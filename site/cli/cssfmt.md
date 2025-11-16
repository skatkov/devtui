---
title: cssfmt
parent: CLI
---

## devtui cssfmt

Format and prettify CSS files

### Synopsis

Format and prettify CSS files with customizable indentation and formatting options.

By default, uses 2-space indentation. Use --tab for tab indentation or --indent to specify
custom spacing. Input can be a string argument or piped from stdin.

```bash
devtui cssfmt [string or file] [flags]
```

### Examples

```bash
# Format CSS from stdin
devtui cssfmt < styles.css
cat minified.css | devtui cssfmt
# Format CSS string argument
devtui cssfmt 'body{margin:0;padding:0}'
# Use tab indentation
devtui cssfmt --tab < styles.css
devtui cssfmt -t < styles.css
# Use custom indent spacing
devtui cssfmt --indent 4 < styles.css
devtui cssfmt -i 4 < styles.css
# Output to file
devtui cssfmt < input.css > formatted.css
# Show results in interactive TUI
devtui cssfmt --tui < styles.css
# Chain with other commands
curl -s https://example.com/styles.css | devtui cssfmt
```

### Options

```
  -h, --help         help for cssfmt
  -i, --indent int   spaces for indentation (default 2)
      --semicolon    always end rule with semicolon, even if not needed (default true)
  -t, --tab          use tabs for indentation
      --tui          present result in a TUI
```
