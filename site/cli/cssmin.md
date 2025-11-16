---
title: cssmin
parent: CLI
---

## devtui cssmin

Minify CSS files by removing whitespace and unnecessary characters

### Synopsis

Minify CSS files by removing whitespace, line breaks, and unnecessary characters.

This reduces file size for production use while maintaining CSS functionality.
Input can be a string argument or piped from stdin.

```bash
devtui cssmin [string or file] [flags]
```

### Examples

```bash
# Minify CSS from stdin
devtui cssmin < styles.css
cat source.css | devtui cssmin
# Minify CSS string argument
devtui cssmin 'body { margin: 0; padding: 0; }'
# Output to file
devtui cssmin < input.css > minified.css
cat styles.css | devtui cssmin > styles.min.css
# Chain with other commands
curl -s https://example.com/styles.css | devtui cssmin
devtui cssfmt < messy.css | devtui cssmin
```

### Options

```
  -h, --help   help for cssmin
```
