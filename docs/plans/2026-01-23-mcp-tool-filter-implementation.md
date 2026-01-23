# MCP Tool Filter Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Hide selected CLI commands from MCP tools and block their execution through `tools/call`.

**Architecture:** Add a simple denylist matcher in the MCP tooling layer. Filter tool schemas in `BuildTools` and reject blocked tool names in `ExecuteTool` to prevent bypasses.

**Tech Stack:** Go 1.25, Cobra, internal MCP tooling.

### Task 1: Filter blocked tools in MCP tool listing

**Files:**
- Modify: `internal/mcp/tools.go`
- Modify: `internal/mcp/tools_test.go`

**Step 1: Write the failing test**

Add a test to assert blocked tools do not appear in `BuildTools` output:

```go
func TestBuildToolsFiltersBlocked(t *testing.T) {
	root := &cobra.Command{Use: "devtui"}
	root.AddCommand(&cobra.Command{Use: "version", Short: "Version", Run: func(*cobra.Command, []string) {}})
	root.AddCommand(&cobra.Command{Use: "completion", Short: "Completion", Run: func(*cobra.Command, []string) {}})
	root.AddCommand(&cobra.Command{Use: "serve", Short: "Serve", Run: func(*cobra.Command, []string) {}})
	root.AddCommand(&cobra.Command{Use: "jsonfmt", Short: "JSON fmt", Run: func(*cobra.Command, []string) {}})

	tools := BuildTools(root)
	for _, tool := range tools {
		if tool.Name == "devtui.completion" || tool.Name == "devtui.version" || tool.Name == "devtui.serve" {
			t.Fatalf("blocked tool present: %s", tool.Name)
		}
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test -v ./internal/mcp -run TestBuildToolsFiltersBlocked`
Expected: FAIL showing blocked tools present.

**Step 3: Write minimal implementation**

In `internal/mcp/tools.go`, add a small matcher and filter in `BuildTools`:

```go
var blockedToolNames = map[string]struct{}{
	"devtui.completion": {},
	"devtui.version": {},
	"devtui.mcp": {},
	"devtui.serve": {},
}

func isBlockedTool(name string) bool {
	if _, ok := blockedToolNames[name]; ok {
		return true
	}
	return false
}
```

Then, when building tools, skip any blocked tool name:

```go
if isBlockedTool(name) {
	return
}
```

**Step 4: Run test to verify it passes**

Run: `go test -v ./internal/mcp -run TestBuildToolsFiltersBlocked`
Expected: PASS.

**Step 5: Commit**

```bash
git add internal/mcp/tools.go internal/mcp/tools_test.go
git commit -m "feat: filter blocked MCP tools"
```

### Task 2: Block execution of filtered tools

**Files:**
- Modify: `internal/mcp/executor.go`
- Modify: `internal/mcp/executor_test.go`

**Step 1: Write the failing test**

Add test cases to ensure `ExecuteTool` rejects blocked tools:

```go
func TestExecuteToolBlocksFilteredTools(t *testing.T) {
	root := &cobra.Command{Use: "devtui"}
	root.AddCommand(&cobra.Command{Use: "version", Short: "Version", Run: func(*cobra.Command, []string) {}})

	_, err := ExecuteTool(root, CallParams{Name: "devtui.version"})
	if err == nil {
		t.Fatalf("expected error for blocked tool")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test -v ./internal/mcp -run TestExecuteToolBlocksFilteredTools`
Expected: FAIL (no error currently).

**Step 3: Write minimal implementation**

In `internal/mcp/executor.go`, before setting args, reject blocked tools:

```go
if isBlockedTool(params.Name) {
	return "", fmt.Errorf("tool not available")
}
```

**Step 4: Run test to verify it passes**

Run: `go test -v ./internal/mcp -run TestExecuteToolBlocksFilteredTools`
Expected: PASS.

**Step 5: Commit**

```bash
git add internal/mcp/executor.go internal/mcp/executor_test.go
git commit -m "feat: block MCP tool execution for filtered tools"
```

### Task 3: Run focused verification

**Files:**
- None

**Step 1: Run relevant tests**

Run: `go test -v ./internal/mcp`
Expected: PASS.

**Step 2: Commit**

No commit needed if nothing changed.
