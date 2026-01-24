package cmd

import (
	"errors"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/skatkov/devtui/internal/cmderror"
	"github.com/skatkov/devtui/internal/input"
	"github.com/skatkov/devtui/internal/ui"
	"github.com/skatkov/devtui/tui/json2toml"
	"github.com/spf13/cobra"
)

var json2tomlCmd = &cobra.Command{
	Use:   "json2toml [string or file]",
	Short: "Convert JSON to TOML format",
	Long: `Convert JSON to TOML format.

Input can be a string argument or piped from stdin. JSON numbers are preserved
as integers when appropriate (not converted to floats). Use --tui flag to view
results in an interactive terminal interface.`,
	Example: `  # Convert JSON from stdin
  devtui json2toml < config.json
  cat app.json | devtui json2toml

  # Convert JSON string argument
  devtui json2toml '{"name": "myapp", "version": "1.0.0"}'

  # Output to file
  devtui json2toml < input.json > output.toml
  cat config.json | devtui json2toml > config.toml

  # Show results in interactive TUI
  devtui json2toml --tui < config.json
  devtui json2toml -t < config.json

  # Chain with other commands
  curl -s https://api.example.com/config | devtui json2toml`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := input.ReadBytesFromArgsOrStdin(cmd, args)
		if err != nil {
			return fmt.Errorf("error reading from stdin: %w", err)
		}

		if len(data) == 0 {
			return errors.New("no input provided. Pipe JSON input to this command")
		}

		if flagTUI {
			common := &ui.CommonModel{
				Width:  0,
				Height: 0,
			}

			model := json2toml.NewJsonTomlModel(common)
			err = model.SetContent(string(data))
			if err != nil {
				return err
			}

			p := tea.NewProgram(model, tea.WithAltScreen())
			if _, err := p.Run(); err != nil {
				return err
			}
			return nil
		}

		inputStr := string(data)
		result, err := json2toml.Convert(inputStr)
		if err != nil {
			return cmderror.FormatParseError("json2toml", inputStr, err)
		}

		_, err = fmt.Fprintln(cmd.OutOrStdout(), result)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(json2tomlCmd)
	json2tomlCmd.Flags().BoolVarP(&flagTUI, "tui", "t", false, "Show output in TUI")
}
