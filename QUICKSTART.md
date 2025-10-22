# Quick Start Guide

## Setup Complete!

Your speech-to-clipboard application is ready to use!

### What's Been Set Up

✓ PortAudio installed (brew)
✓ pkg-config installed (required for building)
✓ Go dependencies downloaded
✓ OpenAI API key configured
✓ Application built successfully (8.3MB binary)
✓ All tests passing (4/4 test suites)

### Running the Application

The application has been built and is ready to run:

```bash
./speech-to-clipboard
```

**Note:** The OPENAI_API_KEY environment variable is set in your current shell session. If you open a new terminal, you'll need to set it again:

```bash
export OPENAI_API_KEY="<KEY>"
```

Or add it to your shell profile (~/.zshrc or ~/.bash_profile):

```bash
echo 'export OPENAI_API_KEY="<KEY>"' >> ~/.zshrc
```

### How to Use

1. Run the application:
   ```bash
   ./speech-to-clipboard
   ```

2. Follow the on-screen instructions:
   - Press ENTER to start recording
   - Speak clearly into your microphone
   - Press ENTER again to stop and transcribe
   - The text will be automatically copied to your clipboard
   - Paste it anywhere with Cmd+V (macOS) or Ctrl+V (Linux/Windows)

3. Press Ctrl+C to exit

### Example Usage

```
Speech-to-Clipboard Application
================================

Instructions:
- Press ENTER to start recording
- Press ENTER again to stop recording and transcribe
- Press Ctrl+C to exit

Press ENTER to start recording: [Press ENTER]
Recording... Press ENTER to stop
[Speak: "Hello, this is a test of the speech to text system"]
[Press ENTER]

Stopping recording...
Captured 48000 samples. Transcribing...

Transcribed text: Hello, this is a test of the speech to text system.

Text copied to clipboard! You can now paste it anywhere.
```

### Useful Commands

```bash
# Rebuild the application
make build

# Run tests
make test

# Run with verbose output
go run ./cmd/speech-to-clipboard

# Clean and rebuild
make clean && make build
```

### Troubleshooting

**Microphone not working?**
- Check System Preferences → Security & Privacy → Microphone
- Ensure the terminal app has microphone access permissions

**API errors?**
- Verify your API key is set: `echo $OPENAI_API_KEY`
- Check your OpenAI account has credits available
- Ensure you have internet connectivity

**Build errors after system restart?**
- Make sure you've set the OPENAI_API_KEY environment variable again

### Next Steps

- Read the full [README.md](README.md) for detailed documentation
- Check out the code in `pkg/` to understand the architecture
- Modify `cmd/speech-to-clipboard/main.go` to customize behavior
- Add new speech-to-text providers by implementing the `Transcriber` interface

Enjoy your speech-to-clipboard application!
