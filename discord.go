package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	ChannelID = "917736643304243230"
)

func StartBot() error {
	cfg, err := LoadConfig()
	if err != nil {
		return err
	}

	dg, err := discordgo.New("Bot " + cfg.DiscordToken)
	if err != nil {
		return err
	}

	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsMessageContent

	dg.AddHandler(listenChannel)

	err = dg.Open()
	if err != nil {
		return err
	}
	defer dg.Close()

	err = Connect()
	if err != nil {
		return err
	}

	select {}
}

func listenChannel(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.ChannelID != ChannelID {
		return
	}

	if strings.HasPrefix(m.Content, "http") {
		err := AddURLToPlaylist(m.Content)
		if err != nil {
			fmt.Println("Erreur lors de l'ajout de l'url ", err)
		}
	}
}

func ListMessages() {
}
