
package services

import (
	"github.com/bwmarrin/discordgo"
	"askeladden/internal/database"
)

// createDailyQuestionEmbed creates the embed for the daily question
func CreateDailyQuestionEmbed(question *database.Question, author *discordgo.User) *discordgo.MessageEmbed {
	var authorName, authorIcon string
	if author != nil {
		authorName = author.Username
		authorIcon = author.AvatarURL("")
	} else {
		authorName = question.AuthorName
		authorIcon = ""
	}
	return &discordgo.MessageEmbed{
		Title:       "🌅 Dagens spørsmål",
		Description: question.Question,
		Color:       0x0099ff, // Blue color
		Author: &discordgo.MessageEmbedAuthor{
			Name:    authorName,
			IconURL: authorIcon,
		},
	}
}

// CreateBotEmbed returns a MessageEmbed with the bot as author (username/avatar from live session)
func CreateBotEmbed(session *discordgo.Session, title, description string, color int) *discordgo.MessageEmbed {
	botUser, _ := session.User("@me")
	return &discordgo.MessageEmbed{
		Title:       title,
		Description: description,
		Color:       color,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    botUser.Username,
			IconURL: botUser.AvatarURL(""),
		},
	}
}



