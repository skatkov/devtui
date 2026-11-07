package cmd

import (
	"errors"
	"strings"
	"testing"
	"testing/iotest"
)

const countTable = "┌──────────┬─┐\n" +
	"│Characters│3│\n" +
	"│Spaces    │1│\n" +
	"│Words     │2│\n" +
	"└──────────┴─┘\n"

func TestCountCmd(t *testing.T) {
	stdout, _, err := runDevtui(strings.NewReader("a b"), "count")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := stdout.String(); got != countTable {
		t.Fatalf("expected %q, got %q", countTable, got)
	}
}

func TestCountCmdArgumentInput(t *testing.T) {
	stdout, _, err := runDevtui(strings.NewReader(""), "count", "a b")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := stdout.String(); got != countTable {
		t.Fatalf("expected %q, got %q", countTable, got)
	}
}

func TestCountCmdReadError(t *testing.T) {
	want := errors.New("read failed")

	_, _, err := runDevtui(iotest.ErrReader(want), "count")
	if !errors.Is(err, want) {
		t.Fatalf("expected %v, got %v", want, err)
	}
}
