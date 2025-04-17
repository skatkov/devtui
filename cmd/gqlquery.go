package cmd

import (
	"io"
	"log"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/skatkov/devtui/internal/ui"
	graphqlquery "github.com/skatkov/devtui/tui/graphql-query"
	"github.com/spf13/cobra"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/formatter"
	"github.com/vektah/gqlparser/v2/parser"
)

var gqlfmtCmd = &cobra.Command{
	Use:   "gqlquery",
	Short: "Format GraphQL queries",
	Long:  "Format GraphQL queries for better readability",
	Example: `
	gqlquery < testdata/query.graphql # Format and output to stdout
 	gqlquery < testdata/query.graphql > formatted.graphql # Output to file
	gqlfmt --indent "    " --with-comments --with-descriptions < testdata/query.graphql # With formatting options
	`,
	Run: func(cmd *cobra.Command, args []string) {
		// Read all input data from stdin
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			log.Printf("ERROR: %s", err)
			return
		}

		if flagTUI {
			// Initialize the TUI
			common := &ui.CommonModel{
				Width:  100, // Default width, will be adjusted by the TUI
				Height: 30,  // Default height, will be adjusted by the TUI
			}

			model := graphqlquery.NewGraphQLQueryModel(common)
			model.SetContent(string(data))

			p := tea.NewProgram(
				model,
				tea.WithAltScreen(),       // Use alternate screen buffer
				tea.WithMouseCellMotion(), // Enable mouse support
			)

			if _, err := p.Run(); err != nil {
				log.Printf("ERROR running TUI: %s", err)
			}
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
	flagTUI             bool
)

func init() {
	rootCmd.AddCommand(gqlfmtCmd)

	gqlfmtCmd.Flags().StringVarP(&gqlIndentString, "indent", "i", "  ",
		"Indent string for nested elements (default is 2 spaces)")

	gqlfmtCmd.Flags().BoolVarP(&gqlWithComments, "with-comments", "c", false,
		"Include comments in the formatted output")

	gqlfmtCmd.Flags().BoolVarP(&gqlWithDescriptions, "with-descriptions", "d", false,
		"Include descriptions in the formatted output (omitted by default)")
	gqlfmtCmd.Flags().BoolVarP(&flagTUI, "tui", "t", false,
		"Open result in TUI")
}
