package reactions

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"roersla.no/askeladden/internal/bot"
	"roersla.no/askeladden/internal/bot/services"
)
// handleApprovalReaction is registered dynamically in InitializeReactions

func handleApprovalReaction(s *discordgo.Session, r *discordgo.MessageReactionAdd, b *bot.Bot) {
	// Get the question by approval message ID
	question, err := b.Database.GetQuestionByApprovalMessageID(r.MessageID)
	if err != nil {
		log.Printf("Could not find question for approval message %s: %v", r.MessageID, err)
		return
	}

	// Approve the question
	err = b.Database.ApproveQuestion(question.ID, r.UserID)
	if err != nil {
		log.Printf("Failed to approve question: %v", err)
		return
	}

	log.Printf("Question approved by opplysar %s: %s", r.UserID, question.Question)

	// Notify the original user
	approvalService := &services.ApprovalService{Bot: b}
	approvalService.NotifyUserApproval(s, question, r.UserID)

	// Update the approval message to indicate it's been processed
	approvedEmbed := services.CreateBotEmbed(s, "✅ GODKJENT", fmt.Sprintf("**Spørsmål:** %s\n**Frå:** %s\n**Godkjent av:** <@%s>", question.Question, question.AuthorName, r.UserID), 0x00ff00)
	s.ChannelMessageEditEmbed(r.ChannelID, r.MessageID, approvedEmbed)
}

