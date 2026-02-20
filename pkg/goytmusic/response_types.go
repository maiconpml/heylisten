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
			Tabs []struct {
				TabRenderer struct {
					Content struct {
						SectionList struct {
							Contents []struct {
								GridRenderer struct {
									Items []musicTowRowItem `json:"items"`
								} `json:"gridRenderer"`
							} `json:"contents"`
						} `json:"sectionListRenderer"`
					} `json:"content"`
				} `json:"tabRenderer"`
			} `json:"tabs"`
		} `json:"singleColumnBrowseResultsRenderer"`
	} `json:"contents"`
}

type musicTowRowItem struct {
	MusicTwoRow struct {
		Title struct {
			Runs []run `json:"runs"`
		} `json:"title"`
		Subtitle struct {
			Runs []run `json:"runs"`
		} `json:"subtitle"`
		NavEndpoint *navigationEndpoint `json:"navigationEndpoint,omitempty"`
	} `json:"musicTwoRowItemRenderer"`
}

type run struct {
	Text        *string             `json:"text"`
	NavEndpoint *navigationEndpoint `json:"navigationEndpoint,omitempty"`
}

type navigationEndpoint struct {
	BrowseEndpoint *struct {
		BrowseID *string `json:"browseId"`
	} `json:"browseEndpoint,omitempty"`
}
