package ui

import (
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
)

func TestIndent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		s    string
		n    int
		want string
	}{
		{"empty string", "", 2, ""},
		{"zero indent", "a", 0, "a"},
		{"negative indent", "a", -1, "a"},
		{"single space", "a", 1, " a\n"},
		{"multiple lines", "a\nb", 2, "  a\n  b\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := Indent(tt.s, tt.n); got != tt.want {
				t.Fatalf("Indent(%q, %d) = %q, want %q", tt.s, tt.n, got, tt.want)
			}
		})
	}
}

func TestAltScreenView(t *testing.T) {
	t.Parallel()

	if !AltScreenView("body").AltScreen {
		t.Fatal("expected AltScreenView to use the alt screen")
	}
	if !WithAltScreen(tea.NewView("body")).AltScreen {
		t.Fatal("expected WithAltScreen to use the alt screen")
	}
}

func TestStatusMessageTimeout(t *testing.T) {
	t.Parallel()

	if StatusMessageTimeout != 3*time.Second {
		t.Fatalf("expected 3s, got %v", StatusMessageTimeout)
	}
}

func TestWaitForStatusMessageTimeout(t *testing.T) {
	t.Parallel()

	msg := WaitForStatusMessageTimeout(time.NewTimer(0))()
	if _, ok := msg.(StatusMessageTimeoutMsg); !ok {
		t.Fatalf("expected StatusMessageTimeoutMsg, got %T", msg)
	}
}
