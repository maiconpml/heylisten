package goytmusic

import "github.com/tidwall/gjson"

type User struct {
	Name     string
	BrowseID *string
}

func extractUser(res gjson.Result) *User {
	u := &User{}
	if buf := res.Get(pText); buf.Exists() {
		u.Name = buf.String()
	} else if buf := res.Get(joinPaths(pFacepileStackView, pText, pContent)); buf.Exists() {
		u.Name = buf.String()
	}

	if buf := res.Get(joinPaths(pNavEndpoint, pBrowseEndID)); buf.Exists() {
		u.BrowseID = Ptr(buf.String())
	} else if buf := res.Get(joinPaths(pFacepileStackView, pRendCtxtInnertubeCommand, pBrowseEndID)); buf.Exists() {
		u.BrowseID = Ptr(buf.String())
	}
	return u
}
