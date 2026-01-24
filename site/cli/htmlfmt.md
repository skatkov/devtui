---
title: htmlfmt
parent: CLI
---

## devtui htmlfmt

Format and prettify HTML

### Synopsis

Format and prettify HTML input with consistent indentation.

Input can be a string argument, piped from stdin, or read from a file.

```bash
devtui htmlfmt [string or file] [flags]
```

### Examples

```bash
# Format HTML from stdin
devtui htmlfmt < page.html
cat page.html | devtui htmlfmt
# Format HTML string argument
devtui htmlfmt '<div><span>hello</span></div>'
# Output to file
devtui htmlfmt < input.html > formatted.html
```

### Options

```
  -h, --help   help for htmlfmt
```
