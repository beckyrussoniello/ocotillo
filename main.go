package main

import (
	"fmt"
	"log"
	"reflect"

	"github.com/zmb3/spotify"
)

const releaseRadarId = "37i9dQZEVXblHXYINKqgaL"
const releaseRadarLength = 30

type song struct {
	ID               string
	Name             string
	ArtistName       string
	AlbumName        string
	Danceability     float32
	Duration         int
	Energy           float32
	Instrumentalness float32
	Liveness         float32
	Popularity       int
	Speechiness      float32
	Tempo            int
	Valence          float32
}

func getPlaylistInfo(client *spotify.Client) []spotify.PlaylistTrack {
	page, err := client.GetPlaylistTracks(spotify.ID(releaseRadarId))
	if err != nil {
		log.Fatalf("couldn't get features playlists: %v", err)
	}

	return page.Tracks
}

func getTracksInfo(client *spotify.Client, trackIDs []spotify.ID) []*spotify.AudioFeatures {
	page, err := client.GetAudioFeatures(trackIDs...)
	if err != nil {
		log.Fatalf("couldn't get tracks info: %v", err)
	}
	return page
}

func printSongsInfo(songsMap map[spotify.ID]song) {
	fmt.Printf("ID | Name | Artist | Album | Danceability | Duration | Energy | Instrumentalness | Liveness | Popularity | Speechiness | Tempo | Valence\n")
	for _, song := range songsMap {
		value := reflect.ValueOf(song)
		for attr := 0; attr < value.NumField(); attr++ {
			fmt.Printf("%v | ", value.Field(attr).Interface())
		}
		fmt.Println()
	}
}

func assembleTrackIDs(songInfo map[spotify.ID]song) []spotify.ID {
	trackIDs := make([]spotify.ID, 0, len(songInfo))
	for key, _ := range songInfo {
		trackIDs = append(trackIDs, key)
	}
	return trackIDs
}

func main() {
	fmt.Println("Ocotillo")
	client := clientCredentialsAuth()
	songInfo := make(map[spotify.ID]song)
	tracklist := getPlaylistInfo(client)

	for _, trackObj := range tracklist {
		track := trackObj.Track
		songInfo[track.ID] = song{Name: track.Name, ArtistName: track.Artists[0].Name,
			AlbumName: track.Album.Name, Popularity: track.Popularity}
	}
	trackIDs := assembleTrackIDs(songInfo)
	tracksInfo := getTracksInfo(client, trackIDs)
	for _, track := range tracksInfo {
		song := songInfo[track.ID]
		song.Danceability = track.Danceability
		song.Duration = track.Duration
		song.Energy = track.Energy
		song.Instrumentalness = track.Instrumentalness
		song.Liveness = track.Liveness
		song.Speechiness = track.Speechiness
		song.Tempo = int(track.Tempo)
		song.Valence = track.Valence
		songInfo[track.ID] = song
	}
	printSongsInfo(songInfo)
}
