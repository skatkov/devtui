# MCP Initialize Handler Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Implement an MCP `initialize` handler so MCP clients can complete the handshake with `devtui mcp`.

**Architecture:** Extend the MCP server config with protocol/version metadata and respond to `initialize` with a minimal handshake payload. Wire CLI to pass version info and add tests for the new handler.

**Tech Stack:** Go 1.25, Cobra, internal MCP server.

### Task 1: Add MCP initialize response metadata types

**Files:**
- Modify: `internal/mcp/types.go`
- Modify: `internal/mcp/server.go`

**Step 1: Write the failing test**

In `internal/mcp/server_test.go`, add a new test that expects `initialize` to return `protocolVersion` and `serverInfo` with name/version.

```go
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
```

**Step 2: Run test to verify it fails**

Run: `go test -v ./internal/mcp -run TestHandleInitialize`
Expected: FAIL with method not found or missing result fields.

**Step 3: Write minimal implementation**

In `internal/mcp/types.go`, add:

```go
type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}
```

In `internal/mcp/server.go`, extend `ServerConfig` and `Server` with:

```go
const ProtocolVersion = "2025-11-25"

type ServerConfig struct {
	Tools          []ToolSchema
	Call           func(name string, args CallParams) (string, error)
	ProtocolVersion string
	Capabilities    map[string]any
	ServerInfo      ServerInfo
}
```

Default these in `NewServer` if empty:

```go
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
	return &Server{tools: cfg.Tools, call: cfg.Call, protocolVersion: protocolVersion, capabilities: capabilities, serverInfo: serverInfo}
}
```

Add fields to `Server`:

```go
protocolVersion string
capabilities    map[string]any
serverInfo      ServerInfo
```

Handle `initialize`:

```go
case "initialize":
	return Response{JSONRPC: "2.0", ID: req.ID, Result: map[string]any{
		"protocolVersion": s.protocolVersion,
		"capabilities":    s.capabilities,
		"serverInfo": map[string]any{
			"name":    s.serverInfo.Name,
			"version": s.serverInfo.Version,
		},
	}}
```

**Step 4: Run test to verify it passes**

Run: `go test -v ./internal/mcp -run TestHandleInitialize`
Expected: PASS.

**Step 5: Commit**

```bash
git add internal/mcp/types.go internal/mcp/server.go internal/mcp/server_test.go
git commit -m "feat: add MCP initialize response metadata"
```

### Task 2: Wire CLI version into MCP server

**Files:**
- Modify: `cmd/version.go`
- Modify: `cmd/mcp.go`
- Modify: `cmd/mcp_test.go`

**Step 1: Write the failing test**

Add a CLI test that uses `initialize` and expects output:

```go
func TestMCPCommandInitialize(t *testing.T) {
	cmd := GetRootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetIn(bytes.NewBufferString("{\"id\":1,\"method\":\"initialize\"}\n"))
	cmd.SetArgs([]string{"mcp"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute failed: %v", err)
	}
	if buf.Len() == 0 {
		t.Fatalf("expected output")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test -v ./cmd -run TestMCPCommandInitialize`
Expected: FAIL with method not found or empty output.

**Step 3: Write minimal implementation**

Expose raw version string in `cmd/version.go`:

```go
// GetVersion returns the raw version string.
func GetVersion() string {
	return version
}
```

Wire server info in `cmd/mcp.go`:

```go
server := mcp.NewServer(mcp.ServerConfig{
	Tools: tools,
	ServerInfo: mcp.ServerInfo{
		Name:    "devtui",
		Version: GetVersion(),
	},
	Call: func(name string, params mcp.CallParams) (string, error) {
		root := GetRootCmd()
		return mcp.ExecuteTool(root, params)
	},
})
```

**Step 4: Run test to verify it passes**

Run: `go test -v ./cmd -run TestMCPCommandInitialize`
Expected: PASS.

**Step 5: Commit**

```bash
git add cmd/version.go cmd/mcp.go cmd/mcp_test.go
git commit -m "feat: add MCP initialize response to CLI"
```

### Task 3: Run focused verification

**Files:**
- None

**Step 1: Run relevant tests**

Run: `go test -v ./internal/mcp ./cmd`
Expected: PASS.

**Step 2: Commit**

No commit needed if nothing changed.
