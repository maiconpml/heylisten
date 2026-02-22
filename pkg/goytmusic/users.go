package goytmusic

import "github.com/tidwall/gjson"

type User struct {
	Name     string
	BrowseID *string
}

func extractUser(res *gjson.Result) *User {
	u := &User{}
	if buf := res.Get(pathTrackName); buf.Exists() {
		u.Name = buf.String()
	} else if buf := res.Get(pathPlaylistAuthor + "." + pathTextContent); buf.Exists() {
		u.Name = buf.String()
	}

	if buf := res.Get(pathNavEndpointBrowseID); buf.Exists() {
		u.BrowseID = Ptr(buf.String())
	} else if buf := res.Get(pathPlaylistAuthor + "." + pathPlaylistAuthorNav); buf.Exists() {
		u.BrowseID = Ptr(buf.String())
	}
	return u
}
