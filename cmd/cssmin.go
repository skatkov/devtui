package cmd

import (
	"bytes"
	"strings"

	"github.com/client9/csstool"
	"github.com/skatkov/devtui/internal/cmderror"
	"github.com/skatkov/devtui/internal/input"
	"github.com/spf13/cobra"
)

var cssminCmd = &cobra.Command{
	Use:   "cssmin [string or file]",
	Short: "Minify CSS files by removing whitespace and unnecessary characters",
	Long: `Minify CSS files by removing whitespace, line breaks, and unnecessary characters.

This reduces file size for production use while maintaining CSS functionality.
Input can be a string argument or piped from stdin.`,
	Example: `  # Minify CSS from stdin
  devtui cssmin < styles.css
  cat source.css | devtui cssmin

  # Minify CSS string argument
  devtui cssmin 'body { margin: 0; padding: 0; }'

  # Output to file
  devtui cssmin < input.css > minified.css
  cat styles.css | devtui cssmin > styles.min.css

  # Chain with other commands
  curl -s https://example.com/styles.css | devtui cssmin
  devtui cssfmt < messy.css | devtui cssmin`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cssformat := csstool.NewCSSFormat(0, false, nil)
		cssformat.AlwaysSemicolon = false

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
			return cmderror.FormatParseError("cssmin", inputStr, err)
		}
		if outputJSON {
			return writeJSONValue(cmd.OutOrStdout(), buffer.String())
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(cssminCmd)
}
