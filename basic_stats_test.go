package main

import (
	"testing"
	"github.com/zmb3/spotify"
)

func mockPlaylistData() PlaylistData {
	playlist := make(PlaylistData)
	playlist[spotify.ID("pics of you")] = Song{ Valence: 50}
	playlist[spotify.ID("fascination st")] = Song{ Valence: 45}
	playlist[spotify.ID("lullaby")] = Song{ Valence: 30}
	return playlist
}

func Test_sortByValence(t *testing.T) {
	sorted := sortByValence(mockPlaylistData())
	assertEqual(t, int(sorted[0].Valence), 30, "")
	assertEqual(t, int(sorted[1].Valence), 45, "")
	assertEqual(t, int(sorted[2].Valence), 50, "")
}