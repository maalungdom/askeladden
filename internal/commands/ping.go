package commands

import (
	"askeladden/internal/bot"
	"askeladden/internal/bot/services"
	"github.com/bwmarrin/discordgo"
)

func init() {
	commands["ping"] = Command{
		name:        "ping",
		description: "Sjekk om boten svarar",
		emoji:       "🏓",
		handler:     Ping,
	}
}

// Ping handsamer ping-kommandoen
//--------------------------------------------------------------------------------

func Ping(s *discordgo.Session, m *discordgo.MessageCreate, bot *bot.Bot) {
	embed := services.CreateBotEmbed(s, "Pong! 🏓", "Bot er oppe og svarar.", services.EmbedTypeSuccess)
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}
