package handlers

import (
	"log"
	"strings"
	"net/http"
	"encoding/json"

	"github.com/bwmarrin/discordgo"
	"roersla.no/askeladden/internal/bot"
	"roersla.no/askeladden/internal/bot/services"
	"roersla.no/askeladden/internal/commands"
)

// Handler struct holds the bot instance and services.
type Handler struct {
	Bot      bot.BotIface
	Services *services.BotServices
}

// New creates a new Handler instance.
func New(b bot.BotIface, s *services.BotServices) *Handler {
	return &Handler{Bot: b, Services: s}
}

// Ready handles the ready event.
func (h *Handler) Ready(s *discordgo.Session, event *discordgo.Ready) {
	log.Println("[BOT] Askeladden is connected and ready.")
	if h.Bot.GetConfig().Discord.LogChannelID != "" {
		s.ChannelMessageSend(h.Bot.GetConfig().Discord.LogChannelID, "Askeladden is online and ready! ✨")
	}
}

// MessageCreate handles new messages.
func (h *Handler) MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if s.State.User.ID == m.Author.ID {
		return
	}

	// Ignore messages that don't start with the prefix
	if !strings.HasPrefix(m.Content, h.Bot.GetConfig().Discord.Prefix) {
		return
	}

	// ---- Custom AI SLOP Warning Section ----
	cfg := h.Bot.GetConfig()
	if cfg.ShowAISlopWarning {
		version := getRemoteVersion() // You need to implement getRemoteVersion()
		warningText := cfg.AISlopWarningText
		if warningText == "" {
			warningText = "Denne botten bruker AI-generert kode ("ai slop"). Inntil alt er manuelt sjekket, kan det vere rare eller feilende ting i svaret!"
		}
		warningEmbed := services.CreateAISlopWarningEmbed(s, version, warningText)
		s.ChannelMessageSendEmbed(m.ChannelID, warningEmbed)
	}
	// ---- End Custom ----

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

// getRemoteVersion fetches the version string from GitHub releases/tag/latest or repo HEAD
func getRemoteVersion() string {
	resp, err := http.Get("https://api.github.com/repos/maalungdom/askeladden/releases/latest")
	if err == nil && resp.StatusCode == 200 {
		defer resp.Body.Close()
		var data struct{ TagName string `json:"tag_name"` }
		if err := json.NewDecoder(resp.Body).Decode(&data); err == nil && data.TagName != "" {
			return data.TagName
		}
	}
	// fallback: try HEAD branch SHA
	resp, err = http.Get("https://api.github.com/repos/maalungdom/askeladden/commits/main")
	if err == nil && resp.StatusCode == 200 {
		defer resp.Body.Close()
		var data struct{ Sha string `json:"sha"` }
		if err := json.NewDecoder(resp.Body).Decode(&data); err == nil && data.Sha != "" {
			return data.Sha[:7] // Shorten
		}
	}
	return "dev"
}

// ReactionAdd handles when a user reacts to a message.
func (h *Handler) ReactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if r.UserID == s.State.User.ID {
		return
	}

	switch r.Emoji.Name {
	case "⭐":
		h.HandleStarReaction(s, r)
	case "👍":
		if r.ChannelID == h.Bot.GetConfig().Approval.QueueChannelID {
			h.HandleApprovalReaction(s, r)
		}
	}
}

