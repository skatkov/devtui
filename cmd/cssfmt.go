package cmd

import (
	"io"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/client9/csstool"
	"github.com/skatkov/devtui/internal/ui"
	"github.com/skatkov/devtui/tui/css"
	"github.com/spf13/cobra"
)

var cssfmtCmd = &cobra.Command{
	Use:     "cssfmt",
	Short:   "Format CSS files",
	Long:    "Format CSS files",
	Example: `cssfmt < testdata/bootstrap.min.css`,
	Run: func(cmd *cobra.Command, args []string) {
		if flagTui {
			// Read input CSS
			input, err := io.ReadAll(os.Stdin)
			if err != nil {
				log.Printf("ERROR reading input: %s", err)
				return
			}

			// Initialize the TUI
			common := &ui.CommonModel{
				Width:  80, // Default width, will be adjusted by the TUI
				Height: 24, // Default height, will be adjusted by the TUI
			}

			model := css.NewCSSFormatterModel(common)
			model.SetContent(string(input))

			p := tea.NewProgram(
				model,
				tea.WithAltScreen(),       // Use alternate screen buffer
				tea.WithMouseCellMotion(), // Enable mouse support
			)

			if _, err := p.Run(); err != nil {
				log.Printf("ERROR running TUI: %s", err)
			}
			return
		}

		if flagTab {
			flagIndent = 1
		}
		cssformat := csstool.NewCSSFormat(flagIndent, flagTab, nil)
		cssformat.AlwaysSemicolon = flagSemicolon
		err := cssformat.Format(os.Stdin, os.Stdout)
		if err != nil {
			log.Printf("ERROR: %s", err)
		}
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
