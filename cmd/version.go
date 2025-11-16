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

// GetVersionString returns a formatted version string with all version information.
func GetVersionString() string {
	return fmt.Sprintf("Version: %s\nCommit:  %s\nDate:    %s", version, commit, date)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  "Print version, commit and date of release for this software",
	RunE: func(_ *cobra.Command, _ []string) error {
		fmt.Print(GetVersionString())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
