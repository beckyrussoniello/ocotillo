package main

import (
	"fmt"
	"testing"

	"github.com/zmb3/spotify"
)

func theCureArtists() []spotify.SimpleArtist {
	artists := make([]spotify.SimpleArtist, 0, 1)
	artists = append(artists, spotify.SimpleArtist{
		Name:         "The Cure",
		ID:           spotify.ID("hello 1"),
		URI:          "",
		Endpoint:     "",
		ExternalURLs: map[string]string{},
	})
	return artists
}

func distintegrationAlbum() spotify.SimpleAlbum {
	return spotify.SimpleAlbum{
		Name:                 "Disintegration",
		Artists:              []spotify.SimpleArtist{},
		AlbumGroup:           "",
		AlbumType:            "",
		ID:                   spotify.ID("hello 3"),
		URI:                  "",
		AvailableMarkets:     []string{},
		Endpoint:             "",
		Images:               []spotify.Image{},
		ExternalURLs:         map[string]string{},
		ReleaseDate:          "",
		ReleaseDatePrecision: "",
	}
}

func spotifyTracks() []spotify.PlaylistTrack {
	tracks := make([]spotify.PlaylistTrack, 0, 1)
	tracks = append(tracks, spotify.PlaylistTrack{
		AddedAt: "",
		AddedBy: spotify.User{},
		IsLocal: false,
		Track: spotify.FullTrack{
			SimpleTrack: spotify.SimpleTrack{
				Artists:          theCureArtists(),
				AvailableMarkets: []string{},
				DiscNumber:       0,
				Duration:         0,
				Explicit:         false,
				ExternalURLs:     map[string]string{},
				Endpoint:         "",
				ID:               spotify.ID("hello 2"),
				Name:             "Pictures Of You",
				PreviewURL:       "",
				TrackNumber:      2,
				URI:              "",
			},
			Album: distintegrationAlbum(),
			ExternalIDs: map[string]string{},
			Popularity:  0,
			IsPlayable:  new(bool),
			LinkedFrom:  &spotify.LinkedFromInfo{},
		},
	})
	return tracks
}

type mockSpotifyClient struct{}

func (m *mockSpotifyClient) GetPlaylistTracks(playlistID spotify.ID) (*spotify.PlaylistTrackPage, error) {
	return &spotify.PlaylistTrackPage{
		Tracks: spotifyTracks(),
	}, nil
}

func (m *mockSpotifyClient) GetAudioFeatures(ids ...spotify.ID) ([]*spotify.AudioFeatures, error) {
	audioFeatures := make([]*spotify.AudioFeatures, 0, 1)
	audioFeatures = append(audioFeatures, &spotify.AudioFeatures{
		Acousticness:     5,
		AnalysisURL:      "",
		Danceability:     0,
		Duration:         0,
		Energy:           35,
		ID:               "hello 2",
		Instrumentalness: 40,
		Key:              0,
		Liveness:         10,
		Loudness:         0,
		Mode:             0,
		Speechiness:      0,
		Tempo:            0,
		TimeSignature:    0,
		TrackURL:         "",
		URI:              "",
		Valence:          0,
	})
	return audioFeatures, nil
}

func basicSongInfo() PlaylistData {
	data := make(PlaylistData, 1)
	data[spotify.ID("hello 2")] = Song{
		ID:         "hello 2",
		Name:       "Pictures Of You",
		ArtistName: "The Cure",
		AlbumName:  "Disintegration",
	}
	return data
}

func assertEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a == b {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v != %v", a, b)
	}
	t.Fatal(message)
}

func Test_buildBasicSongInfo(t *testing.T) {
	client := mockSpotifyClient{}
	pg := PlaylistGetter{client: &client}
	songInfo := pg.buildBasicSongInfo("fake playlist id")
	assertEqual(t, len(songInfo), len(spotifyTracks()), "")
	assertEqual(t, songInfo[spotify.ID("hello 2")].Name, "Pictures Of You", "")
	assertEqual(t, songInfo[spotify.ID("hello 2")].AlbumName, "Disintegration", "")
	assertEqual(t, songInfo[spotify.ID("hello 2")].ArtistName, "The Cure", "")
}

func Test_addAudioFeatures(t *testing.T) {
	client := mockSpotifyClient{}
	pg := PlaylistGetter{client: &client}
	songInfo := pg.addAudioFeatures(basicSongInfo())
	fmt.Println(songInfo)
	assertEqual(t, len(songInfo), len(spotifyTracks()), "")
	assertEqual(t, songInfo[spotify.ID("hello 2")].Name, "Pictures Of You", "")
	assertEqual(t, int(songInfo[spotify.ID("hello 2")].Energy), 35, "")
	assertEqual(t, int(songInfo[spotify.ID("hello 2")].Instrumentalness), 40, "")
}
