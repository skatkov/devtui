---
title: csv2md
parent: CLI
---

## devtui csv2md

Convert CSV to Markdown table format

### Synopsis

Convert CSV (Comma-Separated Values) to Markdown table format for documentation.

Input can be piped from stdin or read from a file. Use --align to align column widths
and --header to add a main heading (h1) to the output.

```bash
devtui csv2md [string or file] [flags]
```

### Examples

```bash
# Convert CSV from stdin
devtui csv2md < example.csv
cat data.csv | devtui csv2md
# Output to file
devtui csv2md < input.csv > output.md
cat data.csv | devtui csv2md > table.md
# Add main header to output
devtui csv2md --header "User Data" < users.csv
devtui csv2md -t "Sales Report" < sales.csv
# Align column widths for better readability
devtui csv2md --align < data.csv
devtui csv2md -a < data.csv
# Combine options
devtui csv2md --header "Results" --align < data.csv
devtui csv2md -t "Results" -a < data.csv
```

### Options

```
  -a, --align           align columns width
  -t, --header string   add main header (h1) to result
  -h, --help            help for csv2md
```
