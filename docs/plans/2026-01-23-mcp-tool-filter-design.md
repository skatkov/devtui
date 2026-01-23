Title: MCP tool filtering for devtui mcp
Date: 2026-01-23

## Context

The MCP server currently exposes every Cobra leaf command as a tool. Some commands should not be callable via MCP, including:

- devtui.completion
- devtui.version
- devtui.mcp
- devtui.serve

## Goals

- Remove the listed tools from `tools/list`.
- Prevent `tools/call` from invoking the removed tools.
- Keep the rest of the MCP tooling behavior unchanged.

## Non-goals

- Switching to an allowlist for all tools.
- Changing Cobra command definitions or CLI behavior.
- Adding new MCP features.

## Proposed Design

### Filtering Strategy

Filter tools at the MCP layer:

- `BuildTools` filters out any tool whose name matches a denylist.
- `ExecuteTool` checks the same denylist and returns an error if a request attempts a blocked tool.

This ensures the tool list and execution path stay consistent even if a client calls a blocked tool directly.

### Matching Rules

- Exact match: `devtui.completion`, `devtui.version`, `devtui.mcp`, `devtui.serve`

### Error Handling

When a blocked tool is invoked via `tools/call`, return an error so the MCP server returns a JSON-RPC error for the call. This should surface as a `method not found` or `invalid params` response depending on existing error handling.

## Testing Plan

- Extend MCP tool list tests to assert blocked tools are not included.
- Add execution tests to ensure blocked tools return errors when called.
