package cmd

import (
	"errors"
	"fmt"

	"github.com/skatkov/devtui/internal/cmderror"
	"github.com/skatkov/devtui/internal/csv2json"
	"github.com/skatkov/devtui/internal/input"
	"github.com/spf13/cobra"
)

var csv2jsonCmd = &cobra.Command{
	Use:   "csv2json [string or file]",
	Short: "Convert CSV to JSON",
	Long: `Convert CSV (Comma-Separated Values) into formatted JSON.

Input can be a string argument or piped from stdin.`,
	Example: `  # Convert CSV from stdin
  devtui csv2json < data.csv
  cat data.csv | devtui csv2json

  # Convert CSV string argument
  devtui csv2json 'name,age\nAlice,30'

  # Output to file
  devtui csv2json < input.csv > output.json`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := input.ReadBytesFromArgsOrStdin(cmd, args)
		if err != nil {
			return fmt.Errorf("error reading from stdin: %w", err)
		}
		if len(data) == 0 {
			return errors.New("no input provided. pipe CSV input to this command")
		}

		inputStr := string(data)
		result, err := csv2json.Convert(inputStr)
		if err != nil {
			return cmderror.FormatParseError("csv2json", inputStr, err)
		}

		_, err = fmt.Fprintln(cmd.OutOrStdout(), result)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(csv2jsonCmd)
}
