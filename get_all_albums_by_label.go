package main

import (
	"fmt"
	"log"

	"github.com/zmb3/spotify"
)

var maxLimit int = 50

const offsetMax int = 951

var yearSpans = []string{"1900-1971", "1972-1986", "1987-2001", "2002-2011", "2012-2021"}

func getAllAlbumsByLabel(recordLabelName string) {
	client := clientCredentialsAuth()
	for yearSpansIndex := 0; yearSpansIndex < len(yearSpans); yearSpansIndex++ {
		gotAllResults := false
		for offset := 0; offset < offsetMax && !gotAllResults; offset += maxLimit {
			years := yearSpans[yearSpansIndex]
			// search for albums from record label
			results, err := client.SearchOpt(fmt.Sprintf("label:\"%v\" year:%v", recordLabelName, years), spotify.SearchTypeAlbum, &spotify.Options{Limit: &maxLimit, Offset: &offset})
			if err != nil {
				log.Fatal(err)
			}
			printSearchResults(years, results)
			if results.Albums.Total <= (offset + maxLimit) {
				gotAllResults = true
			}
		}
	}

}

func printSearchResults(years string, results *spotify.SearchResult) {
	// handle album results
	if results.Albums != nil {
		fmt.Println("Albums (", years, "):")
		for _, item := range results.Albums.Albums {
			fmt.Println("   ", item.Artists[0].Name, "-", item.Name, " || ", item.ReleaseDate)
		}
	}

	// handle artist results
	if results.Artists != nil {
		fmt.Println("Artists:")
		for _, item := range results.Artists.Artists {
			fmt.Println("   ", item.Name, "(", item.ID.String(), ")", item.Genres)
		}
	}
}
