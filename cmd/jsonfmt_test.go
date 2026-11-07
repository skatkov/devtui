package cmd

import (
	"strings"
	"testing"
)

func TestJsonfmtCmd(t *testing.T) {
	stdout, _, err := runDevtui(strings.NewReader(`{"a":1}`), "jsonfmt")
	if err != nil {
		t.Fatalf("jsonfmt failed: %v", err)
	}

	if want := "{\n  \"a\": 1\n}\n\n"; stdout.String() != want {
		t.Fatalf("expected %q, got %q", want, stdout.String())
	}
}

func TestJsonfmtCmdArgumentInput(t *testing.T) {
	stdout, _, err := runDevtui(nil, "jsonfmt", `{"a":1}`)
	if err != nil {
		t.Fatalf("jsonfmt failed: %v", err)
	}

	if want := "{\n  \"a\": 1\n}\n\n"; stdout.String() != want {
		t.Fatalf("expected %q, got %q", want, stdout.String())
	}
}
