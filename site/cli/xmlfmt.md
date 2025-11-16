---
title: xmlfmt
parent: CLI
---

## devtui xmlfmt

Format and prettify XML files

### Synopsis

Format and prettify XML files with customizable indentation and formatting options.

By default, uses 2-space indentation. Customize with --indent, --prefix, and --nested flags.
Input can be a string argument or piped from stdin.

```bash
devtui xmlfmt [string or file] [flags]
```

### Examples

```bash
# Format XML from stdin
devtui xmlfmt < document.xml
cat unformatted.xml | devtui xmlfmt
# Format XML string argument
devtui xmlfmt '<root><item>value</item></root>'
# Output to file
devtui xmlfmt < input.xml > formatted.xml
cat document.xml | devtui xmlfmt > pretty.xml
# Custom indentation
devtui xmlfmt --indent "    " < document.xml
devtui xmlfmt -i "\t" < document.xml
# Add prefix to each line
devtui xmlfmt --prefix "  " < document.xml
devtui xmlfmt -p "  " < document.xml
# Handle nested tags in comments
devtui xmlfmt --nested < document.xml
devtui xmlfmt -n < document.xml
# Show results in interactive TUI
devtui xmlfmt --tui < document.xml
devtui xmlfmt -t < document.xml
# Chain with other commands
curl -s https://example.com/feed.xml | devtui xmlfmt
```

### Options

```
  -h, --help            help for xmlfmt
  -i, --indent string   Indent string for nested elements (default "  ")
  -n, --nested          Nested tags in comments
  -p, --prefix string   Each element begins on a new line and this prefix
  -t, --tui             Show output in TUI
```
