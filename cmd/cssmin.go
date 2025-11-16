package cmd

import (
	"github.com/client9/csstool"
	"github.com/spf13/cobra"
)

var cssminCmd = &cobra.Command{
	Use:     "cssmin",
	Short:   "Minify CSS files",
	Long:    "Minify CSS files",
	Example: "cssmin < testdata/bootstrap.min.css",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cssformat := csstool.NewCSSFormat(0, false, nil)
		cssformat.AlwaysSemicolon = false
		err := cssformat.Format(cmd.InOrStdin(), cmd.OutOrStdout())
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(cssminCmd)
}
