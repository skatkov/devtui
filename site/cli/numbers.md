---
title: numbers
parent: CLI
---

## devtui numbers

Convert numbers between bases

### Synopsis

Convert numbers between binary, octal, decimal, and hexadecimal.

Input can be a number argument or piped from stdin.

```bash
devtui numbers [number] [flags]
```

### Examples

```bash
# Convert a decimal number
devtui numbers 42
# Convert a binary number
devtui numbers --base 2 101010
# Convert from stdin
echo "ff" | devtui numbers --base 16
# Output as JSON
devtui numbers --json 42
```

### Options

```
  -b, --base int   input number base (2, 8, 10, 16) (default 10)
  -h, --help       help for numbers
      --json       output conversions as JSON
```
