package mcp

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestStdioServerHandlesLine(t *testing.T) {
	server := NewServer(ServerConfig{Tools: []ToolSchema{{Name: "devtui.jsonfmt"}}})
	input := bytes.NewBufferString(`{"id":1,"method":"tools/list"}` + "\n")
	output := &bytes.Buffer{}

	if err := ServeStdio(server, input, output); err != nil {
		t.Fatalf("serve failed: %v", err)
	}

	lines := bytes.Split(output.Bytes(), []byte("\n"))
	if len(lines[0]) == 0 {
		t.Fatalf("expected output")
	}
	if !json.Valid(lines[0]) {
		t.Fatalf("response not valid json")
	}
}
