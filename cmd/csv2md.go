package cmd

import (
	"encoding/csv"
	"strings"

	"github.com/skatkov/devtui/internal/cmderror"
	"github.com/skatkov/devtui/internal/csv2md"
	"github.com/skatkov/devtui/internal/input"
	"github.com/spf13/cobra"
)

// Heavily inspired by csv2md
//   @see https://git.axenov.dev/anthony/csv2md/src/branch/master

var csv2mdCmd = &cobra.Command{
	Use:   "csv2md [string or file]",
	Short: "Convert CSV to Markdown table format",
	Long: `Convert CSV (Comma-Separated Values) to Markdown table format for documentation.

Input can be piped from stdin or read from a file. Use --align to align column widths
and --header to add a main heading (h1) to the output.`,
	Example: `  # Convert CSV from stdin
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
  devtui csv2md -t "Results" -a < data.csv`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := input.ReadBytesFromArgsOrStdin(cmd, args)
		if err != nil {
			return err
		}

		inputStr := string(data)
		csvReader := csv.NewReader(strings.NewReader(inputStr))
		records, err := csvReader.ReadAll()
		if err != nil {
			return cmderror.FormatParseError("csv2md", inputStr, err)
		}

		csv2md.Print(csv2md.Convert(csv2mdHeader, records, csv2mdAlignColumns))

		return nil
	},
}

var (
	csv2mdAlignColumns bool   // align columns width
	csv2mdHeader       string // add main header (h1) to result
)

func init() {
	rootCmd.AddCommand(csv2mdCmd)

	csv2mdCmd.Flags().BoolVarP(&csv2mdAlignColumns, "align", "a", false, "align columns width")
	csv2mdCmd.Flags().StringVarP(&csv2mdHeader, "header", "t", "", "add main header (h1) to result")
}
