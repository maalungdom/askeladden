package services

import (
	"strings"

	"askeladden/internal/bot"
)

// GetPratsamRoleID retrieves the ID of the "pratsam" role for the given guild
func GetPratsamRoleID(bot *bot.Bot, guildID string) (string, error) {
	roles, err := bot.Session.GuildRoles(guildID)
	if err != nil {
		return "", err
	}

	for _, role := range roles {
		if strings.EqualFold(role.Name, "pratsam") {
			return role.ID, nil
		}
	}

	return "", nil // Role not found
}
