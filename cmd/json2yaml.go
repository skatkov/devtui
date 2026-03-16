package cmd

import (
	"errors"
	"fmt"

	"github.com/skatkov/devtui/internal/cmderror"
	"github.com/skatkov/devtui/internal/converter"
	"github.com/skatkov/devtui/internal/input"
	"github.com/spf13/cobra"
)

var json2yamlCmd = &cobra.Command{
	Use:   "json2yaml [string or file]",
	Short: "Convert JSON to YAML format",
	Long: `Convert JSON to YAML format.

Input can be a string argument or piped from stdin.`,
	Example: `  # Convert JSON from stdin
  devtui json2yaml < config.json
  cat app.json | devtui json2yaml

  # Convert JSON string argument
  devtui json2yaml '{"name": "myapp", "version": "1.0.0"}'

  # Output to file
  devtui json2yaml < input.json > output.yaml

  # Chain with other commands
  curl -s https://api.example.com/config | devtui json2yaml`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := input.ReadBytesFromArgsOrStdin(cmd, args)
		if err != nil {
			return fmt.Errorf("error reading from stdin: %w", err)
		}

		if len(data) == 0 {
			return errors.New("no input provided. Pipe JSON input to this command")
		}

		inputStr := string(data)
		result, err := converter.JSONToYAML(inputStr)
		if err != nil {
			return cmderror.FormatParseError("json2yaml", inputStr, err)
		}
		if outputJSON {
			return writeJSONValue(cmd.OutOrStdout(), result)
		}

		_, err = fmt.Fprintln(cmd.OutOrStdout(), result)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(json2yamlCmd)
}
