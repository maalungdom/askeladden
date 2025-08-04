package main

import (
	"log"
	"time"
	"roersla.no/askeladden/internal/bot"
	"roersla.no/askeladden/internal/bot/services"
)

// scheduleDailyQuestion sets up the daily trigger for the question.
func scheduleDailyQuestion(b bot.BotIface) *time.Ticker {
	ticker := time.NewTicker(24 * time.Hour)
	go func() {
		for {
			select {
			case <-ticker.C:
				triggerDailyQuestion(b)
			}
		}
	}()
	return ticker
}

// triggerDailyQuestion handles the daily question logic.
func triggerDailyQuestion(b bot.BotIface) {
	// Retrieve least asked approved question
	question, err := b.GetDatabase().GetLeastAskedApprovedQuestion()
	if err != nil {
		log.Printf("[SCHEDULER] Failed to retrieve daily question: %v", err)
		return
	}

	if question == nil {
		log.Println("[SCHEDULER] No approved questions available for the day.")
		return
	}

	// Increment usage for the question
	err = b.GetDatabase().IncrementQuestionUsage(question.ID)
	if err != nil {
		log.Printf("[SCHEDULER] Failed to update question usage: %v", err)
		return
	}

	// Send the question to the default channel
	if b.GetConfig().Discord.DefaultChannelID != "" {
		services.SendDailyQuestion(b, question, "@everyone")
		log.Printf("[SCHEDULER] Daily question sent: %s", question.Question)
	} else {
		log.Println("[SCHEDULER] Default channel not configured.")
	}
}
