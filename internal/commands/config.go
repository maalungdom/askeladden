package commands

import (
	"fmt"

	"askeladden/internal/bot"
	"askeladden/internal/bot/services"
	"github.com/bwmarrin/discordgo"
)

func init() {
	commands["config"] = Command{
		name:        "config",
		description: "Vis gjeldende bot-konfigurasjon (uten hemmelige opplysninger)",
		emoji:       "🔧",
		handler:     handleConfigCommand,
		adminOnly:   true,
	}
}

func handleConfigCommand(s *discordgo.Session, m *discordgo.MessageCreate, b *bot.Bot) {
	cfg := b.Config

	// Helper to get channel name from ID
	getChannelMention := func(id string) string {
		if id == "" {
			return "[ingen]"
		}
		ch, err := s.Channel(id)
		if err == nil {
			return fmt.Sprintf("<#%s> `%s`", id, ch.Name)
		}
		return fmt.Sprintf("<#%s>", id)
	}
	// Helper to get role name from ID
	getRoleMention := func(guildID, roleID string) string {
		if roleID == "" || guildID == "" {
			return "[ingen]"
		}
		role, err := s.State.Role(guildID, roleID)
		if err == nil {
			return fmt.Sprintf("<@&%s> `%s`", roleID, role.Name)
		}
		return fmt.Sprintf("<@&%s>", roleID)
	}
	// Mask secrets
	mask := "[hidden]"

	configInfo := "**🔧 Bot Configuration**\n\n"
	configInfo += "**Discord Settings:**\n"
	configInfo += fmt.Sprintf("• Prefix: `%s`\n", cfg.Discord.Prefix)
	configInfo += fmt.Sprintf("• Log Channel: %s\n", getChannelMention(cfg.Discord.LogChannelID))
	configInfo += fmt.Sprintf("• Default Channel: %s\n\n", getChannelMention(cfg.Discord.DefaultChannelID))

	configInfo += "**Approval Settings:**\n"
	configInfo += fmt.Sprintf("• Queue Channel: %s\n", getChannelMention(cfg.Approval.QueueChannelID))
	configInfo += fmt.Sprintf("• Admin Role: %s\n\n", getRoleMention(m.GuildID, cfg.Approval.OpplysarRoleID))

	configInfo += "**Starboard Settings:**\n"
	configInfo += fmt.Sprintf("• Channel: %s\n", getChannelMention(cfg.Starboard.ChannelID))
	configInfo += fmt.Sprintf("• Threshold: %d reactions\n", cfg.Starboard.Threshold)
	configInfo += fmt.Sprintf("• Emoji: %s\n\n", cfg.Starboard.Emoji)

	configInfo += "**Reaction Emojis:**\n"
	configInfo += fmt.Sprintf("• Question: %s\n", cfg.Reactions.Question)
	configInfo += "• Approval: 👍\n"
	configInfo += "• Reject: 👎\n\n"

	configInfo += "**Database Settings:**\n"
	configInfo += fmt.Sprintf("• Host: %s\n", cfg.Database.Host)
	configInfo += fmt.Sprintf("• Port: %d\n", cfg.Database.Port)
	configInfo += fmt.Sprintf("• Database: %s\n", cfg.Database.DBName)
	configInfo += fmt.Sprintf("• User: %s\n", cfg.Database.User)
	configInfo += fmt.Sprintf("• Password: %s\n\n", mask)

	if cfg.Environment != "" {
		configInfo += fmt.Sprintf("\n**Environment Settings:**\n• Mode: %s", cfg.Environment)
		if cfg.TableSuffix != "" {
			configInfo += fmt.Sprintf("\n• Table Suffix: %s", cfg.TableSuffix)
		}
	}

	if cfg.Scheduler.Enabled {
		configInfo += fmt.Sprintf("\n\n**Scheduler:**\n• Status: %s\n• Timezone: %s\n• Morning Time: %s\n• Evening Time: %s\n• Inactivity Threshold: %d hours",
			map[bool]string{true: "✅ Enabled", false: "❌ Disabled"}[cfg.Scheduler.Enabled],
			cfg.Scheduler.Timezone,
			cfg.Scheduler.MorningTime,
			cfg.Scheduler.EveningTime,
			cfg.Scheduler.InactivityHours)
		if cfg.Scheduler.CronString != "" {
			configInfo += fmt.Sprintf("\n• Fallback Cron: `%s`", cfg.Scheduler.CronString)
		}
	} else if cfg.Scheduler.CronString != "" {
		configInfo += fmt.Sprintf("\n\n**Scheduler:**\n• Status: ❌ Disabled\n• Fallback Cron: `%s`", cfg.Scheduler.CronString)
	}

	embed := services.CreateBotEmbed(s, "🔧 Configuration", configInfo, services.EmbedTypeInfo)
	_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Kunne ikkje sende konfigurasjonsinformasjon.")
	}
}
