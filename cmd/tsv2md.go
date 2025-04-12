package cmd

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/skatkov/devtui/internal/csv2md"
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
	Short: "Convert TSV to Markdown Table",
	Long:  "Convert TSV to Markdown Table",
	Run: func(cmd *cobra.Command, args []string) {
		tsvReader := csv.NewReader(os.Stdin)
		tsvReader.Comma = '\t'

		records, err := tsvReader.ReadAll()
		if err != nil {
			fmt.Printf("Failed to parse input from stdin: %v\n", err)
			return
		}

		csv2md.Print(csv2md.Convert(tsv2mdHeader, records, tsv2mdAlignColumns))
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
