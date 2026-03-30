package player

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/maiconpml/yt-music-tui/internal/audio"
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

type QueueLoadedMsg struct {
	Tracks       []*goytmusic.Track
	CurTrack     int
	Continuation string
}

type PrefetchCompleteMsg struct {
	Index int
	Err   error
}

type (
	NextTrackMsg struct{}
	PrevTrackMsg struct{}
)

type playerStatus int

const (
	standby playerStatus = iota
	playing
	paused
	downloading
)

type Model struct {
	status           playerStatus
	tracks           []*goytmusic.Track
	curTrack         int
	continuation     string
	continuationType int
	width            int
	progress         progress.Model
	position         float64
	duration         float64
	err              error
}

func New() Model {
	prog := progress.New(progress.WithFillCharacters('󱘹', '󰍴'), progress.WithoutPercentage())
	return Model{
		status:   standby,
		curTrack: -1,
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

func (m *Model) prefetchTrack(index int) tea.Cmd {
	if index <= 0 || index > len(m.tracks) {
		return nil
	}

	tr := m.tracks[index]
	if tr.VideoID == nil {
		return nil
	}
	return func() tea.Msg {
		slog.Info("prefetching track", "name", tr.Name)
		_, err := ytdlp.GetAudioPath(*tr.VideoID)
		if err != nil {
			slog.Error("error prefetching track", "error", err.Error())
		}
		slog.Info("track prefetched sucessfully")
		return PrefetchCompleteMsg{Index: index, Err: err}
	}
}

func (m *Model) playTrack(index int) tea.Cmd {
	if index < 0 || index > len(m.tracks) {
		return nil
	}
	m.curTrack = index
	m.status = downloading
	m.position = 0
	m.duration = 0
	m.err = nil

	tr := m.tracks[m.curTrack]
	if tr.VideoID == nil {
		m.err = fmt.Errorf("track without videoID")
		m.status = standby
		return nil
	}

	videoID := *tr.VideoID
	return func() tea.Msg {
		slog.Info("downloading track", "name", tr.Name)
		path, err := ytdlp.GetAudioPath(videoID)
		if err != nil {
			slog.Error("error downloading track", "error", err.Error())
		}
		slog.Info("track downloaded sucessfully")
		return DownloadCompleteMsg{Path: path, Err: err}
	}
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
		cmds = append(cmds, tickCmd())
		if m.status == playing || m.status == paused {
			pos, dur, pausd, err := audio.GetProgress()
			if err == nil {
				m.position = pos
				m.duration = dur
				if pausd {
					m.status = paused
				} else {
					m.status = playing
				}

				if dur > 0 {
					cmd := m.progress.SetPercent(pos / dur)
					cmds = append(cmds, cmd)

					if pos >= dur {
						cmds = append(cmds, func() tea.Msg { return NextTrackMsg{} })
						return m, tea.Batch(cmds...)
					}
				}
			}
		}

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		cmds = append(cmds, cmd)

	case TogglePauseMsg:
		if m.status == playing || m.status == paused {
			audio.TogglePause()
			if m.status == playing {
				m.status = paused
			} else {
				m.status = playing
			}
		}

	case NextTrackMsg:
		if m.curTrack+1 < len(m.tracks) {
			cmds = append(cmds, m.playTrack(m.curTrack+1))
		}
	case PrevTrackMsg:
		if m.curTrack-1 >= 0 {
			cmds = append(cmds, m.playTrack(m.curTrack-1))
		}
	case QueueLoadedMsg:
		m.tracks = msg.Tracks
		m.continuation = msg.Continuation

		cmds = append(cmds, m.playTrack(msg.CurTrack))

	case DownloadCompleteMsg:
		if msg.Err == nil {
			if err := audio.Play(msg.Path); err == nil {
				m.status = playing
				cmds = append(cmds, m.prefetchTrack(m.curTrack+1))
			} else {
				m.status = standby
				m.err = err
			}
		} else {
			m.err = msg.Err
			m.status = standby
		}
	case PrefetchCompleteMsg:
		if msg.Err == nil {
			if msg.Index < m.curTrack+2 {
				cmds = append(cmds, m.prefetchTrack(msg.Index+1))
			}
		}

	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.width == 0 {
		return ""
	}

	var line1, line2 string
	if m.err != nil {
		line1 = fmt.Sprintf("Error: %v", m.err)
	} else if m.status == standby || m.curTrack == -1 {
		line1 = "▶ Not Playing - No track selected"
	} else if m.status == downloading {
		line1 = fmt.Sprintf("⬇ Downloading: %s...", m.tracks[m.curTrack].Name)
	} else {
		tr := m.tracks[m.curTrack]
		artist := ""
		if len(tr.Artists) > 0 {
			artist = tr.Artists[0].Name
		}
		status := "⏸"
		if m.status == playing {
			status = "▶"
		}
		line1 = fmt.Sprintf("%s %s - %s", status, tr.Name, artist)
		line2 = fmt.Sprintf("%s %s %s", formatDuration(m.position), m.progress.View(), formatDuration(m.duration))
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		line1,
		line2,
	)
}

func (m Model) Height() int {
	return lipgloss.Height(m.View()) + styles.ContainerFrameHeight()
}
