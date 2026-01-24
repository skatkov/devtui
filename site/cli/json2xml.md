---
title: json2xml
parent: CLI
---

## devtui json2xml

Convert JSON to XML format

### Synopsis

Convert JSON to XML format.

Input can be a string argument or piped from stdin.

```bash
devtui json2xml [string or file] [flags]
```

### Examples

```bash
# Convert JSON from stdin
devtui json2xml < data.json
cat feed.json | devtui json2xml
# Convert JSON string argument
devtui json2xml '{"item": "value"}'
# Output to file
devtui json2xml < input.json > output.xml
# Chain with other commands
curl -s https://api.example.com/data.json | devtui json2xml
```

### Options

```
  -h, --help   help for json2xml
```
