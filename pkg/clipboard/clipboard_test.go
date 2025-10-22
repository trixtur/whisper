package clipboard

import (
	"fmt"
	"testing"
)

func TestMockManager_Write(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		wantErr bool
	}{
		{
			name:    "write simple text",
			text:    "Hello, World!",
			wantErr: false,
		},
		{
			name:    "write empty text",
			text:    "",
			wantErr: false,
		},
		{
			name:    "write multiline text",
			text:    "Line 1\nLine 2\nLine 3",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mgr := NewMockManager()
			err := mgr.Write(tt.text)

			if (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && mgr.GetContent() != tt.text {
				t.Errorf("Write() content = %v, want %v", mgr.GetContent(), tt.text)
			}
		})
	}
}

func TestMockManager_Read(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantErr bool
	}{
		{
			name:    "read simple text",
			content: "Hello, World!",
			wantErr: false,
		},
		{
			name:    "read empty text",
			content: "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mgr := NewMockManager()
			mgr.Write(tt.content)

			got, err := mgr.Read()
			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.content {
				t.Errorf("Read() = %v, want %v", got, tt.content)
			}
		})
	}
}

func TestMockManager_Error(t *testing.T) {
	mgr := NewMockManager()
	expectedErr := fmt.Errorf("test error")
	mgr.SetError(expectedErr)

	if err := mgr.Write("test"); err == nil {
		t.Error("Write() expected error, got nil")
	}

	if _, err := mgr.Read(); err == nil {
		t.Error("Read() expected error, got nil")
	}
}
