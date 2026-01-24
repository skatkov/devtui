---
title: tomlfmt
parent: CLI
---

## devtui tomlfmt

Format and prettify TOML files

### Synopsis

Format and prettify TOML files with proper indentation.

Input can be a string argument or piped from stdin. Use --tui flag to view results
in an interactive terminal interface.

```bash
devtui tomlfmt [string or file] [flags]
```

### Examples

```bash
# Format TOML from stdin
devtui tomlfmt < config.toml
cat app.toml | devtui tomlfmt
# Format TOML string argument
devtui tomlfmt '[package]\nname = "myapp"'
# Output to file
devtui tomlfmt < input.toml > formatted.toml
cat config.toml | devtui tomlfmt > pretty.toml
# Show results in interactive TUI
devtui tomlfmt --tui < config.toml
devtui tomlfmt -t < config.toml
# Chain with other commands
curl -s https://example.com/config.toml | devtui tomlfmt
```

### Options

```
  -h, --help   help for tomlfmt
  -t, --tui    Show output in TUI
```
