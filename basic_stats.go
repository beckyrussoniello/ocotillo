package main

import (
	"fmt"
	"sort"
)

type ByValence []Song

func (a ByValence) Len() int           { return len(a) }
func (a ByValence) Less(i, j int) bool { return a[i].Valence < a[j].Valence }
func (a ByValence) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func sortByValence(playlistData PlaylistData) []Song {
	songSlice := make([]Song, 0, len(playlistData))
	for _, v := range playlistData {
		songSlice = append(songSlice, v)
	}
	sort.Sort(ByValence(songSlice))
	return songSlice
}

type ValenceReport struct {
	songSlice []Song
}

func (vr *ValenceReport) min() float32 {
	return vr.songSlice[0].Valence
}

func (vr *ValenceReport) max() float32 {
	return vr.songSlice[len(vr.songSlice)-1].Valence
}

func (vr *ValenceReport) sum() float32 {
	var result float32
	for _, song := range vr.songSlice {
		result += song.Valence
	}
	return result
}

func (vr *ValenceReport) median() float32 {
	length := len(vr.songSlice)
	if length%2 == 0 {
		firstValence := vr.songSlice[(length/2)-1].Valence
		secondValence := vr.songSlice[length/2].Valence
		return (firstValence + secondValence) / 2.0
	} else {
		return vr.songSlice[(length-1)/2].Valence
	}
}

func (vr *ValenceReport) mean() float32 {
	return vr.sum() / float32(len(vr.songSlice))
}

func (vr *ValenceReport) print() {
	fmt.Printf("Valence Report | Min: %v | Max: %v | Mean: %v | Median: %v\n", vr.min(), vr.max(), vr.mean(), vr.median())
	printSongInfo(vr.songSlice)
}

func printValenceReport(playlistData PlaylistData) {
	vr := ValenceReport{songSlice: sortByValence(playlistData)}
	vr.print()
}
