package bot

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"roersla.no/askeladden/internal/config"
	"roersla.no/askeladden/internal/database"
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
	// Open connection
	err := b.Session.Open()
	if err != nil {
		return err
	}

	log.Println("[BOT] Askeladden is running and ready to handle messages.")
	return nil
}

// Stop stops the bot.
func (b *Bot) Stop() error {
	log.Println("[BOT] Askeladden is logging off.")
	if b.Config.Discord.LogChannelID != "" {
		b.Session.ChannelMessageSend(b.Config.Discord.LogChannelID, "Askeladden is logging off. Goodbye! 👋")
	}

	// Close database connection
	if b.Database != nil {
		b.Database.Close()
	}

	return b.Session.Close()
}

// GetConfig returns the bot's config.
func (b *Bot) GetConfig() *config.Config {
	return b.Config
}

// GetDatabase returns the bot's database connection.
func (b *Bot) GetDatabase() *database.DB {
	return b.Database
}


// GetSession returns the bot's Discord session.
func (b *Bot) GetSession() *discordgo.Session {
	return b.Session
}

// BotIface provides an interface for interacting with the main bot instance.
type BotIface interface {
	GetConfig() *config.Config
	GetDatabase() *database.DB
	GetSession() *discordgo.Session
	Stop() error
}

