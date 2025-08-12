package reactions

import (
	"fmt"
	"log"

	"askeladden/internal/bot"
	"askeladden/internal/bot/services"
	"github.com/bwmarrin/discordgo"
)

// handleRejectReaction is registered dynamically in InitializeReactions

func handleRejectReaction(s *discordgo.Session, r *discordgo.MessageReactionAdd, b *bot.Bot) {
	// Get the question by approval message ID
	question, err := b.Database.GetQuestionByApprovalMessageID(r.MessageID)
	if err != nil {
		log.Printf("Could not find question for approval message %s: %v", r.MessageID, err)
		return
	}

	// Reject the question
	err = b.Database.RejectQuestion(question.ID, r.UserID)
	if err != nil {
		log.Printf("Failed to reject question: %v", err)
		return
	}

	log.Printf("Question rejected by opplysar %s: %s", r.UserID, question.Question)

	// Notify the original user
	approvalService := &services.ApprovalService{Bot: b}
	approvalService.NotifyUserRejection(s, question, r.UserID)

	// Update the approval message to indicate it's been processed
	rejectedEmbed := services.CreateBotEmbed(s, "❌ AVVIST", fmt.Sprintf("**Spørsmål:** %s\n**Frå:** %s\n**Avvist av:** <@%s>", question.Question, question.AuthorName, r.UserID), services.EmbedTypeError)
	s.ChannelMessageEditEmbed(r.ChannelID, r.MessageID, rejectedEmbed)
}
