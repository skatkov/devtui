package cmd

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/skatkov/devtui/tui/json"
	"github.com/spf13/cobra"
)

var jsonfmtCmd = &cobra.Command{
	Use:   "jsonfmt",
	Short: "Format JSON",
	Long:  "Format JSON",
	Example: `
	devtui jsonfmt < testdata/example.json # Format and output to stdout
 	devtui jsonfmt < testdata/example.json > formatted.json # Output to file
	`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			log.Printf("ERROR: %s", err)
			return
		}

		// if flagTUI {
		// 	common := &ui.CommonModel{
		// 		Width:  80,
		// 		Height: 80,
		// 	}

		// 	model := json.NewJsonModel(common)
		// 	model.SetContent(string(data))

		// 	p := tea.NewProgram(model, tea.WithAltScreen())
		// 	if _, err := p.Run(); err != nil {
		// 		log.Printf("ERROR: %s", err)
		// 	}
		// 	return
		// }

		result := json.FormatJSON(string(data))

		_, err = fmt.Fprintln(cmd.OutOrStdout(), result)
		if err != nil {
			log.Printf("ERROR: %s", err)
		}
	},
}

// TODO: implement setIndent flag
// TODO: implement EscapeHTML flag

func init() {
	rootCmd.AddCommand(jsonfmtCmd)
	// jsonfmtCmd.Flags().BoolVarP(&flagTUI, "tui", "t", false, "Show output in TUI")
}
