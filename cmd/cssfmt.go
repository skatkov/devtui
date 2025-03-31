package cmd

import (
	"log"
	"os"

	"github.com/client9/csstool"
	"github.com/spf13/cobra"
)

// ---
// Example usage:
// go run . cssfmt < testdata/bootstrap.min.css
// ---

var cssfmtCmd = &cobra.Command{
	Use:   "cssfmt",
	Short: "Format CSS files",
	Long:  "Format CSS files",
	Run: func(cmd *cobra.Command, args []string) {
		if flagTab {
			flagIndent = 1
		}
		cssformat := csstool.NewCSSFormat(flagIndent, flagTab, nil)
		cssformat.AlwaysSemicolon = flagSemicolon
		err := cssformat.Format(os.Stdin, os.Stdout)
		if err != nil {
			log.Printf("ERROR: %s", err)
		}
	},
}

var (
	flagTab       bool
	flagIndent    int
	flagSemicolon bool
)

func init() {
	rootCmd.AddCommand(cssfmtCmd)
	cssfmtCmd.Flags().BoolVarP(&flagTab, "tab", "t", false, "use tabs for indentation")
	cssfmtCmd.Flags().IntVarP(&flagIndent, "indent", "i", 2, "spaces for indentation")
	cssfmtCmd.Flags().BoolVarP(&flagSemicolon, "semicolon", "", true, "always end rule with semicolon, even if not needed")
}
