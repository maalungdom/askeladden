package reactions

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"roersla.no/askeladden/internal/bot"
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
		// Send message to starboard channel
		embed := &discordgo.MessageEmbed{
			Author: &discordgo.MessageEmbedAuthor{
				Name:    msg.Author.Username,
				IconURL: msg.Author.AvatarURL(""),
			},
			Description: msg.Content,
			Footer: &discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("%s %d | #%s", b.Config.Starboard.Emoji, stars, getChannelName(s, r.ChannelID)),
			},
			Timestamp: string(msg.Timestamp.Format(time.RFC3339)),
			Color:     0xFFD700, // Gold color
		}

		// Add original message link
		embed.Fields = []*discordgo.MessageEmbedField{
			{
				Name:   "Opphaveleg melding",
				Value:  fmt.Sprintf("[Hopp til melding](https://discord.com/channels/%s/%s/%s)", r.GuildID, r.ChannelID, r.MessageID),
				Inline: false,
			},
		}

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

