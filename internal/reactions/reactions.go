package reactions

import (
	"askeladden/internal/bot"
	"github.com/bwmarrin/discordgo"
)

// Reaction defines the structure for a reaction handler.
type Reaction struct {
	emoji         string
	description   string
	handler       func(s *discordgo.Session, r *discordgo.MessageReactionAdd, b *bot.Bot)
	removeHandler func(s *discordgo.Session, r *discordgo.MessageReactionRemove, b *bot.Bot)
	adminOnly     bool
}

// reactions holds all registered reaction handlers.
var reactions = make(map[string]Reaction)

// Register registers a new reaction handler.
func Register(emoji string, description string, handler func(s *discordgo.Session, r *discordgo.MessageReactionAdd, b *bot.Bot)) Reaction {
	r := Reaction{
		emoji:         emoji,
		description:   description,
		handler:       handler,
		removeHandler: nil,
		adminOnly:     false,
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

// SetRemoveHandler adds a handler for when the reaction is removed.
func (r Reaction) SetRemoveHandler(handler func(s *discordgo.Session, r *discordgo.MessageReactionRemove, b *bot.Bot)) Reaction {
	r.removeHandler = handler
	reactions[r.emoji] = r // Update in map
	return r
}

// MatchAndRunReaction finds and executes a reaction based on its emoji.
func MatchAndRunReaction(emoji string, s *discordgo.Session, r *discordgo.MessageReactionAdd, b *bot.Bot) {
	if reaction, exists := reactions[emoji]; exists {
		reaction.handler(s, r, b)
		return
	}
}

// MatchAndRunReactionRemove finds and executes a reaction removal handler based on its emoji.
func MatchAndRunReactionRemove(emoji string, s *discordgo.Session, r *discordgo.MessageReactionRemove, b *bot.Bot) {
	if reaction, exists := reactions[emoji]; exists && reaction.removeHandler != nil {
		reaction.removeHandler(s, r, b)
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

// InitializeReactions registers all reactions with their configured emojis
func InitializeReactions(b *bot.Bot) {
	// Register starboard reaction
	RegisterStarboardReaction(b)

	// Register question reaction
	RegisterQuestionReaction(b)

	// Register approval reaction (static emoji)
	Register("👍", "Godkjenn eit spørsmål.", handleApprovalReaction).SetAdminOnly()

	// Register reject reaction (static emoji)
	Register("👎", "Avvis eit spørsmål.", handleRejectReaction).SetAdminOnly()
}
