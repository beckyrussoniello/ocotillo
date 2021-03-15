package main

import (
	"sort"
	"fmt"
	"github.com/zmb3/spotify"
)
type ByValence []Song

func (a ByValence) Len() int           { return len(a) }
func (a ByValence) Less(i, j int) bool { return a[i].Valence < a[j].Valence }
func (a ByValence) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func sortByValence(playlistData map[spotify.ID]Song) []Song {
	songSlice := make([]Song, 0, len(playlistData))
	for _, v := range playlistData {
		songSlice = append(songSlice, v)
	}
    sort.Sort(ByValence(songSlice))
    fmt.Println(songSlice)
	return songSlice
}