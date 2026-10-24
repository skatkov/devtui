package cmd

import (
	"errors"
	"strings"
	"testing"

	"github.com/skatkov/devtui/tui/xml"
)

func TestXmlfmtCmd(t *testing.T) {
	stdout, _, err := runDevtui(strings.NewReader("<a><!-- <c/> --></a>"), "xmlfmt")
	if err != nil {
		t.Fatalf("xmlfmt failed: %v", err)
	}

	if want := "\n<a>\n  <!-- <c/> -->\n</a>\n"; stdout.String() != want {
		t.Fatalf("expected %q, got %q", want, stdout.String())
	}
}

func TestXmlfmtCmdArgumentInput(t *testing.T) {
	stdout, _, err := runDevtui(nil, "xmlfmt", "<a><b>x</b></a>")
	if err != nil {
		t.Fatalf("xmlfmt failed: %v", err)
	}

	if want := "\n<a>\n  <b>x</b>\n</a>\n"; stdout.String() != want {
		t.Fatalf("expected %q, got %q", want, stdout.String())
	}
}

func TestXmlfmtCmdTUI(t *testing.T) {
	errQuit := errors.New("quit")
	started := stubTUIProgram(t, errQuit)

	_, _, err := runDevtui(strings.NewReader("<a><b>x</b></a>"), "xmlfmt", "--tui")
	if !errors.Is(err, errQuit) {
		t.Fatalf("expected the program error, got %v", err)
	}

	model, ok := (*started).(xml.XMLFormatterModel)
	if !ok {
		t.Fatalf("expected an xml model, got %T", *started)
	}
	if want := "<a><b>x</b></a>"; model.Content != want {
		t.Fatalf("expected %q, got %q", want, model.Content)
	}
}
