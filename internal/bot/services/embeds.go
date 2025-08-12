package services

import (
	"fmt"
	"time"

	"askeladden/internal/database"
	"github.com/bwmarrin/discordgo"
)

// Embed color constants for consistent styling
const (
	ColorSuccess   = 0x00ff00 // Green
	ColorError     = 0xff0000 // Red
	ColorWarning   = 0xffa500 // Orange
	ColorInfo      = 0x0099ff // Blue
	ColorPrimary   = 0x7289da // Discord Blurple
	ColorStarboard = 0xFFD700 // Gold
)

// EmbedType represents different types of embeds
type EmbedType int

const (
	EmbedTypeSuccess EmbedType = iota
	EmbedTypeError
	EmbedTypeWarning
	EmbedTypeInfo
	EmbedTypePrimary
)

// EmbedBuilder provides a fluent interface for building Discord embeds
type EmbedBuilder struct {
	embed *discordgo.MessageEmbed
}

// NewEmbedBuilder creates a new embed builder
func NewEmbedBuilder() *EmbedBuilder {
	return &EmbedBuilder{
		embed: &discordgo.MessageEmbed{},
	}
}

// SetTitle sets the embed title
func (eb *EmbedBuilder) SetTitle(title string) *EmbedBuilder {
	eb.embed.Title = title
	return eb
}

// SetDescription sets the embed description
func (eb *EmbedBuilder) SetDescription(description string) *EmbedBuilder {
	eb.embed.Description = description
	return eb
}

// SetColor sets the embed color
func (eb *EmbedBuilder) SetColor(color int) *EmbedBuilder {
	eb.embed.Color = color
	return eb
}

// SetColorByType sets the color based on embed type
func (eb *EmbedBuilder) SetColorByType(embedType EmbedType) *EmbedBuilder {
	switch embedType {
	case EmbedTypeSuccess:
		eb.embed.Color = ColorSuccess
	case EmbedTypeError:
		eb.embed.Color = ColorError
	case EmbedTypeWarning:
		eb.embed.Color = ColorWarning
	case EmbedTypeInfo:
		eb.embed.Color = ColorInfo
	case EmbedTypePrimary:
		eb.embed.Color = ColorPrimary
	default:
		eb.embed.Color = ColorInfo
	}
	return eb
}

// SetAuthor sets the embed author
func (eb *EmbedBuilder) SetAuthor(name, iconURL string) *EmbedBuilder {
	eb.embed.Author = &discordgo.MessageEmbedAuthor{
		Name:    name,
		IconURL: iconURL,
	}
	return eb
}

// SetAuthorFromUser sets the author from a Discord user
func (eb *EmbedBuilder) SetAuthorFromUser(user *discordgo.User) *EmbedBuilder {
	if user != nil {
		eb.embed.Author = &discordgo.MessageEmbedAuthor{
			Name:    user.Username,
			IconURL: user.AvatarURL(""),
		}
	}
	return eb
}

// SetAuthorFromBot sets the author from the bot user
func (eb *EmbedBuilder) SetAuthorFromBot(session *discordgo.Session) *EmbedBuilder {
	if botUser, err := session.User("@me"); err == nil {
		eb.embed.Author = &discordgo.MessageEmbedAuthor{
			Name:    botUser.Username,
			IconURL: botUser.AvatarURL(""),
		}
	}
	return eb
}

// AddField adds a field to the embed
func (eb *EmbedBuilder) AddField(name, value string, inline bool) *EmbedBuilder {
	field := &discordgo.MessageEmbedField{
		Name:   name,
		Value:  value,
		Inline: inline,
	}
	eb.embed.Fields = append(eb.embed.Fields, field)
	return eb
}

// SetFooter sets the embed footer
func (eb *EmbedBuilder) SetFooter(text, iconURL string) *EmbedBuilder {
	eb.embed.Footer = &discordgo.MessageEmbedFooter{
		Text:    text,
		IconURL: iconURL,
	}
	return eb
}

// SetTimestamp sets the embed timestamp to current time
func (eb *EmbedBuilder) SetTimestamp() *EmbedBuilder {
	eb.embed.Timestamp = time.Now().Format(time.RFC3339)
	return eb
}

// Build returns the final embed
func (eb *EmbedBuilder) Build() *discordgo.MessageEmbed {
	return eb.embed
}

// Convenience functions for common embed patterns

// CreateSuccessEmbed creates a standardized success embed
func CreateSuccessEmbed(title, description string) *discordgo.MessageEmbed {
	return NewEmbedBuilder().
		SetTitle("✅ " + title).
		SetDescription(description).
		SetColorByType(EmbedTypeSuccess).
		Build()
}

// CreateErrorEmbed creates a standardized error embed
func CreateErrorEmbed(title, description string) *discordgo.MessageEmbed {
	return NewEmbedBuilder().
		SetTitle("❌ " + title).
		SetDescription(description).
		SetColorByType(EmbedTypeError).
		Build()
}

