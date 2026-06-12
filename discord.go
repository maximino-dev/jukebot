package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var ChannelID string

func StartBot() error {
	cfg, err := LoadConfig()
	if err != nil {
		return err
	}

	dg, err := discordgo.New("Bot " + cfg.DiscordToken)
	if err != nil {
		return err
	}

	ChannelID = cfg.ChannelId

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

	if m.Author.Bot {
		return
	}

	if m.ChannelID != ChannelID {
		return
	}

	if len(m.Embeds) > 0 {
		err := AddURLToPlaylist(m.Embeds[0].URL)
		if err != nil {
			fmt.Println("Erreur lors de l'ajout de l'url ", err)
			s.ChannelMessageSend(ChannelID, "Erreur lors de l'ajout de la musique "+err.Error())
		} else {
			s.ChannelMessageSend(ChannelID, "Musique ajoutée dans la playlist")
		}
	}
}

func ListMessages() {
}
