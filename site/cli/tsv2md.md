---
title: tsv2md
parent: CLI
---

## devtui tsv2md

Convert TSV to Markdown table format

### Synopsis

Convert TSV (Tab-Separated Values) to Markdown table format for documentation.

Input can be piped from stdin or read from a file. Use --align to align column widths
and --header to add a main heading (h1) to the output.

```bash
devtui tsv2md [string or file] [flags]
```

### Examples

```bash
# Convert TSV from stdin
devtui tsv2md < example.tsv
cat data.tsv | devtui tsv2md
# Output to file
devtui tsv2md < input.tsv > output.md
cat data.tsv | devtui tsv2md > table.md
# Add main header to output
devtui tsv2md --header "User Data" < users.tsv
devtui tsv2md -t "Sales Report" < sales.tsv
# Align column widths for better readability
devtui tsv2md --align < data.tsv
devtui tsv2md -a < data.tsv
# Combine options
devtui tsv2md --header "Results" --align < data.tsv
devtui tsv2md -t "Results" -a < data.tsv
```

### Options

```
  -a, --align           align columns width
  -t, --header string   add main header (h1) to result
  -h, --help            help for tsv2md
```
