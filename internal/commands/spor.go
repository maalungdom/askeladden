package commands

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"roersla.no/askeladden/internal/bot"
	"roersla.no/askeladden/internal/bot/services"
)

func init() {
	commands["!spør"] = Command{
		name:        "!spør",
		description: "Legg til eit spørsmål for daglege spørsmål",
		emoji:       "❓",
		handler:     Spor,
		aliases:     []string{"spor"},
	}
}

// Spor handsamer spør-kommandoen
func Spor(s *discordgo.Session, m *discordgo.MessageCreate, bot bot.BotIface) {
	db := bot.GetDatabase()
	// Parse kommandoen for å hente spørsmålet
	parts := strings.SplitN(m.Content, " ", 2)
	if len(parts) < 2 {
			embed := services.CreateBotEmbed(s, "❓ Feil", "Du må skrive eit spørsmål! Eksempel: `!spør Kva er din yndlingsmat?`", 0xff0000)
			s.ChannelMessageSendEmbed(m.ChannelID, embed)
			return
	}

	question := strings.TrimSpace(parts[1])
	if question == "" {
			embed := services.CreateBotEmbed(s, "❓ Feil", "Spørsmålet kan ikkje vere tomt!", 0xff0000)
			s.ChannelMessageSendEmbed(m.ChannelID, embed)
			return
	}

	// Send bekreftelse til brukaren
	embed := services.CreateBotEmbed(s, "📝 Spørsmål motteke!", fmt.Sprintf("Takk! Spørsmålet ditt er sendt til godkjenning: \"%s\"\n\n*Du vil få ei melding når det blir godkjent av opplysarane våre! ✨*", question), 0x0099ff)
	response, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
	if err != nil {
		log.Printf("Feil ved sending av melding: %v", err)
		return
	}

	// Lagre spørsmålet i databasen med meldings-ID
	questionID, err := db.AddQuestion(question, m.Author.ID, m.Author.Username, response.ID, m.ChannelID)
	if err != nil {
		log.Printf("Feil ved lagring av spørsmål: %v", err)
			embed := services.CreateBotEmbed(s, "❌ Feil", "Det oppstod ein feil ved lagring av spørsmålet.", 0xff0000)
			s.ChannelMessageSendEmbed(m.ChannelID, embed)
			return
	}

	// Send DM bekreftelse til brukaren
	privateChannel, err := s.UserChannelCreate(m.Author.ID)
	if err == nil {
		embed := services.CreateBotEmbed(s, "📝 Spørsmål motteke!", fmt.Sprintf("Hei %s! 👋\n\nSpørsmålet ditt har blitt sendt til godkjenning:\n\n**\"%s\"**\n\nDu vil få beskjed når det blir godkjent av opplysarane våre! 📝✨", m.Author.Username, question), 0x0099ff)
		s.ChannelMessageSendEmbed(privateChannel.ID, embed)
	}

	// Send question to the approval queue channel
	approvalService := &services.ApprovalService{Bot: bot}
	approvalService.PostNewQuestionToApprovalQueue(questionID)
}
