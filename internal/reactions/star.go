package reactions

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"askeladden/internal/bot"
	"askeladden/internal/bot/services"
)

// RegisterStarboardReaction registers the starboard reaction with the configured emoji
func RegisterStarboardReaction(b *bot.Bot) {
	emoji := b.Config.Starboard.Emoji
	Register(emoji, "Legg til ei melding på stjernebrettet", handleStarReaction).SetRemoveHandler(handleStarReactionRemove)
}

func handleStarReaction(s *discordgo.Session, r *discordgo.MessageReactionAdd, b *bot.Bot) {
	if r.UserID == s.State.User.ID { // Ignore bot's own reactions
		return
	}
	// Don't process reactions in the starboard channel itself
	if r.ChannelID == b.Config.Starboard.ChannelID {
		return
	}

	handleStarboardUpdate(s, r.ChannelID, r.MessageID, r.GuildID, b)
}

func handleStarReactionRemove(s *discordgo.Session, r *discordgo.MessageReactionRemove, b *bot.Bot) {
	if r.UserID == s.State.User.ID { // Ignore bot's own reactions
		return
	}
	// Don't process reactions in the starboard channel itself
	if r.ChannelID == b.Config.Starboard.ChannelID {
		return
	}

	handleStarboardUpdate(s, r.ChannelID, r.MessageID, r.GuildID, b)
}

func handleStarboardUpdate(s *discordgo.Session, channelID, messageID, guildID string, b *bot.Bot) {
	// Fetch message
	msg, err := s.ChannelMessage(channelID, messageID)
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
	log.Printf("Message %s in channel %s has %d stars (threshold: %d)", messageID, channelID, stars, b.Config.Starboard.Threshold)

	// Check if a starboard message already exists for this original message
	existingStarboardMessageID, err := b.Database.GetStarboardMessage(messageID)
	if err != nil {
		log.Printf("Error checking for existing starboard message: %v", err)
		return
	}

	if stars >= b.Config.Starboard.Threshold {
		// Create updated embed
		embed := services.CreateStarboardEmbed(msg, stars, getChannelName(s, channelID), b.Config.Starboard.Emoji, guildID)

		if existingStarboardMessageID != "" {
			// Update existing starboard message
			log.Printf("Updating existing starboard message %s with %d stars", existingStarboardMessageID, stars)
			_, err := s.ChannelMessageEditEmbed(b.Config.Starboard.ChannelID, existingStarboardMessageID, embed)
			if err != nil {
				log.Printf("Error updating starboard message: %v", err)
			}
		} else {
			// Create new starboard message
			log.Printf("Creating new starboard message for original message %s with %d stars", messageID, stars)
			starboardMsg, err := s.ChannelMessageSendEmbed(b.Config.Starboard.ChannelID, embed)
			if err != nil {
				log.Printf("Error sending starboard message: %v", err)
				return
			}

			// Record the mapping in the database
			err = b.Database.AddStarboardMessage(messageID, starboardMsg.ID, channelID)
			if err != nil {
				log.Printf("Error recording starboard message mapping: %v", err)
			}
		}
	} else if existingStarboardMessageID != "" {
		// Stars dropped below threshold and there's an existing starboard message - delete it
		log.Printf("Stars dropped below threshold (%d < %d), deleting starboard message %s", stars, b.Config.Starboard.Threshold, existingStarboardMessageID)
		err := s.ChannelMessageDelete(b.Config.Starboard.ChannelID, existingStarboardMessageID)
		if err != nil {
			log.Printf("Error deleting starboard message: %v", err)
		} else {
			// Remove the mapping from the database since the message was deleted
			err = b.Database.RemoveStarboardMessage(messageID)
			if err != nil {
				log.Printf("Error removing starboard message mapping: %v", err)
			}
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

