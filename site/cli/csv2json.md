---
title: csv2json
parent: CLI
---

## devtui csv2json

Convert CSV to JSON

### Synopsis

Convert CSV into formatted JSON.

Input can be a string argument or piped from stdin.

```bash
devtui csv2json [string or file] [flags]
```

### Examples

```bash
# Convert CSV from stdin
devtui csv2json < data.csv
cat data.csv | devtui csv2json
# Convert CSV string argument
devtui csv2json 'name,age\nAlice,30'
# Output to file
devtui csv2json < input.csv > output.json
```

### Options

```
  -h, --help   help for csv2json
```
