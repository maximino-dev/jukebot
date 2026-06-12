package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

var playlistID string
var service *youtube.Service

type Config struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RefreshToken string `json:"refresh_token"`
	PlaylistID   string `json:"playlist_id"`
	DiscordToken string `json:"discord_token"`
	ChannelId    string `json:"channel_id"`
}

func LoadConfig() (*Config, error) {
	data, err := os.ReadFile("config.json")
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if cfg.ClientID == "" || cfg.ClientSecret == "" {
		return nil, fmt.Errorf("Veuillez saisir les champs client_secret et client_id dans le fichier de config")
	}

	if cfg.ChannelId == "" {
		return nil, fmt.Errorf("Veuillez saisir un id de salon Discord dans le fichier de configuration")
	}

	if cfg.DiscordToken == "" {
		return nil, fmt.Errorf("Veuillez saisir un token Discord dans le fichier de config")
	}

	return &cfg, nil
}

func saveConfig(cfg *Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile("config.json", data, 0644)
}

func getYoutubeService(cfg *Config) (*youtube.Service, error) {
	ctx := context.Background()

	config := &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://localhost:0",
		Scopes:       []string{youtube.YoutubeForceSslScope},
	}

	token := &oauth2.Token{
		RefreshToken: cfg.RefreshToken,
	}

	client := config.Client(ctx, token)

	return youtube.NewService(ctx, option.WithHTTPClient(client))
}
func authFlow(cfg *Config) error {
	ctx := context.Background()

	oauthConfig := &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://localhost:0",
		Scopes:       []string{"https://www.googleapis.com/auth/youtube"},
	}

	authURL := oauthConfig.AuthCodeURL(
		"state",
		oauth2.AccessTypeOffline,
	)

	fmt.Println("Connection URL :")
	fmt.Println(authURL)

	fmt.Print("Please paste the redirected URL here : ")

	var redirectURL string
	fmt.Scan(&redirectURL)

	u, err := url.Parse(redirectURL)
	if err != nil {
		return err
	}

	code := u.Query().Get("code")

	token, err := oauthConfig.Exchange(ctx, code)
	if err != nil {
		return err
	}

	cfg.RefreshToken = token.RefreshToken

	return saveConfig(cfg)
}

func Connect() error {
	cfg, err := LoadConfig()
	if err != nil {
		return err
	}

	if cfg.PlaylistID == "" {
		return fmt.Errorf("Please set a playlist Id in the config.json file")
	}

	playlistID = cfg.PlaylistID

	if cfg.ClientID == "" || cfg.ClientSecret == "" {
		return fmt.Errorf("Please set client Id and client secret in the config.json file")
	}

	if cfg.RefreshToken == "" {
		if err := authFlow(cfg); err != nil {
			return err
		}
	}

	s, err := getYoutubeService(cfg)
	if err != nil {
		return err
	}

	service = s
	return nil
}

func extractVideoID(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	if v := u.Query().Get("v"); v != "" {
		return v, nil
	}

	if strings.Contains(u.Host, "youtu.be") {
		return strings.TrimPrefix(u.Path, "/"), nil
	}

	return "", fmt.Errorf("invalid youtube url")
}

func AddURLToPlaylist(url string) error {
	id, err := extractVideoID(url)
	if err != nil {
		return err
	}

	return addToPlaylist(id)
}

func addToPlaylist(videoID string) error {
	item := &youtube.PlaylistItem{
		Snippet: &youtube.PlaylistItemSnippet{
			PlaylistId: playlistID,
			ResourceId: &youtube.ResourceId{
				Kind:    "youtube#video",
				VideoId: videoID,
			},
		},
	}

	_, err := service.PlaylistItems.Insert([]string{"snippet"}, item).Do()

	return err
}

func ListPlaylists() error {
	call := service.Playlists.List([]string{"snippet", "contentDetails"}).Mine(true).MaxResults(50)

	response, err := call.Do()
	if err != nil {
		return err
	}

	for _, item := range response.Items {
		fmt.Println("Title:", item.Snippet.Title)
		fmt.Println("ID:", item.Id)
	}

	return nil
}