// CreateWarningEmbed creates a standardized warning embed
func CreateWarningEmbed(title, description string) *discordgo.MessageEmbed {
	return NewEmbedBuilder().
		SetTitle("⚠️ " + title).
		SetDescription(description).
		SetColorByType(EmbedTypeWarning).
		Build()
}

// CreateInfoEmbed creates a standardized info embed
func CreateInfoEmbed(title, description string) *discordgo.MessageEmbed {
	return NewEmbedBuilder().
		SetTitle("ℹ️ " + title).
		SetDescription(description).
		SetColorByType(EmbedTypeInfo).
		Build()
}

// CreateBotEmbed creates an embed with the bot as author
func CreateBotEmbed(session *discordgo.Session, title, description string, embedType EmbedType) *discordgo.MessageEmbed {
	return NewEmbedBuilder().
		SetTitle(title).
		SetDescription(description).
		SetColorByType(embedType).
		SetAuthorFromBot(session).
		Build()
}

// Legacy function maintained for compatibility
func CreateBotEmbedLegacy(session *discordgo.Session, title, description string, color int) *discordgo.MessageEmbed {
	return NewEmbedBuilder().
		SetTitle(title).
		SetDescription(description).
		SetColor(color).
		SetAuthorFromBot(session).
		Build()
}

// CreateDailyQuestionEmbed creates the embed for daily questions
func CreateDailyQuestionEmbed(question *database.Question, author *discordgo.User) *discordgo.MessageEmbed {
	builder := NewEmbedBuilder().
		SetTitle("🌅 Dagens spørsmål").
		SetDescription(question.Question).
		SetColorByType(EmbedTypeInfo)

	if author != nil {
		builder.SetAuthorFromUser(author)
	} else {
		builder.SetAuthor(question.AuthorName, "")
	}

	return builder.Build()
}

// CreateApprovalEmbed creates standardized approval embeds
func CreateApprovalEmbed(title, description string, author *discordgo.User) *discordgo.MessageEmbed {
	builder := NewEmbedBuilder().
		SetTitle(title).
		SetDescription(description).
		SetColorByType(EmbedTypeError) // Red for pending approval

	if author != nil {
		builder.SetAuthorFromUser(author)
	}

	return builder.Build()
}

// CreateBannedWordWarningEmbed creates standardized banned word warning embeds
func CreateBannedWordWarningEmbed(bannedWords []string, forumThreads []string) *discordgo.MessageEmbed {
	var warningText string
	if len(bannedWords) == 1 {
		warningText = fmt.Sprintf("⚠️ **Grammatisk merknad**\n\nOrdet **\"%s\"** er markert som feilaktig i norsk.", bannedWords[0])
	} else {
		warningText = fmt.Sprintf("⚠️ **Grammatisk merknad**\n\nDesse orda er markerte som feilaktige i norsk: **%s**",
			fmt.Sprintf("%v", bannedWords))
	}

	// Add forum thread references if available
	if len(forumThreads) > 0 {
		// Remove duplicates
		uniqueThreads := make(map[string]bool)
		var uniqueThreadList []string
		for _, threadID := range forumThreads {
			if !uniqueThreads[threadID] {
				uniqueThreads[threadID] = true
				uniqueThreadList = append(uniqueThreadList, threadID)
			}
		}

		if len(uniqueThreadList) == 1 {
			warningText += fmt.Sprintf("\n\nSjå diskusjon: <#%s>", uniqueThreadList[0])
		} else if len(uniqueThreadList) > 1 {
			threadLinks := make([]string, len(uniqueThreadList))
			for i, threadID := range uniqueThreadList {
				threadLinks[i] = fmt.Sprintf("<#%s>", threadID)
			}
			warningText += fmt.Sprintf("\n\nSjå diskusjonar: %s", fmt.Sprintf("%v", threadLinks))
		}
	} else {
		warningText += "\n\nSjå grammatikkforumet for meir informasjon."
	}

	return NewEmbedBuilder().
		SetTitle("📝 Språkrettleiing").
		SetDescription(warningText).
		SetColorByType(EmbedTypeWarning).
		Build()
}

// CreateStarboardEmbed creates standardized starboard embeds
func CreateStarboardEmbed(msg *discordgo.Message, stars int, channelName, emoji, guildID string) *discordgo.MessageEmbed {
	builder := NewEmbedBuilder().
		SetDescription(msg.Content).
		SetColor(ColorStarboard).
		SetAuthorFromUser(msg.Author).
		SetFooter(fmt.Sprintf("%s %d | #%s", emoji, stars, channelName), "").
		AddField("Opphaveleg melding",
			fmt.Sprintf("[Hopp til melding](https://discord.com/channels/%s/%s/%s)",
				guildID, msg.ChannelID, msg.ID), false)

	// Set timestamp from original message
	builder.embed.Timestamp = msg.Timestamp.Format(time.RFC3339)

	return builder.Build()
}
