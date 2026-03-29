package player

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/maiconpml/yt-music-tui/internal/audio"
	"github.com/maiconpml/yt-music-tui/internal/tui/components/tracks"
	"github.com/maiconpml/yt-music-tui/internal/tui/styles"
	"github.com/maiconpml/yt-music-tui/internal/ytdlp"
	"github.com/maiconpml/yt-music-tui/pkg/goytmusic"
)

type TickMsg time.Time

type TogglePauseMsg struct{}

type DownloadCompleteMsg struct {
	Path string
	Err  error
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

type Model struct {
	playing     bool
	downloading bool
	track       *goytmusic.Track
	width       int
	progress    progress.Model
	position    float64
	duration    float64
	err         error
}

func New() Model {
	prog := progress.New(progress.WithDefaultGradient())
	return Model{
		playing:  false,
		progress: prog,
	}
}

func (m Model) Init() tea.Cmd {
	return tickCmd()
}

func formatDuration(seconds float64) string {
	mins := int(seconds) / 60
	secs := int(seconds) % 60
	return fmt.Sprintf("%02d:%02d", mins, secs)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		fw := styles.ContainerFrameWidth()
		m.progress.Width = msg.Width - fw - 30
		if m.progress.Width < 10 {
			m.progress.Width = 10
		}

	case TickMsg:
		if m.track != nil && !m.downloading {
			pos, dur, paused, err := audio.GetProgress()
			if err == nil {
				m.position = pos
				m.duration = dur
				m.playing = !paused

				if dur > 0 {
					cmd := m.progress.SetPercent(pos / dur)
					cmds = append(cmds, cmd)
				}
			}
		}
		cmds = append(cmds, tickCmd())

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		cmds = append(cmds, cmd)

	case TogglePauseMsg:
		if m.track != nil && !m.downloading {
			audio.TogglePause()
			m.playing = !m.playing
		}

	case tracks.TrackSelectedMsg:
		m.track = msg.Track
		m.playing = false
		m.downloading = true
		m.position = 0
		m.duration = 0
		m.err = nil
		if msg.Track.VideoID != nil {
			videoID := *msg.Track.VideoID
			cmds = append(cmds, func() tea.Msg {
				path, err := ytdlp.GetAudioPath(videoID)
				return DownloadCompleteMsg{Path: path, Err: err}
			})
		}

	case DownloadCompleteMsg:
		m.downloading = false
		if msg.Err == nil {
			if err := audio.Play(msg.Path); err == nil {
				m.playing = true
			} else {
				m.err = err
			}
		} else {
			m.err = msg.Err
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.width == 0 {
		return ""
	}

	var content string
	if m.err != nil {
		content = fmt.Sprintf("Error: %v", m.err)
	} else if m.track == nil {
		content = "▶ Not Playing - No track selected"
	} else if m.downloading {
		content = fmt.Sprintf("⬇ Downloading: %s...", m.track.Name)
	} else {
		artist := ""
		if len(m.track.Artists) > 0 {
			artist = m.track.Artists[0].Name
		}
		status := "⏸"
		if m.playing {
			status = "▶"
		}
		info := fmt.Sprintf("%s %s - %s", status, m.track.Name, artist)
		timeInfo := fmt.Sprintf("%s / %s", formatDuration(m.position), formatDuration(m.duration))

		content = lipgloss.JoinVertical(lipgloss.Left,
			info,
			fmt.Sprintf("%s %s", m.progress.View(), timeInfo),
		)
	}

	return content
}

func (m Model) Height() int {
	return lipgloss.Height(m.View()) + styles.ContainerFrameHeight()
}
