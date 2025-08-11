package commands

import (
	"fmt"
	"log"
	"strings"

	"askeladden/internal/bot"
	"askeladden/internal/bot/services"

	"github.com/bwmarrin/discordgo"
)

func init() {
	commands["kjeften"] = Command{
		name:        "kjeften",
		description: "Si Askeladden han må teie for dei upratsame.",
		emoji:       "🤐",
		handler:     Kjeften,
	}
}

// Kjeften toggles the "pratsam" role on the invoking user. If the role does not
// exist it reports an error. It will also check bot role hierarchy and return
// a friendly embed explaining why the action failed if the bot cannot modify the role.
func Kjeften(s *discordgo.Session, m *discordgo.MessageCreate, bot *bot.Bot) {
	// Must be used in a guild
	if m.GuildID == "" {
		embed := services.CreateBotEmbed(s, "Feil", "Denne kommandoen må brukast i ein server (ikkje PM).", services.EmbedTypeError)
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
		return
	}

	// Load guild roles
	guildRoles, err := s.GuildRoles(m.GuildID)
	if err != nil {
		log.Printf("failed to fetch guild roles: %v", err)
		embed := services.CreateBotEmbed(s, "Feil", "Klarte ikkje hente roller i guilden.", services.EmbedTypeError)
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
		return
	}

	// Build a lookup map of roles by ID and find the role named "pratsam" (case-insensitive)
	rolesByID := make(map[string]*discordgo.Role, len(guildRoles))
	var pratsamRoleID string
	for _, r := range guildRoles {
		rolesByID[r.ID] = r
		if strings.EqualFold(r.Name, "pratsam") {
			pratsamRoleID = r.ID
		}
	}

	if pratsamRoleID == "" {
		embed := services.CreateBotEmbed(s, "Feil", "Fann ikkje rolla 'pratsam' i guilden.", services.EmbedTypeError)
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
		return
	}

	// Fetch the member invoking the command
	member, err := s.State.Member(m.GuildID, m.Author.ID)
	if err != nil {
		member, err = s.GuildMember(m.GuildID, m.Author.ID)
		if err != nil {
			log.Printf("failed to fetch member: %v", err)
			embed := services.CreateBotEmbed(s, "Feil", "Klarte ikkje hente medlem sin informasjon.", services.EmbedTypeError)
			s.ChannelMessageSendEmbed(m.ChannelID, embed)
			return
		}
	}

	// Check whether the member already has the role
	hasRole := false
	for _, rid := range member.Roles {
		if rid == pratsamRoleID {
			hasRole = true
			break
		}
	}

	// Determine bot identity and its highest role position to check hierarchy
	var botUserID string
	if s.State != nil && s.State.User != nil {
		botUserID = s.State.User.ID
	} else {
		u, err := s.User("@me")
		if err == nil {
			botUserID = u.ID
		}
	}

	// Try to fetch the bot's guild member to inspect its roles (may fail)
	var botMember *discordgo.Member
	if botUserID != "" {
		botMember, _ = s.State.Member(m.GuildID, botUserID)
		if botMember == nil {
			botMember, _ = s.GuildMember(m.GuildID, botUserID)
		}
	}

	// If we have botMember and role info, check hierarchy: bot must be higher than the target role
	if botMember != nil {
		botHighest := -9999
		targetRole := rolesByID[pratsamRoleID]
		if targetRole != nil {
			for _, rid := range botMember.Roles {
				if r := rolesByID[rid]; r != nil {
					if r.Position > botHighest {
						botHighest = r.Position
					}
				}
			}
			if botHighest <= targetRole.Position {
				msg := fmt.Sprintf("Botens rolle er ikkje høg nok til å endre rolla '%s'. Flytt boten sin rolle over '%s' i serverinnstillingane.", targetRole.Name, targetRole.Name)
				embed := services.CreateBotEmbed(s, "Feil", msg, services.EmbedTypeError)
				s.ChannelMessageSendEmbed(m.ChannelID, embed)
				return
			}
		}
	}

	// Toggle the role: remove if present, add if absent
	if hasRole {
		if err := s.GuildMemberRoleRemove(m.GuildID, m.Author.ID, pratsamRoleID); err != nil {
			log.Printf("failed to remove role: %v", err)
			embed := services.CreateBotEmbed(s, "Feil", fmt.Sprintf("Klarte ikkje fjerne rolla 'pratsam': %v", err), services.EmbedTypeError)
			s.ChannelMessageSendEmbed(m.ChannelID, embed)
			return
		}
		embed := services.CreateBotEmbed(s, "Orsak! 🤐", "Eg visste ikkje at du ikkje var ein pratsam type. Eg skal lata vere å plaga deg.", services.EmbedTypeSuccess)
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
		log.Printf("Removed role 'pratsam' from %s (%s)", m.Author.Username, m.Author.ID)
		return
	}

	// Add role
	if err := s.GuildMemberRoleAdd(m.GuildID, m.Author.ID, pratsamRoleID); err != nil {
		log.Printf("failed to add role: %v", err)
		embed := services.CreateBotEmbed(s, "Feil", fmt.Sprintf("Klarte ikkje legge til rolla 'pratsam': %v", err), services.EmbedTypeError)
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
		return
	}

	embed := services.CreateBotEmbed(s, "Hei du! 📢", "Eg trur vi kjem til å vere gode venar!", services.EmbedTypeSuccess)
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
	log.Printf("Added role 'pratsam' to %s (%s)", m.Author.Username, m.Author.ID)
}
