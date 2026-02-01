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
	Use:   "xmlfmt [string or file]",
	Short: "Format and prettify XML files",
	Long: `Format and prettify XML files with customizable indentation and formatting options.

By default, uses 2-space indentation. Customize with --indent, --prefix, and --nested flags.
Input can be a string argument or piped from stdin.`,
	Example: `  # Format XML from stdin
  devtui xmlfmt < document.xml
  cat unformatted.xml | devtui xmlfmt

  # Format XML string argument
  devtui xmlfmt '<root><item>value</item></root>'

  # Output to file
  devtui xmlfmt < input.xml > formatted.xml
  cat document.xml | devtui xmlfmt > pretty.xml

  # Custom indentation
  devtui xmlfmt --indent "    " < document.xml
  devtui xmlfmt -i "\t" < document.xml

  # Add prefix to each line
  devtui xmlfmt --prefix "  " < document.xml
  devtui xmlfmt -p "  " < document.xml

  # Handle nested tags in comments
  devtui xmlfmt --nested < document.xml
  devtui xmlfmt -n < document.xml

  # Show results in interactive TUI
  devtui xmlfmt --tui < document.xml
  devtui xmlfmt -t < document.xml

  # Chain with other commands
  curl -s https://example.com/feed.xml | devtui xmlfmt`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
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

		result := xmlfmt.FormatXML(string(data),
			xmlPrefix, xmlIndent, xmlNested)
		if outputJSON {
			return writeJSONValue(cmd.OutOrStdout(), result)
		}

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
