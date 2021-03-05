package main

import (
	"fmt"
	"log"

	"github.com/zmb3/spotify"
)

const releaseRadarId = "37i9dQZEVXblHXYINKqgaL"
const releaseRadarLength = 30

func getPlaylistInfo(client *spotify.Client) []spotify.PlaylistTrack {
	page, err := client.GetPlaylistTracks(spotify.ID(releaseRadarId))
	if err != nil {
		log.Fatalf("couldn't get features playlists: %v", err)
	}

	return page.Tracks
}

func printPlaylistInfo(tracklist []spotify.PlaylistTrack) {
	fmt.Printf("ID | Name | Popularity | Artists | Album | Duration")
	for _, trackObj := range tracklist {
		track := trackObj.Track
		fmt.Printf("%v | %v | %v | %v | %v | %v\n", track.ID, track.Name, track.Popularity, track.Artists[0].Name, track.Album.Name, track.TimeDuration())
	}
}

func assembleTrackIDs(tracklist []spotify.PlaylistTrack) []spotify.ID {
	var trackIDs = make([]spotify.ID, releaseRadarLength)
	for _, trackObj := range tracklist {
		trackIDs = append(trackIDs, trackObj.Track.ID)
	}
	return trackIDs
}

func main() {
	fmt.Println("Ocotillo")
	client := clientCredentialsAuth()
	tracklist := getPlaylistInfo(client)
	trackIDs := assembleTrackIDs(tracklist)
	printPlaylistInfo(tracklist)
	fmt.Println(trackIDs)
}
