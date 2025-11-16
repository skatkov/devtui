---
title: base64
parent: CLI
---

## devtui base64

Encode or decode base64 strings and files

### Synopsis

Encode or decode base64 strings and files.

By default, input is encoded to base64. Use the --decode flag to decode base64 input.
Input can be a string argument or piped from stdin.

```bash
devtui base64 [string or file] [flags]
```

### Examples

```bash
# Encode a string
devtui base64 "hello world"
# Decode a base64 string
devtui base64 "aGVsbG8gd29ybGQ=" --decode
devtui base64 "aGVsbG8gd29ybGQ=" -d
# Output to file
devtui base64 "hello world" > encoded.txt
devtui base64 "aGVsbG8gd29ybGQ=" --decode > decoded.txt
# Pipe input from other commands
echo -n "hello world" | devtui base64
echo -n "aGVsbG8gd29ybGQ=" | devtui base64 --decode
cat file.txt | devtui base64
# Chain with other commands
cat file.txt | devtui base64 | devtui base64 --decode
```

### Options

```
  -d, --decode   decode base64 input instead of encoding
  -h, --help     help for base64
```
