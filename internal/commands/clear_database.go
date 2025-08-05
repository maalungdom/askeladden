package commands

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"roersla.no/askeladden/internal/bot"
)

func init() {
	commands["tøm-db"] = Command{
		name:        "tøm-db",
		description: "Tømmer databasen for alle spørsmål. Dette kan ikkje angrast.",
		emoji:       "🗑️",
		handler:     ClearDatabase,
		adminOnly:   true,
	}
}

// ClearDatabase handles the command to clear the database
func ClearDatabase(s *discordgo.Session, m *discordgo.MessageCreate, bot *bot.Bot) {
	// Send a confirmation message with a button
	confirmationEmbed := &discordgo.MessageEmbed{
		Title:       "🗑️ Stadfesting av databasetømming",
		Description: "Er du sikker på at du vil slette **alle** data frå databasen? Dette kan ikkje angrast.",
		Color:       0xff0000, // Red color
	}

	msg, err := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Embed: confirmationEmbed,
		Components: []discordgo.MessageComponent{
			&discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					&discordgo.Button{
						Label:    "Ja, slett alt!",
						Style:    discordgo.DangerButton,
						CustomID: "confirm_clear_database",
					},
				},
			},
		},
	})
	if err != nil {
		log.Printf("Failed to send confirmation message: %v", err)
		return
	}

	// Store the message ID to verify the button click later
	// This is a simplified example; a more robust solution would store this mapping
	log.Printf("Sent confirmation message with ID: %s", msg.ID)
}

