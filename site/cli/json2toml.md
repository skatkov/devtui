---
title: json2toml
parent: CLI
---

## devtui json2toml

Convert JSON to TOML format

### Synopsis

Convert JSON (JavaScript Object Notation) to TOML (Tom's Obvious Minimal Language) format.

Input can be a string argument or piped from stdin. JSON numbers are preserved
as integers when appropriate (not converted to floats). Use --tui flag to view
results in an interactive terminal interface.

```bash
devtui json2toml [string or file] [flags]
```

### Examples

```bash
# Convert JSON from stdin
devtui json2toml < config.json
cat app.json | devtui json2toml
# Convert JSON string argument
devtui json2toml '{"name": "myapp", "version": "1.0.0"}'
# Output to file
devtui json2toml < input.json > output.toml
cat config.json | devtui json2toml > config.toml
# Show results in interactive TUI
devtui json2toml --tui < config.json
devtui json2toml -t < config.json
# Chain with other commands
curl -s https://api.example.com/config | devtui json2toml
```

### Options

```
  -h, --help   help for json2toml
  -t, --tui    Show output in TUI
```
