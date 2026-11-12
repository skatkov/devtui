package cmd

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"testing/iotest"
)

func TestHTMLFmtCmd(t *testing.T) {
	input := "<div><span>hi</span></div>"

	cmd := GetRootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetIn(strings.NewReader(input))
	cmd.SetArgs([]string{"htmlfmt"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("htmlfmt command failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "<div>") || !strings.Contains(output, "<span>") || !strings.Contains(output, "hi") || !strings.Contains(output, "</span>") {
		t.Fatalf("htmlfmt output missing expected tags: %s", output)
	}
}

func TestHTMLFmtCmdNoInput(t *testing.T) {
	cmd := GetRootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetIn(strings.NewReader(""))
	cmd.SetArgs([]string{"htmlfmt"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("htmlfmt command should return error when no input provided")
	}
}

func TestHTMLFmtCmdArgumentInput(t *testing.T) {
	stdout, _, err := runDevtui(strings.NewReader(""), "htmlfmt", "<p>x</p>")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := "<p>\n  x\n</p>\n"; stdout.String() != want {
		t.Fatalf("expected %q, got %q", want, stdout.String())
	}
}

func TestHTMLFmtCmdReadError(t *testing.T) {
	want := errors.New("read failed")

	_, _, err := runDevtui(iotest.ErrReader(want), "htmlfmt")
	if !errors.Is(err, want) {
		t.Fatalf("expected %v, got %v", want, err)
	}
}
