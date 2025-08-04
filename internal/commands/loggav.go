
package commands

import (
	"os"

	"github.com/bwmarrin/discordgo"
	"roersla.no/askeladden/internal/bot"
)

func init() {
	commands["!loggav"] = Command{
		name:        "!loggav",
		description: "Loggar av boten og avsluttar programmet (kun for admin)",
		emoji:       "👋",
		handler:     Loggav,
		adminOnly:   true,
	}
}

// Loggav handsamar loggav-kommandoen
func Loggav(s *discordgo.Session, m *discordgo.MessageCreate, bot bot.BotIface) {
	bot.Stop()
	os.Exit(0)
}

