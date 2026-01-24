---
title: yamlfmt
parent: CLI
---

## devtui yamlfmt

Format and prettify YAML

### Synopsis

Format and prettify YAML input with proper indentation.

Input can be a string argument, piped from stdin, or read from a file.

```bash
devtui yamlfmt [string or file] [flags]
```

### Examples

```bash
# Format YAML from stdin
devtui yamlfmt < config.yaml
cat config.yaml | devtui yamlfmt
# Format YAML string argument
devtui yamlfmt 'name: myapp\nversion: 1.0.0'
# Output to file
devtui yamlfmt < input.yaml > formatted.yaml
```

### Options

```
  -h, --help   help for yamlfmt
```
