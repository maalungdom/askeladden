package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"askeladden/internal/bot"
	"askeladden/internal/bot/services"
)

func init() {
	commands["info"] = Command{
		name:        "info",
		description: "Syn opplysingar om boten",
		emoji:       "📊",
		handler:     Info,
	}
}

// Info handsamer info-kommandoen
//--------------------------------------------------------------------------------
func Info(s *discordgo.Session, m *discordgo.MessageCreate, bot *bot.Bot) {
	guildCount := len(s.State.Guilds)
	infoText := fmt.Sprintf("**Om Askeladden:**\n" +
		"🤖 Ein norsk Discord-bot\n" +
		"💻 Skriven i Go\n" +
		"🏠 Laga av rørsla\n" +
		"🖥️ Køyrer på %d servarar\n" +
		"🤖 Bot-brukar: %s#%s", 
		guildCount, s.State.User.Username, s.State.User.Discriminator)
	embed := services.CreateBotEmbed(s, "📊 Om Askeladden", infoText, 0x3399ff)
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}
