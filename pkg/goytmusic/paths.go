package goytmusic

const (
	// json paths to extract data from API responses
	pathRootSingleColumnRenderer = "contents.singleColumnBrowseResultsRenderer"
	pathGridRendererItems        = "gridRenderer.items"
	pathNavEndpointBrowseID      = "navigationEndpoint.browseEndpoint.browseId"
	pathItemTitle                = "title.runs.0"
	pathItemSubtitle             = "subtitle.runs.0"

	pathTab0Contents0 = "tabs.0.tabRenderer.content.sectionListRenderer.contents.0"

	pathMusicResponsiveHeader       = "musicResponsiveHeaderRenderer"
	pathMusicEditablePlaylistHeader = "musicEditablePlaylistDetailHeaderRenderer.header"

	pathRootTwoColumnRenderer = "contents.twoColumnBrowseResultsRenderer"
	pathTracks                = "secondaryContents.sectionListRenderer.contents.0.musicPlaylistShelfRenderer.contents"
	pathPlaylistTrack         = "musicResponsiveListItemRenderer.flexColumns"
	pathTrackAttribute        = "musicResponsiveListItemFlexColumnRenderer.text.runs"
	pathTrackName             = "text"
	pathNavEndpointVideoID    = "navigationEndpoint.watchEndpoint.videoId"
	pathPlaylistAuthor        = "facepile.avatarStackViewModel"
	pathPlaylistAuthorNav     = "rendererContext.commandContext.onTap.innertubeCommand.browseEndpoint.browseId"
	pathTextContent           = "text.content"
)
