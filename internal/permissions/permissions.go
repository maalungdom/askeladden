package permissions

import (
	"fmt"
	"log"
	"strings"

	"askeladden/internal/config"
	"github.com/bwmarrin/discordgo"
)

// PermissionManager handles role-based permissions and approvals
type PermissionManager struct {
	Config *config.Config
}

// NewPermissionManager creates a new permission manager
func NewPermissionManager(cfg *config.Config) *PermissionManager {
	return &PermissionManager{
		Config: cfg,
	}
}

// UserRole represents the roles a user can have
type UserRole int

const (
	RoleNone UserRole = iota
	RoleOpplysar
	RoleRettskrivar
	RoleBoth
)

// GetUserRole returns the user's role(s)
func (pm *PermissionManager) GetUserRole(s *discordgo.Session, guildID, userID string) UserRole {
	member, err := s.GuildMember(guildID, userID)
	if err != nil {
		log.Printf("Failed to get guild member: %v", err)
		return RoleNone
	}

	hasOpplysar := false
	hasRettskrivar := false

	for _, roleID := range member.Roles {
		if roleID == pm.Config.Approval.OpplysarRoleID {
			hasOpplysar = true
		}
		if roleID == pm.Config.BannedWords.RettskrivarRoleID {
			hasRettskrivar = true
		}
	}

	if hasOpplysar && hasRettskrivar {
		return RoleBoth
	} else if hasOpplysar {
		return RoleOpplysar
	} else if hasRettskrivar {
		return RoleRettskrivar
	}

	return RoleNone
}

// HasOpplysarRole checks if user has opplysar role
func (pm *PermissionManager) HasOpplysarRole(s *discordgo.Session, guildID, userID string) bool {
	role := pm.GetUserRole(s, guildID, userID)
	return role == RoleOpplysar || role == RoleBoth
}

// HasRettskrivarRole checks if user has rettskrivar role
func (pm *PermissionManager) HasRettskrivarRole(s *discordgo.Session, guildID, userID string) bool {
	role := pm.GetUserRole(s, guildID, userID)
	return role == RoleRettskrivar || role == RoleBoth
}

// ApprovalState represents the approval state for banned words
type ApprovalState struct {
	HasOpplysarApproval    bool
	HasRettskrivarApproval bool
	OpplysarApprovers      []string
	RettskrivarApprovers   []string
}

// CheckCombinedApproval checks all reactions on a message to see if both roles are represented
func (pm *PermissionManager) CheckCombinedApproval(s *discordgo.Session, channelID, messageID, emoji string) (*ApprovalState, error) {
	// Get all users who reacted with the approval emoji
	users, err := s.MessageReactions(channelID, messageID, emoji, 100, "", "")
	if err != nil {
		return nil, err
	}

	state := &ApprovalState{
		OpplysarApprovers:    make([]string, 0),
		RettskrivarApprovers: make([]string, 0),
	}

	// Get guild ID from channel
	channel, err := s.Channel(channelID)
	if err != nil {
		return nil, err
	}

	// Check each user's roles
	for _, user := range users {
		// Skip bot users
		if user.Bot {
			continue
		}

		role := pm.GetUserRole(s, channel.GuildID, user.ID)

		switch role {
		case RoleOpplysar:
			state.HasOpplysarApproval = true
			state.OpplysarApprovers = append(state.OpplysarApprovers, user.ID)
		case RoleRettskrivar:
			state.HasRettskrivarApproval = true
			state.RettskrivarApprovers = append(state.RettskrivarApprovers, user.ID)
		case RoleBoth:
			// User has both roles, count for both
			state.HasOpplysarApproval = true
			state.HasRettskrivarApproval = true
			state.OpplysarApprovers = append(state.OpplysarApprovers, user.ID)
			state.RettskrivarApprovers = append(state.RettskrivarApprovers, user.ID)
		}
	}

	return state, nil
}

// IsFullyApproved checks if both required roles have approved
func (state *ApprovalState) IsFullyApproved() bool {
	return state.HasOpplysarApproval && state.HasRettskrivarApproval
}

// GetApprovalSummary returns a simple approval status summary
func (state *ApprovalState) GetApprovalSummary(s *discordgo.Session) string {
	var parts []string

	// Show opplysar approvals
	if state.HasOpplysarApproval {
		var opplysarNames []string
		for _, userID := range state.OpplysarApprovers {
			user, err := s.User(userID)
			if err == nil {
				opplysarNames = append(opplysarNames, user.Username)
			} else {
				opplysarNames = append(opplysarNames, "Ukjend")
			}
		}
		parts = append(parts, fmt.Sprintf("🧘‍♀️ Opplysar-godkjenning: %s", strings.Join(opplysarNames, ", ")))
	} else {
		parts = append(parts, "⏳ Opplysar-godkjenning: ventar")
	}

	// Show rettskrivar approvals
	if state.HasRettskrivarApproval {
		var rettskrivarNames []string
		for _, userID := range state.RettskrivarApprovers {
			user, err := s.User(userID)
			if err == nil {
				rettskrivarNames = append(rettskrivarNames, user.Username)
			} else {
				rettskrivarNames = append(rettskrivarNames, "Ukjend")
			}
		}
		parts = append(parts, fmt.Sprintf("📝 Rettskrivar-godkjenning: %s", strings.Join(rettskrivarNames, ", ")))
	} else {
		parts = append(parts, "⏳ Rettskrivar-godkjenning: ventar")
	}

	return strings.Join(parts, "\n")
}
