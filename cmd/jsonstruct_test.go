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
	if !strings.Contains(output, "Name string `json:\"name\"`") {
		t.Fatalf("jsonstruct output missing Name field with a json tag: %s", output)
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

func TestJSONStructCmdArgumentInput(t *testing.T) {
	stdout, _, err := runDevtui(nil, "jsonstruct", `{"a":1}`)
	if err != nil {
		t.Fatalf("jsonstruct failed: %v", err)
	}

	if want := "package main\n\ntype Root struct {\n\tA int `json:\"a\"`\n}\n\n"; stdout.String() != want {
		t.Fatalf("expected %q, got %q", want, stdout.String())
	}
}
