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

var csv2mdCmd = &cobra.Command{
	Use:   "csv2md",
	Short: "Convert CSV to Markdown Table",
	Long:  "Convert CSV to Markdown Table",
	Example: `  devtui csv2md -t < example.tsv          - convert tsv from stdin and view result in stdout
	devtui csv2md < example.tsv > output.md - convert tsv from stdin and write result in new file
	cat example.tsv | devtui csv2md         - convert tsv from stdin and view result in stdout`,
	Run: func(cmd *cobra.Command, args []string) {
		csvReader := csv.NewReader(os.Stdin)

		records, err := csvReader.ReadAll()
		if err != nil {
			fmt.Printf("Failed to parse input from stdin: %v\n", err)
			return
		}

		csv2md.Print(csv2md.Convert(csv2mdHeader, records, csv2mdAlignColumns))
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
