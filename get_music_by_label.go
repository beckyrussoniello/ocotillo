package main

import (
	"fmt"
	"log"

	"github.com/zmb3/spotify"
)

var maxLimit int = 50

const offsetMax int = 950

var yearSpans = []string{"1900-1971", "1972-1986", "1987-2001", "2002-2011", "2012-2021"}

type AlbumsByYearSpan struct {
	years         string
	albumIDs      []spotify.ID
	label         string
	api           SpotifyAPI
	gotAllResults bool
}

func (sp SpotifyAPI) getAllAlbumsByLabel(recordLabelName string) *SongSet {
	var allTracks_SongSet SongSet = make(SongSet)
	for yearSpansIndex := 0; yearSpansIndex < len(yearSpans); yearSpansIndex++ {
		albumsByYearSpan := AlbumsByYearSpan{
			years:         yearSpans[yearSpansIndex],
			albumIDs:      make([]spotify.ID, 0, 1000),
			label:         recordLabelName,
			api:           sp,
			gotAllResults: false,
		}
		albumsByYearSpan.getAlbumsData()
		sp.getTracksForAlbums(&albumsByYearSpan.albumIDs, allTracks_SongSet)
	}

	sp.addAudioFeatures(allTracks_SongSet)
	return &allTracks_SongSet
}

func (albs AlbumsByYearSpan) getAlbumsData() {
	for offset := 0; offset < offsetMax && !albs.gotAllResults; offset += maxLimit {
		searchQuery := fmt.Sprintf("label:\"%v\" year:%v", albs.label, albs.years)
		options := &spotify.Options{Limit: &maxLimit, Offset: &offset}

		results, err := albs.api.client.SearchOpt(searchQuery, spotify.SearchTypeAlbum, options)
		if err != nil {
			log.Fatal(err)
		}
		recordAlbumIDs(albs.years, &albs.albumIDs, results)
		if results.Albums.Total <= (offset + maxLimit) {
			albs.gotAllResults = true
		}
	}
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

		for _, album := range albumsData {
			for _, track := range album.Tracks.Tracks {
				trackInfo[track.ID] = Song{ReleaseDate: album.ReleaseDate}
			}
		}
	}
	return trackInfo
}

func recordAlbumIDs(years string, albumIDs *[]spotify.ID, results *spotify.SearchResult) {
	if results.Albums != nil {
		for _, item := range results.Albums.Albums {
			*albumIDs = append(*albumIDs, item.ID)
		}
	}
}
