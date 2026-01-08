package cmd

import (
	"errors"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/skatkov/devtui/internal/cmderror"
	"github.com/skatkov/devtui/internal/input"
	"github.com/skatkov/devtui/internal/ui"
	"github.com/skatkov/devtui/tui/toml"
	"github.com/spf13/cobra"
)

var toml2yamlCmd = &cobra.Command{
	Use:   "toml2yaml [string or file]",
	Short: "Convert TOML to YAML format",
	Long: `Convert TOML (Tom's Obvious Minimal Language) to YAML (YAML Ain't Markup Language) format.

Input can be a string argument or piped from stdin. Use --tui flag to view
results in an interactive terminal interface.`,
	Example: `  # Convert TOML from stdin
  devtui toml2yaml < config.toml
  cat app.toml | devtui toml2yaml

  # Convert TOML string argument
  devtui toml2yaml 'name = "myapp"\nversion = "1.0.0"'

  # Output to file
  devtui toml2yaml < input.toml > output.yaml

  # Show results in interactive TUI
  devtui toml2yaml --tui < config.toml

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

		if flagTUI {
			common := &ui.CommonModel{
				Width:  0,
				Height: 0,
			}

			model := toml.NewTomlFormatModel(common)
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
		result, err := toml.TOMLToYAML(inputStr)
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
	toml2yamlCmd.Flags().BoolVarP(&flagTUI, "tui", "t", false, "Show output in TUI")
}
