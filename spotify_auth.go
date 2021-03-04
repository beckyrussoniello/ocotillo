package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
)

func clientCredentialsAuth() []spotify.SimplePlaylist {
	config := &clientcredentials.Config{
		ClientID:       os.Getenv("SPOTIFY_CLIENT"),
		ClientSecret:   os.Getenv("SPOTIFY_SECRET"),
		TokenURL:       spotify.TokenURL,
		Scopes:         []string{},
		EndpointParams: map[string][]string{},
		AuthStyle:      0,
	}
	token, err := config.Token(context.Background())
	if err != nil {
		log.Fatalf("couldn't get token: %v", err)
	}

	client := spotify.Authenticator{}.NewClient(token)
	msg, page, err := client.FeaturedPlaylists()
	if err != nil {
		log.Fatalf("couldn't get features playlists: %v", err)
	}

	fmt.Println(msg)
	for _, playlist := range page.Playlists {
		fmt.Println("  ", playlist.Name)
	}

	return page.Playlists
}
