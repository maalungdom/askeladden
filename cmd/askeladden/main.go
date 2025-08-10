package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"askeladden/internal/bot"
	"askeladden/internal/bot/handlers"
	"askeladden/internal/bot/services"
	"askeladden/internal/config"
	"askeladden/internal/database"
	"askeladden/internal/reactions"

	"github.com/bwmarrin/discordgo"
)

func main() {
	// Load configuration using environment variables for file paths
	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		configFile = "config/config.yaml"
	}

	secretsFile := os.Getenv("SECRETS_FILE")
	if secretsFile == "" {
		secretsFile = "config/secrets.yaml"
	}

	cfg, err := config.LoadWithFiles(configFile, secretsFile)
	if err != nil {
		log.Fatalf("[MAIN] Could not load configuration: %v", err)
	}

	// Opprett database-tilkobling
	db, err := database.New(cfg)
	if err != nil {
		log.Fatalf("[MAIN] Could not connect to the database: %v", err)
	}

	// Opprett Discord-sesjon
	session, err := discordgo.New("Bot " + cfg.Discord.Token)
	if err != nil {
		log.Fatalf("[MAIN] Could not create Discord session: %v", err)
	}

	// Enable necessary intents for message content
	session.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsMessageContent | discordgo.IntentsGuildMessageReactions

	// Opprett bot
	askeladden := bot.New(cfg, db, session)

	// Initialize reactions with configured emojis
	reactions.InitializeReactions(askeladden)

	// Opprett tenester og handterarar
	botServices := services.New(askeladden)
	botHandlers := handlers.New(askeladden)
	botHandlers.Services = botServices

	// Set opp hendingshandterarar
	session.AddHandler(botHandlers.Ready)
	session.AddHandler(botHandlers.MessageCreate)
	session.AddHandler(botHandlers.ReactionAdd)
	session.AddHandler(botHandlers.InteractionCreate)

	// Start bot
	if err := askeladden.Start(); err != nil {
		log.Fatalf("[MAIN] Error running bot: %v", err)
	}

	// Scheduler for daily question trigger
	ticker := scheduleDailyQuestion(askeladden)
	defer ticker.Stop()

	// Vent på avslutningssignal
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Send goodbye message before stopping
	if askeladden.Config.Discord.LogChannelID != "" {
		embed := services.CreateBotEmbed(session, "🔴 Offline", "Askeladden is logging off. Goodbye! 👋", 0xff0000)
		session.ChannelMessageSendEmbed(askeladden.Config.Discord.LogChannelID, embed)
	}

	// Stopp bot
	if err := askeladden.Stop(); err != nil {
		log.Fatalf("[MAIN] Error stopping bot: %v", err)
	}
}
