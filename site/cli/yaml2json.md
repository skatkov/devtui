---
title: yaml2json
parent: CLI
---

## devtui yaml2json

Convert YAML to JSON format

### Synopsis

Convert YAML to JSON format.

Input can be a string argument or piped from stdin.

```bash
devtui yaml2json [string or file] [flags]
```

### Examples

```bash
# Convert YAML from stdin
devtui yaml2json < config.yaml
cat app.yaml | devtui yaml2json
# Convert YAML string argument
devtui yaml2json 'name: myapp\nversion: 1.0.0'
# Output to file
devtui yaml2json < input.yaml > output.json
# Chain with other commands
curl -s https://api.example.com/config.yaml | devtui yaml2json
```

### Options

```
  -h, --help   help for yaml2json
```
