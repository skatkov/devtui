---
title: uuiddecode
parent: CLI
---

## devtui uuiddecode

Decode a UUID into its components

### Synopsis

Decode a UUID and show its components, including version and variant.

Input can be provided as an argument or piped from stdin.

```bash
devtui uuiddecode [uuid] [flags]
```

### Examples

```bash
# Decode a UUID argument
devtui uuiddecode 4326ff5f-774d-4506-a18c-4bc50c761863
# Decode a UUID from stdin
echo "4326ff5f-774d-4506-a18c-4bc50c761863" | devtui uuiddecode
# Output as JSON
devtui uuiddecode --json 4326ff5f-774d-4506-a18c-4bc50c761863
```

### Options

```
  -h, --help   help for uuiddecode
      --json   output decoded fields as JSON
```
