package handlers

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

// HandleApprovalReaction handles reactions in the approval queue channel.
func (h *Handler) HandleApprovalReaction(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	// Check if user has the opplysar role
	if !h.Services.Approval.UserHasOpplysarRole(s, r.GuildID, r.UserID) {
		log.Printf("User %s does not have opplysar role, ignoring reaction", r.UserID)
		return
	}

	// Get the question by approval message ID
	question, err := h.Bot.GetDatabase().GetQuestionByApprovalMessageID(r.MessageID)
	if err != nil {
		log.Printf("Could not find question for approval message %s: %v", r.MessageID, err)
		return
	}

	// Approve the question
	err = h.Bot.GetDatabase().ApproveQuestion(question.ID, r.UserID)
	if err != nil {
		log.Printf("Failed to approve question: %v", err)
		return
	}

	log.Printf("Question approved by opplysar %s: %s", r.UserID, question.Question)

	// Notify the original user
	h.Services.Approval.NotifyUserApproval(s, question, r.UserID)

	// Update the approval message to indicate it's been processed
	s.ChannelMessageEdit(r.ChannelID, r.MessageID,
		fmt.Sprintf("✅ **GODKJENT**\\n\\n**Spørsmål:** %s\\n**Frå:** %s\\n**Godkjent av:** <@%s>",
			question.Question, question.AuthorName, r.UserID))
}
