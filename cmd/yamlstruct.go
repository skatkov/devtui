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

var yamlstructCmd = &cobra.Command{
	Use:   "yamlstruct [string or file]",
	Short: "Convert YAML to Go struct",
	Long: `Convert YAML input into a Go struct definition.

Input can be a string argument or piped from stdin.`,
	Example: `  # Convert YAML from stdin
  devtui yamlstruct < config.yaml
  cat config.yaml | devtui yamlstruct

  # Convert YAML string argument
  devtui yamlstruct 'name: Alice\nage: 30'

  # Output to file
  devtui yamlstruct < input.yaml > struct.go`,
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
		result, err := structgen.YAMLToGoStruct(strings.NewReader(inputStr))
		if err != nil {
			return cmderror.FormatParseError("yamlstruct", inputStr, err)
		}

		_, err = fmt.Fprintln(cmd.OutOrStdout(), result)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(yamlstructCmd)
}
