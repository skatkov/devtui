---
title: json2yaml
parent: CLI
---

## devtui json2yaml

Convert JSON to YAML format

### Synopsis

Convert JSON (JavaScript Object Notation) to YAML (YAML Ain't Markup Language) format.

Input can be a string argument or piped from stdin.

```bash
devtui json2yaml [string or file] [flags]
```

### Examples

```bash
# Convert JSON from stdin
devtui json2yaml < config.json
cat app.json | devtui json2yaml
# Convert JSON string argument
devtui json2yaml '{"name": "myapp", "version": "1.0.0"}'
# Output to file
devtui json2yaml < input.json > output.yaml
# Chain with other commands
curl -s https://api.example.com/config | devtui json2yaml
```

### Options

```
  -h, --help   help for json2yaml
```
