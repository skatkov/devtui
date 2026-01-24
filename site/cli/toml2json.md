---
title: toml2json
parent: CLI
---

## devtui toml2json

Convert TOML to JSON format

### Synopsis

Convert TOML to JSON format.

Input can be a string argument or piped from stdin. Use --tui flag to view results
in an interactive terminal interface.

```bash
devtui toml2json [string or file] [flags]
```

### Examples

```bash
# Convert TOML from stdin
devtui toml2json < config.toml
cat app.toml | devtui toml2json
# Convert TOML string argument
devtui toml2json '[package]\nname = "myapp"'
# Output to file
devtui toml2json < input.toml > output.json
cat config.toml | devtui toml2json > data.json
# Show results in interactive TUI
devtui toml2json --tui < config.toml
devtui toml2json -t < config.toml
# Chain with other commands
curl -s https://example.com/config.toml | devtui toml2json
```

### Options

```
  -h, --help   help for toml2json
  -t, --tui    Show output in TUI
```
