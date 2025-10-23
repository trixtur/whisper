@echo off
REM Windows Build Script for speech-to-clipboard
REM This script helps set up the build environment and compile the application

echo ========================================
echo Speech-to-Clipboard Windows Build Script
echo ========================================
echo.

REM Check if Go is installed
where go >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo [ERROR] Go is not installed or not in PATH
    echo Please install Go from https://go.dev/dl/
    exit /b 1
)

echo [OK] Go is installed
go version
echo.

REM Check if MSYS2 is installed
if not exist "C:\msys64" (
    echo [ERROR] MSYS2 not found at C:\msys64
    echo.
    echo Please install MSYS2 first:
    echo 1. Download from https://www.msys2.org/
    echo 2. Run the installer
    echo 3. Run this script again
    echo.
    echo See docs\WINDOWS_BUILD.md for detailed instructions
    exit /b 1
)

echo [OK] MSYS2 found at C:\msys64
echo.

REM Check if MinGW gcc is installed
if not exist "C:\msys64\mingw64\bin\gcc.exe" (
    echo [WARNING] MinGW GCC not found
    echo Installing required MSYS2 packages...
    echo.
    C:\msys64\usr\bin\bash -lc "pacman -Syu --noconfirm"
    C:\msys64\usr\bin\bash -lc "pacman -S --noconfirm mingw-w64-x86_64-gcc mingw-w64-x86_64-pkg-config mingw-w64-x86_64-portaudio"
    if %ERRORLEVEL% NEQ 0 (
        echo [ERROR] Failed to install MSYS2 packages
        exit /b 1
    )
    echo [OK] MSYS2 packages installed
) else (
    echo [OK] MinGW GCC is installed
)
echo.

REM Add MinGW to PATH for this session
echo Adding C:\msys64\mingw64\bin to PATH...
set PATH=C:\msys64\mingw64\bin;%PATH%
echo.

REM Verify GCC
where gcc >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo [ERROR] GCC not found in PATH after adding MinGW
    echo Please add C:\msys64\mingw64\bin to your system PATH manually
    exit /b 1
)

echo [OK] GCC is accessible
gcc --version | findstr "gcc"
echo.

REM Verify pkg-config
where pkg-config >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo [ERROR] pkg-config not found
    echo Installing pkg-config via MSYS2...
    C:\msys64\usr\bin\bash -lc "pacman -S --noconfirm mingw-w64-x86_64-pkg-config"
    if %ERRORLEVEL% NEQ 0 (
        echo [ERROR] Failed to install pkg-config
        exit /b 1
    )
)

echo [OK] pkg-config is accessible
pkg-config --version
echo.

REM Verify PortAudio
echo Checking for PortAudio library...
pkg-config --exists portaudio-2.0
if %ERRORLEVEL% NEQ 0 (
    echo [WARNING] PortAudio not found
    echo Installing PortAudio via MSYS2...
    C:\msys64\usr\bin\bash -lc "pacman -S --noconfirm mingw-w64-x86_64-portaudio"
    if %ERRORLEVEL% NEQ 0 (
        echo [ERROR] Failed to install PortAudio
        exit /b 1
    )
    echo [OK] PortAudio installed
) else (
    echo [OK] PortAudio is installed
    pkg-config --cflags --libs portaudio-2.0
)
echo.

REM Enable CGO
echo Enabling CGO...
set CGO_ENABLED=1
echo.

REM Download Go dependencies
echo Downloading Go dependencies...
go mod download
if %ERRORLEVEL% NEQ 0 (
    echo [ERROR] Failed to download Go dependencies
    exit /b 1
)
echo [OK] Dependencies downloaded
echo.

REM Build the application
echo Building application with static linking...
go build -ldflags "-extldflags=-static" -o speech-to-clipboard.exe .\cmd\speech-to-clipboard
if %ERRORLEVEL% NEQ 0 (
    echo.
    echo [ERROR] Build failed
    echo.
    echo Common issues:
    echo - Make sure C:\msys64\mingw64\bin is in your PATH
    echo - Try closing and reopening your terminal
    echo - See docs\WINDOWS_BUILD.md for troubleshooting
    exit /b 1
)

echo.
echo ========================================
echo [SUCCESS] Build completed successfully!
echo ========================================
echo.
echo The application has been built: speech-to-clipboard.exe
echo.
echo To run the application:
echo   1. Set your OpenAI API key:
echo      set OPENAI_API_KEY=your-api-key-here
echo.
echo   2. Run the application:
echo      speech-to-clipboard.exe
echo.
echo Note: You may need to add C:\msys64\mingw64\bin to your
echo system PATH permanently for future builds.
echo.
echo See docs\WINDOWS_BUILD.md for more information.
echo.

exit /b 0
