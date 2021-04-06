package main

import (
	"fmt"

	"log"

	"github.com/zmb3/spotify"
)

var userID = "1211800245"
var chunkSize = 100

type LabelPlaylist struct {
	labelName    string
	playlistType string
	api          SpotifyAPI
	spotifyID    spotify.ID
	trackIDs     []spotify.ID
}

func (sp *SpotifyAPI) createLabelPlaylist(labelName string, playlistType string) *LabelPlaylist {
	labelPlaylist := LabelPlaylist{
		labelName:    labelName,
		playlistType: playlistType,
		api:          *sp,
	}

	playlist, err := labelPlaylist.api.client.CreatePlaylistForUser(userID, labelPlaylist.title(), labelPlaylist.description(), true)
	if err != nil {
		log.Fatal(err)
	}
	labelPlaylist.spotifyID = playlist.ID
	allTracks := labelPlaylist.api.getAllAlbumsByLabel(labelName)
	statReport := StatReport{sortByField(*allTracks, "Popularity"), "Popularity"}
	sortedTracks := statReport.toSongSet()
	labelPlaylist.addTracksToPlaylist(sortedTracks)
	return &labelPlaylist
}

func (labelPlaylist LabelPlaylist) title() string {
	title := fmt.Sprintf("%v Records", labelPlaylist.labelName)
	if labelPlaylist.playlistType == "Top" {
		title = title + " - Top Tracks"
	}
	return title
}

func (labelPlaylist LabelPlaylist) description() string {
	return fmt.Sprintf("%v releases from %v Records", labelPlaylist.playlistType, labelPlaylist.labelName)
}

func (labelPlaylist *LabelPlaylist) addTracksToPlaylist(allTracks *SongSet) {
	trackChunks := trackIDChunks(*allTracks, chunkSize)
	for _, chunk := range trackChunks {
		labelPlaylist.api.client.AddTracksToPlaylist(labelPlaylist.spotifyID, chunk...)
		labelPlaylist.trackIDs = append(labelPlaylist.trackIDs, chunk...)
	}
}
