package main

import (
	"fmt"

	"log"

	"github.com/zmb3/spotify"
)

var userID = "1211800245"
var chunkSize = 100

func createLabelPlaylist(labelName string, playlistType string) {
	sp := SpotifyAPI{client: userAuth()}
	title := fmt.Sprintf("%v Records", labelName)
	if playlistType == "Top" {
		title = title + " - Top Tracks"
	}
	// note: should handle error case below, don't just throw away
	playlist, err := sp.client.CreatePlaylistForUser(userID, title, playlistDescription(labelName, playlistType), true)
	if err != nil {
		log.Fatal(err)
	}
	allTracks := sp.getAllAlbumsByLabel(labelName)
	trackChunks := trackIDChunks(*allTracks, chunkSize)
	for _, chunk := range trackChunks {
		for i := 0; i < len(chunk); i += chunkSize { // max trackIDs to add to playlist in one req
			trackIDsForReq := make([]spotify.ID, 0, chunkSize)
			maxJ := chunkSize + i
			if len(chunk) < maxJ {
				maxJ = len(chunk)
			}
			for j := 0 + i; j < maxJ; j++ {
				trackIDsForReq = append(trackIDsForReq, chunk[j])
			}
			fmt.Println(playlist)
			fmt.Println(trackIDsForReq)
			sp.client.AddTracksToPlaylist(playlist.ID, trackIDsForReq...)
		}
	}

}

func playlistDescription(labelName string, playlistType string) string {
	return fmt.Sprintf("%v releases from %v Records", playlistType, labelName)
}
