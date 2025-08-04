
package services

import (
	"github.com/bwmarrin/discordgo"
	"roersla.no/askeladden/internal/database"
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
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Sendt inn av brukeren",
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

// CreateAISlopWarningEmbed returns the AI slop warning embed message
func CreateAISlopWarningEmbed(session *discordgo.Session, version string, message string) *discordgo.MessageEmbed {
	botUser, _ := session.User("@me")
	return &discordgo.MessageEmbed{
		Title:       "⚠️ AI Slop Warning",
		Description: message + "\n\nVersion: " + version + "\n➔ Vil du bidra? Send en pull request! https://github.com/maalungdom/askeladden",
		Color:       0xFFA500, // Orange
		Author: &discordgo.MessageEmbedAuthor{
			Name:    botUser.Username,
			IconURL: botUser.AvatarURL(""),
		},
	}
}

