package reactions

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"roersla.no/askeladden/internal/bot"
	"roersla.no/askeladden/internal/bot/services"
)

func init() {
	Register("❓", "Spør eit spørsmål.", handleQuestionReaction)
}

func handleQuestionReaction(s *discordgo.Session, r *discordgo.MessageReactionAdd, bot bot.BotIface) {
	// Fetch the message
	msg, err := s.ChannelMessage(r.ChannelID, r.MessageID)
	if err != nil {
		log.Printf("Failed to fetch message: %v", err)
		return
	}

	// Add the message as a question
	db := bot.GetDatabase()
	questionID, err := db.AddQuestion(msg.Content, msg.Author.ID, msg.Author.Username, msg.ID, msg.ChannelID)
	if err != nil {
		log.Printf("Failed to add question from message: %v", err)
		// Optionally, react with an error emoji
		s.MessageReactionAdd(r.ChannelID, r.MessageID, "❌")
		return
	}

	// Send question to the approval queue channel
	approvalService := &services.ApprovalService{Bot: bot}
	approvalService.PostNewQuestionToApprovalQueue(questionID)

	// React with a success emoji
	s.MessageReactionAdd(r.ChannelID, r.MessageID, "✅")
}

