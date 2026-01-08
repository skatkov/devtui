package cmd

import (
	"errors"
	"fmt"

	"github.com/skatkov/devtui/internal/cmderror"
	"github.com/skatkov/devtui/internal/input"
	"github.com/skatkov/devtui/internal/yaml"
	"github.com/spf13/cobra"
)

var yaml2jsonCmd = &cobra.Command{
	Use:   "yaml2json [string or file]",
	Short: "Convert YAML to JSON format",
	Long: `Convert YAML (YAML Ain't Markup Language) to JSON (JavaScript Object Notation) format.

Input can be a string argument or piped from stdin.`,
	Example: `  # Convert YAML from stdin
  devtui yaml2json < config.yaml
  cat app.yaml | devtui yaml2json

  # Convert YAML string argument
  devtui yaml2json 'name: myapp\nversion: 1.0.0'

  # Output to file
  devtui yaml2json < input.yaml > output.json

  # Chain with other commands
  curl -s https://api.example.com/config.yaml | devtui yaml2json`,
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
		result, err := yaml.YAMLToJSON(inputStr)
		if err != nil {
			return cmderror.FormatParseError("yaml2json", inputStr, err)
		}

		_, err = fmt.Fprintln(cmd.OutOrStdout(), result)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(yaml2jsonCmd)
}
