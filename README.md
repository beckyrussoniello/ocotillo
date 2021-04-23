# ocotillo
Exploratory app as I learn Go, using the Spotify API.

Current Features:
1) Label Playlists
- Search Spotify for all albums released by a particular label
- Sort all tracks on these albums by popularity (or another attribute)
- Create a new playlist, either with all tracks from the label, or top tracks (top 10% by attribute)
2) Stat Report
- Take a set of songs (either from an existing playlist, or the label search above)
- Print out basic stats about a given audio feature (valence, energy, liveness, etc)
- Information includes min, max, sum, mean, median
3) Song Set
- SongSet is a struct which stores information for a set of tracks
- A track can be looked up by its Spotify ID. This key is associated with a Song struct containing data for the track
- Stores data from several different Spotify API endpoints
