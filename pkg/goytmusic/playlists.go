package goytmusic

import (
	"errors"

	"github.com/tidwall/gjson"
)

const (
	brIDLikedPlaylists = "FEmusic_liked_playlists"

	// json paths to extract data from API responses
	pathRootSingleColumnRenderer = "contents.singleColumnBrowseResultsRenderer"
	pathGridRendererItems        = "gridRenderer.items"
	pathNavEndpointBrowseID      = "navigationEndpoint.browseEndpoint.browseId"
	pathItemTitle                = "title.runs.0"
	pathItemSubtitle             = "subtitle.runs.0"

	pathTab0Contents0 = "tabs.0.tabRenderer.content.sectionListRenderer.contents.0"
)

type PlaylistsService service

type Playlist struct {
	Name           string
	BrowseID       string
	Tracks         []*Track
	Author   *User
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
	results := gjson.GetBytes(b, pathRootSingleColumnRenderer+"."+pathTab0Contents0+"."+pathGridRendererItems)

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
		pl.Author = extractUser(&author)
	}
	return pl
}
