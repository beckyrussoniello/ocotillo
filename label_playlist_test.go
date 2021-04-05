package main

import (
	"fmt"
	"testing"

	"github.com/zmb3/spotify"
)

func mockLabelPlaylist(playlistType string) LabelPlaylist {
	return LabelPlaylist{
		labelName:    "Fake",
		playlistType: playlistType,
		api:          SpotifyAPI{client: &mockSpotifyClient{}},
	}
}

func Test_title_Top(t *testing.T) {
	lpl := mockLabelPlaylist("Top")
	title := lpl.title()
	assertEqual(t, title, "Fake Records - Top Tracks", "")
}

func Test_title_All(t *testing.T) {
	lpl := mockLabelPlaylist("All")
	title := lpl.title()
	assertEqual(t, title, "Fake Records", "")
}

func Test_description_Top(t *testing.T) {
	lpl := mockLabelPlaylist("Top")
	desc := lpl.description()
	assertEqual(t, desc, "Top releases from Fake Records", "")
}

func Test_description_All(t *testing.T) {
	lpl := mockLabelPlaylist("All")
	desc := lpl.description()
	assertEqual(t, desc, "All releases from Fake Records", "")
}

func Test_addChunkTracksToPlaylist_256(t *testing.T) {
	lpl := mockLabelPlaylist("All")
	songset := SongSet{}
	for i := 0; i < 256; i++ {
		songset[spotify.ID(fmt.Sprintf("%v", i))] = Song{}
	}
	lpl.addTracksToPlaylist(&songset)
	assertEqual(t, len(lpl.trackIDs), 256, "")
}

func Test_addChunkTracksToPlaylist_1(t *testing.T) {
	lpl := mockLabelPlaylist("All")
	songset := SongSet{}
	songset[spotify.ID("1")] = Song{}
	lpl.addTracksToPlaylist(&songset)
	assertEqual(t, len(lpl.trackIDs), 1, "")
	assertEqual(t, (lpl.trackIDs)[0], spotify.ID("1"), "")
}

func Test_addTracksToPlaylist_100(t *testing.T) {
	lpl := mockLabelPlaylist("All")
	songset := SongSet{}
	for i := 0; i < 100; i++ {
		songset[spotify.ID(fmt.Sprintf("%v", i))] = Song{}
	}
	lpl.addTracksToPlaylist(&songset)
	assertEqual(t, len(lpl.trackIDs), 100, "")
}

func Test_createLabelPlaylist(t *testing.T) {
	labelName := "Just Testing"
	playlistType := "Top"
	api := SpotifyAPI{client: &mockSpotifyClient{}}
	lpl := api.createLabelPlaylist(labelName, playlistType)
	fmt.Println(lpl)
	assertEqual(t, lpl.labelName, labelName, "")
	assertEqual(t, lpl.playlistType, playlistType, "")
	assertEqual(t, lpl.spotifyID, spotify.ID("hello"), "")
	assertEqual(t, len(lpl.trackIDs), 20, "")
}
