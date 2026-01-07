package cmd

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/skatkov/devtui/internal/cmderror"
	"github.com/skatkov/devtui/internal/input"
	"github.com/skatkov/devtui/internal/ui"
	"github.com/skatkov/devtui/tui/tomljson"
	"github.com/spf13/cobra"
)

var toml2jsonCmd = &cobra.Command{
	Use:   "toml2json [string or file]",
	Short: "Convert TOML to JSON format",
	Long: `Convert TOML (Tom's Obvious Minimal Language) to JSON format.

Input can be a string argument or piped from stdin. Use --tui flag to view results
in an interactive terminal interface.`,
	Example: `  # Convert TOML from stdin
  devtui toml2json < config.toml
  cat app.toml | devtui toml2json

  # Convert TOML string argument
  devtui toml2json '[package]\nname = "myapp"'

  # Output to file
  devtui toml2json < input.toml > output.json
  cat config.toml | devtui toml2json > data.json

  # Show results in interactive TUI
  devtui toml2json --tui < config.toml
  devtui toml2json -t < config.toml

  # Chain with other commands
  curl -s https://example.com/config.toml | devtui toml2json`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := input.ReadBytesFromArgsOrStdin(cmd, args)
		if err != nil {
			return err
		}

		if flagTUI {
			common := &ui.CommonModel{
				Width:  0,
				Height: 0,
			}

			model := tomljson.NewTomlJsonModel(common)
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
		result, err := tomljson.Convert(inputStr)
		if err != nil {
			return cmderror.FormatParseError("TOML", "toml2json", inputStr, err)
		}

		_, err = fmt.Fprintln(cmd.OutOrStdout(), result)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(toml2jsonCmd)
	toml2jsonCmd.Flags().BoolVarP(&flagTUI, "tui", "t", false, "Show output in TUI")
}
