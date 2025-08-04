
package commands

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"roersla.no/askeladden/internal/bot"
	"roersla.no/askeladden/internal/database"
)

func init() {
	commands["!godkjenn"] = Command{
		name:        "!godkjenn",
		description: "Godkjenn eit spørsmål manuelt (kun for opplysarar)",
		emoji:       "✅",
		handler:   Godkjenn,
		aliases:     []string{},
		adminOnly:   true,
	}
}

// Godkjenn handsamer godkjenn-kommandoen
func Godkjenn(s *discordgo.Session, m *discordgo.MessageCreate, bot bot.BotIface) {
	db := bot.GetDatabase()
	// Parse kommandoen for å hente spørsmål ID eller søkeord
	parts := strings.SplitN(m.Content, " ", 2)
	if len(parts) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Bruk: `!godkjenn [spørsmål-ID]` eller `!godkjenn next` for neste ventande spørsmål")
		return
	}

	arg := strings.TrimSpace(parts[1])

	if arg == "alle" {
		err := db.ApproveAllPendingQuestions(m.Author.ID)
		if err != nil {
			log.Printf("Failed to approve all pending questions: %v", err)
			s.ChannelMessageSend(m.ChannelID, "Feil ved godkjenning av alle spørsmål.")
			return
		}
		s.ChannelMessageSend(m.ChannelID, "Alle ventande spørsmål er no godkjent! ✅")
		return
	}

	var question *database.Question
	var err error

	if arg == "next" || arg == "neste" {
		// Get next pending question
		question, err = db.GetPendingQuestion()
		if err != nil {
			log.Printf("Failed to get next pending question: %v", err)
			s.ChannelMessageSend(m.ChannelID, "Feil ved henting av neste spørsmål.")
			return
		}
		if question == nil {
			s.ChannelMessageSend(m.ChannelID, "Ingen ventande spørsmål! 🎉")
			return
		}
	} else {
		// Try to parse as question ID
		questionID, parseErr := strconv.Atoi(arg)
		if parseErr != nil {
			s.ChannelMessageSend(m.ChannelID, "Ugyldig spørsmål-ID. Bruk eit tal eller 'next' for neste ventande spørsmål.")
			return
		}

		// Find pending question by ID
		question, err = db.GetPendingQuestionByID(questionID)
		if err != nil {
			log.Printf("Failed to get pending question by ID %d: %v", questionID, err)
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Kunne ikkje finne ventande spørsmål med ID %d.", questionID))
			return
		}
	}

	// Approve the question
	err = db.ApproveQuestion(question.ID, m.Author.ID)
	if err != nil {
		log.Printf("Failed to approve question: %v", err)
		s.ChannelMessageSend(m.ChannelID, "Feil ved godkjenning av spørsmålet.")
		return
	}

	// Send confirmation
	confirmationEmbed := &discordgo.MessageEmbed{
		Title:       "✅ Spørsmål godkjent!",
		Description: fmt.Sprintf("**Spørsmål:** %s\n**Frå:** %s\n**Godkjent av:** %s", question.Question, question.AuthorName, m.Author.Username),
		Color:       0x00ff00, // Green color
	}
	s.ChannelMessageSendEmbed(m.ChannelID, confirmationEmbed)

	// Notify the original user
	privateChannel, err := s.UserChannelCreate(question.AuthorID)
	if err != nil {
		log.Printf("Failed to create private channel for approval notification: %v", err)
	} else {
		approver, err := s.User(m.Author.ID)
		var approverName string
		if err != nil {
			approverName = "ein opplysar"
		} else {
			approverName = approver.Username
		}

		embed := services.CreateBotEmbed(s, "🎉 Gratulerer! 🎉", fmt.Sprintf("Spørsmålet ditt har blitt godkjent av %s!\n\n**\"%s\"**\n\nDet er no tilgjengeleg for daglege spørsmål! ✨", approverName, question.Question), 0x00ff00)
		s.ChannelMessageSendEmbed(privateChannel.ID, embed)
	}

	log.Printf("Question manually approved by %s: %s", m.Author.Username, question.Question)
}

