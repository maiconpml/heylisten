package goytmusic

import "github.com/tidwall/gjson"

type TracksService service

type Track struct {
	Name     string
	Artists  []*User
	VideoID  *string
	Duration string
	Album    *Album
}

// Parses res into a Track struct
// Expects the track contained in browseId=VLPL... endpoint JSON response
func extractTrack(res *gjson.Result) *Track {
	tr := &Track{}

	tr.Name = res.Get(pathPlaylistTrack + ".0." + pathTrackAttribute + ".0." + pathTrackName).String()
	buf := res.Get(pathPlaylistTrack + ".0." + pathTrackAttribute + ".0." + pathNavEndpointVideoID)
	if buf.Exists() {
		tr.VideoID = Ptr(buf.String())
	}

	artists := res.Get(pathPlaylistTrack + ".1." + pathTrackAttribute)
	artists.ForEach(func(key, value gjson.Result) bool {
		if u := extractUser(&value); u != nil {
			tr.Artists = append(tr.Artists, u)
		}
		return true
	})

	album := res.Get(pathPlaylistTrack + ".2." + pathTrackAttribute + ".0")
	if album.Exists() {
		tr.Album = extractAlbum(&album)
	}

	return tr
}
