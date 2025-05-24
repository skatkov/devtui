---
title: cssfmt
parent: CLI
---

## devtui cssfmt

Format CSS files

### Synopsis

Format CSS files

```
devtui cssfmt [flags]
```

### Examples

```

	cssfmt < testdata/bootstrap.min.css
	cssfmt < testdata/bootstrap.min.css --tui # Show results in a TUI
	
```

### Options

```
  -h, --help         help for cssfmt
  -i, --indent int   spaces for indentation (default 2)
      --semicolon    always end rule with semicolon, even if not needed (default true)
  -t, --tab          use tabs for indentation
      --tui          present result in a TUI
```
