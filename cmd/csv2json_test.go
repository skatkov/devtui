package cmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestCSV2JSONCmd(t *testing.T) {
	input := "name,age\nAlice,30\nBob,25\n"

	cmd := GetRootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetIn(strings.NewReader(input))
	cmd.SetArgs([]string{"csv2json"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("csv2json command failed: %v", err)
	}

	var data []map[string]any
	if err := json.Unmarshal(buf.Bytes(), &data); err != nil {
		t.Fatalf("csv2json output should be valid JSON: %v", err)
	}
	if len(data) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(data))
	}
	if data[0]["name"] != "Alice" {
		t.Fatalf("expected first row name Alice, got %v", data[0]["name"])
	}
}

func TestCSV2JSONCmdNoInput(t *testing.T) {
	cmd := GetRootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetIn(strings.NewReader(""))
	cmd.SetArgs([]string{"csv2json"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("csv2json command should return error when no input provided")
	}
}
