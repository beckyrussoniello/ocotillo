package main

import (
	"fmt"
	"log"
	"reflect"

	"github.com/zmb3/spotify"
)

const releaseRadarId = "37i9dQZEVXblHXYINKqgaL"
const discoverWeeklyId = "37i9dQZEVXcHTLYCsuh2O5"
const meditationId = "1xQ9JBNlCSza7iZ4AXhnIL"

var defaultChunkSize = 100

type SpotifyAPI struct {
	client spotifyClient
}

type spotifyClient interface {
	GetPlaylistTracks(playlistID spotify.ID) (*spotify.PlaylistTrackPage, error)
	GetAudioFeatures(ids ...spotify.ID) ([]*spotify.AudioFeatures, error)
	CreatePlaylistForUser(userID, playlistName, description string, public bool) (*spotify.FullPlaylist, error)
	AddTracksToPlaylist(playlistID spotify.ID, trackIDs ...spotify.ID) (snapshotID string, err error)
	SearchOpt(query string, t spotify.SearchType, opt *spotify.Options) (*spotify.SearchResult, error)
	GetAlbumsOpt(opt *spotify.Options, ids ...spotify.ID) ([]*spotify.FullAlbum, error)
}

type Song struct {
	ID               string
	Name             string
	ArtistName       string
	AlbumName        string
	ReleaseDate      string
	Danceability     float32
	Duration         int
	Energy           float32
	Instrumentalness float32
	Liveness         float32
	Popularity       int
	Speechiness      float32
	Tempo            int
	Valence          float32
}

type SongSet map[spotify.ID]Song

func (sp *SpotifyAPI) getPlaylistInfo(playlistID string) []spotify.PlaylistTrack {
	page, err := sp.client.GetPlaylistTracks(spotify.ID(playlistID))
	if err != nil {
		log.Fatalf("couldn't get features playlists: %v", err)
	}

	return page.Tracks
}

func (sp *SpotifyAPI) buildBasicSongInfo(playlistID string) SongSet {
	tracklist := sp.getPlaylistInfo(playlistID)
	songInfo := make(SongSet)

	for _, trackObj := range tracklist {
		track := trackObj.Track
		songInfo[track.ID] = Song{Name: track.Name, ArtistName: track.Artists[0].Name,
			AlbumName: track.Album.Name, Popularity: track.Popularity}
	}
	return songInfo
}

func (sp *SpotifyAPI) addAudioFeatures(songInfo SongSet) SongSet {
	trackIDChunks := trackIDChunks(songInfo, defaultChunkSize)
	for _, trackIDs := range trackIDChunks {
		tracksData, err := sp.client.GetAudioFeatures(trackIDs...)

		if err != nil {
			fmt.Println("--ERROR--")
			fmt.Println(len(trackIDs), "track ids passed in")
			log.Fatalf("couldn't get audio features: %v", err)
		}
		for _, track := range tracksData {
			song := songInfo[track.ID]
			song.Danceability = track.Danceability
			song.Duration = track.Duration
			song.Energy = track.Energy
			song.Instrumentalness = track.Instrumentalness
			song.Liveness = track.Liveness
			song.Speechiness = track.Speechiness
			song.Tempo = int(track.Tempo)
			song.Valence = track.Valence
			songInfo[track.ID] = song
			fmt.Println(song)
		}
	}

	return songInfo
}

func printSongInfo(songsMap []Song) {
	fmt.Printf("ID | Name | Artist | Album | Release Date | Danceability | Duration | Energy | Instrumentalness | Liveness | Popularity | Speechiness | Tempo | Valence\n")
	for _, song := range songsMap {
		value := reflect.ValueOf(song)
		for attr := 0; attr < value.NumField(); attr++ {
			fmt.Printf("%v | ", value.Field(attr).Interface())
		}
		fmt.Println()

	}
}

func trackIDChunks(songInfo SongSet, chunkSize int) [][]spotify.ID {
	allChunks := make([][]spotify.ID, 0, (len(songInfo)/chunkSize)+1)
	batchKeys := make([]spotify.ID, 0, chunkSize)

	for k := range songInfo {
		fmt.Println(k)
		batchKeys = append(batchKeys, k)
		if len(batchKeys) == chunkSize {
			allChunks = append(allChunks, batchKeys)
			batchKeys = batchKeys[:0]
		}
	}

	if len(batchKeys) > 0 {
		allChunks = append(allChunks, batchKeys)
	}

	return allChunks
}

func main() {
	api := SpotifyAPI{client: userAuth()}
	api.createLabelPlaylist("Kompakt", "All")
}
