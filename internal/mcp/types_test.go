package mcp

import (
	"encoding/json"
	"testing"
)

func TestToolSchemaJSON(t *testing.T) {
	schema := ToolSchema{
		Name:        "devtui.jsonfmt",
		Description: "Format JSON",
		InputSchema: JSONSchema{
			Type: "object",
			Properties: map[string]JSONSchema{
				"input":  {Type: "string"},
				"indent": {Type: "integer", Default: 2},
			},
		},
	}

	data, err := json.Marshal(schema)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	if len(data) == 0 {
		t.Fatalf("expected JSON output")
	}
}
