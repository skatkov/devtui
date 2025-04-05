package cmd

import (
	"io"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/formatter"
	"github.com/vektah/gqlparser/v2/parser"
)

// ---
// Example usage:
// go run . gqlquery < testdata/query.graphql
//
// Output to file:
// go run . gqlquery < testdata/query.graphql > formatted.graphql
//
// With options:
// go run . gqlfmt --indent "    " --with-comments --with-descriptions < testdata/query.graphql
// ---

var gqlfmtCmd = &cobra.Command{
	Use:   "gqlquery",
	Short: "Format GraphQL queries",
	Long:  "Format GraphQL queries for better readability",
	Run: func(cmd *cobra.Command, args []string) {
		// Read all input data from stdin
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			log.Printf("ERROR: %s", err)
			return
		}

		// Parse the GraphQL query
		query, err := parser.ParseQuery(&ast.Source{
			Input: string(data),
			Name:  "stdin",
		})
		if err != nil {
			log.Printf("ERROR parsing GraphQL: %s", err)
			return
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

		// Include descriptions if requested
		if !gqlWithDescriptions {
			opts = append(opts, formatter.WithoutDescription())
		}

		// Format the query
		var buf strings.Builder
		f := formatter.NewFormatter(&buf, opts...)
		f.FormatQueryDocument(query)
		result := buf.String()

		// Write the result to stdout
		_, err = os.Stdout.WriteString(result)
		if err != nil {
			log.Printf("ERROR: %s", err)
		}
	},
}

var (
	gqlIndentString     string
	gqlWithComments     bool
	gqlWithDescriptions bool
)

func init() {
	rootCmd.AddCommand(gqlfmtCmd)

	gqlfmtCmd.Flags().StringVarP(&gqlIndentString, "indent", "i", "  ",
		"Indent string for nested elements (default is 2 spaces)")

	gqlfmtCmd.Flags().BoolVarP(&gqlWithComments, "with-comments", "c", false,
		"Include comments in the formatted output")

	gqlfmtCmd.Flags().BoolVarP(&gqlWithDescriptions, "with-descriptions", "d", false,
		"Include descriptions in the formatted output (omitted by default)")
}
