package bot

import (
	"log"

	"askeladden/internal/config"
	"askeladden/internal/database"
	"github.com/bwmarrin/discordgo"
)

// Bot represents the main bot structure.
type Bot struct {
	Session  *discordgo.Session
	Config   *config.Config
	Database *database.DB
}

// New creates a new Bot instance.
func New(cfg *config.Config, db *database.DB, session *discordgo.Session) *Bot {
	return &Bot{
		Session:  session,
		Config:   cfg,
		Database: db,
	}
}

// Start starts the bot.
func (b *Bot) Start() error {
	log.Println("[BOT] Attempting to connect to Discord...")
	// Open connection
	err := b.Session.Open()
	if err != nil {
		log.Printf("[BOT] Failed to open Discord session: %v", err)
		return err
	}

	log.Println("[BOT] Discord session opened successfully")
	log.Println("[BOT] Askeladden is running and ready to handle messages.")
	return nil
}

// Stop stops the bot.
func (b *Bot) Stop() error {
	log.Println("[BOT] Askeladden is logging off.")
	// Log channel message will be sent from main.go before calling Stop()

	// Close database connection
	if b.Database != nil {
		b.Database.Close()
	}

	return b.Session.Close()
}

// Note: Direct field access is preferred in Go for simplicity
// Bot fields (Session, Config, Database) are exported for direct access
