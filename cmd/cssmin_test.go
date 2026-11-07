package cmd

import (
	"errors"
	"strings"
	"testing"
	"testing/iotest"
)

func TestCssminCmd(t *testing.T) {
	stdout, _, err := runDevtui(strings.NewReader("a { color: red; margin: 0; }\nb { x: y; }"), "cssmin")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := "a{color:red;margin:0}b{x:y}"; stdout.String() != want {
		t.Fatalf("expected %q, got %q", want, stdout.String())
	}
}

func TestCssminCmdArgumentInput(t *testing.T) {
	stdout, _, err := runDevtui(strings.NewReader(""), "cssmin", "a { b: c; }")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := "a{b:c}"; stdout.String() != want {
		t.Fatalf("expected %q, got %q", want, stdout.String())
	}
}

func TestCssminCmdReadError(t *testing.T) {
	want := errors.New("read failed")

	_, _, err := runDevtui(iotest.ErrReader(want), "cssmin")
	if !errors.Is(err, want) {
		t.Fatalf("expected %v, got %v", want, err)
	}
}
