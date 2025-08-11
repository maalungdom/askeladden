# Discord Embed Guidelines

This document provides standardized guidelines for creating Discord embeds in the Askeladden bot.

## Overview

The embed system in Askeladden provides consistent styling, reduces code duplication, and ensures a professional appearance across all bot interactions.

## Embed Types

### Standard Embed Types

The bot supports several standardized embed types with consistent colors and styling:

| Type | Color | Use Case | Example |
|------|--------|----------|---------|
| **Success** | Green (`0x00ff00`) | Successful operations, confirmations | Command completed successfully |
| **Error** | Red (`0xff0000`) | Error messages, failures | Command failed to execute |
| **Warning** | Orange (`0xffa500`) | Warnings, cautions | Banned word detected |
| **Info** | Blue (`0x0099ff`) | General information, help text | Daily question, bot info |
| **Primary** | Discord Blurple (`0x7289da`) | Primary bot branding | Help commands, main features |

### Special Purpose Colors

| Purpose | Color | Use Case |
|---------|--------|----------|
| **Starboard** | Gold (`0xFFD700`) | Starred messages |
| **Approval Pending** | Red (`0xff0000`) | Items awaiting approval |

## Usage

### Basic Embed Creation

```go
// Using convenience functions
embed := services.CreateSuccessEmbed("Operation Complete", "Your request was processed successfully")
embed := services.CreateErrorEmbed("Error", "Something went wrong")
embed := services.CreateWarningEmbed("Warning", "Please check your input")
embed := services.CreateInfoEmbed("Information", "Here's what you need to know")
```

### Using the Embed Builder

For more complex embeds, use the fluent builder interface:

```go
embed := services.NewEmbedBuilder().
    SetTitle("Complex Embed").
    SetDescription("This embed has multiple fields").
    SetColorByType(services.EmbedTypeInfo).
    SetAuthorFromBot(session).
    AddField("Field 1", "Value 1", true).
    AddField("Field 2", "Value 2", true).
    SetFooter("Footer text", "").
    SetTimestamp().
    Build()
```

### Bot Author Embeds

For embeds that should show the bot as the author:

```go
// Using the new typed approach
embed := services.CreateBotEmbed(session, "Bot Message", "Message content", services.EmbedTypeSuccess)

// For legacy compatibility (will be phased out)
embed := services.CreateBotEmbedLegacy(session, "Bot Message", "Message content", 0x00ff00)
```

### User Author Embeds

For embeds that should show a user as the author:

```go
embed := services.NewEmbedBuilder().
    SetTitle("User Action").
    SetDescription("User performed an action").
    SetAuthorFromUser(user).
    SetColorByType(services.EmbedTypeInfo).
    Build()
```

## Specialized Embed Functions

### Daily Question Embeds

```go
embed := services.CreateDailyQuestionEmbed(question, author)
```

### Approval Embeds

```go
embed := services.CreateApprovalEmbed("Question Pending", "Waiting for approval", author)
```

### Banned Word Warning Embeds

```go
embed := services.CreateBannedWordWarningEmbed(bannedWords, forumThreads)
```

### Starboard Embeds

```go
embed := services.CreateStarboardEmbed(message, starCount, channelName, emoji, guildID)
```

## Best Practices

### 1. Use Consistent Colors

Always use the predefined color constants or embed types rather than hardcoding colors:

```go
// ✅ Good
embed := services.NewEmbedBuilder().SetColorByType(services.EmbedTypeSuccess).Build()

// ❌ Avoid
embed := &discordgo.MessageEmbed{Color: 0x00ff00} // Hardcoded color
```

### 2. Use Appropriate Icons

Standard emoji prefixes for different embed types:

- Success: ✅
- Error: ❌  
- Warning: ⚠️
- Info: ℹ️
- Primary: varies by context

### 3. Consistent Author Patterns

- Use `SetAuthorFromBot()` for bot-initiated messages
- Use `SetAuthorFromUser()` for user-initiated actions
- Use `SetAuthor()` for custom authors

### 4. Field Guidelines

- Use inline fields for related information that can be displayed side-by-side
- Use non-inline fields for detailed information that needs full width
- Limit to 25 fields maximum per embed
- Keep field names concise but descriptive

### 5. Footer Usage

- Use footers for metadata, timestamps, or contextual information
- Keep footer text brief
- Consider using timestamps for time-sensitive information

## Migration Guide

### From Legacy Embed Creation

**Before:**
```go
embed := &discordgo.MessageEmbed{
    Title:       "Success",
    Description: "Operation completed",
    Color:       0x00ff00,
    Author: &discordgo.MessageEmbedAuthor{
        Name:    session.State.User.Username,
        IconURL: session.State.User.AvatarURL(""),
    },
}
```

**After:**
```go
embed := services.CreateBotEmbed(session, "✅ Success", "Operation completed", services.EmbedTypeSuccess)
```

### From Manual Field Creation

**Before:**
```go
embed := &discordgo.MessageEmbed{
    Title: "Complex Embed",
    Fields: []*discordgo.MessageEmbedField{
        {Name: "Field 1", Value: "Value 1", Inline: true},
        {Name: "Field 2", Value: "Value 2", Inline: true},
    },
}
```

**After:**
```go
embed := services.NewEmbedBuilder().
    SetTitle("Complex Embed").
    AddField("Field 1", "Value 1", true).
    AddField("Field 2", "Value 2", true).
    Build()
```

## Examples

### Command Help Embed

```go
embed := services.NewEmbedBuilder().
    SetTitle("🤖 Askeladden - Kommandoer").
    SetDescription("Available commands for the bot").
    SetColorByType(services.EmbedTypePrimary).
    SetAuthorFromBot(session).
    AddField("General Commands", commandList, false).
    SetFooter("Use command prefix: " + bot.Config.Discord.Prefix, "").
    Build()
```

### Error Response Embed

```go
embed := services.CreateErrorEmbed("Command Failed", 
    "The command could not be executed. Please check your permissions and try again.")
```

### User Notification Embed

```go
embed := services.NewEmbedBuilder().
    SetTitle("🎉 Question Approved").
    SetDescription(fmt.Sprintf("Your question has been approved: \"%s\"", question.Question)).
    SetColorByType(services.EmbedTypeSuccess).
    SetAuthorFromUser(approver).
    SetTimestamp().
    Build()
```

## Validation

All embed functions include automatic validation for:

- Title length (256 characters max)
- Description length (4096 characters max)
- Field count (25 max)
- Field name length (256 characters max)
- Field value length (1024 characters max)
- Footer text length (2048 characters max)

The builder pattern helps prevent common embed creation errors and ensures consistency across the application.

## Future Considerations

- Consider adding embed templates for common patterns
- Implement embed preview functionality for development
- Add embed analytics to track which types are most effective
- Consider localization support for multi-language communities