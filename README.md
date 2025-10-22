# Speech-to-Clipboard

A Go application that listens to your microphone, converts speech to text using OpenAI's Whisper API, and automatically copies the transcribed text to your clipboard for easy pasting.

## Features

- Real-time microphone audio capture
- Speech-to-text transcription using OpenAI Whisper API
- Automatic clipboard integration
- Cross-platform support (macOS, Linux, Windows)
- Comprehensive unit tests for business logic
- Clean, modular architecture

## Architecture

The project follows a clean architecture pattern with clear separation of concerns:

```
speech-to-clipboard/
├── cmd/
│   └── speech-to-clipboard/    # Main application entry point
├── pkg/
│   ├── audio/                  # Microphone capture and WAV encoding
│   ├── stt/                    # Speech-to-text transcription
│   └── clipboard/              # Clipboard operations
├── internal/
│   └── config/                 # Configuration management
└── go.mod
```

## Prerequisites

### System Dependencies

1. **PortAudio** - Required for microphone access
   ```bash
   # macOS
   brew install portaudio

   # Ubuntu/Debian
   sudo apt-get install portaudio19-dev

   # Fedora
   sudo dnf install portaudio-devel

   # Windows
   # Download from http://www.portaudio.com/download.html
   ```

2. **Clipboard support** - Usually built-in on most systems
   - macOS: Uses `pbcopy`/`pbpaste`
   - Linux: Requires `xclip` or `xsel`
   - Windows: Native support

### Go Dependencies

The project uses the following Go libraries:
- `github.com/gordonklaus/portaudio` - Audio capture
- `github.com/atotto/clipboard` - Clipboard operations
- Standard library for HTTP, JSON, and other utilities

## Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd speech-to-clipboard
   ```

2. Install Go dependencies:
   ```bash
   go mod download
   ```

3. Set up your OpenAI API key:
   ```bash
   export OPENAI_API_KEY="your-api-key-here"
   ```

## Configuration

The application is configured via environment variables:

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `OPENAI_API_KEY` | OpenAI API key for Whisper | - | Yes |
| `STT_MODEL` | Whisper model to use | `whisper-1` | No |
| `STT_LANGUAGE` | Language code for transcription | `en` | No |

### Example Configuration

```bash
export OPENAI_API_KEY="sk-..."
export STT_MODEL="whisper-1"
export STT_LANGUAGE="en"
```

## Usage

1. Build the application:
   ```bash
   go build -o speech-to-clipboard ./cmd/speech-to-clipboard
   ```

2. Run the application:
   ```bash
   ./speech-to-clipboard
   ```

3. Follow the on-screen instructions:
   - Press ENTER to start recording
   - Speak into your microphone
   - Press ENTER again to stop recording
   - Wait for transcription (usually 2-5 seconds)
   - The transcribed text will be automatically copied to your clipboard
   - Press Ctrl+C to exit

### Example Session

```
Speech-to-Clipboard Application
================================

Instructions:
- Press ENTER to start recording
- Press ENTER again to stop recording and transcribe
- Press Ctrl+C to exit

Press ENTER to start recording:
Recording... Press ENTER to stop

Stopping recording...
Captured 48000 samples. Transcribing...

Transcribed text: Hello, this is a test of the speech to text system.

Text copied to clipboard! You can now paste it anywhere.

Press ENTER to start recording:
```

## Running Tests

Run all unit tests:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```

Run tests with verbose output:
```bash
go test -v ./...
```

Generate coverage report:
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Project Structure

### `pkg/audio`
Handles microphone audio capture using PortAudio. Features:
- 16kHz mono audio capture
- WAV file format encoding
- Clean start/stop interface

Key files:
- `capture.go` - Audio capture implementation
- `capture_test.go` - Unit tests for audio utilities

### `pkg/stt`
Speech-to-text transcription. Features:
- OpenAI Whisper API integration
- Mock implementation for testing
- Pluggable transcriber interface

Key files:
- `transcriber.go` - Transcriber interface and mock
- `whisper.go` - Whisper API client
- `transcriber_test.go` - Unit tests

### `pkg/clipboard`
Clipboard operations. Features:
- Cross-platform clipboard access
- Mock implementation for testing
- Simple read/write interface

Key files:
- `clipboard.go` - Clipboard manager
- `clipboard_test.go` - Unit tests

### `internal/config`
Configuration management. Features:
- Environment variable loading
- Default values
- Validation

Key files:
- `config.go` - Configuration loader
- `config_test.go` - Unit tests

## Development

### Adding a New Speech-to-Text Provider

To add a new STT provider:

1. Implement the `Transcriber` interface in `pkg/stt`:
   ```go
   type Transcriber interface {
       Transcribe(ctx context.Context, audioData io.Reader) (string, error)
   }
   ```

2. Create a new file (e.g., `pkg/stt/google.go`)

3. Update the main application to use your new transcriber

### Testing

The project includes comprehensive unit tests for all business logic:
- Clipboard operations (mock-based)
- Configuration loading
- Audio format encoding
- Speech-to-text processing (mock-based)

Note: Audio capture tests are limited due to hardware dependencies. In production, you may want to add integration tests with mock audio devices.

## Troubleshooting

### "Failed to initialize portaudio"
- Ensure PortAudio is installed on your system
- On Linux, ensure you have ALSA/PulseAudio configured
- Try running with `sudo` (may be needed for microphone access)

### "OPENAI_API_KEY environment variable is required"
- Set your OpenAI API key: `export OPENAI_API_KEY="sk-..."`
- Ensure the API key is valid and has access to the Whisper API

### "Failed to write to clipboard"
- On Linux, install `xclip`: `sudo apt-get install xclip`
- On macOS, ensure accessibility permissions are granted
- On Windows, run as administrator if needed

### "No audio captured"
- Check your microphone is connected and working
- Ensure your system's default input device is correct
- Try speaking louder or closer to the microphone
- Check system audio settings/permissions

## License

[Specify your license here]

## Contributing

Contributions are welcome! Please:
1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

## Acknowledgments

- [PortAudio](http://www.portaudio.com/) - Cross-platform audio I/O
- [OpenAI Whisper](https://openai.com/research/whisper) - Speech recognition
- [atotto/clipboard](https://github.com/atotto/clipboard) - Clipboard operations
