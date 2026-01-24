package cmd

import (
	"errors"
	"fmt"

	"github.com/skatkov/devtui/internal/htmlfmt"
	"github.com/skatkov/devtui/internal/input"
	"github.com/spf13/cobra"
)

var htmlfmtCmd = &cobra.Command{
	Use:   "htmlfmt [string or file]",
	Short: "Format and prettify HTML",
	Long: `Format and prettify HTML input with consistent indentation.

Input can be a string argument, piped from stdin, or read from a file.`,
	Example: `  # Format HTML from stdin
  devtui htmlfmt < page.html
  cat page.html | devtui htmlfmt

  # Format HTML string argument
  devtui htmlfmt '<div><span>hello</span></div>'

  # Output to file
  devtui htmlfmt < input.html > formatted.html`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := input.ReadBytesFromArgsOrStdin(cmd, args)
		if err != nil {
			return fmt.Errorf("error reading from stdin: %w", err)
		}
		if len(data) == 0 {
			return errors.New("no input provided. pipe HTML input to this command")
		}

		result := htmlfmt.Format(string(data))
		_, err = fmt.Fprintln(cmd.OutOrStdout(), result)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(htmlfmtCmd)
}
