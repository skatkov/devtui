package ui

import (
	"strings"
	"testing"
)

func TestAppNameStyleIsBold(t *testing.T) {
	t.Parallel()

	if got := AppNameStyle("x"); !strings.HasPrefix(got, "\x1b[1;") {
		t.Fatalf("expected bold app name, got %q", got)
	}
}

func TestNewStyle(t *testing.T) {
	t.Parallel()

	s := NewStyle()

	if top, right, bottom, left := s.Base.GetPadding(); top != 1 || right != 4 || bottom != 1 || left != 2 {
		t.Fatalf("Base padding = %d %d %d %d, want 1 4 1 2", top, right, bottom, left)
	}
	if top, right, bottom, left := s.Title.GetPadding(); top != 0 || right != 1 || bottom != 0 || left != 1 {
		t.Fatalf("Title padding = %d %d %d %d, want 0 1 0 1", top, right, bottom, left)
	}
}
