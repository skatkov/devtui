package cmd

import (
	"github.com/skatkov/devtui/cmd/license"
	"github.com/spf13/cobra"
)

var LicenseCmd = &cobra.Command{
	Use:   "license",
	Short: "License management commands",
	Long:  "Commands for activating, validating, and deactivating licenses",
}

func init() {
	rootCmd.AddCommand(LicenseCmd)
	LicenseCmd.AddCommand(license.ActivateCmd)
	LicenseCmd.AddCommand(license.DeactivateCmd)
	LicenseCmd.AddCommand(license.ValidateCmd)
}
