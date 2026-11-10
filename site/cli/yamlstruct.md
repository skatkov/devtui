---
title: yamlstruct
parent: CLI
description: "devtui yamlstruct: convert YAML to Go struct from the command line. Usage, flags and examples for piping and scripts."
---

# devtui yamlstruct

Convert YAML to Go struct

### Synopsis

Convert YAML input into a Go struct definition.

Input can be a string argument or piped from stdin.

```bash
devtui yamlstruct [string or file] [flags]
```

### Examples

```bash
# Convert YAML from stdin
devtui yamlstruct < config.yaml
cat config.yaml | devtui yamlstruct
# Convert YAML string argument
devtui yamlstruct 'name: Alice\nage: 30'
# Output to file
devtui yamlstruct < input.yaml > struct.go
```

### Options

```
  -h, --help   help for yamlstruct
```
