package goytmusic

import (
	"github.com/tidwall/gjson"
)

type TracksService service

type Track struct {
	Name     string
	Artists  []*User
	VideoID  *string
	Duration string
	Album    *Album
}

// Parses res into a Track struct
// Expects the track contained in browseId=VLPL... endpoint JSON response
func extractTrack(res gjson.Result) *Track {
	tr := &Track{}

	name := res.Get(joinPaths(pRespListItem, pFlexColumn0, pRespListItemFlexColumn, pText, pRun))
	tr.Name = name.Get(joinPaths(pText)).String()
	buf := name.Get(joinPaths(pNavEndpoint, pWatchEndID))
	if buf.Exists() {
		tr.VideoID = Ptr(buf.String())
	}

	artists := res.Get(joinPaths(pRespListItem, pFlexColumn1, pRespListItemFlexColumn, pText, pRuns))
	artists.ForEach(func(key, value gjson.Result) bool {
		if u := extractUser(value); u != nil {
			tr.Artists = append(tr.Artists, u)
		}
		return true
	})

	album := res.Get(joinPaths(pRespListItem, pFlexColumn2, pRespListItemFlexColumn, pText, pRun))
	if album.Exists() {
		tr.Album = extractAlbum(&album)
	}

	return tr
}
