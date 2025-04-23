package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var activateCmd = &cobra.Command{
	Use:     "activate",
	Short:   "Activate a license",
	Long:    "Activate a license",
	Example: "devtui activate",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Hey!")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(activateCmd)
}
