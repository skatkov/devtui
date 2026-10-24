package cmd

import (
	"strings"
	"testing"
)

func TestTsv2mdCmd(t *testing.T) {
	stdout, _, err := runDevtui(strings.NewReader("name\tage\nAlice\t30\n"), "tsv2md")
	if err != nil {
		t.Fatalf("tsv2md failed: %v", err)
	}

	want := "| name | age | \n| --- | --- | \n| Alice | 30 | \n"
	if got := stdout.String(); got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestTsv2mdCmdArgumentInput(t *testing.T) {
	stdout, _, err := runDevtui(nil, "tsv2md", "x\ty")
	if err != nil {
		t.Fatalf("tsv2md failed: %v", err)
	}

	if want := "| x | y | \n| --- | --- | \n"; stdout.String() != want {
		t.Fatalf("expected %q, got %q", want, stdout.String())
	}
}

func TestTsv2mdCmdInvalidInput(t *testing.T) {
	_, _, err := runDevtui(strings.NewReader("a\"b\tc\n"), "tsv2md")
	if err == nil {
		t.Fatal("expected an error for a bare quote")
	}
}
