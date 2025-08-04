package reactions

import (
	"github.com/bwmarrin/discordgo"
	"roersla.no/askeladden/internal/bot"
)

// Reaction defines the structure for a reaction handler.
type Reaction struct {
	emoji       string
	description string
	handler     func(s *discordgo.Session, r *discordgo.MessageReactionAdd, b bot.BotIface)
	adminOnly   bool
}

// reactions holds all registered reaction handlers.
var reactions = make(map[string]Reaction)

// Register registers a new reaction handler.
func Register(emoji string, description string, handler func(s *discordgo.Session, r *discordgo.MessageReactionAdd, b bot.BotIface)) Reaction {
	r := Reaction{
		emoji:       emoji,
		description: description,
		handler:     handler,
		adminOnly:   false,
	}
	reactions[emoji] = r
	return r
}

// SetAdminOnly marks a reaction as admin-only.
func (r Reaction) SetAdminOnly() Reaction {
	r.adminOnly = true
	reactions[r.emoji] = r // Update in map
	return r
}

// MatchAndRunReaction finds and executes a reaction based on its emoji.
func MatchAndRunReaction(emoji string, s *discordgo.Session, r *discordgo.MessageReactionAdd, b bot.BotIface) {
	if reaction, exists := reactions[emoji]; exists {
		reaction.handler(s, r, b)
		return
	}
}

// IsAdminReaction checks if a reaction is admin-only.
func IsAdminReaction(emoji string) bool {
	if r, exists := reactions[emoji]; exists {
		return r.adminOnly
	}
	return false
}

