package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestJSONStructCmd(t *testing.T) {
	input := `{"name":"Alice","age":30}`

	cmd := GetRootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetIn(strings.NewReader(input))
	cmd.SetArgs([]string{"jsonstruct"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("jsonstruct command failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "type Root struct") {
		t.Fatalf("jsonstruct output missing struct definition: %s", output)
	}
	if !strings.Contains(output, "Name") {
		t.Fatalf("jsonstruct output missing Name field: %s", output)
	}
}

func TestJSONStructCmdNoInput(t *testing.T) {
	cmd := GetRootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetIn(strings.NewReader(""))
	cmd.SetArgs([]string{"jsonstruct"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("jsonstruct command should return error when no input provided")
	}
}
