package main

import (
	"fmt"
	"reflect"
	"sort"
)

type StatReport struct {
	songSlice []Song
	fieldName string
}

func (sr StatReport) Len() int { return len(sr.songSlice) }

func (sr StatReport) Less(i, j int) bool {
	iFieldValue := reflect.ValueOf(sr.songSlice[i]).FieldByName(sr.fieldName).Float()
	jFieldValue := reflect.ValueOf(sr.songSlice[j]).FieldByName(sr.fieldName).Float()
	return iFieldValue < jFieldValue
}

func (sr StatReport) Swap(i, j int) {
	sr.songSlice[i], sr.songSlice[j] = sr.songSlice[j], sr.songSlice[i]
}

func sortByField(songSet SongSet, field string) []Song {
	songSlice := make([]Song, 0, len(songSet))
	for _, v := range songSet {
		songSlice = append(songSlice, v)
	}
	sort.Sort(StatReport{songSlice, field})
	return songSlice
}

func (sr *StatReport) toSongSet() *SongSet {
	songSet := make(SongSet)
	for _, song := range sr.songSlice {
		songSet[song.ID] = song
	}
	return &songSet
}

func (sr *StatReport) min() float32 {
	return float32(reflect.ValueOf(sr.songSlice[0]).FieldByName(sr.fieldName).Float())
}

func (sr *StatReport) max() float32 {
	return float32(reflect.ValueOf(sr.songSlice[len(sr.songSlice)-1]).FieldByName(sr.fieldName).Float())
}

func (sr *StatReport) sum() float32 {
	var result float32
	for _, song := range sr.songSlice {
		result += float32(reflect.ValueOf(song).FieldByName(sr.fieldName).Float())
	}
	return result
}

func (sr *StatReport) median() float32 {
	length := len(sr.songSlice)
	if length%2 == 0 {
		firstMiddle := reflect.ValueOf(sr.songSlice[(length/2)-1]).FieldByName(sr.fieldName).Float()
		secondMiddle := reflect.ValueOf(sr.songSlice[length/2]).FieldByName(sr.fieldName).Float()
		return float32((firstMiddle + secondMiddle) / 2.0)
	} else {
		return float32(reflect.ValueOf(sr.songSlice[(length-1)/2]).FieldByName(sr.fieldName).Float())
	}
}

func (sr *StatReport) mean() float32 {
	return sr.sum() / float32(len(sr.songSlice))
}

func (sr *StatReport) print() {
	fmt.Printf("%v Report | Min: %v | Max: %v | Mean: %v | Median: %v\n", sr.fieldName, sr.min(), sr.max(), sr.mean(), sr.median())
	printSongInfo(sr.songSlice)
}

func printStatReport(songSet SongSet, fieldName string) {
	vr := StatReport{sortByField(songSet, fieldName), fieldName}
	vr.print()
}
