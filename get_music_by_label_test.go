package main

import (
	"fmt"
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
	albumSlice = append(albumSlice, spotify.SimpleAlbum{Name: "Disintegration", ID: spotify.ID("iddd")})
	albumPage := spotify.SimpleAlbumPage{
		Albums: albumSlice,
	}
	return &albumPage
}

func Test_getAllAlbumsByLabel(t *testing.T) {
	sp := SpotifyAPI{client: &mockSpotifyClient{}}
	recordLabelName := "Test-it-all"
	testSongSet := sp.getAllAlbumsByLabel(recordLabelName)
	assertEqual(t, len(*testSongSet), 20, "")
	assertEqual(t, (*testSongSet)[spotify.ID("0")].ReleaseDate, "3/1/2000", "")
	assertEqual(t, (*testSongSet)[spotify.ID("0")].Energy, float32(35), "")
	assertEqual(t, (*testSongSet)[spotify.ID("0")].Liveness, float32(10), "")
	assertEqual(t, (*testSongSet)[spotify.ID("0")].ID, spotify.ID("0"), "")
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

func Test_getTracksForAlbums(t *testing.T) {
	albs := AlbumsByYearSpan{
		api:      SpotifyAPI{client: &mockSpotifyClient{}},
		albumIDs: make([]spotify.ID, 0, 20),
	}
	for i := 0; i < 20; i++ {
		albs.albumIDs = append(albs.albumIDs, spotify.ID(fmt.Sprintf("%v", i)))
	}
	songSet := SongSet{}
	albs.getTracksForAlbums(&songSet)
	assertEqual(t, len(songSet), 20, "")
	assertEqual(t, songSet[spotify.ID("0")].ReleaseDate, "3/1/2000", "")
}

func Test_recordAlbumIDs(t *testing.T) {
	albs := AlbumsByYearSpan{albumIDs: make([]spotify.ID, 0, 3)}
	results := spotify.SearchResult{Albums: &spotify.SimpleAlbumPage{
		Albums: make([]spotify.SimpleAlbum, 0, 3),
	}}
	results.Albums.Albums = append(
		results.Albums.Albums,
		spotify.SimpleAlbum{ID: spotify.ID("heeyy 1"), Artists: simpleArtistSlice()},
		spotify.SimpleAlbum{ID: spotify.ID("heeyy 2"), Artists: simpleArtistSlice()},
		spotify.SimpleAlbum{ID: spotify.ID("heeyy 3"), Artists: simpleArtistSlice()},
	)

	albs.recordAlbumIDs(&results)
	assertEqual(t, len(albs.albumIDs), 3, "")
	assertEqual(t, albs.albumIDs[0], spotify.ID("heeyy 1"), "")
}
