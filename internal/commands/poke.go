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
	commands["!poke"] = Command{
		name:        "!poke",
		description: "Trigger dagens spørsmål manuelt (kun admin)",
		emoji:       "👉",
		handler:   handlePoke,
		adminOnly:   true,
	}
}

func handlePoke(s *discordgo.Session, m *discordgo.MessageCreate, bot bot.BotIface) {
	db := bot.GetDatabase()
	log.Printf("Manual daily question trigger requested by %s", m.Author.Username)

	// Support !poke alle
	pokeAlle := false
	args := strings.Fields(m.Content)
	if len(args) > 1 && args[1] == "alle" {
		pokeAlle = true
	}

	question, err := db.GetLeastAskedApprovedQuestion()
	if err != nil {
		log.Printf("Failed to get least asked question: %v", err)
		embed := services.CreateBotEmbed(s, "❌ Feil", "Feil ved henting av spørsmål frå databasen.", 0xff0000)
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
		return
	}

	if question == nil {
		log.Println("No approved questions available")
		embed := services.CreateBotEmbed(s, "😔 Ingen godkjente spørsmål", "Ingen godkjente spørsmål tilgjengelege for no.", 0xffff00)
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
		return
	}

	// Increment the usage count for this question
	err = db.IncrementQuestionUsage(question.ID)
	if err != nil {
		log.Printf("Failed to increment question usage: %v", err)
		embed := services.CreateBotEmbed(s, "❌ Feil", "Feil ved oppdatering av spørsmål-statistikk.", 0xff0000)
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
		return
	}

	// Tag everyone if poke alle, else tag the question submitter
	mention := ""
	if pokeAlle {
		mention = "@everyone"
	} else {
		mention = fmt.Sprintf("<@%s>", question.AuthorID)
	}

	services.SendDailyQuestion(bot, question, mention)

	log.Printf("Daily question manually triggered: %s (asked %d times total)", question.Question, question.TimesAsked+1)
	
	// Get stats for confirmation message
	totalApproved, totalAsked, minAsked, err := db.GetApprovedQuestionStats()
	if err != nil {
		log.Printf("[DATABASE] Failed to get question stats: %v", err)
	} else {
			statsMessage := fmt.Sprintf(`📊 **Statistikk**: %d godkjente spørsmål, %d gonger stilt totalt, minst stilt: %d gonger`, 
				totalApproved, totalAsked+1, minAsked)
			embed := services.CreateBotEmbed(s, "📊 Statistikk", statsMessage, 0x3399ff)
			s.ChannelMessageSendEmbed(bot.GetConfig().Discord.LogChannelID, embed)
		}
}
