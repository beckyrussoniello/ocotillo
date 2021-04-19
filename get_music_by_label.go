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
	allTracks_SongSet := SongSet{}
	allTracks_SongSet.data = make(map[spotify.ID]Song)
	allTracks_SongSet.orderedKeys = make([]spotify.ID, 0, 10000)
	for yearSpansIndex := 0; yearSpansIndex < len(yearSpans); yearSpansIndex++ {
		albumsByYearSpan := AlbumsByYearSpan{
			years:         yearSpans[yearSpansIndex],
			albumIDs:      make([]spotify.ID, 0, 1000),
			label:         recordLabelName,
			api:           sp,
			gotAllResults: false,
		}
		albumsByYearSpan.getAlbumsData()
		albumsByYearSpan.getTracksForAlbums(&allTracks_SongSet)
	}

	sp.addTracksInfo(allTracks_SongSet)
	return &allTracks_SongSet
}

func (albs *AlbumsByYearSpan) getAlbumsData() {
	for offset := 0; offset < offsetMax && !albs.gotAllResults; offset += maxLimit {
		searchQuery := fmt.Sprintf("label:\"%v\" year:%v", albs.label, albs.years)
		options := &spotify.Options{Limit: &maxLimit, Offset: &offset}

		results, err := albs.api.client.SearchOpt(searchQuery, spotify.SearchTypeAlbum, options)
		if err != nil {
			log.Fatal(err)
		}
		albs.recordAlbumIDs(results)
		if results.Albums.Total <= (offset + maxLimit) {
			albs.gotAllResults = true
		}
	}
}

func (albs *AlbumsByYearSpan) getTracksForAlbums(trackInfo *SongSet) *SongSet {
	var albumsPerRequest = 20
	for offset := 0; offset < len(albs.albumIDs); offset += albumsPerRequest {
		var endIndex = offset + albumsPerRequest
		if len(albs.albumIDs) < endIndex {
			endIndex = len(albs.albumIDs)
		}
		var albumIDsForReq = albs.albumIDs[offset:endIndex]
		options := &spotify.Options{Limit: &albumsPerRequest, Offset: &offset}
		var albumsData, _ = albs.api.client.GetAlbumsOpt(options, albumIDsForReq...)
		for _, album := range albumsData {
			for _, track := range album.Tracks.Tracks {
				(*trackInfo).data[track.ID] = Song{ReleaseDate: album.ReleaseDate, ID: track.ID}
				trackInfo.orderedKeys = append(trackInfo.orderedKeys, track.ID)
			}
		}
	}
	return trackInfo
}

func (albs *AlbumsByYearSpan) recordAlbumIDs(results *spotify.SearchResult) {
	if results.Albums != nil {
		for _, item := range results.Albums.Albums {
			albs.albumIDs = append(albs.albumIDs, item.ID)
		}
	}
}
