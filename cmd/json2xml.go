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

var json2xmlCmd = &cobra.Command{
	Use:   "json2xml [string or file]",
	Short: "Convert JSON to XML format",
	Long: `Convert JSON (JavaScript Object Notation) to XML (Extensible Markup Language) format.

Input can be a string argument or piped from stdin. Use --tui flag to view
results in an interactive terminal interface.`,
	Example: `  # Convert JSON from stdin
  devtui json2xml < data.json
  cat feed.json | devtui json2xml

  # Convert JSON string argument
  devtui json2xml '{"item": "value"}'

  # Output to file
  devtui json2xml < input.json > output.xml

  # Show results in interactive TUI
  devtui json2xml --tui < data.json

  # Chain with other commands
  curl -s https://api.example.com/data.json | devtui json2xml`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := input.ReadBytesFromArgsOrStdin(cmd, args)
		if err != nil {
			return fmt.Errorf("error reading from stdin: %w", err)
		}

		if len(data) == 0 {
			return errors.New("no input provided. Pipe JSON input to this command")
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
		result, err := xml.JSONToXML(inputStr)
		if err != nil {
			return cmderror.FormatParseError("json2xml", inputStr, err)
		}

		_, err = fmt.Fprintln(cmd.OutOrStdout(), result)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(json2xmlCmd)
	json2xmlCmd.Flags().BoolVarP(&flagTUI, "tui", "t", false, "Show output in TUI")
}
