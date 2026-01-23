Title: MCP initialize handler for devtui mcp
Date: 2026-01-23

## Context

`devtui mcp` currently implements only `tools/list` and `tools/call`. MCP clients such as mcpjam initiate a handshake with `initialize`. Without this method, the server returns `method not found` and the connection fails.

## Goals

- Accept the MCP `initialize` request.
- Respond with protocol version, capabilities, and server info so clients can complete the handshake.
- Keep the implementation minimal, focused on handshake compatibility.

## Non-goals

- Implementing additional MCP methods beyond `tools/list` and `tools/call`.
- Enforcing strict validation of client protocol versions.

## Proposed Design

### Server Configuration

Extend MCP server configuration to include defaults for initialization response:

- `ProtocolVersion` (default `2025-11-25`)
- `Capabilities` (default empty map)
- `ServerInfo` (default `name=devtui`, `version=dev`)

Add a small `ServerInfo` type and store these fields in the MCP `Server`.

### Initialize Handler

Handle `initialize` in `Server.HandleRequest`:

- Return JSON-RPC `Result` object with:
  - `protocolVersion`
  - `capabilities`
  - `serverInfo` (name + version)
- Accept parameters but ignore them for now.
- Keep error handling consistent with existing methods.

### Version Source

Expose a `GetVersion()` accessor in `cmd/version.go` and pass that value into `cmd/mcp.go` when constructing the MCP server config.

## Data Flow

1. MCP client sends `initialize` over stdio.
2. `ServeStdio` scans the line and dispatches it to `Server.HandleRequest`.
3. The server returns a JSON-RPC response with handshake metadata.
4. Client proceeds to `tools/list` and `tools/call` as before.

## Error Handling

- Return `method not found` for unsupported methods.
- Keep parameter validation minimal to avoid false negatives.

## Testing Plan

- Add `internal/mcp/server_test.go` coverage for `initialize` response shape.
- Update `cmd/mcp_test.go` with an initialize request to ensure CLI wiring works.

## Compatibility Notes

- This is backward compatible for existing clients using only `tools/list` and `tools/call`.
- Future protocol versions can be added by expanding server config or validating client requests.
