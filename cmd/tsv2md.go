package cmd

import (
	"encoding/csv"

	"github.com/skatkov/devtui/internal/csv2md"
	"github.com/spf13/cobra"
)

// Heavily inspired by csv2md
//   @see https://git.axenov.dev/anthony/csv2md/src/branch/master

var tsv2mdCmd = &cobra.Command{
	Use:   "tsv2md",
	Short: "Convert TSV to Markdown Table",
	Long:  "Convert TSV to Markdown Table",
	Example: `  devtui tsv2md -t < example.tsv          # convert tsv from stdin and view result in stdout
	devtui tsv2md < example.tsv > output.md # convert tsv from stdin and write result in new file
	cat example.tsv | devtui tsv2md         # convert tsv from stdin and view result in stdout`,
	RunE: func(cmd *cobra.Command, args []string) error {
		tsvReader := csv.NewReader(cmd.InOrStdin())
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
