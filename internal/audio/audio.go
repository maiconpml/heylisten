package audio

import (
	"fmt"
	"io"
	"log/slog"
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
	streamerSR  beep.SampleRate
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

func play(s beep.StreamSeekCloser, f beep.Format) error {
	speaker.Clear()
	if streamer != nil {
		streamer.Close()
	}

	resampled := beep.Resample(4, f.SampleRate, sampleRate, s)

	ctrl = &beep.Ctrl{Streamer: resampled, Paused: false}
	volumeCtrl = &effects.Volume{Streamer: ctrl, Base: 2, Volume: volumeCtrl.Volume}
	streamer = s
	streamerSR = f.SampleRate

	speaker.Play(volumeCtrl)
	return nil
}

func PlayStream(r io.Reader) error {
	slog.Info("Iniciando decodificação do stream MP3...")

	var rc io.ReadCloser
	if c, ok := r.(io.ReadCloser); ok {
		rc = c
	} else {
		rc = io.NopCloser(r)
	}

	s, format, err := mp3.Decode(rc)
	if err != nil {
		return fmt.Errorf("erro ao decodificar stream: %v", err)
	}

	slog.Info("Áudio pronto, iniciando reprodução", "sampleRate", format.SampleRate)
	return play(s, format)
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

// GetProgress returns the current position, duration, paused state and if it finished.
func GetProgress() (position float64, duration float64, paused bool, finished bool, err error) {
	mu.Lock()
	defer mu.Unlock()

	if streamer == nil || ctrl == nil {
		return 0, 0, false, false, fmt.Errorf("no track playing")
	}

	speaker.Lock()
	pos := streamer.Position()
	len := streamer.Len()
	isPaused := ctrl.Paused
	speaker.Unlock()

	sr := float64(streamerSR)
	if sr == 0 {
		sr = float64(sampleRate)
	}

	position = float64(pos) / sr
	duration = float64(len) / sr
	finished = len > 0 && pos >= len

	return position, duration, isPaused, finished, nil
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
