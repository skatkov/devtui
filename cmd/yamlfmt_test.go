package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestYamlfmtCmd(t *testing.T) {
	input := "b: 2\na: 1\n"

	cmd := GetRootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetIn(strings.NewReader(input))
	cmd.SetArgs([]string{"yamlfmt"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("yamlfmt command failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "a: 1") || !strings.Contains(output, "b: 2") {
		t.Fatalf("yamlfmt output missing expected keys: %s", output)
	}
}

func TestYamlfmtCmdNoInput(t *testing.T) {
	cmd := GetRootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetIn(strings.NewReader(""))
	cmd.SetArgs([]string{"yamlfmt"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("yamlfmt command should return error when no input provided")
	}
}
