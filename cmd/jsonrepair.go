package cmd

import (
	"errors"
	"fmt"
	"io"

	"github.com/skatkov/devtui/tui/jsonrepair"
	"github.com/spf13/cobra"
)

var jsonrepairCmd = &cobra.Command{
	Use:   "jsonrepair",
	Short: "Repair malformed JSON",
	Long: `Repair malformed JSON, particularly useful for fixing JSON output from LLMs.

This tool can fix various JSON issues including:
- Single quotes instead of double quotes
- Unclosed arrays and objects
- Mixed quotes
- Uppercase TRUE/FALSE/Null values
- Trailing commas
- JSON wrapped in markdown code blocks
- And many more LLM-specific JSON issues`,
	Example: `
	# Repair JSON from stdin
	echo "{'key': 'value'}" | devtui jsonrepair

	# Repair JSON from file
	devtui jsonrepair < broken.json

	# Output to file
	devtui jsonrepair < broken.json > fixed.json

	# Chain with other commands
	cat llm-output.txt | devtui jsonrepair | devtui jsonfmt
	`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := io.ReadAll(cmd.InOrStdin())
		if err != nil {
			return fmt.Errorf("error reading from stdin: %w", err)
		}

		if len(data) == 0 {
			return errors.New("no input provided. Pipe JSON input to this command")
		}

		result, err := jsonrepair.RepairJSON(string(data))
		if err != nil {
			return fmt.Errorf("failed to repair JSON: %w", err)
		}

		_, err = fmt.Fprintln(cmd.OutOrStdout(), result)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(jsonrepairCmd)
}
