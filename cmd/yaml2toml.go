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

var yaml2tomlCmd = &cobra.Command{
	Use:   "yaml2toml [string or file]",
	Short: "Convert YAML to TOML format",
	Long: `Convert YAML (YAML Ain't Markup Language) to TOML (Tom's Obvious Minimal Language) format.

Input can be a string argument or piped from stdin. Use --tui flag to view
results in an interactive terminal interface.`,
	Example: `  # Convert YAML from stdin
  devtui yaml2toml < config.yaml
  cat app.yaml | devtui yaml2toml

  # Convert YAML string argument
  devtui yaml2toml 'name: myapp\nversion: 1.0.0'

  # Output to file
  devtui yaml2toml < input.yaml > output.toml

  # Show results in interactive TUI
  devtui yaml2toml --tui < config.yaml

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
		result, err := yaml.YAMLToTOML(inputStr)
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
	yaml2tomlCmd.Flags().BoolVarP(&flagTUI, "tui", "t", false, "Show output in TUI")
}
