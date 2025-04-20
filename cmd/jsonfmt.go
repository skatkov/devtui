package cmd

import (
	"io"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/skatkov/devtui/internal/ui"
	"github.com/skatkov/devtui/tui/json"
	"github.com/spf13/cobra"
)

var jsonfmtCmd = &cobra.Command{
	Use:   "jsonfmt",
	Short: "Format JSON",
	Long:  "Format JSON",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			log.Printf("ERROR: %s", err)
			return
		}

		if flagTUI {
			common := &ui.CommonModel{
				Width:  80,
				Height: 80,
			}

			model := json.NewJsonModel(common)
			model.SetContent(string(data))

			p := tea.NewProgram(model, tea.WithAltScreen())
			if _, err := p.Run(); err != nil {
				log.Printf("ERROR: %s", err)
			}
			return
		}

		result := json.FormatJSON(string(data))

		_, err = os.Stdout.WriteString(result)
		if err != nil {
			log.Printf("ERROR: %s", err)
		}

	},
}

// TODO: implement setIndent flag
// TODO: implement EscapeHTML flag

func init() {
	rootCmd.AddCommand(jsonfmtCmd)
	jsonfmtCmd.Flags().BoolVarP(&flagTUI, "tui", "t", false, "Show output in TUI")
}
