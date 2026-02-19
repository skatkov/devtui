package cmd

import (
	"encoding/json"
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
	Example: `devtui version
	devtui version --json`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		if flagJSON {
			versionInfo := map[string]string{
				"version": version,
				"commit":  commit,
				"date":    date,
			}
			jsonBytes, err := json.Marshal(versionInfo)
			if err != nil {
				return err
			}
			fmt.Println(string(jsonBytes))
			return nil
		}

		versionInfo := fmt.Sprintf("Version: %s\nCommit:  %s\nDate:    %s", version, commit, date)
		fmt.Print(versionInfo)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
