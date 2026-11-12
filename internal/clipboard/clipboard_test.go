package clipboard

import (
	"errors"
	"testing"
)

func TestEnhanceError(t *testing.T) {
	t.Parallel()

	original := errors.New("exec: not found")
	tests := []struct {
		name        string
		goos        string
		sessionType string
		want        string
	}{
		{"darwin keeps original error", "darwin", "", "clipboard copy failed: exec: not found"},
		{"windows keeps original error", "windows", "x11", "clipboard copy failed: exec: not found"},
		{"linux wayland", "linux", "wayland", "clipboard copy failed: Install 'wl-clipboard' for clipboard support"},
		{"linux x11", "linux", "x11", "clipboard copy failed: Install 'xclip' for clipboard support"},
		{"linux unknown session", "linux", "tty", "clipboard copy failed: clipboard manager was not found on this system"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := enhanceError(original, "copy", tt.goos, tt.sessionType)
			if err.Error() != tt.want {
				t.Fatalf("enhanceError() = %q, want %q", err.Error(), tt.want)
			}
		})
	}
}

func TestEnhanceErrorWrapsOriginalOutsideLinux(t *testing.T) {
	t.Parallel()

	original := errors.New("exec: not found")
	if err := enhanceError(original, "paste", "darwin", ""); !errors.Is(err, original) {
		t.Fatalf("expected %v to wrap %v", err, original)
	}
}

func stubClipboard(t *testing.T, copyFn func(string) error, pasteFn func() (string, error)) {
	t.Helper()

	origCopy, origPaste := copyText, pasteText
	copyText, pasteText = copyFn, pasteFn
	t.Cleanup(func() { copyText, pasteText = origCopy, origPaste })
}

func TestCopy(t *testing.T) {
	var copied string
	stubClipboard(t, func(text string) error { copied = text; return nil }, nil)

	if err := Copy("hello"); err != nil {
		t.Fatalf("Copy returned error: %v", err)
	}
	if copied != "hello" {
		t.Fatalf("expected %q to be copied, got %q", "hello", copied)
	}
}

func TestCopyFailureReturnsError(t *testing.T) {
	stubClipboard(t, func(string) error { return errors.New("exec: not found") }, nil)

	err := Copy("hello")
	if err == nil {
		t.Fatal("expected error when clipboard copy fails")
	}
}

func TestPaste(t *testing.T) {
	stubClipboard(t, nil, func() (string, error) { return "pasted", nil })

	got, err := Paste()
	if err != nil {
		t.Fatalf("Paste returned error: %v", err)
	}
	if got != "pasted" {
		t.Fatalf("expected %q, got %q", "pasted", got)
	}
}

func TestPasteFailureReturnsError(t *testing.T) {
	stubClipboard(t, nil, func() (string, error) { return "partial", errors.New("exec: not found") })

	got, err := Paste()
	if err == nil {
		t.Fatal("expected error when clipboard paste fails")
	}
	if got != "" {
		t.Fatalf("expected empty text on error, got %q", got)
	}
}
