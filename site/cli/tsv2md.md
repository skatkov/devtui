---
title: tsv2md
parent: CLI
---

## devtui tsv2md

Convert TSV to Markdown Table

### Synopsis

Convert TSV to Markdown Table

```bash
devtui tsv2md [flags]
```

### Examples

```bash
devtui tsv2md -t < example.tsv          # convert tsv from stdin and view result in stdout
devtui tsv2md < example.tsv > output.md # convert tsv from stdin and write result in new file
cat example.tsv | devtui tsv2md         # convert tsv from stdin and view result in stdout
```

### Options

```
  -a, --align           align columns width
  -t, --header string   add main header (h1) to result
  -h, --help            help for tsv2md
```
