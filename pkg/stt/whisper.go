package stt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

// WhisperTranscriber uses OpenAI's Whisper API
type WhisperTranscriber struct {
	apiKey string
	model  string
	client *http.Client
}

// NewWhisperTranscriber creates a new Whisper API transcriber
func NewWhisperTranscriber(apiKey string) Transcriber {
	return &WhisperTranscriber{
		apiKey: apiKey,
		model:  "whisper-1",
		client: &http.Client{},
	}
}

type whisperResponse struct {
	Text string `json:"text"`
}

// Transcribe sends audio to Whisper API and returns transcribed text
func (w *WhisperTranscriber) Transcribe(ctx context.Context, audioData io.Reader) (string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add file field
	part, err := writer.CreateFormFile("file", "audio.wav")
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %w", err)
	}

	if _, err := io.Copy(part, audioData); err != nil {
		return "", fmt.Errorf("failed to copy audio data: %w", err)
	}

	// Add model field
	if err := writer.WriteField("model", w.model); err != nil {
		return "", fmt.Errorf("failed to write model field: %w", err)
	}

	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("failed to close writer: %w", err)
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/audio/transcriptions", body)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+w.apiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send request
	resp, err := w.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Parse response
	var result whisperResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Text, nil
}
