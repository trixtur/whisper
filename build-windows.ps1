# Windows Build Script for speech-to-clipboard
# PowerShell version - run with: powershell -ExecutionPolicy Bypass -File build-windows.ps1

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Speech-to-Clipboard Windows Build Script" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# Check if Go is installed
Write-Host "Checking for Go..." -NoNewline
try {
    $goVersion = go version 2>$null
    if ($LASTEXITCODE -eq 0) {
        Write-Host " [OK]" -ForegroundColor Green
        Write-Host "  $goVersion"
    } else {
        throw
    }
} catch {
    Write-Host " [ERROR]" -ForegroundColor Red
    Write-Host "Go is not installed or not in PATH" -ForegroundColor Red
    Write-Host "Please install Go from https://go.dev/dl/"
    exit 1
}
Write-Host ""

# Check if MSYS2 is installed
Write-Host "Checking for MSYS2..." -NoNewline
$msys2Path = "C:\msys64"
if (Test-Path $msys2Path) {
    Write-Host " [OK]" -ForegroundColor Green
    Write-Host "  Found at $msys2Path"
} else {
    Write-Host " [ERROR]" -ForegroundColor Red
    Write-Host ""
    Write-Host "MSYS2 not found at $msys2Path" -ForegroundColor Red
    Write-Host ""
    Write-Host "Please install MSYS2 first:"
    Write-Host "1. Download from https://www.msys2.org/"
    Write-Host "2. Run the installer"
    Write-Host "3. Run this script again"
    Write-Host ""
    Write-Host "See docs\WINDOWS_BUILD.md for detailed instructions"
    exit 1
}
Write-Host ""

# Check if MinGW gcc is installed
Write-Host "Checking for MinGW GCC..." -NoNewline
$gccPath = "$msys2Path\mingw64\bin\gcc.exe"
if (-not (Test-Path $gccPath)) {
    Write-Host " [WARNING]" -ForegroundColor Yellow
    Write-Host "Installing required MSYS2 packages..." -ForegroundColor Yellow
    Write-Host ""

    # Update MSYS2
    & "$msys2Path\usr\bin\bash" -lc "pacman -Syu --noconfirm"

    # Install packages
    & "$msys2Path\usr\bin\bash" -lc "pacman -S --noconfirm mingw-w64-x86_64-gcc mingw-w64-x86_64-pkg-config mingw-w64-x86_64-portaudio"

    if ($LASTEXITCODE -ne 0) {
        Write-Host "[ERROR] Failed to install MSYS2 packages" -ForegroundColor Red
        exit 1
    }
    Write-Host "[OK] MSYS2 packages installed" -ForegroundColor Green
} else {
    Write-Host " [OK]" -ForegroundColor Green
}
Write-Host ""

# Add MinGW to PATH for this session
Write-Host "Adding MinGW to PATH..." -NoNewline
$mingwBin = "$msys2Path\mingw64\bin"
$env:PATH = "$mingwBin;$env:PATH"
Write-Host " [OK]" -ForegroundColor Green
Write-Host ""

# Verify GCC
Write-Host "Verifying GCC..." -NoNewline
try {
    $gccVersion = & gcc --version 2>$null | Select-Object -First 1
    if ($LASTEXITCODE -eq 0) {
        Write-Host " [OK]" -ForegroundColor Green
        Write-Host "  $gccVersion"
    } else {
        throw
    }
} catch {
    Write-Host " [ERROR]" -ForegroundColor Red
    Write-Host "GCC not found in PATH after adding MinGW" -ForegroundColor Red
    Write-Host "Please add C:\msys64\mingw64\bin to your system PATH manually"
    exit 1
}
Write-Host ""

