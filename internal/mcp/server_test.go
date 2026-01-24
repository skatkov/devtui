package mcp

import (
	"encoding/json"
	"testing"
)

func TestHandleToolsList(t *testing.T) {
	server := NewServer(ServerConfig{
		Tools: []ToolSchema{{Name: "devtui.jsonfmt"}},
	})

	req := Request{ID: 1, Method: "tools/list"}
	resp := server.HandleRequest(req)

	data, _ := json.Marshal(resp)
	if !json.Valid(data) {
		t.Fatalf("response not valid json")
	}
	if resp.Error != nil {
		t.Fatalf("expected no error")
	}
}

func TestHandleToolsCall(t *testing.T) {
	server := NewServer(ServerConfig{
		Tools: []ToolSchema{{Name: "devtui.base64"}},
		Call: func(name string, args CallParams) (string, error) {
			if name != "devtui.base64" {
				t.Fatalf("unexpected tool name: %s", name)
			}
			if args.Input != "hello" {
				t.Fatalf("unexpected input")
			}
			return "aGVsbG8=", nil
		},
	})

	toolCallParams := ToolCallParams{
		Name:      "devtui.base64",
		Arguments: ToolCallArguments{Input: "hello"},
	}
	data, _ := json.Marshal(toolCallParams)
	resp := server.HandleRequest(Request{ID: 2, Method: "tools/call", Params: data})

	if resp.Error != nil {
		t.Fatalf("expected no error")
	}
}

func TestHandleInitialize(t *testing.T) {
	server := NewServer(ServerConfig{})
	resp := server.HandleRequest(Request{ID: 3, Method: "initialize"})
	if resp.Error != nil {
		t.Fatalf("expected no error")
	}
	result, ok := resp.Result.(map[string]any)
	if !ok {
		t.Fatalf("expected result map")
	}
	if result["protocolVersion"] == "" {
		t.Fatalf("expected protocolVersion")
	}
	info, ok := result["serverInfo"].(map[string]any)
	if !ok {
		t.Fatalf("expected serverInfo map")
	}
	if info["name"] == "" || info["version"] == "" {
		t.Fatalf("expected serverInfo name and version")
	}
}
