package services

import (
	"log"

	"askeladden/internal/bot"
	"askeladden/internal/database"
	"github.com/bwmarrin/discordgo"
)

// SendDailyQuestion sends the daily question to the appropriate channel
// mention may be "@everyone", "<@user_id>", or blank
func SendDailyQuestion(bot *bot.Bot, question *database.Question, mention string) {
	// Use configured default channel ID instead of hardcoded
	channelID := bot.Config.Discord.DefaultChannelID
	if channelID == "" {
		log.Printf("[MESSAGING] No default channel ID configured, cannot send daily question")
		return
	}

	// Try to fetch pretty channel name
	chanObj, chanErr := bot.Session.State.Channel(channelID)
	channelName := channelID
	if chanErr == nil {
		channelName = "#" + chanObj.Name
	}
	// Fetch Discord user for embed author
	authorObj, _ := bot.Session.User(question.AuthorID)
	// Embed
	embed := CreateDailyQuestionEmbed(question, authorObj)
	// Use @mention string if provided; else, empty
	msg := &discordgo.MessageSend{
		Content: mention,
		Embeds:  []*discordgo.MessageEmbed{embed},
	}
	log.Printf("[MESSAGING] Sending daily question to %s for %s: \"%s\" [mention:'%s']", channelName, embed.Author.Name, question.Question, mention)
	_, err := bot.Session.ChannelMessageSendComplex(channelID, msg)
	if err != nil {
		log.Printf("[MESSAGING] Failed to send daily question: %v", err)
	}
}
