package services

import (
	"roersla.no/askeladden/internal/bot"
)

// BotServices holds all the services the bot can use.
type BotServices struct {
	Approval *ApprovalService
}

// New creates a new BotServices instance.
func New(b bot.BotIface) *BotServices {
	return &BotServices{
		Approval: &ApprovalService{Bot: b},
	}
}

