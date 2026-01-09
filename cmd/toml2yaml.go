package cmd

import (
	"errors"
	"fmt"

	"github.com/skatkov/devtui/internal/cmderror"
	"github.com/skatkov/devtui/internal/converter"
	"github.com/skatkov/devtui/internal/input"
	"github.com/spf13/cobra"
)

var toml2yamlCmd = &cobra.Command{
	Use:   "toml2yaml [string or file]",
	Short: "Convert TOML to YAML format",
	Long: `Convert TOML (Tom's Obvious Minimal Language) to YAML (YAML Ain't Markup Language) format.

Input can be a string argument or piped from stdin.`,
	Example: `  # Convert TOML from stdin
  devtui toml2yaml < config.toml
  cat app.toml | devtui toml2yaml

  # Convert TOML string argument
  devtui toml2yaml 'name = "myapp"\nversion = "1.0.0"'

  # Output to file
  devtui toml2yaml < input.toml > output.yaml

  # Chain with other commands
  curl -s https://api.example.com/config.toml | devtui toml2yaml`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := input.ReadBytesFromArgsOrStdin(cmd, args)
		if err != nil {
			return fmt.Errorf("error reading from stdin: %w", err)
		}

		if len(data) == 0 {
			return errors.New("no input provided. Pipe TOML input to this command")
		}

		inputStr := string(data)
		result, err := converter.TOMLToYAML(inputStr)
		if err != nil {
			return cmderror.FormatParseError("toml2yaml", inputStr, err)
		}

		_, err = fmt.Fprintln(cmd.OutOrStdout(), result)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(toml2yamlCmd)
}
