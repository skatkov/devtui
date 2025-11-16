package cmd

import (
	"strings"

	"github.com/client9/csstool"
	"github.com/skatkov/devtui/internal/input"
	"github.com/spf13/cobra"
)

var cssminCmd = &cobra.Command{
	Use:     "cssmin",
	Short:   "Minify CSS files",
	Long:    "Minify CSS files",
	Example: "cssmin < testdata/bootstrap.min.css",
	Args:    cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cssformat := csstool.NewCSSFormat(0, false, nil)
		cssformat.AlwaysSemicolon = false

		data, err := input.ReadBytesFromArgsOrStdin(cmd, args)
		if err != nil {
			return err
		}

		err = cssformat.Format(strings.NewReader(string(data)), cmd.OutOrStdout())
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(cssminCmd)
}
