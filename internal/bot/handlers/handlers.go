package handlers

import (
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"roersla.no/askeladden/internal/bot"
	"roersla.no/askeladden/internal/bot/services"
	"roersla.no/askeladden/internal/commands"
	"roersla.no/askeladden/internal/reactions"
)

// Handler struct holds the bot instance and services.
type Handler struct {
	Bot            *bot.Bot
	Services       *services.BotServices
	warnedChannels map[string]bool
}

// New creates a new Handler instance.
func New(b *bot.Bot) *Handler {
return &Handler{Bot: b}
}

// Ready handles the ready event.
func (h *Handler) Ready(s *discordgo.Session, event *discordgo.Ready) {
	log.Println("[BOT] Askeladden is connected and ready.")
	if h.Bot.Config.Discord.LogChannelID != "" {
		embed := services.CreateBotEmbed(s, "🟢 Online", "Askeladden is online and ready! ✨", 0x00ff00)
		s.ChannelMessageSendEmbed(h.Bot.Config.Discord.LogChannelID, embed)
	}
}

// MessageCreate handles new messages.
func (h *Handler) MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if s.State.User.ID == m.Author.ID {
		return
	}

	// Ignore messages that don't start with the prefix
	if !strings.HasPrefix(m.Content, h.Bot.Config.Discord.Prefix) {
		return
	}

	// Extract command and arguments
	commandWithPrefix := strings.Split(m.Content, " ")[0]

	// Check if the command is admin-only
	if commands.IsAdminCommand(commandWithPrefix) {
		if !h.Services.Approval.UserHasOpplysarRole(s, m.GuildID, m.Author.ID) {
			return // Silently ignore admin commands from non-admins
		}
	}

	// Run the command
	commands.MatchAndRunCommand(commandWithPrefix, s, m, h.Bot)

}


// ReactionAdd handles when a user reacts to a message.
func (h *Handler) ReactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if r.UserID == s.State.User.ID {
		return
	}

	// Check if the reaction is admin-only
	if reactions.IsAdminReaction(r.Emoji.Name) {
		if !h.Services.Approval.UserHasOpplysarRole(s, r.GuildID, r.UserID) {
			return // Silently ignore admin reactions from non-admins
		}
	}

	// Run the reaction handler
	reactions.MatchAndRunReaction(r.Emoji.Name, s, r, h.Bot)
}

// InteractionCreate handles button clicks and other interactions
func (h *Handler) InteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionMessageComponent {
		customID := i.MessageComponentData().CustomID

		if customID == "confirm_clear_database" {
			// Check if the user is an admin
			if !h.Services.Approval.UserHasOpplysarRole(s, i.GuildID, i.Member.User.ID) {
				// Respond to the interaction with an error message
				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Du har ikkje tilgang til å tømme databasen.",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})
				if err != nil {
					log.Printf("Failed to send interaction response: %v", err)
				}
				return
			}

			// Clear the database
			if err := h.Bot.Database.ClearDatabase(); err != nil {
				log.Printf("Failed to clear database: %v", err)
				// Let the user know something went wrong
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Ein feil oppstod under tømming av databasen.",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})
				return
			}

			// Respond to the interaction
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "✅ Databasen har blitt tømt.",
				},
			})
			if err != nil {
				log.Printf("Failed to send interaction response: %v", err)
			}

			// Delete the original confirmation message
			s.ChannelMessageDelete(i.ChannelID, i.Message.ID)
		}
	}
}


