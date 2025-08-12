package reactions

import (
	"fmt"
	"log"
	"strings"

	"askeladden/internal/bot"
	"askeladden/internal/bot/services"
	"askeladden/internal/database"
	"askeladden/internal/permissions"
	"github.com/bwmarrin/discordgo"
)

// handleBannedWordApprovalReaction handles reactions for banned word approval process.
func handleBannedWordApprovalReaction(s *discordgo.Session, r *discordgo.MessageReactionAdd, b *bot.Bot) {
	bannedWord, err := b.Database.GetBannedWordByApprovalMessageID(r.MessageID)
	if err != nil {
		log.Printf("Could not find banned word for approval message %s: %v", r.MessageID, err)
		return
	}

	// Skip if the reacting user is a bot
	if r.Member != nil && r.Member.User.Bot {
		return
	}

	// Initialize permission manager
	permManager := permissions.NewPermissionManager(b.Config)

	// Check if the user has any required role
	userRole := permManager.GetUserRole(s, r.GuildID, r.UserID)
	if userRole == permissions.RoleNone {
		log.Printf("User %s does not have required roles for approval", r.UserID)
		return
	}

	// Check combined approval state from all reactions
	approvalState, err := permManager.CheckCombinedApproval(s, r.ChannelID, r.MessageID, r.Emoji.Name)
	if err != nil {
		log.Printf("Failed to check combined approval: %v", err)
		return
	}

	// Update the embed with current approval state
	var embedColor int
	var embedTitle string
	var embedDescription string

	if approvalState.IsFullyApproved() {
		// Both roles have approved - finalize the approval
		err = b.Database.ApproveBannedWordCombined(int(bannedWord.ID), approvalState.OpplysarApprovers, approvalState.RettskrivarApprovers)
		if err != nil {
			log.Printf("Failed to approve banned word: %v", err)
			return
		}

		// Create forum thread for discussion
		approvalService := services.ApprovalService{Bot: b}
		originalMessageID := ""
		originalChannelID := ""
		if bannedWord.OriginalMessageID != nil {
			// Parse the combined channel|message format
			parts := strings.Split(*bannedWord.OriginalMessageID, "|")
			if len(parts) == 2 {
				originalChannelID = parts[0]
				originalMessageID = parts[1]
			} else {
				// Fallback for old format (just message ID)
				originalMessageID = *bannedWord.OriginalMessageID
				originalChannelID = ""
			}
		}
		thread := approvalService.PostBannedWordReport(s, []string{bannedWord.Word}, bannedWord.AuthorID, r.GuildID, originalChannelID, originalMessageID)
		if thread != nil {
			log.Printf("Created forum thread %s for banned word %s", thread.ID, bannedWord.Word)
			// Update the banned word with the forum thread ID
			b.Database.UpdateBannedWordForumThreadID(int(bannedWord.ID), thread.ID)
		}

		embedColor = services.ColorSuccess // Green
		embedTitle = bannedWord.Word
		embedDescription = approvalState.GetApprovalSummary(s)
		log.Printf("Banned word %s fully approved by combined roles", bannedWord.Word)
	} else {
		// Partial approval - update status but don't finalize
		embedColor = services.ColorWarning // Yellow
		embedTitle = bannedWord.Word
		embedDescription = approvalState.GetApprovalSummary(s)
		log.Printf("Banned word %s partially approved - waiting for additional roles", bannedWord.Word)
	}

	// Get hammer user info for embed author
	hammerUser, err := s.User(bannedWord.AuthorID)
	var authorName, avatarURL string
	if err == nil {
		authorName = hammerUser.Username
		avatarURL = hammerUser.AvatarURL("")
	} else {
		authorName = bannedWord.AuthorName
		avatarURL = ""
	}

	// Create updated embed with hammer user as author
	updatedEmbed := &discordgo.MessageEmbed{
		Title:       embedTitle,
		Description: embedDescription,
		Color:       embedColor,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    authorName,
			IconURL: avatarURL,
		},
	}
	s.ChannelMessageEditEmbed(r.ChannelID, r.MessageID, updatedEmbed)
}

func handleApprovalReaction(s *discordgo.Session, r *discordgo.MessageReactionAdd, b *bot.Bot) {
	// Try to find a banned word first
	_, err := b.Database.GetBannedWordByApprovalMessageID(r.MessageID)
	if err == nil {
		// Found a banned word - handle banned word approval
		handleBannedWordApprovalReaction(s, r, b)
		return
	}

	// If no banned word found, try to find a question
	question, err := b.Database.GetQuestionByApprovalMessageID(r.MessageID)
	if err != nil {
		log.Printf("Could not find question or banned word for approval message %s: %v", r.MessageID, err)
		return
	}

	// Handle question approval
	handleQuestionApprovalReaction(s, r, b, question)
}

func handleQuestionApprovalReaction(s *discordgo.Session, r *discordgo.MessageReactionAdd, b *bot.Bot, question *database.Question) {
	// Approve the question
	err := b.Database.ApproveQuestion(question.ID, r.UserID)
	if err != nil {
		log.Printf("Failed to approve question: %v", err)
		return
	}

	log.Printf("Question approved by opplysar %s: %s", r.UserID, question.Question)

	// Notify the original user
	approvalService := &services.ApprovalService{Bot: b}
	approvalService.NotifyUserApproval(s, question, r.UserID)

	// Get the approver's info for the approval message
	approver, err := s.User(r.UserID)
	var approverName string
	if err == nil {
		approverName = approver.Username
	} else {
		approverName = "Ukjend"
	}

	// Get the question author's info for embed author
	questionAuthor, err := s.User(question.AuthorID)
	var authorName, avatarURL string
	if err == nil {
		authorName = questionAuthor.Username
		avatarURL = questionAuthor.AvatarURL("")
	} else {
		authorName = question.AuthorName
		avatarURL = ""
	}

	// Update the approval message to match banned word format
	approvedEmbed := &discordgo.MessageEmbed{
		Title:       question.Question,
		Description: fmt.Sprintf("🧘‍♀️ Opplysar-godkjenning: %s", approverName),
		Color:       services.ColorSuccess, // Green
		Author: &discordgo.MessageEmbedAuthor{
			Name:    authorName,
			IconURL: avatarURL,
		},
	}
	s.ChannelMessageEditEmbed(r.ChannelID, r.MessageID, approvedEmbed)
}
