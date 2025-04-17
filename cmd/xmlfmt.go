package cmd

import (
	"io"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-xmlfmt/xmlfmt"
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
	Run: func(cmd *cobra.Command, args []string) {
		// Read all input data from stdin
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

			model := xml.NewXMLFormatterModel(common)
			model.SetContent(string(data))

			p := tea.NewProgram(model, tea.WithAltScreen())
			if _, err := p.Run(); err != nil {
				log.Printf("ERROR: %s", err)
			}
			return
		}

		// Process the XML
		result := xmlfmt.FormatXML(string(data),
			xmlPrefix, xmlIndent, xmlNested)

		// Write the result to stdout
		_, err = os.Stdout.WriteString(result)
		if err != nil {
			log.Printf("ERROR: %s", err)
		}
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
