package stt

import (
	"bytes"
	"context"
	"fmt"
	"io"
)

// Transcriber converts audio data to text
type Transcriber interface {
	Transcribe(ctx context.Context, audioData io.Reader) (string, error)
}

// TranscriberConfig holds configuration for transcription services
type TranscriberConfig struct {
	APIKey   string
	Model    string
	Language string
}

// MockTranscriber is a mock implementation for testing
type MockTranscriber struct {
	Response string
	Error    error
}

// NewMockTranscriber creates a mock transcriber
func NewMockTranscriber(response string, err error) Transcriber {
	return &MockTranscriber{
		Response: response,
		Error:    err,
	}
}

// Transcribe returns the mock response
func (m *MockTranscriber) Transcribe(ctx context.Context, audioData io.Reader) (string, error) {
	if m.Error != nil {
		return "", m.Error
	}
	return m.Response, nil
}

// ProcessAudioBuffer is a helper function that converts int16 audio to a format ready for transcription
func ProcessAudioBuffer(data []int16) (*bytes.Buffer, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty audio data")
	}

	buf := new(bytes.Buffer)

	// Write WAV header (simplified for demonstration)
	// In production, you'd use a proper WAV encoder
	for _, sample := range data {
		buf.WriteByte(byte(sample))
		buf.WriteByte(byte(sample >> 8))
	}

	return buf, nil
}
