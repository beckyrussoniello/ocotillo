# ocotillo
Exploratory app as I learn Go, using the Spotify API.

Right now, it:
1) Takes the ID of a Spotify Playlist
2) Makes an API call to get the tracklist, including some basic attributes of each song
3) Creates a map (data structure) to store this information, where each key is an individual song's Spotify ID, and the value is a struct
4) Makes a second API call to get more detailed "audio features" of all the songs
5) Further populates the data structure with these additional fields
6) Prints out all of the playlist's track data with some formatting
7) Creates a ValenceReport struct, with basic stats methods (min, max, sum, mean, median) and prints it out
