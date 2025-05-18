---
title: csv2md
parent: CLI
---

## devtui csv2md

Convert CSV to Markdown Table

### Synopsis

Convert CSV to Markdown Table

```
devtui csv2md [flags]
```

### Examples

```
  devtui csv2md -t < example.tsv          - convert tsv from stdin and view result in stdout
	devtui csv2md < example.tsv > output.md - convert tsv from stdin and write result in new file
	cat example.tsv | devtui csv2md         - convert tsv from stdin and view result in stdout
```

### Options

```
  -a, --align           align columns width
  -t, --header string   add main header (h1) to result
  -h, --help            help for csv2md
```
