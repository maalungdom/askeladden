package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"roersla.no/askeladden/internal/bot"
	"roersla.no/askeladden/internal/bot/services"
)

func init() {
	commands["config"] = Command{
		name:        "config",
		description: "Vis noverande bot-konfigurasjon (ikkje hemmelege opplysningar)",
		emoji:       "🔧",
		handler:     handleConfigCommand,
		adminOnly:   true,
	}
}

func handleConfigCommand(s *discordgo.Session, m *discordgo.MessageCreate, b *bot.Bot) {
	cfg := b.Config
	
	// Build configuration display (excluding secrets)
	configInfo := fmt.Sprintf("**🔧 Bot Configuration**\n\n"+
		"**Discord Settings:**\n"+
		"• Prefix: `%s`\n"+
		"• Log Channel: <#%s>\n"+
		"• Default Channel: <#%s>\n\n"+
		"**Approval Settings:**\n"+
		"• Queue Channel: <#%s>\n"+
		"• Admin Role: <@&%s>\n\n"+
		"**Starboard Settings:**\n"+
		"• Channel: <#%s>\n"+
		"• Threshold: %d reactions\n"+
		"• Emoji: %s\n\n"+
		"**Reaction Emojis:**\n"+
		"• Question: %s\n"+
		"• Approval: 👍\n"+
		"• Reject: 👎\n\n"+
		"**Database Settings:**\n"+
		"• Host: %s\n"+
		"• Port: %d\n"+
		"• Database: %s",
		cfg.Discord.Prefix,
		cfg.Discord.LogChannelID,
		cfg.Discord.DefaultChannelID,
		cfg.Approval.QueueChannelID,
		cfg.Approval.OpplysarRoleID,
		cfg.Starboard.ChannelID,
		cfg.Starboard.Threshold,
		cfg.Starboard.Emoji,
		cfg.Reactions.Question,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName)

	// Add environment-specific info if present
	if cfg.Environment != "" {
		configInfo += fmt.Sprintf("\n\n**Environment Settings:**\n• Mode: %s", cfg.Environment)
		
		if cfg.TableSuffix != "" {
			configInfo += fmt.Sprintf("\n• Table Suffix: %s", cfg.TableSuffix)
		}
	}

	// Add scheduler info
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

	// Send as embed
	embed := services.CreateBotEmbed(s, "Configuration", configInfo, 0x0099ff)
	_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Fekk ikkje sendt konfigurasjonsinformasjon.")
	}
}
