---
title: jsonrepair
parent: CLI
---

## devtui jsonrepair

Repair malformed JSON

### Synopsis

Repair malformed JSON, particularly useful for fixing JSON output from LLMs.

This tool can fix various JSON issues including:
- Single quotes instead of double quotes
- Unclosed arrays and objects
- Mixed quotes
- Uppercase TRUE/FALSE/Null values
- Trailing commas
- JSON wrapped in markdown code blocks
- And many more LLM-specific JSON issues

```bash
devtui jsonrepair [flags]
```

### Examples

```bash
# Repair JSON from stdin
echo "{'key': 'value'}" | devtui jsonrepair
# Repair JSON from file
devtui jsonrepair < broken.json
# Output to file
devtui jsonrepair < broken.json > fixed.json
# Chain with other commands
cat llm-output.txt | devtui jsonrepair | devtui jsonfmt
```

### Options

```
  -h, --help   help for jsonrepair
```
