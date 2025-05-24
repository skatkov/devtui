---
title: gqlquery
parent: CLI
---

## devtui gqlquery

Format GraphQL queries

### Synopsis

Format GraphQL queries for better readability

```
devtui gqlquery [flags]
```

### Examples

```

	gqlquery < testdata/query.graphql # Format and output to stdout
 	gqlquery < testdata/query.graphql > formatted.graphql # Output to file
	gqlquery --indent "    " --with-comments --with-descriptions < testdata/query.graphql # With formatting options
	gqlquery < testdata/query.graphql --tui # Show results in a TUI
	
```

### Options

```
  -h, --help                help for gqlquery
  -i, --indent string       Indent string for nested elements (default is 2 spaces) (default "  ")
  -t, --tui                 Open result in TUI
  -c, --with-comments       Include comments in the formatted output
  -d, --with-descriptions   Include descriptions in the formatted output (omitted by default)
```
