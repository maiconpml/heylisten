package goytmusic

// The InnerTube API returns complex JSON objects that coordinate the rendering
// steps of the web page, rather than just raw data. To extract relevant
// information from these large responses, this package uses "shadow structs"
// that mirror the original JSON structure while including only the necessary
// fields. These nested structures are used to filter and clean the data
// received from the API responses.

type libraryCollectionListResponse struct {
	Contents struct {
		SingleColumn struct {
			Tabs []tab `json:"tabs"`
		} `json:"singleColumnBrowseResultsRenderer"`
	} `json:"contents"`
}

type tab struct {
	TabRenderer struct {
		Content struct {
			SectionList struct {
				Contents []sectionContent `json:"contents"`
			} `json:"sectionListRenderer"`
		} `json:"content"`
	} `json:"tabRenderer"`
}

type sectionContent struct {
	GridRenderer               *gridRenderer               `json:"gridRenderer,omitempty"`
}

type gridRenderer struct {
	Items []musicTowRowItem `json:"items"`
}
type musicTowRowItem struct {
	MusicTwoRow struct {
		Title    text                `json:"title"`
		Subtitle text                `json:"subtitle"`
		Endpoint *navigationEndpoint `json:"navigationEndpoint,omitempty"`
	} `json:"musicTwoRowItemRenderer"`
}

type text struct {
	Runs []run `json:"runs"`
}
type run struct {
	Text        string              `json:"text"`
	NavEndpoint *navigationEndpoint `json:"navigationEndpoint,omitempty"`
}

type navigationEndpoint struct {
	BrowseEndpoint *struct {
		BrowseID *string `json:"browseId"`
	} `json:"browseEndpoint,omitempty"`

// ExtractPlaylists extracts playlists from a library response
func (r *libraryCollectionListResponse) ExtractPlaylists() []*Playlist {
	var playlists []*Playlist
	for _, tab := range r.Contents.SingleColumn.Tabs {
		for _, section := range tab.TabRenderer.Content.SectionList.Contents {
			if section.GridRenderer != nil {
				for _, item := range section.GridRenderer.Items {
					playlists = append(playlists, item.toPlaylist())
				}
			}
		}
	}
	return playlists
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

