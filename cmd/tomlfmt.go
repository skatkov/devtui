package cmd

import (
	"io"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/skatkov/devtui/internal/ui"
	"github.com/skatkov/devtui/tui/toml"
	"github.com/spf13/cobra"
)

var tomlfmtCmd = &cobra.Command{
	Use:   "tomlfmt",
	Short: "Format TOML files",
	Long:  "Format TOML files",
	Run: func(cmd *cobra.Command, args []string) {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			log.Printf("ERROR: %s", err)
			return
		}

		if flagTUI {
			common := &ui.CommonModel{
				Width:  0, // Will be set by tea.WindowSizeMsg
				Height: 0,
			}

			model := toml.NewTomlFormatModel(common)
			err = model.SetContent(string(data))
			if err != nil {
				log.Printf("ERROR reading input: %s", err)
				return
			}

			p := tea.NewProgram(model, tea.WithAltScreen())
			if _, err := p.Run(); err != nil {
				log.Printf("ERROR: %s", err)
			}
			return
		}

		result, err := toml.Convert(string(data))
		if err != nil {
			log.Printf("ERROR: %s", err)
			return
		}

		_, err = os.Stdout.WriteString(result)
		if err != nil {
			log.Printf("ERROR: %s", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(tomlfmtCmd)
	tomlfmtCmd.Flags().BoolVarP(&flagTUI, "tui", "t", false, "Show output in TUI")
}
