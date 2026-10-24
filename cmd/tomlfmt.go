package cmd

import (
	"fmt"

	"github.com/skatkov/devtui/internal/cmderror"
	"github.com/skatkov/devtui/internal/input"
	"github.com/skatkov/devtui/internal/ui"
	"github.com/skatkov/devtui/tui/toml"
	"github.com/spf13/cobra"
)

var tomlfmtCmd = &cobra.Command{
	Use:   "tomlfmt [string or file]",
	Short: "Format and prettify TOML files",
	Long: `Format and prettify TOML files with proper indentation.

Input can be a string argument or piped from stdin. Use --tui flag to view results
in an interactive terminal interface.`,
	Example: `  # Format TOML from stdin
  devtui tomlfmt < config.toml
  cat app.toml | devtui tomlfmt

  # Format TOML string argument
  devtui tomlfmt '[package]\nname = "myapp"'

  # Output to file
  devtui tomlfmt < input.toml > formatted.toml
  cat config.toml | devtui tomlfmt > pretty.toml

  # Show results in interactive TUI
  devtui tomlfmt --tui < config.toml
  devtui tomlfmt -t < config.toml

  # Chain with other commands
  curl -s https://example.com/config.toml | devtui tomlfmt`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := input.ReadBytesFromArgsOrStdin(cmd, args)
		if err != nil {
			return err
		}

		if tomlfmtTUI {
			common := &ui.CommonModel{}

			model := toml.NewTomlFormatModel(common)
			err = model.SetContent(string(data))
			if err != nil {
				return err
			}

			return runProgram(model)
		}

		inputStr := string(data)
		result, err := toml.Convert(inputStr)
		if err != nil {
			return cmderror.FormatParseError("tomlfmt", inputStr, err)
		}

		_, err = fmt.Fprintln(cmd.OutOrStdout(), result)
		return err
	},
}

var tomlfmtTUI bool

func init() {
	rootCmd.AddCommand(tomlfmtCmd)
	tomlfmtCmd.Flags().BoolVarP(&tomlfmtTUI, "tui", "t", false, "Show output in TUI")
}
