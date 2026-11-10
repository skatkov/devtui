---
title: jsonstruct
parent: CLI
description: "devtui jsonstruct: convert JSON to Go struct from the command line. Usage, flags and examples for piping and scripts."
---

# devtui jsonstruct

Convert JSON to Go struct

### Synopsis

Convert JSON input into a Go struct definition.

Input can be a string argument or piped from stdin.

```bash
devtui jsonstruct [string or file] [flags]
```

### Examples

```bash
# Convert JSON from stdin
devtui jsonstruct < data.json
cat data.json | devtui jsonstruct
# Convert JSON string argument
devtui jsonstruct '{"name":"Alice","age":30}'
# Output to file
devtui jsonstruct < input.json > struct.go
```

### Options

```
  -h, --help   help for jsonstruct
```
