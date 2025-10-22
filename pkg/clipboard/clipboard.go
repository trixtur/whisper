package clipboard

import (
	"fmt"

	"github.com/atotto/clipboard"
)

// Manager handles clipboard operations
type Manager interface {
	Write(text string) error
	Read() (string, error)
}

type clipboardManager struct{}

// NewManager creates a new clipboard manager
func NewManager() Manager {
	return &clipboardManager{}
}

// Write writes text to the system clipboard
func (c *clipboardManager) Write(text string) error {
	if err := clipboard.WriteAll(text); err != nil {
		return fmt.Errorf("failed to write to clipboard: %w", err)
	}
	return nil
}

// Read reads text from the system clipboard
func (c *clipboardManager) Read() (string, error) {
	text, err := clipboard.ReadAll()
	if err != nil {
		return "", fmt.Errorf("failed to read from clipboard: %w", err)
	}
	return text, nil
}

// MockManager is a mock implementation for testing
type MockManager struct {
	content string
	err     error
}

// NewMockManager creates a mock clipboard manager
func NewMockManager() *MockManager {
	return &MockManager{
		content: "",
		err:     nil,
	}
}

// Write stores text in the mock clipboard
func (m *MockManager) Write(text string) error {
	if m.err != nil {
		return m.err
	}
	m.content = text
	return nil
}

// Read retrieves text from the mock clipboard
func (m *MockManager) Read() (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.content, nil
}

// SetError sets an error to be returned by operations
func (m *MockManager) SetError(err error) {
	m.err = err
}

// GetContent returns the current content (for testing)
func (m *MockManager) GetContent() string {
	return m.content
}
