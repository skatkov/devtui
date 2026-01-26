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
	Tools           []ToolSchema
	Call            func(name string, args CallParams) (string, error)
	ProtocolVersion string
	Capabilities    map[string]any
	ServerInfo      ServerInfo
}

type Server struct {
	tools           []ToolSchema
	call            func(name string, args CallParams) (string, error)
	protocolVersion string
	capabilities    map[string]any
	serverInfo      ServerInfo
}

const ProtocolVersion = "2025-11-25"

func NewServer(cfg ServerConfig) *Server {
	protocolVersion := cfg.ProtocolVersion
	if protocolVersion == "" {
		protocolVersion = ProtocolVersion
	}
	capabilities := cfg.Capabilities
	if capabilities == nil {
		capabilities = map[string]any{}
	}
	serverInfo := cfg.ServerInfo
	if serverInfo.Name == "" {
		serverInfo.Name = "devtui"
	}
	if serverInfo.Version == "" {
		serverInfo.Version = "dev"
	}
	return &Server{
		tools:           cfg.Tools,
		call:            cfg.Call,
		protocolVersion: protocolVersion,
		capabilities:    capabilities,
		serverInfo:      serverInfo,
	}
}

func (s *Server) HandleRequest(req Request) Response {
	switch req.Method {
	case "tools/list":
		return Response{JSONRPC: "2.0", ID: req.ID, Result: map[string]any{"tools": s.tools}}
	case "initialize":
		return Response{JSONRPC: "2.0", ID: req.ID, Result: map[string]any{
			"protocolVersion": s.protocolVersion,
			"capabilities":    s.capabilities,
			"serverInfo": map[string]any{
				"name":    s.serverInfo.Name,
				"version": s.serverInfo.Version,
			},
		}}
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
