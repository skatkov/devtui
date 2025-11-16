// devtui/internal/clipboard/clipboard_test.go
package clipboard

import (
	"os"
	"strings"
	"testing"
)

func TestEnhanceError_NonLinux(t *testing.T) {
	// This test verifies that on non-Linux systems, we just return a simple error
	originalErr := os.ErrNotExist
	enhanced := enhanceError(originalErr, "copy")

	if enhanced == nil {
		t.Fatal("Expected error, got nil")
	}

	if !strings.Contains(enhanced.Error(), "clipboard copy failed") {
		t.Errorf("Expected error message to contain 'clipboard copy failed', got: %s", enhanced.Error())
	}

	// Should not contain Linux-specific help messages
	if strings.Contains(enhanced.Error(), "wl-clipboard") {
		t.Errorf("Expected no Linux-specific help on non-Linux, got: %s", enhanced.Error())
	}
}

func TestEnhanceError_WaylandSession(t *testing.T) {
	// Skip if not on Linux
	if os.Getenv("CI") != "" {
		t.Skip("Skipping session-dependent test in CI")
	}

	// Save and restore original session type
	originalSession := os.Getenv("XDG_SESSION_TYPE")
	defer func() {
		if originalSession == "" {
			os.Unsetenv("XDG_SESSION_TYPE")
		} else {
			os.Setenv("XDG_SESSION_TYPE", originalSession)
		}
	}()

	// Set Wayland session
	os.Setenv("XDG_SESSION_TYPE", "wayland")

	// Create an error that looks like a missing executable
	originalErr := &os.PathError{
		Op:   "exec",
		Path: "wl-copy",
		Err:  os.ErrNotExist,
	}

	enhanced := enhanceError(originalErr, "copy")

	if enhanced == nil {
		t.Fatal("Expected error, got nil")
	}

	errMsg := enhanced.Error()
	if !strings.Contains(errMsg, "clipboard copy failed") {
		t.Errorf("Expected error message to contain 'clipboard copy failed', got: %s", errMsg)
	}

	if !strings.Contains(errMsg, "wl-clipboard") {
		t.Errorf("Expected error message to mention wl-clipboard, got: %s", errMsg)
	}
}

func TestEnhanceError_X11Session(t *testing.T) {
	// Skip if not on Linux
	if os.Getenv("CI") != "" {
		t.Skip("Skipping session-dependent test in CI")
	}

	// Save and restore original session type
	originalSession := os.Getenv("XDG_SESSION_TYPE")
	defer func() {
		if originalSession == "" {
			os.Unsetenv("XDG_SESSION_TYPE")
		} else {
			os.Setenv("XDG_SESSION_TYPE", originalSession)
		}
	}()

	// Set X11 session
	os.Setenv("XDG_SESSION_TYPE", "x11")

	// Create an error that looks like a missing executable
	originalErr := &os.PathError{
		Op:   "exec",
		Path: "xclip",
		Err:  os.ErrNotExist,
	}

	enhanced := enhanceError(originalErr, "paste")

	if enhanced == nil {
		t.Fatal("Expected error, got nil")
	}

	errMsg := enhanced.Error()
	if !strings.Contains(errMsg, "clipboard paste failed") {
		t.Errorf("Expected error message to contain 'clipboard paste failed', got: %s", errMsg)
	}

	if !strings.Contains(errMsg, "xclip") {
		t.Errorf("Expected error message to mention xclip, got: %s", errMsg)
	}
}

func TestEnhanceError_UnknownSession(t *testing.T) {
	// Skip if not on Linux
	if os.Getenv("CI") != "" {
		t.Skip("Skipping session-dependent test in CI")
	}

	// Save and restore original session type
	originalSession := os.Getenv("XDG_SESSION_TYPE")
	defer func() {
		if originalSession == "" {
			os.Unsetenv("XDG_SESSION_TYPE")
		} else {
			os.Setenv("XDG_SESSION_TYPE", originalSession)
		}
	}()

	// Unset session type
	os.Unsetenv("XDG_SESSION_TYPE")

	// Create an error that looks like a missing executable
	originalErr := &os.PathError{
		Op:   "exec",
		Path: "xclip",
		Err:  os.ErrNotExist,
	}

	enhanced := enhanceError(originalErr, "copy")

	if enhanced == nil {
		t.Fatal("Expected error, got nil")
	}

	errMsg := enhanced.Error()

	// Should mention that clipboard manager was not found
	if !strings.Contains(errMsg, "clipboard manager was not found") {
		t.Errorf("Expected error message to mention clipboard manager not found, got: %s", errMsg)
	}

	if !strings.Contains(errMsg, "clipboard copy failed") {
		t.Errorf("Expected error message to contain operation type, got: %s", errMsg)
	}
}

func TestEnhanceError_AllOperations(t *testing.T) {
	// Skip if not on Linux
	if os.Getenv("CI") != "" {
		t.Skip("Skipping session-dependent test in CI")
	}

	// Save and restore original session type
	originalSession := os.Getenv("XDG_SESSION_TYPE")
	defer func() {
		if originalSession == "" {
			os.Unsetenv("XDG_SESSION_TYPE")
		} else {
			os.Setenv("XDG_SESSION_TYPE", originalSession)
		}
	}()

	testCases := []struct {
		name        string
		sessionType string
		operation   string
		wantContain string
	}{
		{
			name:        "wayland copy",
			sessionType: "wayland",
			operation:   "copy",
			wantContain: "wl-clipboard",
		},
		{
			name:        "wayland paste",
			sessionType: "wayland",
			operation:   "paste",
			wantContain: "wl-clipboard",
		},
		{
			name:        "x11 copy",
			sessionType: "x11",
			operation:   "copy",
			wantContain: "xclip",
		},
		{
			name:        "x11 paste",
			sessionType: "x11",
			operation:   "paste",
			wantContain: "xclip",
		},
		{
			name:        "unknown copy",
			sessionType: "",
			operation:   "copy",
			wantContain: "clipboard manager was not found",
		},
		{
			name:        "unknown paste",
			sessionType: "",
			operation:   "paste",
			wantContain: "clipboard manager was not found",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.sessionType == "" {
				os.Unsetenv("XDG_SESSION_TYPE")
			} else {
				os.Setenv("XDG_SESSION_TYPE", tc.sessionType)
			}

			originalErr := &os.PathError{
				Op:   "exec",
				Path: "xclip",
				Err:  os.ErrNotExist,
			}

			enhanced := enhanceError(originalErr, tc.operation)

			if enhanced == nil {
				t.Fatal("Expected error, got nil")
			}

			errMsg := enhanced.Error()

			if !strings.Contains(errMsg, tc.wantContain) {
				t.Errorf("Expected error to contain %q, got: %s", tc.wantContain, errMsg)
			}

			expectedOp := "clipboard " + tc.operation + " failed"
			if !strings.Contains(errMsg, expectedOp) {
				t.Errorf("Expected error to contain %q, got: %s", expectedOp, errMsg)
			}
		})
	}
}
