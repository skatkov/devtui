---
title: gqlquery
parent: CLI
---

## devtui gqlquery

Format and prettify GraphQL queries

### Synopsis

Format and prettify GraphQL queries for better readability.

By default, uses 2-space indentation and omits descriptions. Use flags to customize
indentation, include comments, or include descriptions. Input can be a string argument
or piped from stdin.

```bash
devtui gqlquery [string or file] [flags]
```

### Examples

```bash
# Format GraphQL query from stdin
devtui gqlquery < query.graphql
cat query.graphql | devtui gqlquery
# Format GraphQL string argument
devtui gqlquery 'query { user(id: 1) { name email } }'
# Output to file
devtui gqlquery < input.graphql > formatted.graphql
cat query.graphql | devtui gqlquery > pretty.graphql
# Custom indentation (4 spaces)
devtui gqlquery --indent "    " < query.graphql
devtui gqlquery -i "    " < query.graphql
# Include comments in output
devtui gqlquery --with-comments < query.graphql
devtui gqlquery -c < query.graphql
# Include descriptions in output
devtui gqlquery --with-descriptions < query.graphql
devtui gqlquery -d < query.graphql
# Combine formatting options
devtui gqlquery -i "    " -c -d < query.graphql
# Show results in interactive TUI
devtui gqlquery --tui < query.graphql
devtui gqlquery -t < query.graphql
# Chain with other commands
curl -s https://api.example.com/schema.graphql | devtui gqlquery
```

### Options

```
  -h, --help                help for gqlquery
  -i, --indent string       Indent string for nested elements (default is 2 spaces) (default "  ")
  -t, --tui                 Open result in TUI
  -c, --with-comments       Include comments in the formatted output
  -d, --with-descriptions   Include descriptions in the formatted output (omitted by default)
```
