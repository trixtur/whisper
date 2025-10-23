package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"speech-to-clipboard/internal/config"
	"speech-to-clipboard/pkg/audio"
	"speech-to-clipboard/pkg/clipboard"
	"speech-to-clipboard/pkg/stt"
	"syscall"
	"time"
)

func main() {
	fmt.Println("Speech-to-Clipboard Application")
	fmt.Println("================================")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize components
	capturer, err := audio.NewCapturer()
	if err != nil {
		log.Fatalf("Failed to initialize audio capturer: %v", err)
	}
	defer func() {
		if err := audio.Cleanup(); err != nil {
			log.Printf("Error cleaning up audio: %v", err)
		}
	}()

	transcriber := stt.NewWhisperTranscriber(cfg.OpenAIAPIKey)
	clipMgr := clipboard.NewManager()

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	fmt.Println("\nInstructions:")
	fmt.Println("- Press ENTER to start recording")
	fmt.Println("- Press ENTER again to stop recording and transcribe")
	fmt.Println("- Press Ctrl+C to exit")
	fmt.Println()

	go func() {
		<-sigChan
		fmt.Println("\n\nShutting down...")
		if capturer.IsRecording() {
			if err := capturer.Stop(); err != nil {
				log.Printf("Error stopping capturer: %v", err)
			}
		}
		os.Exit(0)
	}()

	// Main loop
	for {
		fmt.Print("Press ENTER to start recording: ")
		_, _ = fmt.Scanln()

		fmt.Println("Recording... Press ENTER to stop")
		if err := capturer.Start(); err != nil {
			log.Printf("Error starting recording: %v", err)
			continue
		}

		// Wait for user to press ENTER again
		_, _ = fmt.Scanln()

		fmt.Println("Stopping recording...")
		if err := capturer.Stop(); err != nil {
			log.Printf("Error stopping recording: %v", err)
			continue
		}

		// Get audio data
		audioData, err := capturer.GetAudioData()
		if err != nil {
			log.Printf("Error getting audio data: %v", err)
			continue
		}

		if len(audioData) == 0 {
			fmt.Println("No audio captured. Please try again.")
			continue
		}

		fmt.Printf("Captured %d samples. Transcribing...\n", len(audioData))

		// Convert to WAV format
		wavBuf := new(bytes.Buffer)
		if err := audio.SaveToWAV(audioData, wavBuf); err != nil {
			log.Printf("Error converting to WAV: %v", err)
			continue
		}

		// Transcribe
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		text, err := transcriber.Transcribe(ctx, wavBuf)
		cancel()

		if err != nil {
			log.Printf("Error transcribing: %v", err)
			continue
		}

		if text == "" {
			fmt.Println("No speech detected. Please try again.")
			continue
		}

		fmt.Printf("\nTranscribed text: %s\n\n", text)

		// Copy to clipboard
		if err := clipMgr.Write(text); err != nil {
			log.Printf("Error writing to clipboard: %v", err)
			continue
		}

		fmt.Println("Text copied to clipboard! You can now paste it anywhere.")
		fmt.Println()
	}
}
