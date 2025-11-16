---
title: urls
parent: CLI
---

## devtui urls

Extract URLs from text, files, or stdin

### Synopsis

Extract URLs from text, files, or stdin.

By default, uses relaxed mode which finds URLs without requiring a scheme.
Use the --strict flag to only find URLs with valid schemes (http, https, ftp, etc.).
Input can be a string argument or piped from stdin.

```bash
devtui urls [string or file] [flags]
```

### Examples

```bash
# Extract URLs from a string
devtui urls "Visit https://google.com and http://example.com"
# Extract URLs in strict mode (requires valid schemes)
devtui urls "Visit google.com and https://example.com" --strict
# Extract URLs from stdin
cat file.html | devtui urls
echo "Check out google.com" | devtui urls
# Chain with other commands
curl -s https://example.com | devtui urls
cat file.txt | devtui urls > extracted_urls.txt
```

### Options

```
  -h, --help     help for urls
  -s, --strict   use strict mode (require valid URL schemes)
```
