package services

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"roersla.no/askeladden/internal/bot"
	"roersla.no/askeladden/internal/database"
)

// SendDailyQuestion sends the daily question to the appropriate channel
// mention may be "@everyone", "<@user_id>", or blank
func SendDailyQuestion(bot bot.BotIface, question *database.Question, mention string) {
	channelID := "1379979709055762518"

	// Try to fetch pretty channel name
	chanObj, chanErr := bot.GetSession().State.Channel(channelID)
	channelName := channelID
	if chanErr == nil {
		channelName = "#" + chanObj.Name
	}
	// Fetch Discord user for embed author
	authorObj, _ := bot.GetSession().User(question.AuthorID)
	// Embed
	embed := CreateDailyQuestionEmbed(question, authorObj)
	// Use @mention string if provided; else, empty
	msg := &discordgo.MessageSend{
		Content: mention,
		Embeds:  []*discordgo.MessageEmbed{embed},
	}
	log.Printf("[MESSAGING] Sending daily question to %s for %s: \"%s\" [mention:'%s']", channelName, embed.Author.Name, question.Question, mention)
	_, err := bot.GetSession().ChannelMessageSendComplex(channelID, msg)
	if err != nil {
		log.Printf("[MESSAGING] Failed to send daily question: %v", err)
	}
}

