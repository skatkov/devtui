package cmd

import (
	"errors"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/skatkov/devtui/internal/cmderror"
	"github.com/skatkov/devtui/internal/input"
	"github.com/skatkov/devtui/internal/ui"
	"github.com/skatkov/devtui/tui/xml"
	"github.com/spf13/cobra"
)

var xml2jsonCmd = &cobra.Command{
	Use:   "xml2json [string or file]",
	Short: "Convert XML to JSON format",
	Long: `Convert XML (Extensible Markup Language) to JSON (JavaScript Object Notation) format.

Input can be a string argument or piped from stdin. Use --tui flag to view
results in an interactive terminal interface.`,
	Example: `  # Convert XML from stdin
  devtui xml2json < data.xml
  cat feed.xml | devtui xml2json

  # Convert XML string argument
  devtui xml2json '<root><item>value</item></root>'

  # Output to file
  devtui xml2json < input.xml > output.json

  # Show results in interactive TUI
  devtui xml2json --tui < data.xml

  # Chain with other commands
  curl -s https://api.example.com/data.xml | devtui xml2json`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := input.ReadBytesFromArgsOrStdin(cmd, args)
		if err != nil {
			return fmt.Errorf("error reading from stdin: %w", err)
		}

		if len(data) == 0 {
			return errors.New("no input provided. Pipe XML input to this command")
		}

		if flagTUI {
			common := &ui.CommonModel{
				Width:  0,
				Height: 0,
			}

			model := xml.NewXMLFormatterModel(common)
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
		result, err := xml.XMLToJSON(inputStr)
		if err != nil {
			return cmderror.FormatParseError("xml2json", inputStr, err)
		}

		_, err = fmt.Fprintln(cmd.OutOrStdout(), result)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(xml2jsonCmd)
	xml2jsonCmd.Flags().BoolVarP(&flagTUI, "tui", "t", false, "Show output in TUI")
}
