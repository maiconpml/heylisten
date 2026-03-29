package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/maiconpml/yt-music-tui/internal/audio"
	"github.com/maiconpml/yt-music-tui/internal/config"
	"github.com/maiconpml/yt-music-tui/internal/tui"
	"github.com/maiconpml/yt-music-tui/internal/ytdlp"
	"github.com/maiconpml/yt-music-tui/pkg/goytmusic"
)

func main() {
	if err := ytdlp.CheckDependencies(); err != nil {
		log.Fatalf("Dependency error: %v", err)
	}

	if err := audio.Init(); err != nil {
		log.Printf("Warning: Could not initialize audio system: %v\nAudio playback will not work.", err)
	} else {
		defer audio.Quit()
	}

	cookiePath, err := config.GetCookiePath()
	if err != nil {
		log.Fatalf("Error getting config path: %v", err)
	}

	if _, err := os.Stat(cookiePath); os.IsNotExist(err) {
		fmt.Printf("No authentication cookie found. Please paste your YouTube Music cookie string:\n> ")
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Error reading input: %v", err)
		}
		input = strings.TrimSpace(input)
		if input == "" {
			log.Fatalf("Cookie cannot be empty")
		}
		if err := config.SaveCookie(input); err != nil {
			log.Fatalf("Error saving cookie: %v", err)
		}
		fmt.Println("Cookie saved successfully!")
	}

	cookieString, err := config.LoadCookie()
	if err != nil {
		log.Fatalf("Error loading cookie: %v", err)
	}

	if err := ytdlp.Init(cookiePath); err != nil {
		log.Fatalf("error on ytdlp initializing: %v", err)
	}

	client := goytmusic.NewClient(&http.Client{}).WithAuthCookie(cookieString)

	liked, err := client.Playlists.ListLiked()
	if err != nil {
		log.Fatalf("failed to list liked playlists: %v", err)
	}

	m := tui.NewModel(client, liked)
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
