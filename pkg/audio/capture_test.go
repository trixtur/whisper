package audio

import (
	"bytes"
	"testing"
)

func TestWriteInt16(t *testing.T) {
	tests := []struct {
		name  string
		value uint16
		want  []byte
	}{
		{
			name:  "write zero",
			value: 0,
			want:  []byte{0x00, 0x00},
		},
		{
			name:  "write small value",
			value: 256,
			want:  []byte{0x00, 0x01},
		},
		{
			name:  "write max value",
			value: 0xFFFF,
			want:  []byte{0xFF, 0xFF},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := make([]byte, 2)
			writeInt16(buf, tt.value)

			if !bytes.Equal(buf, tt.want) {
				t.Errorf("writeInt16() = %v, want %v", buf, tt.want)
			}
		})
	}
}

func TestWriteInt32(t *testing.T) {
	tests := []struct {
		name  string
		value uint32
		want  []byte
	}{
		{
			name:  "write zero",
			value: 0,
			want:  []byte{0x00, 0x00, 0x00, 0x00},
		},
		{
			name:  "write small value",
			value: 256,
			want:  []byte{0x00, 0x01, 0x00, 0x00},
		},
		{
			name:  "write sample rate",
			value: 16000,
			want:  []byte{0x80, 0x3E, 0x00, 0x00},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := make([]byte, 4)
			writeInt32(buf, tt.value)

			if !bytes.Equal(buf, tt.want) {
				t.Errorf("writeInt32() = %v, want %v", buf, tt.want)
			}
		})
	}
}

func TestSaveToWAV(t *testing.T) {
	tests := []struct {
		name    string
		data    []int16
		wantErr bool
	}{
		{
			name:    "save valid audio data",
			data:    []int16{100, 200, 300, 400, 500},
			wantErr: false,
		},
		{
			name:    "save empty audio data",
			data:    []int16{},
			wantErr: false,
		},
		{
			name:    "save single sample",
			data:    []int16{42},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			err := SaveToWAV(tt.data, buf)

			if (err != nil) != tt.wantErr {
				t.Errorf("SaveToWAV() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Check that WAV header is present (44 bytes)
				if buf.Len() < 44 {
					t.Errorf("SaveToWAV() output too small, got %d bytes, want at least 44", buf.Len())
				}

				// Check RIFF header
				header := buf.Bytes()
				if string(header[0:4]) != "RIFF" {
					t.Errorf("SaveToWAV() missing RIFF header")
				}

				if string(header[8:12]) != "WAVE" {
					t.Errorf("SaveToWAV() missing WAVE header")
				}

				// Check expected file size
				expectedSize := 44 + len(tt.data)*2
				if buf.Len() != expectedSize {
					t.Errorf("SaveToWAV() size = %d, want %d", buf.Len(), expectedSize)
				}
			}
		})
	}
}

func TestConstants(t *testing.T) {
	if SampleRate != 16000 {
		t.Errorf("SampleRate = %d, want 16000", SampleRate)
	}

	if Channels != 1 {
		t.Errorf("Channels = %d, want 1", Channels)
	}

	if FramesPerBuffer != 1024 {
		t.Errorf("FramesPerBuffer = %d, want 1024", FramesPerBuffer)
	}
}
