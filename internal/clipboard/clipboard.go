// Package clipboard provides clipboard operations with helpful error messages
// that detect the Linux session type and suggest appropriate dependencies.
package clipboard

import (
	"fmt"
	"os"
	"runtime"

	"github.com/tiagomelo/go-clipboard/clipboard"
)

// copyText and pasteText reach the system clipboard; tests replace them.
var (
	copyText  = func(text string) error { return clipboard.New().CopyText(text) }
	pasteText = func() (string, error) { return clipboard.New().PasteText() }
)

// Copy copies text to the clipboard with enhanced error messages
func Copy(text string) error {
	err := copyText(text)
	if err != nil {
		return enhanceError(err, "copy", runtime.GOOS, os.Getenv("XDG_SESSION_TYPE"))
	}
	return nil
}

// Paste retrieves text from the clipboard with enhanced error messages
func Paste() (string, error) {
	text, err := pasteText()
	if err != nil {
		return "", enhanceError(err, "paste", runtime.GOOS, os.Getenv("XDG_SESSION_TYPE"))
	}
	return text, nil
}

// enhanceError adds helpful context to clipboard errors based on the OS and session type
func enhanceError(originalErr error, operation, goos, sessionType string) error {
	if goos != "linux" {
		// On non-Linux systems, return the original error
		return fmt.Errorf("clipboard %s failed: %w", operation, originalErr)
	}

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
