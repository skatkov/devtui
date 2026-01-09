---
title: yaml2toml
parent: CLI
---

## devtui yaml2toml

Convert YAML to TOML format

### Synopsis

Convert YAML (YAML Ain't Markup Language) to TOML (Tom's Obvious Minimal Language) format.

Input can be a string argument or piped from stdin.

```bash
devtui yaml2toml [string or file] [flags]
```

### Examples

```bash
# Convert YAML from stdin
devtui yaml2toml < config.yaml
cat app.yaml | devtui yaml2toml
# Convert YAML string argument
devtui yaml2toml 'name: myapp\nversion: 1.0.0'
# Output to file
devtui yaml2toml < input.yaml > output.toml
# Chain with other commands
curl -s https://api.example.com/config.yaml | devtui yaml2toml
```

### Options

```
  -h, --help   help for yaml2toml
```
