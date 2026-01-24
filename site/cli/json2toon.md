---
title: json2toon
parent: CLI
---

## devtui json2toon

Convert JSON to TOON

### Synopsis

Convert JSON to TOON - a compact, human-readable
format designed for passing structured data to Large Language Models with significantly
reduced token usage (typically 30-60% fewer tokens than JSON).

```bash
devtui json2toon [flags]
```

### Examples

```bash
devtui json2toon < example.json                    # Convert with defaults
devtui json2toon -i 4 < example.json               # Use 4-space indent
devtui json2toon -l '#' < example.json             # Add length marker prefix
cat example.json | devtui json2toon > output.toon  # Pipe and save to file
```

### Options

```
  -h, --help                   help for json2toon
  -i, --indent int             Number of spaces per indentation level (default 2)
  -l, --length-marker string   Optional marker to prefix array lengths (e.g., '#')
```
