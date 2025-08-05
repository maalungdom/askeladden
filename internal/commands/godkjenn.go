package commands
import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"askeladden/internal/bot"
	"askeladden/internal/database"
	"askeladden/internal/bot/services"
)
func init() {
	commands["godkjenn"] = Command{
		name:        "godkjenn",
		description: "Godkjenn eit spørsmål for hand (kun for opplysarar)",
		emoji:       "✅",
		handler:   Godkjenn,
		aliases:     []string{},
		adminOnly:   true,
	}
}

// Godkjenn handsamer godkjenn-kommandoen
func Godkjenn(s *discordgo.Session, m *discordgo.MessageCreate, bot *bot.Bot) {
	db := bot.Database
	// Parse kommandoen for å hente spørsmål ID eller søkeord
	parts := strings.SplitN(m.Content, " ", 2)
	if len(parts) < 2 {
		embed := services.CreateBotEmbed(s, "❓ Feil", "Bruk: `!godkjenn [spørsmål-ID]` eller `!godkjenn next` for neste ventande spørsmål", 0xff0000)
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
		return
	}

	arg := strings.TrimSpace(parts[1])

	if arg == "alle" {
		// TODO: Implement ApproveAllPendingQuestions functionality
		embed := services.CreateBotEmbed(s, "⚠️ Ikkje implementert", "Godkjenning av alle spørsmål er ikkje enno implementert.", 0xffa500)
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
		return
	}

	var question *database.Question
	var err error

	if arg == "next" || arg == "neste" {
		// Get next pending question
		question, err = db.GetPendingQuestion()
		if err != nil {
			log.Printf("Failed to get next pending question: %v", err)
			embed := services.CreateBotEmbed(s, "❌ Feil", "Mislukkast i å hente neste spørsmål.", 0xff0000)
			s.ChannelMessageSendEmbed(m.ChannelID, embed)
			
			return
		}
		if question == nil {
			embed := services.CreateBotEmbed(s, "🎉 Ingen ventande spørsmål!", "", 0x00ff00)
			s.ChannelMessageSendEmbed(m.ChannelID, embed)
			
			return
		}
	} else {
		// Try to parse as question ID
		questionID, parseErr := strconv.Atoi(arg)
		if parseErr != nil {
			embed := services.CreateBotEmbed(s, "❓ Feil", "Ugyldig spørsmål-ID. Bruk eit tal eller «next» for neste ventande spørsmål.", 0xff0000)
			s.ChannelMessageSendEmbed(m.ChannelID, embed)
			
			return
		}

		// Find pending question by ID
		question, err = db.GetPendingQuestionByID(questionID)
		if err != nil {
			log.Printf("Failed to get pending question by ID %d: %v", questionID, err)
			embed := services.CreateBotEmbed(s, "❌ Feil", fmt.Sprintf("Kunne ikkje finne ventande spørsmål med ID %d.", questionID), 0xff0000)
			s.ChannelMessageSendEmbed(m.ChannelID, embed)
			
			return
		}
	}

	// Approve the question
	err = db.ApproveQuestion(question.ID, m.Author.ID)
	if err != nil {
		log.Printf("Failed to approve question: %v", err)
		embed := services.CreateBotEmbed(s, "❌ Feil", "Feil ved godkjenning av spørsmålet.", 0xff0000)
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
		
		return
	}

	// Send confirmation
	confirmationEmbed := services.CreateBotEmbed(s, "✅ Spørsmål godkjent!", fmt.Sprintf("**Spørsmål:** %s\n**Frå:** %s\n**Godkjent av:** %s", question.Question, question.AuthorName, m.Author.Username), 0x00ff00)
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

		embed := services.CreateBotEmbed(s, "🎉 Gratulerer! 🎉", fmt.Sprintf("Spørsmålet ditt er vorte godkjent av %s!\n\n**\"%s\"**\n\nDet er no tilgjengeleg for daglege spørsmål! ✨", approverName, question.Question), 0x00ff00)
		s.ChannelMessageSendEmbed(privateChannel.ID, embed)
		
	}

	log.Printf("Question manually approved by %s: %s", m.Author.Username, question.Question)
}

