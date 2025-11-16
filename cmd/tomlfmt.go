package cmd

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/skatkov/devtui/internal/input"
	"github.com/skatkov/devtui/internal/ui"
	"github.com/skatkov/devtui/tui/toml"
	"github.com/spf13/cobra"
)

var tomlfmtCmd = &cobra.Command{
	Use:   "tomlfmt",
	Short: "Format TOML files",
	Long:  "Format TOML files",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := input.ReadFromStdin(cmd)
		if err != nil {
			return err
		}

		if flagTUI {
			common := &ui.CommonModel{
				Width:  0, // Will be set by tea.WindowSizeMsg
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

		result, err := toml.Convert(string(data))
		if err != nil {
			return err
		}

		_, err = fmt.Fprintln(cmd.OutOrStdout(), result)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(tomlfmtCmd)
	tomlfmtCmd.Flags().BoolVarP(&flagTUI, "tui", "t", false, "Show output in TUI")
}
