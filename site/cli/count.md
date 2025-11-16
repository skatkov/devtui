---
title: count
parent: CLI
---

## devtui count

Character, spaces and word counter

### Synopsis

Count characters, spaces and words in a string

```bash
devtui count [flags]
```

### Examples

```bash
# Count text from a string
devtui count "test me please"
# Count text from stdin
cat testdata/example.csv | devtui count
echo "hello world" | devtui count
```

### Options

```
  -h, --help   help for count
```
