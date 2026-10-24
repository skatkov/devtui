package cmd

import (
	"errors"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/skatkov/devtui/tui/css"
)

const cssfmtInput = "a{color:red;margin:0}"

func TestCssfmtCmd(t *testing.T) {
	captureProgram(t, nil)

	tests := []struct {
		name string
		args []string
		want string
	}{
		{"default two spaces and semicolons", nil, "a {\n  color: red;\n  margin: 0;\n}\n"},
		{"tabs", []string{"--tab"}, "a {\n\tcolor: red;\n\tmargin: 0;\n}\n"},
		{"tabs ignore indent", []string{"--tab", "--indent", "4"}, "a {\n\tcolor: red;\n\tmargin: 0;\n}\n"},
		{"custom indent", []string{"--indent", "4"}, "a {\n    color: red;\n    margin: 0;\n}\n"},
		{"without last semicolon", []string{"--semicolon=false"}, "a {\n  color: red;\n  margin: 0\n}\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, _, err := runDevtui(strings.NewReader(cssfmtInput), append([]string{"cssfmt"}, tt.args...)...)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := stdout.String(); got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

func TestCssfmtCmdArgumentInput(t *testing.T) {
	stdout, _, err := runDevtui(strings.NewReader(""), "cssfmt", "a{b:c}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := "a {\n  b: c;\n}\n"; stdout.String() != want {
		t.Fatalf("expected %q, got %q", want, stdout.String())
	}
}

func TestCssfmtCmdReadError(t *testing.T) {
	want := errors.New("read failed")

	for _, args := range [][]string{{"cssfmt"}, {"cssfmt", "--tui"}} {
		_, _, err := runDevtui(iotest.ErrReader(want), args...)
		if !errors.Is(err, want) {
			t.Fatalf("%v: expected %v, got %v", args, want, err)
		}
	}
}

func TestCssfmtCmdTUI(t *testing.T) {
	model := captureProgram(t, nil)

	_, _, err := runDevtui(strings.NewReader(cssfmtInput), "cssfmt", "--tui")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m, ok := (*model).(css.CSSFormatterModel)
	if !ok {
		t.Fatalf("expected css.CSSFormatterModel, got %T", *model)
	}
	if want := "a {\n  color: red;\n  margin: 0;\n}\n"; m.FormattedContent != want {
		t.Fatalf("expected formatted content %q, got %q", want, m.FormattedContent)
	}
}

func TestCssfmtCmdTUIError(t *testing.T) {
	want := errors.New("program failed")
	captureProgram(t, want)

	_, _, err := runDevtui(strings.NewReader(cssfmtInput), "cssfmt", "--tui")
	if !errors.Is(err, want) {
		t.Fatalf("expected %v, got %v", want, err)
	}
}
