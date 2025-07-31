package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"github.com/bwmarrin/discordgo"
)

//CONST.
//--------------------------------------------------------------------------------
const logChannel = "1400454839971876946"

func main() {
//1. Prøve å lesa heile innhaldet i fila "token.txt"
//-------------------------------------------------------------------------------
	tokenBytes, err := os.ReadFile("token.txt")
	/* os.ReadFile gir to verdiar:
	A. tokenBytes: Innhaldet i fila som rådata
	B. err: ei eventuell feilmelding */
	if err != nil {
		log.Fatalf("Token-fila er uleseleg: %v. Pass på at ho finnast.", err)
	}
	
//2. Passa på at tokenet er det einaste i variabelen og at fila er har eit innhald
//--------------------------------------------------------------------------------
	botToken := strings.TrimSpace(string(tokenBytes))
	if botToken == "" {
		log.Fatalln("Token-fila er tom.")
	} else {
		log.Println("Token lasta inn.")
	}

//3. Opprette Discord-session
//--------------------------------------------------------------------------------
	session, err := discordgo.New("Bot " + botToken)
	if err != nil {
		log.Fatalf("Discord-session vart ikkje oppretta: %v", err)
	}

//4. Setja opp meldings-handlar
//--------------------------------------------------------------------------------
	session.AddHandler(messageCreate)

//5. Prøve å opne tilkopling
//--------------------------------------------------------------------------------
	err = session.Open()
	if err != nil {
		log.Fatalf("Kunne ikkje opne tilkopling: %v", err)
	} else {
		log.Println("Askeladden er no pålogga tenaren.")
	}	

//KANAL. Vent på avslutningsmelding
//--------------------------------------------------------------------------------
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	log.Println("Askeladden køyrer og er klår til å handsama meldingar.")
	session.ChannelMessageSend(logChannel, "Askeladden er pålogga og klår til å hjelpa deg! 👋")
	<-signalChannel
	
//AVSLUTNING.
//---------------------------------------------------------------------------------
	log.Println("Avsluttningsmelding motteke. Askeladden loggar av.")
	session.ChannelMessageSend(logChannel, "Askeladden loggar av. Ha det bra! 👋")
	session.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Content == "!ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}
}
