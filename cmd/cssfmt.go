package cmd

import (
	"bytes"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/client9/csstool"
	"github.com/skatkov/devtui/internal/cmderror"
	"github.com/skatkov/devtui/internal/input"
	"github.com/skatkov/devtui/internal/ui"
	"github.com/skatkov/devtui/tui/css"
	"github.com/spf13/cobra"
)

var cssfmtCmd = &cobra.Command{
	Use:   "cssfmt [string or file]",
	Short: "Format and prettify CSS files",
	Long: `Format and prettify CSS files with customizable indentation and formatting options.

By default, uses 2-space indentation. Use --tab for tab indentation or --indent to specify
custom spacing. Input can be a string argument or piped from stdin.`,
	Example: `  # Format CSS from stdin
  devtui cssfmt < styles.css
  cat minified.css | devtui cssfmt

  # Format CSS string argument
  devtui cssfmt 'body{margin:0;padding:0}'

  # Use tab indentation
  devtui cssfmt --tab < styles.css
  devtui cssfmt -t < styles.css

  # Use custom indent spacing
  devtui cssfmt --indent 4 < styles.css
  devtui cssfmt -i 4 < styles.css

  # Output to file
  devtui cssfmt < input.css > formatted.css

  # Show results in interactive TUI
  devtui cssfmt --tui < styles.css

  # Chain with other commands
  curl -s https://example.com/styles.css | devtui cssfmt`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if flagTui {
			data, err := input.ReadBytesFromArgsOrStdin(cmd, args)
			if err != nil {
				return err
			}

			common := &ui.CommonModel{
				Width:  100, // Default width, will be adjusted by the TUI
				Height: 30,  // Default height, will be adjusted by the TUI
			}

			model := css.NewCSSFormatterModel(common)
			err = model.SetContent(string(data))
			if err != nil {
				return err
			}
			p := tea.NewProgram(
				model,
				tea.WithAltScreen(),       // Use alternate screen buffer
				tea.WithMouseCellMotion(), // Enable mouse support
			)

			if _, err := p.Run(); err != nil {
				return err
			}

			return nil
		}

		if flagTab {
			flagIndent = 1
		}
		cssformat := csstool.NewCSSFormat(flagIndent, flagTab, nil)
		cssformat.AlwaysSemicolon = flagSemicolon

		data, err := input.ReadBytesFromArgsOrStdin(cmd, args)
		if err != nil {
			return err
		}

		inputStr := string(data)
		var buffer bytes.Buffer
		output := cmd.OutOrStdout()
		if outputJSON {
			output = &buffer
		}
		err = cssformat.Format(strings.NewReader(inputStr), output)
		if err != nil {
			return cmderror.FormatParseError("cssfmt", inputStr, err)
		}
		if outputJSON {
			return writeJSONValue(cmd.OutOrStdout(), buffer.String())
		}

		return nil
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
