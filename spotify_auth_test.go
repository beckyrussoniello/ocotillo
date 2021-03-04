package main

import "testing"

func TestAuth(t *testing.T) {
	if clientCredentialsAuth() == nil {
		t.Fatal("Should return playlists")
	}
}
