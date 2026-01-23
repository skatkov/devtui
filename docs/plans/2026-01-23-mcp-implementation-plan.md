# DevTUI MCP Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add a `devtui mcp` subcommand that runs an MCP stdio server exposing all Cobra CLI commands as MCP tools.

**Architecture:** Implement a small MCP server in `internal/mcp` that discovers tools from Cobra, handles `tools/list` and `tools/call`, and runs commands via fresh `GetRootCmd()` instances. Add a stdio transport loop that reads newline-delimited JSON-RPC messages and writes responses to stdout only.

**Tech Stack:** Go 1.25, Cobra, JSON-RPC 2.0, standard library `encoding/json`.

### Task 1: Define MCP protocol types and tool schemas

**Files:**
- Create: `internal/mcp/types.go`
- Test: `internal/mcp/types_test.go`

**Step 1: Write the failing test**

```go
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
                "input": {Type: "string"},
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
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/mcp -run TestToolSchemaJSON -v`
Expected: FAIL with “undefined: ToolSchema/JSONSchema”.

**Step 3: Write minimal implementation**

```go
package mcp

type JSONSchema struct {
    Type        string                 `json:"type,omitempty"`
    Description string                 `json:"description,omitempty"`
    Default     any                    `json:"default,omitempty"`
    Enum        []string               `json:"enum,omitempty"`
    Properties  map[string]JSONSchema  `json:"properties,omitempty"`
    Required    []string               `json:"required,omitempty"`
}

type ToolSchema struct {
    Name        string     `json:"name"`
    Description string     `json:"description,omitempty"`
    InputSchema JSONSchema `json:"inputSchema"`
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/mcp -run TestToolSchemaJSON -v`
Expected: PASS.

**Step 5: Commit**

```bash
git add internal/mcp/types.go internal/mcp/types_test.go
git commit -m "feat: add MCP schema types"
```

### Task 2: Implement tool discovery from Cobra

**Files:**
- Create: `internal/mcp/tools.go`
- Test: `internal/mcp/tools_test.go`

**Step 1: Write the failing test**

```go
package mcp

import (
    "testing"

    "github.com/skatkov/devtui/cmd"
)

func TestBuildToolsFromCobra(t *testing.T) {
    root := cmd.GetRootCmd()
    tools := BuildTools(root)

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
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/mcp -run TestBuildToolsFromCobra -v`
Expected: FAIL with “undefined: BuildTools”.

**Step 3: Write minimal implementation**

```go
package mcp

import (
    "strings"

    "github.com/spf13/cobra"
)

func BuildTools(root *cobra.Command) []ToolSchema {
    var tools []ToolSchema
    var walk func(cmd *cobra.Command, path []string)
    walk = func(cmd *cobra.Command, path []string) {
        if cmd.IsAvailableCommand() && !cmd.HasSubCommands() {
            name := "devtui." + strings.Join(path, ".")
            tools = append(tools, ToolSchema{
                Name:        name,
                Description: cmd.Short,
                InputSchema: buildSchema(cmd),
            })
        }

        for _, child := range cmd.Commands() {
            if child.IsAvailableCommand() {
                walk(child, append(path, child.Name()))
            }
        }
    }

    walk(root, []string{root.Name()})
    return tools
}

func buildSchema(cmd *cobra.Command) JSONSchema {
    schema := JSONSchema{
        Type:       "object",
        Properties: map[string]JSONSchema{},
    }
    schema.Properties["input"] = JSONSchema{Type: "string"}
    schema.Properties["args"] = JSONSchema{Type: "array"}

    cmd.Flags().VisitAll(func(flag *cobra.Flag) {
        schema.Properties[flag.Name] = JSONSchema{
            Type:        flagType(flag.Value.Type()),
            Description: flag.Usage,
            Default:     flag.DefValue,
        }
    })

    return schema
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/mcp -run TestBuildToolsFromCobra -v`
Expected: PASS.

