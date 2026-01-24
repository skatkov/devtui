package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestYAMLStructCmd(t *testing.T) {
	input := "name: Alice\nage: 30\n"

	cmd := GetRootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetIn(strings.NewReader(input))
	cmd.SetArgs([]string{"yamlstruct"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("yamlstruct command failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "type Root struct") {
		t.Fatalf("yamlstruct output missing struct definition: %s", output)
	}
	if !strings.Contains(output, "Name") {
		t.Fatalf("yamlstruct output missing Name field: %s", output)
	}
}

func TestYAMLStructCmdNoInput(t *testing.T) {
	cmd := GetRootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetIn(strings.NewReader(""))
	cmd.SetArgs([]string{"yamlstruct"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("yamlstruct command should return error when no input provided")
	}
}
