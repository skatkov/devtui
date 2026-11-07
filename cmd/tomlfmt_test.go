package cmd

import (
	"errors"
	"strings"
	"testing"

	"github.com/skatkov/devtui/tui/toml"
)

func TestTomlfmtCmd(t *testing.T) {
	stdout, _, err := runDevtui(strings.NewReader("a=1"), "tomlfmt")
	if err != nil {
		t.Fatalf("tomlfmt failed: %v", err)
	}

	if want := "a = 1\n\n"; stdout.String() != want {
		t.Fatalf("expected %q, got %q", want, stdout.String())
	}
}

func TestTomlfmtCmdArgumentInput(t *testing.T) {
	stdout, _, err := runDevtui(nil, "tomlfmt", "a=1")
	if err != nil {
		t.Fatalf("tomlfmt failed: %v", err)
	}

	if want := "a = 1\n\n"; stdout.String() != want {
		t.Fatalf("expected %q, got %q", want, stdout.String())
	}
}

func TestTomlfmtCmdInvalidInput(t *testing.T) {
	_, _, err := runDevtui(strings.NewReader("a ="), "tomlfmt")
	if err == nil {
		t.Fatal("expected an error for invalid TOML")
	}
}

func TestTomlfmtCmdTUI(t *testing.T) {
	errQuit := errors.New("quit")
	started := stubTUIProgram(t, errQuit)

	_, _, err := runDevtui(strings.NewReader("a=1"), "tomlfmt", "--tui")
	if !errors.Is(err, errQuit) {
		t.Fatalf("expected the program error, got %v", err)
	}

	model, ok := (*started).(toml.TomlFormatModel)
	if !ok {
		t.Fatalf("expected a toml model, got %T", *started)
	}
	if want := "a = 1\n"; model.FormattedContent != want {
		t.Fatalf("expected %q, got %q", want, model.FormattedContent)
	}
}

func TestTomlfmtCmdTUIInvalidInput(t *testing.T) {
	started := stubTUIProgram(t, nil)

	_, _, err := runDevtui(strings.NewReader("a ="), "tomlfmt", "--tui")
	if err == nil {
		t.Fatal("expected an error for invalid TOML")
	}
	if *started != nil {
		t.Fatal("expected no TUI to start")
	}
}
