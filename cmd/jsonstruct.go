package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/skatkov/devtui/internal/cmderror"
	"github.com/skatkov/devtui/internal/input"
	"github.com/skatkov/devtui/internal/structgen"
	"github.com/spf13/cobra"
)

var jsonstructCmd = &cobra.Command{
	Use:   "jsonstruct [string or file]",
	Short: "Convert JSON to Go struct",
	Long: `Convert JSON input into a Go struct definition.

Input can be a string argument or piped from stdin.`,
	Example: `  # Convert JSON from stdin
  devtui jsonstruct < data.json
  cat data.json | devtui jsonstruct

  # Convert JSON string argument
  devtui jsonstruct '{"name":"Alice","age":30}'

  # Output to file
  devtui jsonstruct < input.json > struct.go`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := input.ReadBytesFromArgsOrStdin(cmd, args)
		if err != nil {
			return fmt.Errorf("error reading from stdin: %w", err)
		}
		if len(data) == 0 {
			return errors.New("no input provided. pipe JSON input to this command")
		}

		inputStr := string(data)
		result, err := structgen.JSONToGoStruct(strings.NewReader(inputStr))
		if err != nil {
			return cmderror.FormatParseError("jsonstruct", inputStr, err)
		}

		_, err = fmt.Fprintln(cmd.OutOrStdout(), result)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(jsonstructCmd)
}
