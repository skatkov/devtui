package cmd

import (
	"fmt"
	"strings"

	"github.com/skatkov/devtui/internal/cmderror"
	"github.com/skatkov/devtui/internal/input"
	"github.com/skatkov/devtui/internal/ui"
	graphqlquery "github.com/skatkov/devtui/tui/graphql-query"
	"github.com/spf13/cobra"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/formatter"
	"github.com/vektah/gqlparser/v2/parser"
)

var gqlfmtCmd = &cobra.Command{
	Use:   "gqlquery [string or file]",
	Short: "Format and prettify GraphQL queries",
	Long: `Format and prettify GraphQL queries for better readability.

By default, uses 2-space indentation and omits descriptions. Use flags to customize
indentation, include comments, or include descriptions. Input can be a string argument
or piped from stdin.`,
	Example: `  # Format GraphQL query from stdin
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
  curl -s https://api.example.com/schema.graphql | devtui gqlquery`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := input.ReadBytesFromArgsOrStdin(cmd, args)
		if err != nil {
			return err
		}

		if gqlTUI {
			// Initialize the TUI
			common := &ui.CommonModel{}

			model := graphqlquery.NewGraphQLQueryModel(common)
			err := model.SetContent(string(data))
			if err != nil {
				return err
			}

			return runProgram(model)
		}

		// Parse the GraphQL query
		inputStr := string(data)
		query, err := parser.ParseQuery(&ast.Source{
			Input: inputStr,
			Name:  "stdin",
		})
		if err != nil {
			return cmderror.FormatParseError("gqlquery", inputStr, err)
		}

		// Configure formatter options
		var opts []formatter.FormatterOption

		// Add indent option
		if gqlIndentString != "" {
			opts = append(opts, formatter.WithIndent(gqlIndentString))
		}

		// Include comments if requested
		if gqlWithComments {
			opts = append(opts, formatter.WithComments())
		}

		// Format the query
		var buf strings.Builder
		f := formatter.NewFormatter(&buf, opts...)
		f.FormatQueryDocument(query)
		result := buf.String()

		// Write the result to stdout
		_, err = fmt.Fprintln(cmd.OutOrStdout(), result)
		return err
	},
}

var (
	gqlIndentString string
	gqlWithComments bool
	gqlTUI          bool
)

func init() {
	rootCmd.AddCommand(gqlfmtCmd)

	gqlfmtCmd.Flags().StringVarP(&gqlIndentString, "indent", "i", "  ",
		"Indent string for nested elements (default is 2 spaces)")

	gqlfmtCmd.Flags().BoolVarP(&gqlWithComments, "with-comments", "c", false,
		"Include comments in the formatted output")

	// Query documents carry no descriptions, so this flag changes nothing; it stays for compatibility.
	gqlfmtCmd.Flags().BoolP("with-descriptions", "d", false,
		"Include descriptions in the formatted output (omitted by default)")
	gqlfmtCmd.Flags().BoolVarP(&gqlTUI, "tui", "t", false,
		"Open result in TUI")
}
