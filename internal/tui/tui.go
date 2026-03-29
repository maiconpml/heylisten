package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/maiconpml/yt-music-tui/internal/tui/components/player"
	"github.com/maiconpml/yt-music-tui/internal/tui/components/playlists"
	"github.com/maiconpml/yt-music-tui/internal/tui/components/tracks"
	"github.com/maiconpml/yt-music-tui/internal/tui/styles"
	"github.com/maiconpml/yt-music-tui/pkg/goytmusic"
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
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.playlists.Init(), m.tracks.Init(), m.player.Init())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "esc", "backspace":
			if m.state == viewTracks {
				m.state = viewPlaylists
				return m, nil
			}
		case " ":
			return m, func() tea.Msg { return player.TogglePauseMsg{} }
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		availWidth := msg.Width - 4
		availHeight := msg.Height - 2

		var cmd tea.Cmd
		m.player, cmd = m.player.Update(tea.WindowSizeMsg{
			Width:  availWidth,
			Height: availHeight,
		})
		cmds = append(cmds, cmd)

		playerHeight := m.player.Height()
		m.playlists.SetSize(availWidth, availHeight-playerHeight)
		m.tracks.SetSize(availWidth, availHeight-playerHeight)
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

	var activeView string
	if m.state == viewPlaylists {
		activeView = styles.RenderContainer("Minhas Playlists", availWidth, m.playlists.View())
	} else {
		activeView = styles.RenderContainer(m.tracks.Title(), availWidth, m.tracks.View())
	}

	playerView := styles.RenderContainer("Player", availWidth, m.player.View())

	ui := lipgloss.JoinVertical(lipgloss.Left,
		activeView,
		playerView,
	)

	return lipgloss.NewStyle().Padding(1, 2).Render(ui)
}
