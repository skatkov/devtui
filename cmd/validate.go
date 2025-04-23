package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:     "validate",
	Short:   "Validate a license",
	Long:    "Validate a license",
	Example: "devtui validate",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Hey!")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
}
