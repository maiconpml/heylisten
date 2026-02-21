package goytmusic

import (
	"errors"

	"github.com/tidwall/gjson"
)

const (
	brIDLikedPlaylists = "FEmusic_liked_playlists"

	pathRootGridRendererItems = "contents.singleColumnBrowseResultsRenderer.tabs.0.tabRenderer.content.sectionListRenderer.contents.0.gridRenderer.items"
	pathNavEndpointBrowseID   = "navigationEndpoint.browseEndpoint.browseId"
	pathItemTitle             = "title.runs.0"
	pathItemSubtitle          = "subtitle.runs.0"
)

type PlaylistsService service

type Playlist struct {
	Name           string
	BrowseID       string
	Author         *string
	AuthorBrowseID *string
	Tracks         []*Track
}

// ListLiked retrieves and returns an array of Playlist. This array
// corresponds the current user's list of liked playlists.
func (s *PlaylistsService) ListLiked() ([]*Playlist, error) {
	if s.client.isGuest {
		return nil, errors.New("Client is not authenticated")
	}
	u := "browse"
	body := s.client.BrowseBody(brIDLikedPlaylists)
	req, err := s.client.NewRequest("POST", u, body)
	if err != nil {
		return nil, err
	}

	respBody, _, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}

	items := extractPlaylists(respBody)

	if len(items) > 2 {
		return items[2:], nil
	}
	return nil, nil
}


func extractPlaylists(b []byte) []*Playlist {
	results := gjson.GetBytes(b, pathRootGridRendererItems)

	var playlists []*Playlist
	results.ForEach(func(key, value gjson.Result) bool {
		if pl := extractPlaylist(&value); pl != nil {
			playlists = append(playlists, pl)
		}
		return true
	})

	return playlists
}

func extractPlaylist(res *gjson.Result) *Playlist {
	render := res.Get("musicTwoRowItemRenderer")
	if !render.Exists() {
		return nil
	}

	pl := &Playlist{
		Name:     render.Get(pathItemTitle + ".text").String(),
		BrowseID: render.Get(pathNavEndpointBrowseID).String(),
	}

	author := render.Get(pathItemSubtitle)
	if author.Exists() {
		pl.Author = Ptr(author.Get("text").String())
		authorBrowseID := author.Get(pathNavEndpointBrowseID).String()
		if authorBrowseID != "" {
			pl.AuthorBrowseID = Ptr(authorBrowseID)
		}
	}
	return pl
}
