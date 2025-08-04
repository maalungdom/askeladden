package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"roersla.no/askeladden/internal/bot"
)

func init() {
	commands["!info"] = Command{
		name:        "!info",
		description: "Vis informasjon om boten",
		emoji:       "📊",
		handler:     Info,
	}
}

// Info handsamer info-kommandoen
//--------------------------------------------------------------------------------
func Info(s *discordgo.Session, m *discordgo.MessageCreate, bot bot.BotIface) {
	guildCount := len(s.State.Guilds)
	infoText := fmt.Sprintf("**Om Askeladden:**\n" +
		"🤖 Ein norsk Discord-bot\n" +
		"💻 Skrive i Go\n" +
		"🏠 Laga av rørsla\n" +
		"🖥️ Køyrer på %d servarar\n" +
		"🤖 Bot-brukar: %s#%s", 
		guildCount, s.State.User.Username, s.State.User.Discriminator)
	embed := services.CreateBotEmbed(s, "📊 Om Askeladden", infoText, 0x3399ff)
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}
