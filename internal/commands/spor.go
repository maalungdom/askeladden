package commands

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"askeladden/internal/bot"
	"askeladden/internal/bot/services"
)

func init() {
	commands["spør"] = Command{
		name:        "spør",
		description: "Legg til eit spørsmål for daglege spørsmål",
		emoji:       "❓",
		handler:     Spor,
		aliases:     []string{"spor"},
	}
}

// Spor handsamer spør-kommandoen
func Spor(s *discordgo.Session, m *discordgo.MessageCreate, bot *bot.Bot) {
	db := bot.Database
	// Parse kommandoen for å hente spørsmålet
	parts := strings.SplitN(m.Content, " ", 2)
	if len(parts) < 2 {
			embed := services.CreateBotEmbed(s, "❓ Feil", "Du må skrive eit spørsmål! Døme: `!spør Kva er din yndlingsmat?`", services.EmbedTypeError)
			s.ChannelMessageSendEmbed(m.ChannelID, embed)
			return
	}

	question := strings.TrimSpace(parts[1])
	if question == "" {
			embed := services.CreateBotEmbed(s, "❓ Feil", "Spørsmålet kan ikkje vere tomt!", services.EmbedTypeError)
			s.ChannelMessageSendEmbed(m.ChannelID, embed)
			return
	}

	// Send bekreftelse til brukaren
	embed := services.CreateBotEmbed(s, "📝 Spørsmål motteke!", fmt.Sprintf("Takk! Spørsmålet ditt er sendt til godkjenning: \"%s\"\n\n*Du får ei melding når det vert godkjent av opplysarane våre! ✨*", question), services.EmbedTypeInfo)
	response, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
	if err != nil {
		log.Printf("Feil ved sending av melding: %v", err)
		return
	}

	// Lagre spørsmålet i databasen med meldings-ID
	questionID, err := db.AddQuestion(question, m.Author.ID, m.Author.Username, response.ID, m.ChannelID)
	if err != nil {
		log.Printf("Feil ved lagring av spørsmål: %v", err)
			embed := services.CreateBotEmbed(s, "❌ Feil", "Det oppstod ein feil ved lagring av spørsmålet.", services.EmbedTypeError)
			s.ChannelMessageSendEmbed(m.ChannelID, embed)
			return
	}

	// Send DM bekreftelse til brukaren
	privateChannel, err := s.UserChannelCreate(m.Author.ID)
	if err == nil {
	embed := services.CreateBotEmbed(s, "📝 Spørsmål motteke!", fmt.Sprintf("Hei %s! 👋\n\nSpørsmålet ditt er vorte sendt til godkjenning:\n\n**\"%s\"**\n\nDu får bod når det vert godkjent av opplysarane våre! 📝✨", m.Author.Username, question), services.EmbedTypeInfo)
		s.ChannelMessageSendEmbed(privateChannel.ID, embed)
	}

	// Send question to the approval queue channel
	approvalService := &services.ApprovalService{Bot: bot}
	approvalService.PostNewQuestionToApprovalQueue(questionID)
}
