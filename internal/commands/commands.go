package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"roersla.no/askeladden/internal/bot"
)

// Command defines the structure for a command
type Command struct {
	name        string
	description string
	emoji       string
	handler     func(s *discordgo.Session, m *discordgo.MessageCreate, bot bot.BotIface)
	aliases     []string
	adminOnly   bool
}

// commands holds all the registered commands
var commands = make(map[string]Command)

// MatchAndRunCommand finds and executes a command based on its name or alias.
func MatchAndRunCommand(input string, s *discordgo.Session, m *discordgo.MessageCreate, bot bot.BotIface) {
	// `input` is the command with prefix, e.g., "!spør"
	if cmd, exists := commands[input]; exists {
		cmd.handler(s, m, bot)
		return
	}

	// Check aliases
	commandWithoutPrefix := strings.TrimPrefix(input, bot.GetConfig().Discord.Prefix)
	for _, cmd := range commands {
		for _, alias := range cmd.aliases {
			if alias == commandWithoutPrefix {
					cmd.handler(s, m, bot)
				return
			}
		}
	}
}

// IsAdminCommand checks if a command is admin-only
func IsAdminCommand(commandName string) bool {
	if cmd, exists := commands[commandName]; exists {
		return cmd.adminOnly
	}
	return false
}

// GetHelpText generates the help text for all commands.
func GetHelpText() string {
	var helpText strings.Builder
	helpText.WriteString("**Askeladden - Kommandoer:**\n")

	for _, cmd := range commands {
		helpText.WriteString(fmt.Sprintf("%s `%s` - %s\n", cmd.emoji, cmd.name, cmd.description))
	}

	return strings.TrimSpace(helpText.String())
}

// ListCommands lists commands based on admin status
func ListCommands(isAdmin bool) *discordgo.MessageEmbed {
	var generalCommands, adminCommands strings.Builder

	for _, cmd := range commands {
		commandLine := fmt.Sprintf("%s `%s` - %s\n", cmd.emoji, cmd.name, cmd.description)
		if cmd.adminOnly {
			adminCommands.WriteString(commandLine)
		} else {
			generalCommands.WriteString(commandLine)
		}
	}

	embed := &discordgo.MessageEmbed{
		Title: "Askeladden - Kommandoer",
		Color: 0x0099ff, // Blue color
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Generelle kommandoer",
				Value: generalCommands.String(),
			},
		},
	}

	if isAdmin && adminCommands.Len() > 0 {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  "Admin-kommandoer",
			Value: adminCommands.String(),
		})
	}

	return embed
}
