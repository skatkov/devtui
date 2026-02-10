//go:build !mcp

package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestMCPDisabledMessage(t *testing.T) {
	cmd := GetRootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"mcp"})

	err := cmd.Execute()
	if err == nil {
		t.Fatalf("expected error when mcp tag disabled")
	}
	if !strings.Contains(buf.String(), "mcp disabled") {
		t.Fatalf("expected disabled message, got: %s", buf.String())
	}
}
