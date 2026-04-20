package cmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestUUIDDecodeCmd(t *testing.T) {
	outputJSON = false
	input := "4326ff5f-774d-4506-a18c-4bc50c761863"

	cmd := GetRootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetIn(strings.NewReader(input))
	cmd.SetArgs([]string{"uuiddecode"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("uuiddecode command failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Standard String Format") {
		t.Fatalf("uuiddecode output missing expected label: %s", output)
	}
	if !strings.Contains(output, input) {
		t.Fatalf("uuiddecode output missing uuid value: %s", output)
	}
}

func TestUUIDDecodeCmdJSON(t *testing.T) {
	outputJSON = false
	input := "4326ff5f-774d-4506-a18c-4bc50c761863"

	cmd := GetRootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetIn(strings.NewReader(input))
	cmd.SetArgs([]string{"uuiddecode", "--json"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("uuiddecode --json command failed: %v", err)
	}

	var payload uuidDecodeJSON
	if err := json.Unmarshal(buf.Bytes(), &payload); err != nil {
		t.Fatalf("uuiddecode --json output invalid JSON: %v", err)
	}
	if payload.UUID != input {
		t.Fatalf("expected uuid %s, got %s", input, payload.UUID)
	}
	if len(payload.Fields) == 0 {
		t.Fatalf("expected decoded fields, got empty")
	}
}

func TestUUIDDecodeCmdNoInput(t *testing.T) {
	outputJSON = false
	cmd := GetRootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetIn(strings.NewReader(""))
	cmd.SetArgs([]string{"uuiddecode"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("uuiddecode command should return error when no input provided")
	}
}
