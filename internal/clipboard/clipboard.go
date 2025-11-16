// Package clipboard provides clipboard operations with helpful error messages
// that detect the Linux session type and suggest appropriate dependencies.
package clipboard

import (
	"fmt"
	"os"
	"runtime"

	"github.com/tiagomelo/go-clipboard/clipboard"
)

// Copy copies text to the clipboard with enhanced error messages
func Copy(text string) error {
	c := clipboard.New()
	err := c.CopyText(text)
	if err != nil {
		return enhanceError(err, "copy")
	}
	return nil
}

// Paste retrieves text from the clipboard with enhanced error messages
func Paste() (string, error) {
	c := clipboard.New()
	text, err := c.PasteText()
	if err != nil {
		return "", enhanceError(err, "paste")
	}
	return text, nil
}

// enhanceError adds helpful context to clipboard errors based on the OS and session type
func enhanceError(originalErr error, operation string) error {
	if runtime.GOOS != "linux" {
		// On non-Linux systems, return the original error
		return fmt.Errorf("clipboard %s failed: %w", operation, originalErr)
	}

	sessionType := os.Getenv("XDG_SESSION_TYPE")

	var helpMsg string
	switch sessionType {
	case "wayland":
		helpMsg = `Install 'wl-clipboard' for clipboard support`
	case "x11":
		helpMsg = `Install 'xclip' for clipboard support`
	default:
		helpMsg = `clipboard manager was not found on this system`
	}

	return fmt.Errorf("clipboard %s failed: %s", operation, helpMsg)
}
