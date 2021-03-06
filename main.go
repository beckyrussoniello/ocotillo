package main

import (
	"fmt"
	"log"
	"reflect"

	"github.com/zmb3/spotify"
)

const releaseRadarId = "37i9dQZEVXblHXYINKqgaL"

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

func getPlaylistInfo(client *spotify.Client, playlistID string) []spotify.PlaylistTrack {
	page, err := client.GetPlaylistTracks(spotify.ID(playlistID))
	if err != nil {
		log.Fatalf("couldn't get features playlists: %v", err)
	}

	return page.Tracks
}

func buildSongInfo(tracklist []spotify.PlaylistTrack) map[spotify.ID]song {
	songInfo := make(map[spotify.ID]song)

	for _, trackObj := range tracklist {
		track := trackObj.Track
		songInfo[track.ID] = song{Name: track.Name, ArtistName: track.Artists[0].Name,
			AlbumName: track.Album.Name, Popularity: track.Popularity}
	}
	return songInfo
}

func addAudioFeatures(client *spotify.Client, songInfo map[spotify.ID]song) map[spotify.ID]song {
	trackIDs := assembleTrackIDs(songInfo)
	tracksData, err := client.GetAudioFeatures(trackIDs...)
	if err != nil {
		log.Fatalf("couldn't get tracks info: %v", err)
	}
	for _, track := range tracksData {
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
	return songInfo
}

func printSongInfo(songsMap map[spotify.ID]song) {
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
	client := clientCredentialsAuth()
	tracklist := getPlaylistInfo(client, releaseRadarId)
	songInfo := buildSongInfo(tracklist)
	songInfo = addAudioFeatures(client, songInfo)

	printSongInfo(songInfo)
}
