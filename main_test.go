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
			Album:       distintegrationAlbum(),
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

func (m *mockSpotifyClient) GetTracks(ids ...spotify.ID) ([]*spotify.FullTrack, error) {
	tracks := make([]*spotify.FullTrack, 0, 1)
	for _, id := range ids {
		tracks = append(tracks, &spotify.FullTrack{
			SimpleTrack: spotify.SimpleTrack{ID: id},
			Album:       spotify.SimpleAlbum{},
			ExternalIDs: map[string]string{},
			Popularity:  44,
			IsPlayable:  new(bool),
			LinkedFrom:  &spotify.LinkedFromInfo{},
		})
	}
	return tracks, nil
}

func (m *mockSpotifyClient) GetAudioFeatures(ids ...spotify.ID) ([]*spotify.AudioFeatures, error) {
	audioFeatures := make([]*spotify.AudioFeatures, 0, 1)
	for _, id := range ids {
		audioFeatures = append(audioFeatures, &spotify.AudioFeatures{
			Acousticness:     5,
			AnalysisURL:      "",
			Danceability:     0,
			Duration:         0,
			Energy:           35,
			ID:               id,
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
	}

	return audioFeatures, nil
}

func (m *mockSpotifyClient) CreatePlaylistForUser(userID, playlistName, description string, public bool) (*spotify.FullPlaylist, error) {
	return &spotify.FullPlaylist{
		SimplePlaylist: spotify.SimplePlaylist{ID: spotify.ID("hello")},
	}, nil
}

func (m *mockSpotifyClient) AddTracksToPlaylist(playlistID spotify.ID, trackIDs ...spotify.ID) (snapshotID string, err error) {
	return "", nil
}

func (m *mockSpotifyClient) SearchOpt(query string, t spotify.SearchType, opt *spotify.Options) (*spotify.SearchResult, error) {
	return &spotify.SearchResult{
		Artists:   &spotify.FullArtistPage{},
		Albums:    simpleAlbumPage(),
		Playlists: &spotify.SimplePlaylistPage{},
		Tracks:    &spotify.FullTrackPage{},
	}, nil
}

func (m *mockSpotifyClient) GetAlbumsOpt(opt *spotify.Options, ids ...spotify.ID) ([]*spotify.FullAlbum, error) {
	return fullAlbumSlice(), nil
}

func fullAlbumSlice() []*spotify.FullAlbum {
	albumSlice := make([]*spotify.FullAlbum, 0, 20)
	for i := 0; i < 20; i++ {
		id := spotify.ID(fmt.Sprintf("%v", i))
		releaseDate := fmt.Sprintf("3/%v/2000", i+1)
		fullAlbum := spotify.FullAlbum{
			SimpleAlbum: spotify.SimpleAlbum{ReleaseDate: releaseDate},
			Tracks:      spotify.SimpleTrackPage{Tracks: []spotify.SimpleTrack{{ID: id}}},
		}
		albumSlice = append(albumSlice, &fullAlbum)
	}
	return albumSlice
}

func basicSongInfo() SongSet {
	songSet := SongSet{}
	songSet.data = make(map[spotify.ID]Song)
	songSet.orderedKeys = make([]spotify.ID, 0, 10000)
	songSet.data[spotify.ID("hello 2")] = Song{
		ID:          "hello 2",
		Name:        "Pictures Of You",
		ArtistName:  "The Cure",
		AlbumName:   "Disintegration",
		ReleaseDate: "5-2-1989",
	}
	songSet.orderedKeys = append(songSet.orderedKeys, spotify.ID("hello 2"))
	return songSet
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
	pg := SpotifyAPI{client: &client}
	songInfo := pg.buildBasicSongInfo("fake playlist id")
	assertEqual(t, len(songInfo.data), len(spotifyTracks()), "")
	assertEqual(t, songInfo.data[spotify.ID("hello 2")].Name, "Pictures Of You", "")
	assertEqual(t, songInfo.data[spotify.ID("hello 2")].AlbumName, "Disintegration", "")
	assertEqual(t, songInfo.data[spotify.ID("hello 2")].ArtistName, "The Cure", "")
}

func Test_getTracksInfo(t *testing.T) {
	client := mockSpotifyClient{}
	pg := SpotifyAPI{client: &client}
	songInfo := pg.addTracksInfo(basicSongInfo())
	assertEqual(t, len(songInfo.data), len(spotifyTracks()), "")
	assertEqual(t, songInfo.data[spotify.ID("hello 2")].Popularity, float32(44), "")
}

func Test_addAudioFeatures(t *testing.T) {
	client := mockSpotifyClient{}
	pg := SpotifyAPI{client: &client}
	songInfo := pg.addAudioFeatures(basicSongInfo())
	assertEqual(t, len(songInfo.data), len(spotifyTracks()), "")
	assertEqual(t, songInfo.data[spotify.ID("hello 2")].Name, "Pictures Of You", "")
	assertEqual(t, int(songInfo.data[spotify.ID("hello 2")].Energy), 35, "")
	assertEqual(t, int(songInfo.data[spotify.ID("hello 2")].Instrumentalness), 40, "")
}

func Test_trackIDChunks(t *testing.T) {
	songSet := SongSet{}
	songSet.data = make(map[spotify.ID]Song)
	songSet.orderedKeys = make([]spotify.ID, 0, 10)
	for i := 0; i < 10; i++ {
		spotifyID := spotify.ID(fmt.Sprintf("hello %v", i))
		songSet.data[spotifyID] = Song{}
		songSet.orderedKeys = append(songSet.orderedKeys, spotifyID)
	}
	trackIDs := trackIDChunks(songSet, 3)
	assertEqual(t, len(trackIDs), 4, "")
	assertEqual(t, len(trackIDs[0]), 3, "")
	assertEqual(t, trackIDs[0][0], spotify.ID("hello 0"), "")
}
