package reactions

import (
	"log"

	"askeladden/internal/bot"
	"askeladden/internal/bot/services"
	"github.com/bwmarrin/discordgo"
)

// RegisterStarboardReaction registers the starboard reaction with the configured emoji
func RegisterStarboardReaction(b *bot.Bot) {
	emoji := b.Config.Starboard.Emoji
	Register(emoji, "Legg til ei melding på stjernebrettet", handleStarReaction)
}

func handleStarReaction(s *discordgo.Session, r *discordgo.MessageReactionAdd, b *bot.Bot) {
	if r.UserID == s.State.User.ID { // Ignore bot's own reactions
		return
	}
	// Don't process reactions in the starboard channel itself
	if r.ChannelID == b.Config.Starboard.ChannelID {
		return
	}

	// Fetch message
	msg, err := s.ChannelMessage(r.ChannelID, r.MessageID)
	if err != nil {
		log.Printf("Error fetching message: %v", err)
		return
	}

	// Count total starboard reactions
	stars := 0
	for _, reaction := range msg.Reactions {
		if reaction.Emoji.Name == b.Config.Starboard.Emoji {
			stars = reaction.Count
			break
		}
	}

	// Log for debugging
	log.Printf("Message %s in channel %s has %d stars (threshold: %d)", r.MessageID, r.ChannelID, stars, b.Config.Starboard.Threshold)

	if stars >= b.Config.Starboard.Threshold {
		// Send message to starboard channel using new embed system
		embed := services.CreateStarboardEmbed(msg, stars, getChannelName(s, r.ChannelID), b.Config.Starboard.Emoji, r.GuildID)

		_, err := s.ChannelMessageSendEmbed(b.Config.Starboard.ChannelID, embed)
		if err != nil {
			log.Printf("Error sending starboard message: %v", err)
		}
	}
}

// getChannelName returns the channel name for a given channel ID
func getChannelName(s *discordgo.Session, channelID string) string {
	channel, err := s.Channel(channelID)
	if err != nil {
		return "unknown-channel"
	}
	return channel.Name
}
