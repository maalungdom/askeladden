package commands

import (
	"fmt"
	"log"
	"reflect"

	"github.com/bwmarrin/discordgo"
	"askeladden/internal/bot"
	"askeladden/internal/bot/services"
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

	configInfo := "**🔧 Bot Configuration**\n\n"
	
	// Use reflection to iterate through the config struct dynamically
	v := reflect.ValueOf(cfg).Elem()
	t := v.Type()
	
	for i := 0; i < v.NumField(); i++ {
		section := v.Field(i)
		sectionType := section.Type()
		
		configInfo += fmt.Sprintf("**%s Settings:**\n", t.Field(i).Name)
		
		for j := 0; j < section.NumField(); j++ {
			field := section.Field(j)
			fieldName := sectionType.Field(j).Name
			configInfo += fmt.Sprintf("• %s: %v\n", fieldName, field.Interface())
		}
		configInfo += "\n"
	}

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
	embed := services.CreateBotEmbed(s, "🔧 Configuration", configInfo, 0x0099ff)
	_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
	if err != nil {
		log.Printf("Failed to send config embed: %v", err)
		s.ChannelMessageSend(m.ChannelID, "Kunne ikke sende konfigurasjonsinformasjon.")
	}
}
