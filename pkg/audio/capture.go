package audio

import (
	"fmt"
	"io"

	"github.com/gordonklaus/portaudio"
)

const (
	SampleRate      = 16000
	FramesPerBuffer = 1024
	Channels        = 1
)

// Capturer handles microphone audio capture
type Capturer interface {
	Start() error
	Stop() error
	GetAudioData() ([]int16, error)
	IsRecording() bool
}

type portAudioCapturer struct {
	stream    *portaudio.Stream
	buffer    []int16
	recording bool
}

// NewCapturer creates a new audio capturer
func NewCapturer() (Capturer, error) {
	if err := portaudio.Initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize portaudio: %w", err)
	}

	return &portAudioCapturer{
		buffer:    make([]int16, 0),
		recording: false,
	}, nil
}

// Start begins capturing audio from the microphone
func (c *portAudioCapturer) Start() error {
	if c.recording {
		return fmt.Errorf("already recording")
	}

	c.buffer = make([]int16, 0)

	stream, err := portaudio.OpenDefaultStream(Channels, 0, float64(SampleRate), FramesPerBuffer, func(in []int16) {
		c.buffer = append(c.buffer, in...)
	})
	if err != nil {
		return fmt.Errorf("failed to open stream: %w", err)
	}

	c.stream = stream
	if err := c.stream.Start(); err != nil {
		return fmt.Errorf("failed to start stream: %w", err)
	}

	c.recording = true
	return nil
}

// Stop stops capturing audio
func (c *portAudioCapturer) Stop() error {
	if !c.recording {
		return fmt.Errorf("not recording")
	}

	if err := c.stream.Stop(); err != nil {
		return fmt.Errorf("failed to stop stream: %w", err)
	}

	if err := c.stream.Close(); err != nil {
		return fmt.Errorf("failed to close stream: %w", err)
	}

	c.recording = false
	return nil
}

// GetAudioData returns the captured audio data
func (c *portAudioCapturer) GetAudioData() ([]int16, error) {
	if c.recording {
		return nil, fmt.Errorf("still recording, call Stop() first")
	}

	return c.buffer, nil
}

// IsRecording returns whether the capturer is currently recording
func (c *portAudioCapturer) IsRecording() bool {
	return c.recording
}

// Cleanup terminates the portaudio library
func Cleanup() {
	portaudio.Terminate()
}

// SaveToWAV saves the audio buffer to a WAV file
func SaveToWAV(data []int16, writer io.Writer) error {
	// WAV header
	dataSize := len(data) * 2 // 2 bytes per int16
	fileSize := 36 + dataSize

	header := make([]byte, 44)

	// RIFF chunk
	copy(header[0:4], "RIFF")
	writeInt32(header[4:8], uint32(fileSize))
	copy(header[8:12], "WAVE")

	// fmt chunk
	copy(header[12:16], "fmt ")
	writeInt32(header[16:20], 16)        // fmt chunk size
	writeInt16(header[20:22], 1)         // PCM format
	writeInt16(header[22:24], Channels)  // channels
	writeInt32(header[24:28], SampleRate) // sample rate
	writeInt32(header[28:32], SampleRate*Channels*2) // byte rate
	writeInt16(header[32:34], Channels*2) // block align
	writeInt16(header[34:36], 16)        // bits per sample

	// data chunk
	copy(header[36:40], "data")
	writeInt32(header[40:44], uint32(dataSize))

	if _, err := writer.Write(header); err != nil {
		return err
	}

	// Write audio data
	buf := make([]byte, 2)
	for _, sample := range data {
		writeInt16(buf, uint16(sample))
		if _, err := writer.Write(buf); err != nil {
			return err
		}
	}

	return nil
}

func writeInt16(b []byte, v uint16) {
	b[0] = byte(v)
	b[1] = byte(v >> 8)
}

func writeInt32(b []byte, v uint32) {
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
}
