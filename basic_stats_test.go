package main

import (
	"testing"

	"github.com/zmb3/spotify"
)

func mockPlaylistData() SongSet {
	playlist := SongSet{}
	playlist.data = make(map[spotify.ID]Song)
	playlist.data[spotify.ID("pics of you")] = Song{Valence: 50}
	playlist.data[spotify.ID("fascination st")] = Song{Valence: 45}
	playlist.data[spotify.ID("lullaby")] = Song{Valence: 30}
	return playlist
}

func mockSongSlice() []Song {
	songSlice := make([]Song, 0, 4) // this slice will have 3 elements; one test appends a 4th
	for i := 0; i < 3; i++ {
		valence := float32(10.0*i + 1)
		songSlice = append(songSlice, Song{Valence: valence})
	}
	return songSlice
}

func Test_sortByValence(t *testing.T) {
	sorted := sortByField(mockPlaylistData(), "Valence")
	assertEqual(t, (*sorted).songSlice[0].Valence, float32(50), "")
	assertEqual(t, (*sorted).songSlice[1].Valence, float32(45), "")
	assertEqual(t, (*sorted).songSlice[2].Valence, float32(30), "")
}

func Test_min(t *testing.T) {
	vr := StatReport{songSlice: mockSongSlice(), fieldName: "Valence"}
	assertEqual(t, vr.min(), float32(1), "")
}

func Test_max(t *testing.T) {
	vr := StatReport{songSlice: mockSongSlice(), fieldName: "Valence"}
	assertEqual(t, vr.max(), float32(21), "")
}

func Test_sum(t *testing.T) {
	vr := StatReport{songSlice: mockSongSlice(), fieldName: "Valence"}
	assertEqual(t, vr.sum(), float32(33), "")
}

func Test_median(t *testing.T) {
	vr := StatReport{songSlice: mockSongSlice(), fieldName: "Valence"}
	assertEqual(t, vr.median(), float32(11), "")
	vr.songSlice = append(vr.songSlice, Song{Valence: 50})
	assertEqual(t, vr.median(), float32(16), "")
}

func Test_mean(t *testing.T) {
	vr := StatReport{songSlice: mockSongSlice(), fieldName: "Valence"}
	assertEqual(t, vr.mean(), float32(11), "")
	vr.songSlice = append(vr.songSlice, Song{Valence: 39})
	assertEqual(t, vr.mean(), float32(18), "")
}

func Example_printStatReport() {
	playlistData := basicSongInfo()
	printStatReport(playlistData, "Valence")
	// Output: Valence Report | Min: 0 | Max: 0 | Mean: 0 | Median: 0
	// ID | Name | Artist | Album | Release Date | Danceability | Duration | Energy | Instrumentalness | Liveness | Popularity | Speechiness | Tempo | Valence
	// hello 2 | Pictures Of You | The Cure | Disintegration | 5-2-1989 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 |
}
