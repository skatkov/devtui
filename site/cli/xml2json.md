---
title: xml2json
parent: CLI
---

## devtui xml2json

Convert XML to JSON format

### Synopsis

Convert XML (Extensible Markup Language) to JSON (JavaScript Object Notation) format.

Input can be a string argument or piped from stdin.

```bash
devtui xml2json [string or file] [flags]
```

### Examples

```bash
# Convert XML from stdin
devtui xml2json < data.xml
cat feed.xml | devtui xml2json
# Convert XML string argument
devtui xml2json '<root><item>value</item></root>'
# Output to file
devtui xml2json < input.xml > output.json
# Chain with other commands
curl -s https://api.example.com/data.xml | devtui xml2json
```

### Options

```
  -h, --help   help for xml2json
```
