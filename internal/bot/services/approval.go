package services

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"roersla.no/askeladden/internal/bot"
	"roersla.no/askeladden/internal/database"
)

// ApprovalService handles logic for question approval.
type ApprovalService struct {
	Bot *bot.Bot
}

// PostNewQuestionToApprovalQueue posts a newly created question to the approval queue.
func (s *ApprovalService) PostNewQuestionToApprovalQueue(questionID int64) {
	// Get the question from database
	question, err := s.Bot.Database.GetPendingQuestionByID(int(questionID))
	if err != nil {
		log.Printf("Failed to get question for approval queue posting: %v", err)
		return
	}

	s.postToApprovalQueue(s.Bot.Session, question)
}

// postToApprovalQueue posts a question to the approval queue channel.
func (s *ApprovalService) postToApprovalQueue(session *discordgo.Session, question *database.Question) {
	if s.Bot.Config.Approval.QueueChannelID == "" {
		log.Println("Approval queue channel not configured")
		return
	}

	approvalEmbed := CreateBotEmbed(session, "📝 Nytt spørsmål til godkjenning", fmt.Sprintf("**Spørsmål:** %s\n**Frå:** %s\n\nReager med 👍 for å godkjenne!", question.Question, question.AuthorName), 0xffa500)

	approvalMessage, err := session.ChannelMessageSendEmbed(s.Bot.Config.Approval.QueueChannelID, approvalEmbed)
	if err != nil {
		log.Printf("Failed to post to approval queue: %v", err)
		return
	}

	// Add thumbs up reaction
	session.MessageReactionAdd(s.Bot.Config.Approval.QueueChannelID, approvalMessage.ID, "👍")

	// Update the database with the approval message ID
	err = s.Bot.Database.UpdateApprovalMessageID(question.ID, approvalMessage.ID)
	if err != nil {
		log.Printf("Failed to update approval message ID: %v", err)
	}
}

// UserHasOpplysarRole checks if a user has the opplysar role.
func (s *ApprovalService) UserHasOpplysarRole(session *discordgo.Session, guildID, userID string) bool {
	if s.Bot.Config.Approval.OpplysarRoleID == "" {
		return false
	}

	member, err := session.GuildMember(guildID, userID)
	if err != nil {
		log.Printf("Failed to get guild member: %v", err)
		return false
	}

	for _, roleID := range member.Roles {
		if roleID == s.Bot.Config.Approval.OpplysarRoleID {
			return true
		}
	}

	return false
}

// NotifyUserApproval notifies the user that their question was approved.
func (s *ApprovalService) NotifyUserApproval(session *discordgo.Session, question *database.Question, approverID string) {
	privateChannel, err := session.UserChannelCreate(question.AuthorID)
	if err != nil {
		log.Printf("Failed to create private channel for approval notification: %v", err)
		return
	}

	approver, err := session.User(approverID)
	var approverName string
	if err != nil {
		approverName = "ein opplysar"
	} else {
		approverName = approver.Username
	}

	embed := CreateBotEmbed(session, "🎉 Gratulerer! 🎉", fmt.Sprintf("Spørsmålet ditt har blitt godkjent av %s!\n\n**\"%s\"**\n\nDet er no tilgjengeleg for daglege spørsmål! ✨", approverName, question.Question), 0x00ff00)
	_, err = session.ChannelMessageSendEmbed(privateChannel.ID, embed)
	if err != nil {
		log.Printf("Failed to send approval notification to user: %v", err)
	}
}
