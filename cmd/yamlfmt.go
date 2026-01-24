package cmd

import (
	"errors"
	"fmt"

	"github.com/skatkov/devtui/internal/cmderror"
	"github.com/skatkov/devtui/internal/input"
	"github.com/skatkov/devtui/internal/yamlfmt"
	"github.com/spf13/cobra"
)

var yamlfmtCmd = &cobra.Command{
	Use:   "yamlfmt [string or file]",
	Short: "Format and prettify YAML",
	Long: `Format and prettify YAML input with proper indentation.

Input can be a string argument, piped from stdin, or read from a file.`,
	Example: `  # Format YAML from stdin
  devtui yamlfmt < config.yaml
  cat config.yaml | devtui yamlfmt

  # Format YAML string argument
  devtui yamlfmt 'name: myapp\nversion: 1.0.0'

  # Output to file
  devtui yamlfmt < input.yaml > formatted.yaml`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := input.ReadBytesFromArgsOrStdin(cmd, args)
		if err != nil {
			return fmt.Errorf("error reading from stdin: %w", err)
		}
		if len(data) == 0 {
			return errors.New("no input provided. pipe YAML input to this command")
		}

		inputStr := string(data)
		result, err := yamlfmt.Format(inputStr)
		if err != nil {
			return cmderror.FormatParseError("yamlfmt", inputStr, err)
		}

		_, err = fmt.Fprintln(cmd.OutOrStdout(), result)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(yamlfmtCmd)
}
