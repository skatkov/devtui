package cmd

import (
	"os"

	"github.com/skatkov/devtui/cmd/license"
	"github.com/spf13/cobra"
)

var LicenseCmd = &cobra.Command{
	Use:   "license",
	Short: "License management commands",
	Long:  "Commands for activating, validating, and deactivating licenses",
}

var (
	licenseKey   string
	activationID string
	orgID        string = "afde3142-5d70-42e3-8214-71c5bbc04e6f"
)

func init() {
	rootCmd.AddCommand(LicenseCmd)
	LicenseCmd.AddCommand(license.ActivateCmd)
	LicenseCmd.AddCommand(license.DeactivateCmd)
	LicenseCmd.AddCommand(license.ValidateCmd)

	LicenseCmd.PersistentFlags().StringVar(&licenseKey, "key", "", "License key to activate")
	LicenseCmd.PersistentFlags().StringVar(&activationID, "id", "", "License activation ID")
}

func GetLicenseKey() string {
	if licenseKey == "" {
		envKey := os.Getenv("DEVTUI_KEY")
		if envKey != "" {
			licenseKey = envKey
		}
	}
	return licenseKey
}

func GetOrgID() string {
	if orgID == "" {
		envKey := os.Getenv("DEVTUI_ORG_ID")
		if envKey != "" {
			orgID = envKey
		}
	}
	return orgID
}

func GetActivationID() string {
	if activationID == "" {
		envKey := os.Getenv("DEVTUI_ID")
		if envKey != "" {
			activationID = envKey
		}
	}
	return activationID
}
