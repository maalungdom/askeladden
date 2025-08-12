package main

import (
	"askeladden/internal/bot"
	"askeladden/internal/bot/services"
	"fmt"
	"log"
	"time"
)

type SchedulerState struct {
	lastActivity    time.Time
	lastDailyPost   time.Time
	timezone        *time.Location
	morningTime     time.Time
	eveningTime     time.Time
	inactivityHours time.Duration
}

// scheduleDailyQuestion sets up the advanced scheduler with timezone and inactivity support
func scheduleDailyQuestion(b *bot.Bot) *time.Ticker {
	if !b.Config.Scheduler.Enabled {
		log.Println("[SCHEDULER] Scheduler is disabled in config")
		return time.NewTicker(time.Hour) // Return a dummy ticker
	}

	// Parse timezone
	timezone, err := time.LoadLocation(b.Config.Scheduler.Timezone)
	if err != nil {
		log.Printf("[SCHEDULER] Invalid timezone '%s', using UTC: %v", b.Config.Scheduler.Timezone, err)
		timezone = time.UTC
	}

	// Parse morning and evening times
	morningTime, err := time.Parse("15:04", b.Config.Scheduler.MorningTime)
	if err != nil {
		log.Printf("[SCHEDULER] Invalid morning time '%s', using 08:00: %v", b.Config.Scheduler.MorningTime, err)
		morningTime, _ = time.Parse("15:04", "08:00")
	}

	eveningTime, err := time.Parse("15:04", b.Config.Scheduler.EveningTime)
	if err != nil {
		log.Printf("[SCHEDULER] Invalid evening time '%s', using 20:00: %v", b.Config.Scheduler.EveningTime, err)
		eveningTime, _ = time.Parse("15:04", "20:00")
	}

	state := &SchedulerState{
		lastActivity:    time.Now(),
		lastDailyPost:   time.Time{}, // Never posted
		timezone:        timezone,
		morningTime:     morningTime,
		eveningTime:     eveningTime,
		inactivityHours: time.Duration(b.Config.Scheduler.InactivityHours) * time.Hour,
	}

	log.Printf("[SCHEDULER] Advanced scheduler started - Timezone: %s, Morning: %s, Evening: %s, Inactivity: %v",
		timezone.String(), b.Config.Scheduler.MorningTime, b.Config.Scheduler.EveningTime, state.inactivityHours)

	// Check every 30 minutes
	ticker := time.NewTicker(30 * time.Minute)
	go func() {
		for {
			select {
			case <-ticker.C:
				checkAndTriggerDailyQuestion(b, state)
			}
		}
	}()

	return ticker
}

// triggerDailyQuestion handles the daily question logic.
func triggerDailyQuestion(b *bot.Bot) {
	// Retrieve least asked approved question
	question, err := b.Database.GetLeastAskedApprovedQuestion()
	if err != nil {
		log.Printf("[SCHEDULER] Failed to retrieve daily question: %v", err)
		return
	}

	if question == nil {
		log.Println("[SCHEDULER] No approved questions available for the day.")
		return
	}

	// Increment usage for the question
	err = b.Database.IncrementQuestionUsage(question.ID)
	if err != nil {
		log.Printf("[SCHEDULER] Failed to update question usage: %v", err)
		return
	}

	// Send the question to the default channel
	if b.Config.Discord.DefaultChannelID != "" {
		services.SendDailyQuestion(b, question, "@pratsam")
		log.Printf("[SCHEDULER] Daily question sent: %s", question.Question)
	} else {
		log.Println("[SCHEDULER] Default channel not configured.")
	}
}

// checkAndTriggerDailyQuestion implements the scheduling logic:
// 1. Post at morning time (08:00)
// 2. Post after 6 hours of inactivity, but only before nighttime (20:00)
// 3. Stop posting once nighttime is reached
func checkAndTriggerDailyQuestion(b *bot.Bot, state *SchedulerState) {
	now := time.Now().In(state.timezone)
	timeSinceLastActivity := now.Sub(state.lastActivity)

	// Get current time components for comparison
	currentTime := time.Date(0, 1, 1, now.Hour(), now.Minute(), 0, 0, time.UTC)
	morningTime := state.morningTime
	eveningTime := state.eveningTime // This is our "nighttime" cutoff

	// Check if we've already posted today
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, state.timezone)
	hasPostedToday := state.lastDailyPost.After(todayStart)

	shouldTrigger := false
	reason := ""

	// Condition 1: It's morning time and we haven't posted today yet
	if currentTime.After(morningTime) && currentTime.Before(morningTime.Add(30*time.Minute)) && !hasPostedToday {
		shouldTrigger = true
		reason = fmt.Sprintf("morning schedule (%s)", b.Config.Scheduler.MorningTime)
	}

	// Condition 2: Inactivity threshold reached, but only if:
	// - It's after morning time
	// - It's before nighttime (evening_time)
	// - We haven't posted today yet
	if timeSinceLastActivity >= state.inactivityHours && !hasPostedToday {
		if currentTime.After(morningTime) && currentTime.Before(eveningTime) {
			shouldTrigger = true
			reason = fmt.Sprintf("inactivity threshold (%v since last activity, before nighttime)", timeSinceLastActivity.Round(time.Minute))
		} else if currentTime.After(eveningTime) {
			// After nighttime - log but don't trigger
			log.Printf("[SCHEDULER] Inactivity threshold reached (%v) but nighttime reached (%s) - waiting until tomorrow morning",
				timeSinceLastActivity.Round(time.Minute), b.Config.Scheduler.EveningTime)
		}
	}

	if shouldTrigger {
		log.Printf("[SCHEDULER] Triggering daily question due to: %s", reason)
		triggerDailyQuestion(b)
		state.lastDailyPost = now

		// Reset activity timer when we post
		state.lastActivity = now
	}
}
