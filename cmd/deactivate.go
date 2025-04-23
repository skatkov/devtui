package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var deactivateCmd = &cobra.Command{
	Use:     "deactivate",
	Short:   "Deactivate a license",
	Long:    "Deactivate a license",
	Example: "devtui deactivate",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Hey!")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deactivateCmd)
}
