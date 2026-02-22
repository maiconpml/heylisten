package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github/maiconpml/yt-music-tui/pkg/ytmusic"
)

func main() {
	cookie := os.Getenv("AUTH_COOKIE")
	if cookie == "" {
		log.Println("WARNING: AUTH_COOKIE environment variable is not set. This will likely fail.")
	}

	client := ytmusic.NewClient(&http.Client{}).WithAuthCookie(cookie)

	liked, err := client.Playlists.ListLiked()
	if err != nil {
		log.Fatalf("failed to list liked playlists: %v", err)
	}

	// for _, pl := range liked {
	// 	fmt.Printf("Name: %s; ID: %s; ", pl.Name, pl.BrowseID)
	// 	if pl.Author == nil {
	// 		fmt.Printf("Author: <nil>;\n")
	// 	} else {
	// 		fmt.Printf("Author: %s;\n", *pl.Author)
	// 	}
	// }

	for _, pl := range liked {
		pl, _ = client.Playlists.Get(&pl.BrowseID)

		fmt.Printf("Playlist %s by %s (BrowseID: %s)\n", pl.Name, pl.Author.Name, pl.BrowseID)
		fmt.Printf("[ ")
		for i := range 5 {
			if len(pl.Tracks) < i {
				break
			}
			fmt.Printf("[ %s, %s", pl.Tracks[i].Name, pl.Tracks[i].Artists[0].Name)
			if pl.Tracks[i].Album != nil {
				fmt.Printf(", %s ]", pl.Tracks[i].Album.Name)
			} else {
				fmt.Printf(" ]")
			}
		}
		fmt.Printf("]\n\n")
	}
}
