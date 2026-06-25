package cmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/skatkov/devtui/internal/numbers"
)

func TestNumbersCmd(t *testing.T) {
	numbersBase = 10
	outputJSON = false

	cmd := GetRootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetIn(strings.NewReader("1010"))
	cmd.SetArgs([]string{"numbers", "--base", "2"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("numbers command failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Base 10 (decimal)") || !strings.Contains(output, "10") {
		t.Fatalf("numbers output missing decimal conversion: %s", output)
	}
}

func TestNumbersCmdJSON(t *testing.T) {
	numbersBase = 10
	outputJSON = false

	cmd := GetRootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"numbers", "--json", "42"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("numbers --json command failed: %v", err)
	}

	var result numbers.Result
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("numbers --json output invalid JSON: %v", err)
	}
	if result.Base != 10 {
		t.Fatalf("expected base 10, got %d", result.Base)
	}
	if len(result.Conversions) == 0 {
		t.Fatalf("expected conversions in JSON output")
	}
}
