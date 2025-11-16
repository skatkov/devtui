package cmd

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-xmlfmt/xmlfmt"
	"github.com/skatkov/devtui/internal/input"
	"github.com/skatkov/devtui/internal/ui"
	"github.com/skatkov/devtui/tui/xml"
	"github.com/spf13/cobra"
)

var xmlfmtCmd = &cobra.Command{
	Use:   "xmlfmt",
	Short: "Format XML",
	Long:  "Format XML",
	Example: `
	xmlfmt < testdata/sample.xml   # Format XML from stdin
	xmlfmt < testdata/sample.xml > output.xml # Output formatted XML to file
	xmlfmt < testdata/sample.xml --tui # Open XML formatter in TUI
	`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Read all input data from stdin or args
		data, err := input.ReadBytesFromArgsOrStdin(cmd, args)
		if err != nil {
			return err
		}

		if flagTUI {
			common := &ui.CommonModel{
				Width:  80,
				Height: 80,
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

		// Process the XML
		result := xmlfmt.FormatXML(string(data),
			xmlPrefix, xmlIndent, xmlNested)

		_, err = fmt.Fprintln(cmd.OutOrStdout(), result)
		if err != nil {
			return err
		}

		return nil
	},
}

var (
	xmlPrefix string
	xmlIndent string
	xmlNested bool
)

func init() {
	rootCmd.AddCommand(xmlfmtCmd)
	xmlfmtCmd.Flags().StringVarP(&xmlPrefix, "prefix", "p", "", "Each element begins on a new line and this prefix")
	xmlfmtCmd.Flags().StringVarP(&xmlIndent, "indent", "i", "  ", "Indent string for nested elements")
	xmlfmtCmd.Flags().BoolVarP(&xmlNested, "nested", "n", false, "Nested tags in comments")
	xmlfmtCmd.Flags().BoolVarP(&flagTUI, "tui", "t", false, "Show output in TUI")
}
