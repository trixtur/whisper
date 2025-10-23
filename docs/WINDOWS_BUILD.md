# Building on Windows

This guide provides step-by-step instructions for building the speech-to-clipboard application on Windows.

## Prerequisites

The application requires CGO (C bindings) due to the PortAudio dependency. You'll need:

1. **Go** (1.24.2 or later)
2. **MSYS2** (provides MinGW-w64 compiler and pkg-config)
3. **PortAudio library**

## Installation Steps

### 1. Install MSYS2

1. Download MSYS2 from https://www.msys2.org/
2. Run the installer (e.g., `msys2-x86_64-20240113.exe`)
3. Follow the installation wizard (default location: `C:\msys64`)
4. After installation, MSYS2 will open a terminal window

### 2. Update MSYS2

In the MSYS2 terminal, run:

```bash
pacman -Syu
```

Close and reopen the MSYS2 terminal if prompted, then run:

```bash
pacman -Su
```

### 3. Install Build Tools and PortAudio

In the MSYS2 terminal (MSYS2 MINGW64), install the required packages:

```bash
# Install MinGW-w64 toolchain
pacman -S mingw-w64-x86_64-gcc

# Install pkg-config
pacman -S mingw-w64-x86_64-pkg-config

# Install PortAudio
pacman -S mingw-w64-x86_64-portaudio

# Install make (optional, for using Makefile)
pacman -S make
```

### 4. Configure Environment Variables

You need to add the MinGW binaries to your Windows PATH so Go can find the C compiler and pkg-config.

**Option A: Temporary (current PowerShell/CMD session only)**

In PowerShell:
```powershell
$env:PATH = "C:\msys64\mingw64\bin;$env:PATH"
```

In CMD:
```cmd
set PATH=C:\msys64\mingw64\bin;%PATH%
```

**Option B: Permanent (recommended)**

1. Open "Edit the system environment variables" from Windows search
2. Click "Environment Variables"
3. Under "User variables" or "System variables", find and edit `PATH`
4. Add new entry: `C:\msys64\mingw64\bin`
5. Click OK to save
6. **Restart your terminal/IDE** for changes to take effect

### 5. Verify Installation

Open a new terminal (PowerShell or CMD) and verify:

```powershell
# Check Go
go version

# Check GCC
gcc --version

# Check pkg-config
pkg-config --version

# Check PortAudio is found
pkg-config --cflags --libs portaudio-2.0
```

Expected output for the last command:
```
-IC:/msys64/mingw64/include -LC:/msys64/mingw64/lib -lportaudio
```

### 6. Set OpenAI API Key

In PowerShell:
```powershell
$env:OPENAI_API_KEY = "your-api-key-here"
```

In CMD:
```cmd
set OPENAI_API_KEY=your-api-key-here
```

To make it permanent, add it to your system environment variables (same as step 4).

### 7. Build the Application

```bash
# Clone the repository (if not already done)
git clone <repository-url>
cd speech-to-clipboard

# Download Go dependencies
go mod download

# Build the application
go build -o speech-to-clipboard.exe ./cmd/speech-to-clipboard
```

### 8. Run the Application

```bash
.\speech-to-clipboard.exe
```

## Running Tests

```bash
go test ./...
```

## Troubleshooting

### Error: "gcc: command not found"

**Problem:** The C compiler is not in your PATH.

**Solution:**
- Make sure you installed `mingw-w64-x86_64-gcc` via MSYS2
- Add `C:\msys64\mingw64\bin` to your PATH (see step 4)
- Restart your terminal

### Error: "Package portaudio-2.0 was not found"

**Problem:** pkg-config cannot find PortAudio.

**Solution:**
- Make sure you installed `mingw-w64-x86_64-portaudio` via MSYS2
- Verify pkg-config works: `pkg-config --version`
- Check PortAudio: `pkg-config --cflags portaudio-2.0`
- Ensure `C:\msys64\mingw64\bin` is in your PATH

### Error: "cgo: C compiler not available"

**Problem:** CGO is disabled or C compiler not found.

**Solution:**
- Enable CGO: `set CGO_ENABLED=1` (CMD) or `$env:CGO_ENABLED = 1` (PowerShell)
- Verify GCC is installed and in PATH
- Restart your terminal/IDE

### Build works in MSYS2 but not in PowerShell/CMD

**Problem:** Environment variables not set correctly in Windows.

**Solution:**
- Make sure `C:\msys64\mingw64\bin` is in your Windows PATH (not just MSYS2)
- Don't build in the MSYS2 terminal - use PowerShell/CMD after setting PATH
- Restart your terminal after changing environment variables

### "Access denied" or permission errors

**Problem:** Windows security or antivirus blocking execution.

**Solution:**
- Run your terminal as Administrator
- Add exceptions in Windows Defender/antivirus for the build directory
- Check that you have write permissions in the project directory

## Alternative: Using MSYS2 Terminal

If you prefer to work entirely within MSYS2 (instead of Windows terminal):

1. Open "MSYS2 MinGW 64-bit" from Start menu
2. Navigate to your project directory
3. Set your API key: `export OPENAI_API_KEY="your-key"`
4. Build: `go build -o speech-to-clipboard.exe ./cmd/speech-to-clipboard`
5. Run: `./speech-to-clipboard.exe`

## Quick Setup Script

Save this as `setup-windows.ps1` and run in PowerShell:

```powershell
# Check if MSYS2 is installed
if (-not (Test-Path "C:\msys64")) {
    Write-Host "ERROR: MSYS2 not found at C:\msys64"
    Write-Host "Please install MSYS2 from https://www.msys2.org/"
    exit 1
}

# Add MSYS2 to PATH for this session
$env:PATH = "C:\msys64\mingw64\bin;$env:PATH"

# Verify tools
Write-Host "Checking prerequisites..."
Write-Host "Go version:"
go version

Write-Host "`nGCC version:"
gcc --version

Write-Host "`npkg-config version:"
pkg-config --version

Write-Host "`nPortAudio check:"
pkg-config --cflags --libs portaudio-2.0

Write-Host "`nAll checks passed! You can now build the application."
Write-Host "Run: go build -o speech-to-clipboard.exe ./cmd/speech-to-clipboard"
```

## Need More Help?

- MSYS2 Documentation: https://www.msys2.org/docs/what-is-msys2/
- PortAudio Documentation: http://www.portaudio.com/docs/
- Go CGO Documentation: https://pkg.go.dev/cmd/cgo

If you continue to have issues, please open a GitHub issue with:
- Your Windows version
- Go version (`go version`)
- Error messages
- Output of `pkg-config --cflags portaudio-2.0`
