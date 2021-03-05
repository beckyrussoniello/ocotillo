package main

import (
	"fmt"
	"log"

	"github.com/zmb3/spotify"
)

var releaseRadarId = "37i9dQZEVXblHXYINKqgaL"

func printPlaylistInfo(client *spotify.Client) []spotify.PlaylistTrack {
	page, err := client.GetPlaylistTracks(spotify.ID(releaseRadarId))
	if err != nil {
		log.Fatalf("couldn't get features playlists: %v", err)
	}

	fmt.Printf("ID | Name | Popularity | Artists | Album | Duration")
	for _, trackObj := range page.Tracks {
		track := trackObj.Track
		fmt.Printf("%v | %v | %v | %v | %v | %v\n", track.ID, track.Name, track.Popularity, track.Artists[0].Name, track.Album.Name, track.TimeDuration())
	}

	return page.Tracks
}

func main() {
	fmt.Println("Ocotillo")
	client := clientCredentialsAuth()
	printPlaylistInfo(client)
	//userAuth()
}
