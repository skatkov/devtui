package cmd

import (
	"errors"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/skatkov/devtui/internal/cmderror"
	"github.com/skatkov/devtui/internal/input"
	"github.com/skatkov/devtui/internal/ui"
	"github.com/skatkov/devtui/tui/yaml"
	"github.com/spf13/cobra"
)

var json2yamlCmd = &cobra.Command{
	Use:   "json2yaml [string or file]",
	Short: "Convert JSON to YAML format",
	Long: `Convert JSON (JavaScript Object Notation) to YAML (YAML Ain't Markup Language) format.

Input can be a string argument or piped from stdin. Use --tui flag to view
results in an interactive terminal interface.`,
	Example: `  # Convert JSON from stdin
  devtui json2yaml < config.json
  cat app.json | devtui json2yaml

  # Convert JSON string argument
  devtui json2yaml '{"name": "myapp", "version": "1.0.0"}'

  # Output to file
  devtui json2yaml < input.json > output.yaml

  # Show results in interactive TUI
  devtui json2yaml --tui < config.json

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

		if flagTUI {
			common := &ui.CommonModel{
				Width:  0,
				Height: 0,
			}

			model := yaml.NewYamlModel(common)
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
		result, err := yaml.JSONToYAML(inputStr)
		if err != nil {
			return cmderror.FormatParseError("json2yaml", inputStr, err)
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
	json2yamlCmd.Flags().BoolVarP(&flagTUI, "tui", "t", false, "Show output in TUI")
}
