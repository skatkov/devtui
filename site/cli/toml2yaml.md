---
title: toml2yaml
parent: CLI
---

## devtui toml2yaml

Convert TOML to YAML format

### Synopsis

Convert TOML (Tom's Obvious Minimal Language) to YAML (YAML Ain't Markup Language) format.

Input can be a string argument or piped from stdin.

```bash
devtui toml2yaml [string or file] [flags]
```

### Examples

```bash
# Convert TOML from stdin
devtui toml2yaml < config.toml
cat app.toml | devtui toml2yaml
# Convert TOML string argument
devtui toml2yaml 'name = "myapp"\nversion = "1.0.0"'
# Output to file
devtui toml2yaml < input.toml > output.yaml
# Chain with other commands
curl -s https://api.example.com/config.toml | devtui toml2yaml
```

### Options

```
  -h, --help   help for toml2yaml
```