**Step 5: Commit**

```bash
git add internal/mcp/tools.go internal/mcp/tools_test.go
git commit -m "feat: derive MCP tools from cobra"
```

### Task 3: Implement JSON-RPC request/response handling

**Files:**
- Create: `internal/mcp/server.go`
- Test: `internal/mcp/server_test.go`

**Step 1: Write the failing test**

```go
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
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/mcp -run TestHandleToolsList -v`
Expected: FAIL with “undefined: NewServer/Request”.

**Step 3: Write minimal implementation**

```go
package mcp

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
    default:
        return Response{JSONRPC: "2.0", ID: req.ID, Error: &ErrorObject{Code: -32601, Message: "method not found"}}
    }
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/mcp -run TestHandleToolsList -v`
Expected: PASS.

**Step 5: Commit**

```bash
git add internal/mcp/server.go internal/mcp/server_test.go
git commit -m "feat: add MCP server request handling"
```

### Task 4: Implement tool execution for tools/call

**Files:**
- Modify: `internal/mcp/server.go`
- Create: `internal/mcp/executor.go`
- Test: `internal/mcp/server_test.go`

**Step 1: Write the failing test**

```go
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

    params := CallParams{ Name: "devtui.base64", Input: "hello" }
    data, _ := json.Marshal(params)
    resp := server.HandleRequest(Request{ID: 2, Method: "tools/call", Params: data})

    if resp.Error != nil {
        t.Fatalf("expected no error")
    }
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/mcp -run TestHandleToolsCall -v`
Expected: FAIL with “undefined: CallParams”.

**Step 3: Write minimal implementation**

```go
package mcp

type CallParams struct {
    Name   string            `json:"name"`
    Args   []string          `json:"args,omitempty"`
    Input  string            `json:"input,omitempty"`
    Flags  map[string]string `json:"flags,omitempty"`
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
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/mcp -run TestHandleToolsCall -v`
Expected: PASS.

**Step 5: Commit**

```bash
git add internal/mcp/server.go internal/mcp/executor.go internal/mcp/server_test.go
git commit -m "feat: handle MCP tools calls"
```

### Task 5: Execute Cobra commands for MCP tool calls

**Files:**
- Modify: `internal/mcp/executor.go`
- Test: `internal/mcp/executor_test.go`

**Step 1: Write the failing test**

```go
package mcp

import (
    "strings"
    "testing"

    "github.com/skatkov/devtui/cmd"
)

func TestExecuteTool(t *testing.T) {
    root := cmd.GetRootCmd()
    out, err := ExecuteTool(root, CallParams{
        Name:  "devtui.base64",
        Input: "hello",
    })
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if strings.TrimSpace(out) != "aGVsbG8gd29ybGQ=" && strings.TrimSpace(out) != "aGVsbG8=" {
        t.Fatalf("unexpected output: %q", out)
    }
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/mcp -run TestExecuteTool -v`
Expected: FAIL with “undefined: ExecuteTool”.

**Step 3: Write minimal implementation**

```go
package mcp

import (
    "bytes"
    "strings"

    "github.com/spf13/cobra"
)

func ExecuteTool(root *cobra.Command, params CallParams) (string, error) {
    cmd, _, err := root.Find(strings.Split(strings.TrimPrefix(params.Name, "devtui."), "."))
    if err != nil {
        return "", err
    }

    bufOut := &bytes.Buffer{}
    bufErr := &bytes.Buffer{}
    cmd.SetOut(bufOut)
    cmd.SetErr(bufErr)
    if params.Input != "" {
        cmd.SetIn(strings.NewReader(params.Input))
    }
    if len(params.Args) > 0 {
        cmd.SetArgs(params.Args)
    }

    if err := cmd.Execute(); err != nil {
        if bufErr.Len() > 0 {
            return "", fmt.Errorf("%s", strings.TrimSpace(bufErr.String()))
        }
        return "", err
    }

    return bufOut.String(), nil
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/mcp -run TestExecuteTool -v`
Expected: PASS.

