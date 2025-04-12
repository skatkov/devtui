package cmd

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// Heavily inspired by csv2md
//   @see https://git.axenov.dev/anthony/csv2md/src/branch/master
// Example usage:
//   devtui tsv2md -t < example.tsv          - convert tsv from stdin and view result in stdout
//   devtui tsv2md < example.tsv > output.md - convert tsv from stdin and write result in new file
//   cat example.tsv | devtui tsv2md         - convert tsv from stdin and view result in stdout

var tsv2mdCmd = &cobra.Command{
	Use:   "tsv2md",
	Short: "Convert TSV to Markdown",
	Long:  "Convert TSV to Markdown",
	Run: func(cmd *cobra.Command, args []string) {
		tsvReader := csv.NewReader(os.Stdin)
		tsvReader.Comma = '\t'

		records, err := tsvReader.ReadAll()
		if err != nil {
			fmt.Printf("Failed to parse input from stdin: %v\n", err)
			return
		}

		Print(Convert(header, records, alignColumns))

	},
}

var (
	alignColumns bool   // align columns width
	header       string // add main header (h1) to result
)

func init() {
	rootCmd.AddCommand(tsv2mdCmd)

	tsv2mdCmd.Flags().BoolVarP(&alignColumns, "align", "a", false, "align columns width")
	tsv2mdCmd.Flags().StringVarP(&header, "header", "t", "", "add main header (h1) to result")
}

// Convert formats data from file or stdin as markdown
func Convert(header string, records [][]string, aligned bool) []string {
	var result []string

	// add h1 if passed
	header = strings.Trim(header, "\t\r\n ")
	if len(header) != 0 {
		result = append(result, "# "+header)
		result = append(result, "")
	}

	// if user wants aligned columns width then we
	// count max length of every value in every column
	widths := make(map[int]int)
	if aligned {
		for _, row := range records {
			for col_idx, col := range row {
				length := len(col)
				if len(widths) == 0 || widths[col_idx] < length {
					widths[col_idx] = length
				}
			}
		}
	}

	// build markdown table
	for row_idx, row := range records {

		// table content
		str := "| "
		for col_idx, col := range row {
			if aligned {
				str += fmt.Sprintf("%-*s | ", widths[col_idx], col)
			} else {
				str += col + " | "
			}
		}
		result = append(result, str)

		// content separator only after first row (header)
		if row_idx == 0 {
			str := "| "
			for col_idx := range row {
				if !aligned || widths[col_idx] < 3 {
					str += "--- | "
				} else {
					str += strings.Repeat("-", widths[col_idx]) + " | "
				}
			}
			result = append(result, str)
		}
	}
	return result
}

func Print(data []string) {
	for _, row := range data {
		fmt.Println(row)
	}
}
