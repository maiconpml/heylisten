package tui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/maiconpml/heylisten/internal/audio"
	"github.com/maiconpml/heylisten/internal/tui/components/player"
	"github.com/maiconpml/heylisten/internal/tui/components/playlists"
	"github.com/maiconpml/heylisten/internal/tui/components/tracks"
	"github.com/maiconpml/heylisten/internal/tui/keys"
	"github.com/maiconpml/heylisten/internal/tui/styles"
	"github.com/maiconpml/heylisten/pkg/goytmusic"
	overlay "github.com/rmhubbert/bubbletea-overlay"
)

type viewState int

const (
	viewPlaylists viewState = iota
	viewTracks
)

type TracksLoadedMsg struct {
	Playlist *goytmusic.Playlist
}

type ErrorMsg struct {
	Err error
}

type Model struct {
	client    *goytmusic.Client
	state     viewState
	playlists playlists.Model
	tracks    tracks.Model
	player    player.Model
	help      help.Model
	keys      keys.KeyMap
	width     int
	height    int
}

func NewModel(client *goytmusic.Client, playlistsData []*goytmusic.Playlist) Model {
	return Model{
		client:    client,
		state:     viewPlaylists,
		playlists: playlists.New(playlistsData),
		tracks:    tracks.New(),
		player:    player.New(),
		help:      help.New(),
		keys:      keys.Keys,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.playlists.Init(), m.tracks.Init(), m.player.Init())
}

func (m *Model) updateSizes() {
	if m.width == 0 || m.height == 0 {
		return
	}
	availWidth := m.width - 4
	availHeight := m.height - 2

	m.help.ShowAll = false
	footerHeight := lipgloss.Height(m.help.View(m.keys))

	playerHeight := m.player.Height()

	listHeight := availHeight - playerHeight - footerHeight
	if listHeight < 0 {
		listHeight = 0
	}

	m.playlists.SetSize(availWidth, listHeight)
	m.tracks.SetSize(availWidth, listHeight)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// If help is shown, intercept keys to close it
		if m.help.ShowAll {
			switch {
			case key.Matches(msg, m.keys.Help),
				key.Matches(msg, m.keys.Back),
				key.Matches(msg, m.keys.Quit):
				m.help.ShowAll = false
				return m, nil
			}
			return m, nil
		}
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Back):
			if m.state == viewTracks {
				m.state = viewPlaylists
				return m, nil
			}
		case key.Matches(msg, m.keys.PlayPause):
			return m, func() tea.Msg { return player.TogglePauseMsg{} }
		case key.Matches(msg, m.keys.NextTrack):
			return m, func() tea.Msg { return player.NextTrackMsg{} }
		case key.Matches(msg, m.keys.PrevTrack):
			return m, func() tea.Msg { return player.PrevTrackMsg{} }
		case key.Matches(msg, m.keys.VolUp):
			audio.IncrVolume()
		case key.Matches(msg, m.keys.VolDown):
			audio.DecrVolume()
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.Width = msg.Width

		availWidth := msg.Width - 4
		availHeight := msg.Height - 2

		var cmd tea.Cmd
		m.player, cmd = m.player.Update(tea.WindowSizeMsg{
			Width:  availWidth,
			Height: availHeight,
		})
		cmds = append(cmds, cmd)

		m.updateSizes()
		return m, tea.Batch(cmds...)

	case playlists.PlaylistSelectedMsg:
		m.state = viewTracks

		id := msg.PlaylistID
		return m, func() tea.Msg {
			pl, err := m.client.Playlists.Get(&id)
			if err != nil {
				return ErrorMsg{Err: err}
			}
			return TracksLoadedMsg{Playlist: pl}
		}

	case TracksLoadedMsg:
		m.tracks.SetTracks(msg.Playlist.Tracks, msg.Playlist.Name)
		return m, nil
	case tracks.TrackSelectedMsg:
		return m, func() tea.Msg {
			cont := ""
			tracks, contin, err := m.client.Tracks.NextTracksByMusicInPlaylist(msg.Track.VideoID, msg.Track.PlaylistSetVideoID, msg.Track.PlaylistID, Ptr(cont))
			if err != nil {
				return ErrorMsg{Err: err}
			}
			return player.QueueLoadedMsg{Tracks: tracks, CurTrack: msg.Index, Continuation: contin}
		}
	}

	// Forward messages to active components
	var cmd tea.Cmd

	switch m.state {
	case viewPlaylists:
		m.playlists, cmd = m.playlists.Update(msg)
		cmds = append(cmds, cmd)
	case viewTracks:
		m.tracks, cmd = m.tracks.Update(msg)
		cmds = append(cmds, cmd)
	}

	m.player, cmd = m.player.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	availWidth := m.width - 4

	isHelpOpen := m.help.ShowAll

	var activeView string
	if m.state == viewPlaylists {
		activeView = styles.RenderContainer("Minhas Playlists", availWidth, m.playlists.View())
	} else {
		activeView = styles.RenderContainer(m.tracks.Title(), availWidth, m.tracks.View())
	}

	playerView := styles.RenderContainer("Player", availWidth, m.player.View())

	m.help.ShowAll = false
	footerHelp := m.help.View(m.keys)

	m.help.ShowAll = isHelpOpen

	ui := lipgloss.JoinVertical(lipgloss.Left,
		activeView,
		playerView,
		footerHelp,
	)

	background := lipgloss.NewStyle().Padding(1, 2).Render(ui)

	if isHelpOpen {
		helpContent := m.help.View(m.keys)
		modal := styles.ModalStyle.Render(helpContent)

		return overlay.Composite(
			modal,
			background,
			overlay.Center,
			overlay.Center,
			0, 0,
		)
	}

	return background
}

// Ptr returns a pointer to value v
func Ptr[T any](v T) *T {
	return &v
}
