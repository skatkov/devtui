package cmd

import (
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/fang/v2"
	"charm.land/lipgloss/v2"
	"github.com/skatkov/devtui/tui/root"
)

// captureProgram replaces runProgram for the test and returns the model it was given.
func captureProgram(t *testing.T, err error) *tea.Model {
	t.Helper()

	var model tea.Model
	original := runProgram
	runProgram = func(m tea.Model) error {
		model = m
		return err
	}
	t.Cleanup(func() { runProgram = original })

	return &model
}

type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) {
	return 0, errors.New("write failed")
}

// failOutput makes the subcommand's writes fail.
func failOutput(t *testing.T, cmd interface{ SetOut(io.Writer) }) {
	t.Helper()

	cmd.SetOut(failingWriter{})
	t.Cleanup(func() { cmd.SetOut(nil) })
}

func TestRootCmdRunsRootScreen(t *testing.T) {
	model := captureProgram(t, nil)

	_, _, err := runDevtui(strings.NewReader(""))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := (*model).(root.RootModel); !ok {
		t.Fatalf("expected root.RootModel, got %T", *model)
	}
}

func TestRootCmdReturnsProgramError(t *testing.T) {
	want := errors.New("program failed")
	captureProgram(t, want)

	_, _, err := runDevtui(strings.NewReader(""))
	if !errors.Is(err, want) {
		t.Fatalf("expected %v, got %v", want, err)
	}
}

func testErrorStyles() fang.Styles {
	return fang.Styles{
		ErrorHeader: lipgloss.NewStyle().SetString("ERROR"),
		ErrorText:   lipgloss.NewStyle().Width(3).Transform(strings.ToUpper),
	}
}

func TestCustomErrorHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  string
		want string
	}{
		{"single line", "devtui failed", "ERROR\nDEVTUI FAILED\n\n"},
		{"lines of one paragraph", "first\nsecond", "ERROR\nFIRST\nsecond\n\n"},
		{"paragraphs", "first\n\nsecond\nthird", "ERROR\nFIRST\n\nsecond\nthird\n\n"},
		{"empty lines are skipped", "first\n\n\nsecond", "ERROR\nFIRST\n\nsecond\n\n"},
		{"empty first line keeps the rest untransformed", "\nsecond", "ERROR\nsecond\n\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			customErrorHandler(&buf, testErrorStyles(), errors.New(tt.err))
			if got := buf.String(); got != tt.want {
				t.Fatalf("customErrorHandler(%q) = %q, want %q", tt.err, got, tt.want)
			}
		})
	}
}

func TestCustomErrorHandlerPlainForNonTerminal(t *testing.T) {
	t.Parallel()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	defer r.Close()

	customErrorHandler(w, testErrorStyles(), errors.New("first\n\nsecond"))
	w.Close()

	got, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("read pipe: %v", err)
	}
	if want := "first\n\nsecond\n"; string(got) != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

// executeDevtui runs Execute with args and returns the exit code it asked for, or -1.
func executeDevtui(t *testing.T, args ...string) int {
	t.Helper()

	cmd := GetRootCmd()
	helpFunc := cmd.HelpFunc()
	commands := cmd.Commands()
	t.Cleanup(func() {
		cmd.SetHelpFunc(helpFunc)
		cmd.ResetCommands()
		cmd.AddCommand(commands...)
		cmd.Version = ""
		cmd.SilenceUsage = false
		cmd.SilenceErrors = false
	})

	code := -1
	original := exit
	exit = func(c int) { code = c }
	t.Cleanup(func() { exit = original })

	cmd.SetIn(strings.NewReader(""))
	cmd.SetOut(new(bytes.Buffer))
	cmd.SetErr(new(bytes.Buffer))
	cmd.SetArgs(args)
	Execute()

	return code
}

func TestExecuteExitsOnError(t *testing.T) {
	if code := executeDevtui(t, "iban"); code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
}

func TestExecuteDoesNotExitOnSuccess(t *testing.T) {
	if code := executeDevtui(t, "iban", "GB"); code != -1 {
		t.Fatalf("expected no exit, got exit code %d", code)
	}
}
