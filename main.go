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

var audioFeaturesChunkSize = 100
var tracksInfoChunkSize = 50
var playlistLengthMax = 10000
var topTracksMaxLength = 500

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
	GetTracks(ids ...spotify.ID) ([]*spotify.FullTrack, error)
}

type Song struct {
	ID               spotify.ID
	Name             string
	ArtistName       string
	AlbumName        string
	ReleaseDate      string
	Danceability     float32
	Duration         int
	Energy           float32
	Instrumentalness float32
	Liveness         float32
	Popularity       float32
	Speechiness      float32
	Tempo            int
	Valence          float32
}

type SongSet struct {
	data        map[spotify.ID]Song
	orderedKeys []spotify.ID
}

func (sp *SpotifyAPI) getPlaylistInfo(playlistID string) []spotify.PlaylistTrack {
	page, err := sp.client.GetPlaylistTracks(spotify.ID(playlistID))
	if err != nil {
		log.Fatalf("couldn't get features playlists: %v", err)
	}

	return page.Tracks
}

func (sp *SpotifyAPI) buildBasicSongInfo(playlistID string) SongSet {
	tracklist := sp.getPlaylistInfo(playlistID)
	songInfo := SongSet{}
	songInfo.data = make(map[spotify.ID]Song)
	songInfo.orderedKeys = make([]spotify.ID, 0, playlistLengthMax)

	for _, trackObj := range tracklist {
		track := trackObj.Track
		songInfo.data[track.ID] = Song{Name: track.Name, ArtistName: track.Artists[0].Name,
			AlbumName: track.Album.Name, Popularity: float32(track.Popularity)}
		songInfo.orderedKeys = append(songInfo.orderedKeys, track.ID)
	}
	return songInfo
}

func (sp *SpotifyAPI) addTracksInfo(songInfo SongSet) SongSet {
	trackIDChunks := trackIDChunks(songInfo, tracksInfoChunkSize)
	for _, trackIDs := range trackIDChunks {
		tracksData, err := sp.client.GetTracks(trackIDs...)

		if err != nil {
			log.Fatalf("couldn't get tracks info: %v", err)
		}
		for _, track := range tracksData {
			song := songInfo.data[track.ID]
			song.Popularity = float32(track.Popularity)
			// what else?
			songInfo.data[track.ID] = song
		}
	}

	return songInfo
}

func (sp *SpotifyAPI) addAudioFeatures(songInfo SongSet) SongSet {
	trackIDChunks := trackIDChunks(songInfo, audioFeaturesChunkSize)
	for _, trackIDs := range trackIDChunks {
		tracksData, err := sp.client.GetAudioFeatures(trackIDs...)

		if err != nil {
			fmt.Println("--ERROR--")
			fmt.Println(len(trackIDs), "track ids passed in")
			log.Fatalf("couldn't get audio features: %v", err)
		}
		for _, track := range tracksData {
			song := songInfo.data[track.ID]
			song.Danceability = track.Danceability
			song.Duration = track.Duration
			song.Energy = track.Energy
			song.Instrumentalness = track.Instrumentalness
			song.Liveness = track.Liveness
			song.Speechiness = track.Speechiness
			song.Tempo = int(track.Tempo)
			song.Valence = track.Valence
			songInfo.data[track.ID] = song
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
	allChunks := make([][]spotify.ID, 0, (len(songInfo.data)/chunkSize)+1)
	batchKeys := make([]spotify.ID, 0, chunkSize)

	for _, k := range songInfo.orderedKeys {
		batchKeys = append(batchKeys, k)
		if len(batchKeys) == chunkSize {
			batchCopy := make([]spotify.ID, len(batchKeys))
			copy(batchCopy, batchKeys)
			allChunks = append(allChunks, batchCopy)
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
	api.createLabelPlaylist("Sacred Bones", "Top")

}
