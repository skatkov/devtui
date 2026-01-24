package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/skatkov/devtui/internal/input"
	"github.com/skatkov/devtui/internal/numbers"
	"github.com/spf13/cobra"
)

var numbersCmd = &cobra.Command{
	Use:   "numbers [number]",
	Short: "Convert numbers between bases",
	Long: `Convert numbers between binary, octal, decimal, and hexadecimal.

Input can be a number argument or piped from stdin.`,
	Example: `  # Convert a decimal number
  devtui numbers 42

  # Convert a binary number
  devtui numbers --base 2 101010

  # Convert from stdin
  echo "ff" | devtui numbers --base 16

  # Output as JSON
  devtui numbers --json 42`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputStr, err := input.ReadFromArgsOrStdin(cmd, args)
		if err != nil {
			return fmt.Errorf("error reading from stdin: %w", err)
		}
		value := strings.TrimSpace(inputStr)
		if value == "" {
			return errors.New("no input provided. pipe number input to this command")
		}

		result, err := numbers.Convert(value, numbersBase)
		if err != nil {
			return err
		}

		if numbersJSONOutput {
			bytes, err := json.MarshalIndent(result, "", "  ")
			if err != nil {
				return err
			}
			_, err = fmt.Fprintln(cmd.OutOrStdout(), string(bytes))
			return err
		}

		output := table.New().Border(lipgloss.NormalBorder())
		for _, conversion := range result.Conversions {
			output.Row(conversion.Label, conversion.Value)
		}

		_, err = fmt.Fprintln(cmd.OutOrStdout(), output.String())
		return err
	},
}

var (
	numbersBase       int
	numbersJSONOutput bool
)

func init() {
	rootCmd.AddCommand(numbersCmd)
	numbersCmd.Flags().IntVarP(&numbersBase, "base", "b", 10, "input number base (2, 8, 10, 16)")
	numbersCmd.Flags().BoolVar(&numbersJSONOutput, "json", false, "output conversions as JSON")
}
