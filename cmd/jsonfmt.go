package cmd

import (
	"fmt"

	"github.com/skatkov/devtui/internal/input"
	"github.com/skatkov/devtui/tui/json"
	"github.com/spf13/cobra"
)

var jsonfmtCmd = &cobra.Command{
	Use:   "jsonfmt [string or file]",
	Short: "Format and prettify JSON",
	Long: `Format and prettify JSON input with proper indentation and syntax highlighting.

Input can be a string argument, piped from stdin, or read from a file.
The output is always valid, properly indented JSON.`,
	Example: `  # Format JSON from stdin
  devtui jsonfmt < example.json
  echo '{"name":"John","age":30}' | devtui jsonfmt

  # Format JSON string argument
  devtui jsonfmt '{"name":"John","age":30}'

  # Output to file
  devtui jsonfmt < input.json > formatted.json
  cat compact.json | devtui jsonfmt > pretty.json

  # Chain with other commands
  curl -s https://api.example.com/data | devtui jsonfmt
  devtui jsonrepair < broken.json | devtui jsonfmt`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := input.ReadBytesFromArgsOrStdin(cmd, args)
		if err != nil {
			return err
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
			return err
		}

		return nil
	},
}

// TODO: implement setIndent flag
// TODO: implement EscapeHTML flag

func init() {
	rootCmd.AddCommand(jsonfmtCmd)
	// jsonfmtCmd.Flags().BoolVarP(&flagTUI, "tui", "t", false, "Show output in TUI")
}
