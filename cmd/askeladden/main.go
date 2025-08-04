package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"roersla.no/askeladden/internal/bot"
	"roersla.no/askeladden/internal/bot/handlers"
	"roersla.no/askeladden/internal/bot/services"
	"roersla.no/askeladden/internal/config"
	"roersla.no/askeladden/internal/database"
)

func main() {
	// Last inn konfigurasjon
	cfg, err := config.Load()
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

	// Opprett bot
	bot := bot.New(cfg, db, session)

	// Opprett tenester og handterarar
	botServices := services.New(bot)
	handlers := handlers.New(bot, botServices)

	// Set opp hendingshandterarar
	session.AddHandler(handlers.Ready)
	session.AddHandler(handlers.MessageCreate)
	session.AddHandler(handlers.ReactionAdd)

	// Start bot
	if err := bot.Start(); err != nil {
		log.Fatalf("[MAIN] Error running bot: %v", err)
	}

	// Scheduler for daily question trigger
	ticker := scheduleDailyQuestion(bot)
	defer ticker.Stop()

	// Vent på avslutningssignal
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Stopp bot
	if err := bot.Stop(); err != nil {
		log.Fatalf("[MAIN] Error stopping bot: %v", err)
	}
}
