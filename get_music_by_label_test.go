package main

import (
	"testing"

	"github.com/zmb3/spotify"
)

func simpleArtistSlice() []spotify.SimpleArtist {
	artistSlice := make([]spotify.SimpleArtist, 0, 1)
	artistSlice = append(artistSlice, spotify.SimpleArtist{Name: "The Cure"})
	return artistSlice
}

func simpleAlbumPage() *spotify.SimpleAlbumPage {
	albumSlice := make([]spotify.SimpleAlbum, 0, 1)
	albumSlice = append(albumSlice, spotify.SimpleAlbum{Name: "Disintegration"})
	albumPage := spotify.SimpleAlbumPage{
		Albums: albumSlice,
	}
	return &albumPage
}

func Test_getAlbumsData(t *testing.T) {
	albumsByYearSpan := AlbumsByYearSpan{
		years:         "2012-2021",
		albumIDs:      make([]spotify.ID, 0, 10),
		label:         "Fake Records",
		api:           SpotifyAPI{client: &mockSpotifyClient{}},
		gotAllResults: false,
	}
	albumsByYearSpan.getAlbumsData()
	assertEqual(t, len(albumsByYearSpan.albumIDs), 1, "")
	assertEqual(t, albumsByYearSpan.gotAllResults, true, "")
}

func Test_recordAlbumIDs(t *testing.T) {
	years := "2012-2021"
	albumIDs := make([]spotify.ID, 0, 3)
	results := spotify.SearchResult{Albums: &spotify.SimpleAlbumPage{
		Albums: make([]spotify.SimpleAlbum, 0, 3),
	}}
	results.Albums.Albums = append(
		results.Albums.Albums,
		spotify.SimpleAlbum{ID: spotify.ID("heeyy 1"), Artists: simpleArtistSlice()},
		spotify.SimpleAlbum{ID: spotify.ID("heeyy 2"), Artists: simpleArtistSlice()},
		spotify.SimpleAlbum{ID: spotify.ID("heeyy 3"), Artists: simpleArtistSlice()},
	)

	recordAlbumIDs(years, &albumIDs, &results)
	assertEqual(t, len(albumIDs), 3, "")
	assertEqual(t, albumIDs[0], spotify.ID("heeyy 1"), "")
}
