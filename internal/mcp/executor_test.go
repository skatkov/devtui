package mcp_test

import (
	"strings"
	"testing"

	"github.com/skatkov/devtui/cmd"
	"github.com/skatkov/devtui/internal/mcp"
)

func TestExecuteTool(t *testing.T) {
	root := cmd.GetRootCmd()
	out, err := mcp.ExecuteTool(root, mcp.CallParams{
		Name:  "devtui.base64",
		Input: "hello",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.TrimSpace(out) != "aGVsbG8=" {
		t.Fatalf("unexpected output: %q", out)
	}
}
