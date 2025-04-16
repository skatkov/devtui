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

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  "Print version, commit and date of release for this software",
	RunE: func(_ *cobra.Command, _ []string) error {
		versionInfo := fmt.Sprintf("Version: %s\nCommit:  %s\nDate:    %s", version, commit, date)
		fmt.Print(versionInfo)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
