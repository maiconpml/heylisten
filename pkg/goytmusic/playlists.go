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

