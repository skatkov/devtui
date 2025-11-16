package cmd

import (
	"encoding/csv"
	"strings"

	"github.com/skatkov/devtui/internal/csv2md"
	"github.com/skatkov/devtui/internal/input"
	"github.com/spf13/cobra"
)

// Heavily inspired by csv2md
//   @see https://git.axenov.dev/anthony/csv2md/src/branch/master

var tsv2mdCmd = &cobra.Command{
	Use:   "tsv2md [string or file]",
	Short: "Convert TSV to Markdown table format",
	Long: `Convert TSV (Tab-Separated Values) to Markdown table format for documentation.

Input can be piped from stdin or read from a file. Use --align to align column widths
and --header to add a main heading (h1) to the output.`,
	Example: `  # Convert TSV from stdin
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
  devtui tsv2md -t "Results" -a < data.tsv`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Read from args or stdin - the function handles both cases
		data, err := input.ReadBytesFromArgsOrStdin(cmd, args)
		if err != nil {
			return err
		}

		tsvReader := csv.NewReader(strings.NewReader(string(data)))
		tsvReader.Comma = '\t'

		records, err := tsvReader.ReadAll()
		if err != nil {
			return err
		}

		csv2md.Print(csv2md.Convert(tsv2mdHeader, records, tsv2mdAlignColumns))

		return nil
	},
}

var (
	tsv2mdAlignColumns bool   // align columns width
	tsv2mdHeader       string // add main header (h1) to result
)

func init() {
	rootCmd.AddCommand(tsv2mdCmd)

	tsv2mdCmd.Flags().BoolVarP(&tsv2mdAlignColumns, "align", "a", false, "align columns width")
	tsv2mdCmd.Flags().StringVarP(&tsv2mdHeader, "header", "t", "", "add main header (h1) to result")
}
