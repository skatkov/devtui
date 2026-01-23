package mcp

import "encoding/json"

type Request struct {
	JSONRPC string          `json:"jsonrpc,omitempty"`
	ID      any             `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type Response struct {
	JSONRPC string       `json:"jsonrpc"`
	ID      any          `json:"id,omitempty"`
	Result  any          `json:"result,omitempty"`
	Error   *ErrorObject `json:"error,omitempty"`
}

type ErrorObject struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type CallParams struct {
	Name  string            `json:"name"`
	Args  []string          `json:"args,omitempty"`
	Input string            `json:"input,omitempty"`
	Flags map[string]string `json:"flags,omitempty"`
}

type ServerConfig struct {
	Tools []ToolSchema
	Call  func(name string, args CallParams) (string, error)
}

type Server struct {
	tools []ToolSchema
	call  func(name string, args CallParams) (string, error)
}

func NewServer(cfg ServerConfig) *Server {
	return &Server{tools: cfg.Tools, call: cfg.Call}
}

func (s *Server) HandleRequest(req Request) Response {
	switch req.Method {
	case "tools/list":
		return Response{JSONRPC: "2.0", ID: req.ID, Result: map[string]any{"tools": s.tools}}
	case "tools/call":
		var params CallParams
		if err := json.Unmarshal(req.Params, &params); err != nil {
			return Response{JSONRPC: "2.0", ID: req.ID, Error: &ErrorObject{Code: -32602, Message: "invalid params"}}
		}
		if s.call == nil {
			return Response{JSONRPC: "2.0", ID: req.ID, Error: &ErrorObject{Code: -32603, Message: "call not configured"}}
		}
		result, err := s.call(params.Name, params)
		if err != nil {
			return Response{JSONRPC: "2.0", ID: req.ID, Error: &ErrorObject{Code: -32000, Message: err.Error()}}
		}
		return Response{JSONRPC: "2.0", ID: req.ID, Result: map[string]any{"content": []map[string]string{{"type": "text", "text": result}}}}
	default:
		return Response{JSONRPC: "2.0", ID: req.ID, Error: &ErrorObject{Code: -32601, Message: "method not found"}}
	}
}