**Step 5: Commit**

```bash
git add internal/mcp/executor.go internal/mcp/executor_test.go
git commit -m "feat: execute cobra commands for MCP tools"
```

### Task 6: Add stdio transport loop

**Files:**
- Create: `internal/mcp/stdio.go`
- Test: `internal/mcp/stdio_test.go`

**Step 1: Write the failing test**

```go
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
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/mcp -run TestStdioServerHandlesLine -v`
Expected: FAIL with “undefined: ServeStdio”.

**Step 3: Write minimal implementation**

```go
package mcp

import (
    "bufio"
    "encoding/json"
    "io"
)

func ServeStdio(server *Server, in io.Reader, out io.Writer) error {
    scanner := bufio.NewScanner(in)
    for scanner.Scan() {
        var req Request
        if err := json.Unmarshal(scanner.Bytes(), &req); err != nil {
            continue
        }
        resp := server.HandleRequest(req)
        data, err := json.Marshal(resp)
        if err != nil {
            continue
        }
        if _, err := out.Write(append(data, '\n')); err != nil {
            return err
        }
    }
    return scanner.Err()
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/mcp -run TestStdioServerHandlesLine -v`
Expected: PASS.

**Step 5: Commit**

```bash
git add internal/mcp/stdio.go internal/mcp/stdio_test.go
git commit -m "feat: add MCP stdio transport"
```

### Task 7: Add mcp command wiring

**Files:**
- Create: `cmd/mcp.go`
- Modify: `main.go` (only if needed for registration)
- Test: `cmd/mcp_test.go`

**Step 1: Write the failing test**

```go
package cmd

import (
    "bytes"
    "testing"
)

func TestMCPCommandListsTools(t *testing.T) {
    cmd := GetRootCmd()
    buf := new(bytes.Buffer)
    cmd.SetOut(buf)
    cmd.SetErr(buf)
    cmd.SetIn(bytes.NewBufferString(`{"id":1,"method":"tools/list"}\n`))
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

Run: `go test ./cmd -run TestMCPCommandListsTools -v`
Expected: FAIL with “unknown command mcp”.

**Step 3: Write minimal implementation**

```go
package cmd

import (
    "github.com/skatkov/devtui/internal/mcp"
    "github.com/spf13/cobra"
)

var mcpCmd = &cobra.Command{
    Use:   "mcp",
    Short: "Run DevTUI as an MCP stdio server",
    RunE: func(cmd *cobra.Command, args []string) error {
        tools := mcp.BuildTools(GetRootCmd())
        server := mcp.NewServer(mcp.ServerConfig{
            Tools: tools,
            Call: func(name string, params mcp.CallParams) (string, error) {
                root := GetRootCmd()
                return mcp.ExecuteTool(root, params)
            },
        })
        return mcp.ServeStdio(server, cmd.InOrStdin(), cmd.OutOrStdout())
    },
}

func init() {
    rootCmd.AddCommand(mcpCmd)
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./cmd -run TestMCPCommandListsTools -v`
Expected: PASS.

**Step 5: Commit**

```bash
git add cmd/mcp.go cmd/mcp_test.go
git commit -m "feat: add mcp subcommand"
```

### Task 8: Update documentation

**Files:**
- Modify: `README.md`

**Step 1: Write the failing test**

Skip (docs-only change).

**Step 2: Update README**

Add a short section describing MCP usage:

```md
## MCP

Run DevTUI as an MCP server over stdio:

```bash
devtui mcp
```
```

**Step 3: Commit**

```bash
git add README.md
git commit -m "docs: add MCP usage"
```

### Task 9: Full verification

**Files:** none

**Step 1: Run full test suite**

Run: `go test ./...`
Expected: PASS.

**Step 2: Report**

Summarize changes and provide next steps.
