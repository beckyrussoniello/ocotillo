package main

import (
	"sort"
	"fmt"
)
type ByValence []Song

func (a ByValence) Len() int           { return len(a) }
func (a ByValence) Less(i, j int) bool { return a[i].Valence < a[j].Valence }
func (a ByValence) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func sum(songSlice []Song) float32 {
	var result float32
	for _, song := range songSlice {
		result += song.Valence
	}
	return result
}

func median(songSlice []Song) float32 {
	length := len(songSlice)
	if length % 2 == 0 {
		return (songSlice[length/2].Valence + songSlice[(length/2)-1].Valence) / 2
	}else{
		return songSlice[(length-1)/2].Valence
	}
}

func sortByValence(playlistData PlaylistData) []Song {
	songSlice := make([]Song, 0, len(playlistData))
	for _, v := range playlistData {
		songSlice = append(songSlice, v)
	}
    sort.Sort(ByValence(songSlice))
	return songSlice
}

func valenceReport(playlistData PlaylistData) {
	sorted := sortByValence(playlistData)
	min := sorted[0].Valence
	max := sorted[len(sorted) - 1].Valence
	mean := sum(sorted) / float32(len(sorted))
	fmt.Printf("Valence Report | Min: %v | Max: %v | Mean: %v | Median: %v\n", min, max, mean, median(sorted))
	printSongInfo(sorted)
}