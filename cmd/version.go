package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// These variables are populated by goreleaser during build.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// GetVersionShort returns a short version string suitable for single-line output.
func GetVersionShort() string {
	return fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  "Print version, commit and date of release for this software",
	RunE: func(_ *cobra.Command, _ []string) error {
		fmt.Print("devtui version " + GetVersionShort())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
