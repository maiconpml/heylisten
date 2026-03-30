package audio

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/effects"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
)

var (
	sampleRate  = beep.SampleRate(44100)
	mu          sync.Mutex
	ctrl        *beep.Ctrl
	streamer    beep.StreamSeekCloser
	initialized bool
	volumeCtrl  *effects.Volume
)

// Init initializes the audio speaker.
func Init() error {
	mu.Lock()
	defer mu.Unlock()
	if initialized {
		return nil
	}
	err := speaker.Init(sampleRate, sampleRate.N(time.Second/20))
	if err == nil {
		initialized = true
	}
	volumeCtrl = &effects.Volume{Volume: -10.5}
	return err
}

// Play plays the audio file at the given path.
func Play(path string) error {
	mu.Lock()
	defer mu.Unlock()

	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}

	s, format, err := mp3.Decode(f)
	if err != nil {
		f.Close()
		return fmt.Errorf("failed to decode mp3: %v", err)
	}

	speaker.Clear()
	if streamer != nil {
		streamer.Close()
	}

	resampled := beep.Resample(4, format.SampleRate, sampleRate, s)

	ctrl = &beep.Ctrl{Streamer: resampled, Paused: false}
	volumeCtrl = &effects.Volume{Streamer: ctrl, Base: 2, Volume: volumeCtrl.Volume}
	streamer = s

	speaker.Play(volumeCtrl)
	return nil
}

// TogglePause toggles the playback between paused and playing states.
func TogglePause() {
	mu.Lock()
	defer mu.Unlock()
	if ctrl != nil {
		speaker.Lock()
		ctrl.Paused = !ctrl.Paused
		speaker.Unlock()
	}
}

// IncrVolume increases the volume by 5%
func IncrVolume() {
	mu.Lock()
	defer mu.Unlock()
	if volumeCtrl != nil {
		speaker.Lock()
		volumeCtrl.Silent = false
		volumeCtrl.Volume += 0.75
		if volumeCtrl.Volume > 0 {
			volumeCtrl.Volume = 0
		}
		speaker.Unlock()
	}
}

// VolumePercent returns the current volume percentage
func VolumePercent() int {
	return int(((volumeCtrl.Volume + 15) / 15) * 100)
}

// DecrVolume decreases the volume by 5%
func DecrVolume() {
	mu.Lock()
	defer mu.Unlock()
	if volumeCtrl != nil {
		speaker.Lock()
		volumeCtrl.Volume -= 0.75
		if volumeCtrl.Volume <= -15 {
			volumeCtrl.Silent = true
			volumeCtrl.Volume = -15
		}
		speaker.Unlock()
	}
}

// GetProgress returns the current position, duration and paused state of the playback.
func GetProgress() (position float64, duration float64, paused bool, err error) {
	mu.Lock()
	defer mu.Unlock()

	if streamer == nil || ctrl == nil {
		return 0, 0, false, fmt.Errorf("no track playing")
	}

	speaker.Lock()
	pos := streamer.Position()
	len := streamer.Len()
	isPaused := ctrl.Paused
	speaker.Unlock()

	position = float64(pos) / float64(sampleRate)
	duration = float64(len) / float64(sampleRate)

	return position, duration, isPaused, nil
}

// Stop stops the playback and releases resources.
func Stop() {
	mu.Lock()
	defer mu.Unlock()
	speaker.Clear()
	if streamer != nil {
		streamer.Close()
		streamer = nil
	}
	ctrl = nil
}

// Quit cleans up the audio system.
func Quit() {
	Stop()
}
