package mcp_test

import (
	"testing"

	"github.com/skatkov/devtui/cmd"
	"github.com/skatkov/devtui/internal/mcp"
)

func TestBuildToolsFromCobra(t *testing.T) {
	root := cmd.GetRootCmd()
	tools := mcp.BuildTools(root)

	if len(tools) != 2 {
		t.Fatalf("expected 2 tools, got %d", len(tools))
	}

	seen := map[string]bool{}
	for _, tool := range tools {
		seen[tool.Name] = true
		if tool.InputSchema.Type != "object" {
			t.Fatalf("expected object schema")
		}
	}

	if !seen["devtui.json2toon"] || !seen["devtui.jsonrepair"] {
		t.Fatalf("expected json2toon and jsonrepair tools")
	}
}

func TestBuildToolsUsesAllowList(t *testing.T) {
	root := cmd.GetRootCmd()
	tools := mcp.BuildTools(root)

	allowed := map[string]struct{}{
		"devtui.json2toon":  {},
		"devtui.jsonrepair": {},
	}

	for _, tool := range tools {
		if _, exists := allowed[tool.Name]; !exists {
			t.Fatalf("unexpected tool %s", tool.Name)
		}
	}
}
