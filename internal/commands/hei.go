package commands

import (
	"github.com/bwmarrin/discordgo"
	"roersla.no/askeladden/internal/bot"
	"roersla.no/askeladden/internal/bot/services"
)

func init() {
	commands["hei"] = Command{
		name:        "hei",
		description: "Sei hei til boten",
		emoji:       "👋",
		handler:     Hei,
		aliases:     []string{"hallo"},
	}
}

// Hei handsamer hei-kommandoen

func Hei(s *discordgo.Session, m *discordgo.MessageCreate, bot *bot.Bot) {
	embed := services.CreateBotEmbed(s, "Hei der! 👋", "Eg er Askeladden, laga av rørsla!", 0x0099ff)
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}
