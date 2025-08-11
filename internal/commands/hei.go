package commands

import (
	"github.com/bwmarrin/discordgo"
	"askeladden/internal/bot"
	"askeladden/internal/bot/services"
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
	embed := services.CreateBotEmbed(s, "Heisann! 👋", "Eg er Askeladden, laga av rørsla!", services.EmbedTypeInfo)
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}
