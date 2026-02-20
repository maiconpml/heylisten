package goytmusic

import (
	"errors"
)

const (
	brIDLikedPlaylists = "FEmusic_liked_playlists"
)

type PlaylistsService service

type Playlist struct {
	Name       string
	BrID       string
	Author     *string
	AuthorBrID *string
}

// toPlaylist parses a musicTwoRowItem struct to a Playlist
func (mt *musicTowRowItem) toPlaylist() *Playlist {
	pl := &Playlist{}
	if len(mt.MusicTwoRow.Title.Runs) > 0 {
		pl.Name = *mt.MusicTwoRow.Title.Runs[0].Text
	}
	if len(mt.MusicTwoRow.Subtitle.Runs) > 0 {
		runs := mt.MusicTwoRow.Subtitle.Runs
		pl.Author = runs[0].Text

		if runs[0].NavEndpoint != nil && runs[0].NavEndpoint.BrowseEndpoint != nil {
			pl.AuthorBrID = runs[0].NavEndpoint.BrowseEndpoint.BrowseID
		}
	}

	if mt.MusicTwoRow.NavEndpoint != nil && mt.MusicTwoRow.NavEndpoint.BrowseEndpoint != nil {
		pl.BrID = *mt.MusicTwoRow.NavEndpoint.BrowseEndpoint.BrowseID
	}

	return pl
}

// toPlaylistCollection parses a libraryCollectionListResponse to a []*Playlist
func (lc *libraryCollectionListResponse) toPlaylistCollection() []*Playlist {
	items := lc.Contents.SingleColumn.Tabs[0].TabRenderer.Content.SectionList.Contents[0].GridRenderer.Items

	var plCollection []*Playlist
	for _, it := range items[1:] {
		plCollection = append(plCollection, it.toPlaylist())
	}
	return plCollection
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

	plColl := raw.toPlaylistCollection()

	return plColl, nil
}
