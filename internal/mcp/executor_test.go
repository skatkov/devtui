package mcp

import (
	"strings"
	"testing"

	"github.com/skatkov/devtui/cmd"
)

func TestExecuteTool(t *testing.T) {
	root := cmd.GetRootCmd()
	out, err := ExecuteTool(root, CallParams{
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
