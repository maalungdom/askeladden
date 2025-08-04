package reactions

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"roersla.no/askeladden/internal/bot"
)

func init() {
	Register("⭐", "Add a message to the starboard", handleStarReaction)
}

func handleStarReaction(s *discordgo.Session, r *discordgo.MessageReactionAdd, b bot.BotIface) {
	if r.UserID == s.State.User.ID { // Ignore bot's own reactions
		return
	}
	// Don't process reactions in the starboard channel itself
	if r.ChannelID == b.GetConfig().Starboard.ChannelID {
		return
	}

	// Fetch message
	msg, err := s.ChannelMessage(r.ChannelID, r.MessageID)
	if err != nil {
		log.Printf("Error fetching message: %v", err)
		return
	}

	// Count total star reactions
	stars := 0
	for _, reaction := range msg.Reactions {
		if reaction.Emoji.Name == "⭐" {
			stars = reaction.Count
			break
		}
	}

	// Log for debugging
	log.Printf("Message %s in channel %s has %d stars (threshold: %d)", r.MessageID, r.ChannelID, stars, b.GetConfig().Starboard.Threshold)

	if stars >= b.GetConfig().Starboard.Threshold {
		// Send message to starboard channel
		embed := &discordgo.MessageEmbed{
			Author: &discordgo.MessageEmbedAuthor{
				Name:    msg.Author.Username,
				IconURL: msg.Author.AvatarURL(""),
			},
			Description: msg.Content,
			Footer: &discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("⭐ %d | #%s", stars, getChannelName(s, r.ChannelID)),
			},
			Timestamp: string(msg.Timestamp.Format(time.RFC3339)),
			Color:     0xFFD700, // Gold color
		}

		// Add original message link
		embed.Fields = []*discordgo.MessageEmbedField{
			{
				Name:   "Original Message",
				Value:  fmt.Sprintf("[Jump to Message](https://discord.com/channels/%s/%s/%s)", r.GuildID, r.ChannelID, r.MessageID),
				Inline: false,
			},
		}

		_, err := s.ChannelMessageSendEmbed(b.GetConfig().Starboard.ChannelID, embed)
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

