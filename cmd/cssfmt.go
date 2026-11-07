package cmd

import (
	"bytes"

	"github.com/client9/csstool"
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

			common := &ui.CommonModel{}

			model := css.NewCSSFormatterModel(common)
			err = model.SetContent(string(data))
			if err != nil {
				return err
			}
			return runProgram(model)
		}

		cssformat := csstool.NewCSSFormat(flagIndent, flagTab, nil)
		cssformat.AlwaysSemicolon = flagSemicolon

		data, err := input.ReadBytesFromArgsOrStdin(cmd, args)
		if err != nil {
			return err
		}

		// csstool reports neither parse nor write errors, so there is no parse error to explain.
		return cssformat.Format(bytes.NewReader(data), cmd.OutOrStdout())
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
