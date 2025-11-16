---
title: jsonfmt
parent: CLI
---

## devtui jsonfmt

Format and prettify JSON

### Synopsis

Format and prettify JSON input with proper indentation and syntax highlighting.

Input can be a string argument, piped from stdin, or read from a file.
The output is always valid, properly indented JSON.

```bash
devtui jsonfmt [string or file] [flags]
```

### Examples

```bash
# Format JSON from stdin
devtui jsonfmt < example.json
echo '{"name":"John","age":30}' | devtui jsonfmt
# Format JSON string argument
devtui jsonfmt '{"name":"John","age":30}'
# Output to file
devtui jsonfmt < input.json > formatted.json
cat compact.json | devtui jsonfmt > pretty.json
# Chain with other commands
curl -s https://api.example.com/data | devtui jsonfmt
devtui jsonrepair < broken.json | devtui jsonfmt
```

### Options

```
  -h, --help   help for jsonfmt
```
