package ytdlp

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	activeCookiePath string
	activeCacheDir   string
)

func Init(cookiePath string) error {
	activeCookiePath = cookiePath

	cacheDir, err := os.UserCacheDir()
	if err != nil {
		cacheDir = os.TempDir()
	}

	activeCacheDir := filepath.Join(cacheDir, "yt-music-tui")
	if err := os.MkdirAll(activeCacheDir, 0o755); err != nil {
		return fmt.Errorf("failed to create cache directory: %v", err)
	}
	return nil
}

// GetAudioPath downloads the audio for the given videoID to a cache directory
// and returns the path to the downloaded MP3 file.
func GetAudioPath(videoID string) (string, error) {
	path := filepath.Join(activeCacheDir, videoID+".mp3")
	if _, err := os.Stat(path); err == nil {
		return path, nil
	}

	url := fmt.Sprintf("https://music.youtube.com/watch?v=%s", videoID)
	// -x: extract audio, --audio-format mp3: convert to mp3, --audio-quality 0: best quality
	args := []string{"-x", "--audio-format", "mp3", "--audio-quality", "0", "--output", path}
	if activeCookiePath != "" {
		args = append(args, "--cookies", activeCookiePath)
	}
	args = append(args, url)

	cmd := exec.Command("yt-dlp", args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("yt-dlp failed: %v, output: %s", err, string(output))
	}

	return path, nil
}

// CheckDependencies checks if yt-dlp and ffmpeg are installed.
func CheckDependencies() error {
	if _, err := exec.LookPath("yt-dlp"); err != nil {
		return fmt.Errorf("yt-dlp is not installed. Please install it to download audio")
	}
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return fmt.Errorf("ffmpeg is not installed. It is required by yt-dlp for audio extraction")
	}
	return nil
}
