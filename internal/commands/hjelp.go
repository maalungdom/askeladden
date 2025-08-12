package commands

import (
	"log"

	"askeladden/internal/bot"
	"askeladden/internal/bot/services"
	"github.com/bwmarrin/discordgo"
)

func init() {
	commands["hjelp"] = Command{
		name:        "hjelp",
		description: "Syn denne hjelpemeldinga",
		emoji:       "❓",
		handler:     Hjelp,
		aliases:     []string{"help", "h"},
	}
}

// Hjelp handsamer hjelp-kommandoen
// --------------------------------------------------------------------------------
func Hjelp(s *discordgo.Session, m *discordgo.MessageCreate, bot *bot.Bot) {
	// Check if user has admin role (we need to implement role checking here)
	// For now, let's use a placeholder implementation
	isAdmin := false

	// Try to get guild member to check roles
	if m.GuildID != "" {
		member, err := s.GuildMember(m.GuildID, m.Author.ID)
		if err == nil {
			// Check for opplysar role (need to get the role ID from config)
			// This is a placeholder - we'll need to pass config or implement differently
			for _, roleID := range member.Roles {
				if roleID == bot.Config.Approval.OpplysarRoleID { // Use config for role ID
					isAdmin = true
					break
				}
			}
		} else {
			log.Printf("Failed to get guild member for role check: %v", err)
		}
	}

	helpEmbed := ListCommands(isAdmin)
	helpBotEmbed := services.CreateBotEmbed(s, helpEmbed.Title, helpEmbed.Description, services.EmbedTypePrimary)
	helpBotEmbed.Fields = helpEmbed.Fields
	if helpEmbed.Footer != nil {
		helpBotEmbed.Footer = helpEmbed.Footer
	}
	s.ChannelMessageSendEmbed(m.ChannelID, helpBotEmbed)
}
