package stt

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"
)

func TestMockTranscriber_Transcribe(t *testing.T) {
	tests := []struct {
		name         string
		mockResponse string
		mockError    error
		wantErr      bool
		wantText     string
	}{
		{
			name:         "successful transcription",
			mockResponse: "Hello, this is a test transcription",
			mockError:    nil,
			wantErr:      false,
			wantText:     "Hello, this is a test transcription",
		},
		{
			name:         "transcription error",
			mockResponse: "",
			mockError:    fmt.Errorf("API error"),
			wantErr:      true,
			wantText:     "",
		},
		{
			name:         "empty transcription",
			mockResponse: "",
			mockError:    nil,
			wantErr:      false,
			wantText:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transcriber := NewMockTranscriber(tt.mockResponse, tt.mockError)
			ctx := context.Background()

			audioData := bytes.NewReader([]byte("fake audio data"))
			got, err := transcriber.Transcribe(ctx, audioData)

			if (err != nil) != tt.wantErr {
				t.Errorf("Transcribe() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.wantText {
				t.Errorf("Transcribe() = %v, want %v", got, tt.wantText)
			}
		})
	}
}

func TestMockTranscriber_WithContext(t *testing.T) {
	transcriber := NewMockTranscriber("test response", nil)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	audioData := bytes.NewReader([]byte("fake audio data"))
	got, err := transcriber.Transcribe(ctx, audioData)

	if err != nil {
		t.Errorf("Transcribe() unexpected error = %v", err)
	}

	if got != "test response" {
		t.Errorf("Transcribe() = %v, want %v", got, "test response")
	}
}

func TestProcessAudioBuffer(t *testing.T) {
	tests := []struct {
		name    string
		data    []int16
		wantErr bool
	}{
		{
			name:    "process valid audio buffer",
			data:    []int16{100, 200, 300, 400, 500},
			wantErr: false,
		},
		{
			name:    "process empty buffer",
			data:    []int16{},
			wantErr: true,
		},
		{
			name:    "process single sample",
			data:    []int16{42},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf, err := ProcessAudioBuffer(tt.data)

			if (err != nil) != tt.wantErr {
				t.Errorf("ProcessAudioBuffer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && buf.Len() == 0 {
				t.Error("ProcessAudioBuffer() returned empty buffer")
			}
		})
	}
}
