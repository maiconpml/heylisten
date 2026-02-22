package goytmusic

import "github.com/tidwall/gjson"

type Album struct {
	Name     string
	BrowseID string
}

// Parses res into a Album struct
// Expects the Album contained in browseId=VLPL... JSON response
func extractAlbum(res *gjson.Result) *Album {
	alb := &Album{}
	alb.Name = res.Get(pathTrackName).String()

	buf := res.Get(pathNavEndpointBrowseID)
	if buf.Exists() {
		alb.BrowseID = buf.String()
	}
	return alb
}
