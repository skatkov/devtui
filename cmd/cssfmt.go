package cmd

import (
	"io"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/client9/csstool"
	"github.com/skatkov/devtui/internal/ui"
	"github.com/skatkov/devtui/tui/css"
	"github.com/spf13/cobra"
)

var cssfmtCmd = &cobra.Command{
	Use:   "cssfmt",
	Short: "Format CSS files",
	Long:  "Format CSS files",
	Example: `
	cssfmt < testdata/bootstrap.min.css
	cssfmt < testdata/bootstrap.min.css --tui # Show results in a TUI
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if flagTui {
			input, err := io.ReadAll(cmd.InOrStdin())
			if err != nil {
				return err
			}

			common := &ui.CommonModel{
				Width:  100, // Default width, will be adjusted by the TUI
				Height: 30,  // Default height, will be adjusted by the TUI
			}

			model := css.NewCSSFormatterModel(common)
			err = model.SetContent(string(input))
			if err != nil {
				return err
			}
			p := tea.NewProgram(
				model,
				tea.WithAltScreen(),       // Use alternate screen buffer
				tea.WithMouseCellMotion(), // Enable mouse support
			)

			if _, err := p.Run(); err != nil {
				return err
			}

			return nil
		}

		if flagTab {
			flagIndent = 1
		}
		cssformat := csstool.NewCSSFormat(flagIndent, flagTab, nil)
		cssformat.AlwaysSemicolon = flagSemicolon
		err := cssformat.Format(os.Stdin, os.Stdout)
		if err != nil {
			return err
		}

		return nil
	},
}

var (
	flagTab       bool
	flagIndent    int
	flagSemicolon bool
	flagTui       bool
)

func init() {
	rootCmd.AddCommand(cssfmtCmd)
	cssfmtCmd.Flags().BoolVarP(&flagTab, "tab", "t", false, "use tabs for indentation")
	cssfmtCmd.Flags().IntVarP(&flagIndent, "indent", "i", 2, "spaces for indentation")
	cssfmtCmd.Flags().BoolVarP(&flagSemicolon, "semicolon", "", true, "always end rule with semicolon, even if not needed")
	cssfmtCmd.Flags().BoolVarP(&flagTui, "tui", "", false, "present result in a TUI")
}
