package cmd

import (
	"errors"
	"fmt"

	"github.com/skatkov/devtui/internal/cmderror"
	"github.com/skatkov/devtui/internal/converter"
	"github.com/skatkov/devtui/internal/input"
	"github.com/spf13/cobra"
)

var yaml2tomlCmd = &cobra.Command{
	Use:   "yaml2toml [string or file]",
	Short: "Convert YAML to TOML format",
	Long: `Convert YAML to TOML format.

Input can be a string argument or piped from stdin.`,
	Example: `  # Convert YAML from stdin
  devtui yaml2toml < config.yaml
  cat app.yaml | devtui yaml2toml

  # Convert YAML string argument
  devtui yaml2toml 'name: myapp\nversion: 1.0.0'

  # Output to file
  devtui yaml2toml < input.yaml > output.toml

  # Chain with other commands
  curl -s https://api.example.com/config.yaml | devtui yaml2toml`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := input.ReadBytesFromArgsOrStdin(cmd, args)
		if err != nil {
			return fmt.Errorf("error reading from stdin: %w", err)
		}

		if len(data) == 0 {
			return errors.New("no input provided. Pipe YAML input to this command")
		}

		inputStr := string(data)
		result, err := converter.YAMLToTOML(inputStr)
		if err != nil {
			return cmderror.FormatParseError("yaml2toml", inputStr, err)
		}

		_, err = fmt.Fprintln(cmd.OutOrStdout(), result)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(yaml2tomlCmd)
}
