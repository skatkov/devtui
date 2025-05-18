---
title: xmlfmt
parent: CLI
---

## devtui xmlfmt

Format XML

### Synopsis

Format XML

```
devtui xmlfmt [flags]
```

### Examples

```

	xmlfmt < testdata/sample.xml   # Format XML from stdin
	xmlfmt < testdata/sample.xml > output.xml # Output formatted XML to file
	xmlfmt < testdata/sample.xml --tui # Open XML formatter in TUI

```

### Options

```
  -h, --help            help for xmlfmt
  -i, --indent string   Indent string for nested elements (default "  ")
  -n, --nested          Nested tags in comments
  -p, --prefix string   Each element begins on a new line and this prefix
  -t, --tui             Show output in TUI
```
