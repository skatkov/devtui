package mcp_test

import (
	"strings"
	"testing"

	"github.com/skatkov/devtui/cmd"
	"github.com/skatkov/devtui/internal/mcp"
	"github.com/spf13/cobra"
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

func TestExecuteToolBlocksFilteredTools(t *testing.T) {
	root := &cobra.Command{Use: "devtui"}
	root.AddCommand(&cobra.Command{Use: "completion", Short: "Completion", Run: func(*cobra.Command, []string) {}})
	root.AddCommand(&cobra.Command{Use: "version", Short: "Version", Run: func(*cobra.Command, []string) {}})
	blocked := []string{"devtui.completion.fish", "devtui.version"}
	for _, name := range blocked {
		_, err := mcp.ExecuteTool(root, mcp.CallParams{Name: name})
		if err == nil {
			t.Fatalf("expected error for blocked tool %s", name)
		}
	}
}