# Verify pkg-config
Write-Host "Checking for pkg-config..." -NoNewline
try {
    $pkgConfigVersion = & pkg-config --version 2>$null
    if ($LASTEXITCODE -eq 0) {
        Write-Host " [OK]" -ForegroundColor Green
        Write-Host "  Version: $pkgConfigVersion"
    } else {
        throw
    }
} catch {
    Write-Host " [WARNING]" -ForegroundColor Yellow
    Write-Host "Installing pkg-config via MSYS2..." -ForegroundColor Yellow
    & "$msys2Path\usr\bin\bash" -lc "pacman -S --noconfirm mingw-w64-x86_64-pkg-config"
    if ($LASTEXITCODE -ne 0) {
        Write-Host "[ERROR] Failed to install pkg-config" -ForegroundColor Red
        exit 1
    }
}
Write-Host ""

# Verify PortAudio
Write-Host "Checking for PortAudio..." -NoNewline
& pkg-config --exists portaudio-2.0 2>$null
if ($LASTEXITCODE -ne 0) {
    Write-Host " [WARNING]" -ForegroundColor Yellow
    Write-Host "Installing PortAudio via MSYS2..." -ForegroundColor Yellow
    & "$msys2Path\usr\bin\bash" -lc "pacman -S --noconfirm mingw-w64-x86_64-portaudio"
    if ($LASTEXITCODE -ne 0) {
        Write-Host "[ERROR] Failed to install PortAudio" -ForegroundColor Red
        exit 1
    }
    Write-Host "[OK] PortAudio installed" -ForegroundColor Green
} else {
    Write-Host " [OK]" -ForegroundColor Green
    $portaudioFlags = & pkg-config --cflags --libs portaudio-2.0 2>$null
    Write-Host "  $portaudioFlags"
}
Write-Host ""

# Enable CGO
Write-Host "Enabling CGO..." -NoNewline
$env:CGO_ENABLED = 1
Write-Host " [OK]" -ForegroundColor Green
Write-Host ""

# Download Go dependencies
Write-Host "Downloading Go dependencies..." -NoNewline
go mod download 2>$null
if ($LASTEXITCODE -ne 0) {
    Write-Host " [ERROR]" -ForegroundColor Red
    Write-Host "Failed to download Go dependencies" -ForegroundColor Red
    exit 1
}
Write-Host " [OK]" -ForegroundColor Green
Write-Host ""

# Build the application
Write-Host "Building application..." -NoNewline
go build -o speech-to-clipboard.exe .\cmd\speech-to-clipboard 2>&1 | Out-Null
if ($LASTEXITCODE -ne 0) {
    Write-Host " [ERROR]" -ForegroundColor Red
    Write-Host ""
    Write-Host "Build failed. Running verbose build to show errors:" -ForegroundColor Yellow
    Write-Host ""
    go build -v -o speech-to-clipboard.exe .\cmd\speech-to-clipboard
    Write-Host ""
    Write-Host "Common issues:" -ForegroundColor Yellow
    Write-Host "- Make sure C:\msys64\mingw64\bin is in your PATH"
    Write-Host "- Try closing and reopening your terminal"
    Write-Host "- See docs\WINDOWS_BUILD.md for troubleshooting"
    exit 1
}
Write-Host " [OK]" -ForegroundColor Green
Write-Host ""

Write-Host "========================================" -ForegroundColor Green
Write-Host "[SUCCESS] Build completed successfully!" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Green
Write-Host ""
Write-Host "The application has been built: speech-to-clipboard.exe"
Write-Host ""
Write-Host "To run the application:" -ForegroundColor Cyan
Write-Host "  1. Set your OpenAI API key:"
Write-Host "     `$env:OPENAI_API_KEY = 'your-api-key-here'" -ForegroundColor Gray
Write-Host ""
Write-Host "  2. Run the application:"
Write-Host "     .\speech-to-clipboard.exe" -ForegroundColor Gray
Write-Host ""
Write-Host "Note: You may need to add C:\msys64\mingw64\bin to your" -ForegroundColor Yellow
Write-Host "system PATH permanently for future builds." -ForegroundColor Yellow
Write-Host ""
Write-Host "To add to PATH permanently, run as Administrator:"
Write-Host "  [Environment]::SetEnvironmentVariable('Path', `$env:Path + ';C:\msys64\mingw64\bin', 'Machine')" -ForegroundColor Gray
Write-Host ""
Write-Host "See docs\WINDOWS_BUILD.md for more information."
Write-Host ""
