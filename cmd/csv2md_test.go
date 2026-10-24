package cmd

import (
	"encoding/csv"
	"errors"
	"strings"
	"testing"
	"testing/iotest"
)

const csv2mdInput = "a,bb\nccc,d\n"

func TestCsv2mdCmd(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{"default", nil, "| a | bb | \n| --- | --- | \n| ccc | d | \n"},
		{"aligned", []string{"--align"}, "| a   | bb | \n| --- | --- | \n| ccc | d  | \n"},
		{"header", []string{"--header", "Data"}, "# Data\n\n| a | bb | \n| --- | --- | \n| ccc | d | \n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, _, err := runDevtui(strings.NewReader(csv2mdInput), append([]string{"csv2md"}, tt.args...)...)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := stdout.String(); got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

func TestCsv2mdCmdArgumentInput(t *testing.T) {
	stdout, _, err := runDevtui(strings.NewReader(""), "csv2md", "a\n1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := "| a | \n| --- | \n| 1 | \n"; stdout.String() != want {
		t.Fatalf("expected %q, got %q", want, stdout.String())
	}
}

func TestCsv2mdCmdInvalidCSV(t *testing.T) {
	_, _, err := runDevtui(strings.NewReader("a,\"b\n"), "csv2md")
	var parseErr *csv.ParseError
	if !errors.As(err, &parseErr) {
		t.Fatalf("expected csv parse error, got %v", err)
	}
}

func TestCsv2mdCmdReadError(t *testing.T) {
	want := errors.New("read failed")

	_, _, err := runDevtui(iotest.ErrReader(want), "csv2md")
	if !errors.Is(err, want) {
		t.Fatalf("expected %v, got %v", want, err)
	}
}

func TestCsv2mdCmdWriteError(t *testing.T) {
	failOutput(t, csv2mdCmd)

	_, _, err := runDevtui(strings.NewReader(csv2mdInput), "csv2md")
	if err == nil || !strings.Contains(err.Error(), "write failed") {
		t.Fatalf("expected write error, got %v", err)
	}
}
