package commands

import (
	"github.com/bwmarrin/discordgo"
	"roersla.no/askeladden/internal/bot"
	"roersla.no/askeladden/internal/bot/services"
)

func init() {
	commands["!ping"] = Command{
		name:        "!ping",
		description: "Sjekk om boten svarar",
		emoji:       "🏓",
		handler:     Ping,
	}
}

// Ping handsamer ping-kommandoen
//--------------------------------------------------------------------------------

func Ping(s *discordgo.Session, m *discordgo.MessageCreate, bot bot.BotIface) {
	embed := services.CreateBotEmbed(s, "Pong! 🏓", "Bot er oppe og svarar.", 0x00ff00)
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}
