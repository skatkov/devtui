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
