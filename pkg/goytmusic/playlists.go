package goytmusic

import (
	"errors"
)

const (
	brIDLikedPlaylists = "FEmusic_liked_playlists"
)

type PlaylistsService service

type Playlist struct {
	Name           string
	BrowseID       string
	Author         *string
	AuthorBrowseID *string
}

// toPlaylist parses a musicTwoRowItem struct to a Playlist
func (mt *musicTowRowItem) toPlaylist() *Playlist {
	pl := &Playlist{}
	if len(mt.MusicTwoRow.Title.Runs) > 0 {
		pl.Name = mt.MusicTwoRow.Title.Runs[0].Text
	}
	if len(mt.MusicTwoRow.Subtitle.Runs) > 0 {
		runs := mt.MusicTwoRow.Subtitle.Runs
		pl.Author = &runs[0].Text

		if runs[0].NavEndpoint != nil && runs[0].NavEndpoint.BrowseEndpoint != nil {
			pl.AuthorBrowseID = runs[0].NavEndpoint.BrowseEndpoint.BrowseID
		}
	}

	if mt.MusicTwoRow.Endpoint != nil && mt.MusicTwoRow.Endpoint.BrowseEndpoint != nil {
		pl.BrowseID = *mt.MusicTwoRow.Endpoint.BrowseEndpoint.BrowseID
	}

	return pl
}

	}
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

	var raw libraryCollectionListResponse
	_, err = s.client.Do(req, &raw)
	if err != nil {
		return nil, err
	}

	items := raw.ExtractPlaylists()

	if len(items) > 2 {
		return items[2:], nil
	}
	return items, nil
}

