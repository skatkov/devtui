package cmd

import (
	"bytes"
	"strings"
	"testing"
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
