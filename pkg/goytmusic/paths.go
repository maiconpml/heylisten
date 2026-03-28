package goytmusic

import "strings"

const (
	// json paths fragments to extract data from API responses
	pSingleColumn                = "contents.singleColumnBrowseResultsRenderer"
	pTwoColumn                   = "contents.twoColumnBrowseResultsRenderer"
	pGridRendererItems           = "gridRenderer.items"
	pTab0                        = "tabs.0"
	pContent0                    = "contents.0"
	pContents                    = "contents"
	pTabRendererContent          = "tabRenderer.content"
	pSectionList                 = "sectionListRenderer"
	pMusicResponsiveHeader       = "musicResponsiveHeaderRenderer"
	pMusicEditablePlaylistHeader = "musicEditablePlaylistDetailHeaderRenderer.header"
	pSecContents                 = "secondaryContents"
	pPlaylistShelf               = "musicPlaylistShelfRenderer"
	pRespListItem                = "musicResponsiveListItemRenderer"
	pRespListItemFlexColumn      = "musicResponsiveListItemFlexColumnRenderer"
	pFlexColumn0                 = "flexColumns.0"
	pFlexColumn1                 = "flexColumns.1"
	pFlexColumn2                 = "flexColumns.2"
	pMusicTwoRow                 = "musicTwoRowItemRenderer"
	pRun                         = "runs.0"
	pRuns                        = "runs"
	pText                        = "text"
	pNavEndpoint                 = "navigationEndpoint"
	pBrowseEnd                   = "browseEndpoint"
	pBrowseID                    = "browseId"
	pWatchEndID                  = "watchEndpoint.videoId"
	pTitle                       = "title"
	pSubtitle                    = "subtitle"
	pFacepileStackView           = "facepile.avatarStackViewModel"
	pRendCtxtInnertubeCommand    = "rendererContext.commandContext.onTap.innertubeCommand"
	pContent                     = "content"
)

// joinPaths takes multiple strings s1, s2, ..., sn and join
// them in s1.s2....sn
func joinPaths(parts ...string) string {
	return strings.Join(parts, ".")
}
