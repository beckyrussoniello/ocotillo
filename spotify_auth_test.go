package main

import "testing"

func TestClientCredentials(t *testing.T) {
	if clientCredentialsAuth() == nil {
		t.Fatal("Should return playlists")
	}
}

func TestUserAuth(t *testing.T) {
	if userAuth() == nil {
		t.Fatal("Should return playlists")
	}
}
