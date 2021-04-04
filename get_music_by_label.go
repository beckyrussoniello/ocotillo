package main

import (
	"fmt"
	"log"

	"github.com/zmb3/spotify"
)

var maxLimit int = 50

const offsetMax int = 950

var yearSpans = []string{"1900-1971", "1972-1986", "1987-2001", "2002-2011", "2012-2021"}

func (sp SpotifyAPI) getAllAlbumsByLabel(recordLabelName string) *SongSet {
	var allTrackIDs SongSet = make(SongSet)
	for yearSpansIndex := 0; yearSpansIndex < len(yearSpans); yearSpansIndex++ {
		gotAllResults := false
		years := yearSpans[yearSpansIndex]
		var albumIDs []spotify.ID = make([]spotify.ID, 0, 1000)

		for offset := 0; offset < offsetMax && !gotAllResults; offset += maxLimit {
			searchQuery := fmt.Sprintf("label:\"%v\" year:%v", recordLabelName, years)
			options := &spotify.Options{Limit: &maxLimit, Offset: &offset}

			results, err := sp.client.SearchOpt(searchQuery, spotify.SearchTypeAlbum, options)
			if err != nil {
				log.Fatal(err)
			}
			recordAlbumIDs(years, &albumIDs, results)
			if results.Albums.Total <= (offset + maxLimit) {
				gotAllResults = true
			}
		}

		sp.getTracksForAlbums(&albumIDs, allTrackIDs)
	}

	sp.addAudioFeatures(allTrackIDs)
	return &allTrackIDs
}

func (sp *SpotifyAPI) getTracksForAlbums(albumIDs *[]spotify.ID, trackInfo SongSet) SongSet {
	var albumsPerRequest = 20
	for offset := 0; offset < len(*albumIDs); offset += albumsPerRequest {
		var endIndex = offset + albumsPerRequest
		if len(*albumIDs) < endIndex {
			endIndex = len(*albumIDs)
		}
		var albumIDsForReq = (*albumIDs)[offset:endIndex]
		options := &spotify.Options{Limit: &albumsPerRequest, Offset: &offset}
		var albumsData, _ = sp.client.GetAlbumsOpt(options, albumIDsForReq...)
		fmt.Println(len(albumIDsForReq), "albums requested.", len(albumsData), "albums received.")

		for _, album := range albumsData {
			fmt.Println("Album", album.Name, "has", len(album.Tracks.Tracks), "tracks.")
			for _, track := range album.Tracks.Tracks {
				trackInfo[track.ID] = Song{ReleaseDate: album.ReleaseDate}
			}
		}
	}
	return trackInfo
}

func recordAlbumIDs(years string, albumIDs *[]spotify.ID, results *spotify.SearchResult) {
	if results.Albums != nil {
		fmt.Println("Albums (", years, "):")
		for _, item := range results.Albums.Albums {
			*albumIDs = append(*albumIDs, item.ID)
			fmt.Println("   ", item.Artists[0].Name, "-", item.Name, " || ", item.ReleaseDate)
		}
	}
}
