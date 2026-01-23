package mcp_test

import (
	"testing"

	"github.com/skatkov/devtui/cmd"
	"github.com/skatkov/devtui/internal/mcp"
)

func TestBuildToolsFromCobra(t *testing.T) {
	root := cmd.GetRootCmd()
	tools := mcp.BuildTools(root)

	if len(tools) == 0 {
		t.Fatalf("expected tools")
	}

	found := false
	for _, tool := range tools {
		if tool.Name == "devtui.jsonfmt" {
			found = true
			if tool.InputSchema.Type != "object" {
				t.Fatalf("expected object schema")
			}
		}
	}

	if !found {
		t.Fatalf("expected devtui.jsonfmt tool")
	}
}
