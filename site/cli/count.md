---
title: count
parent: CLI
---

## devtui count

Count characters, spaces, and words in text

### Synopsis

Count characters, spaces, and words in text input.

Provides detailed statistics including character count, space count, and word count
in a formatted table. Input can be a string argument or piped from stdin.

```bash
devtui count [string or file] [flags]
```

### Examples

```bash
# Count text from a string
devtui count "test me please"
devtui count "hello world"
# Count text from stdin
echo "hello world" | devtui count
cat document.txt | devtui count
# Count text from file
devtui count < document.txt
cat README.md | devtui count
# Output to file
devtui count "sample text" > stats.txt
# Chain with other commands
curl -s https://example.com | devtui count
cat article.txt | devtui count
```

### Options

```
  -h, --help   help for count
```
