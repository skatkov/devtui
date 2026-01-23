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

	params := CallParams{Name: "devtui.base64", Input: "hello"}
	data, _ := json.Marshal(params)
	resp := server.HandleRequest(Request{ID: 2, Method: "tools/call", Params: data})

	if resp.Error != nil {
		t.Fatalf("expected no error")
	}
}
