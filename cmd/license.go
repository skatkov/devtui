package cmd

import (
	"github.com/skatkov/devtui/cmd/license"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var LicenseCmd = &cobra.Command{
	Use:   "license",
	Short: "License management commands",
	Long:  "Commands for activating, validating, and deactivating licenses",
}

var (
	licenseKey   string
	activationID string
)

func init() {
	rootCmd.AddCommand(LicenseCmd)
	LicenseCmd.AddCommand(license.ActivateCmd)
	LicenseCmd.AddCommand(license.DeactivateCmd)
	LicenseCmd.AddCommand(license.ValidateCmd)

	LicenseCmd.PersistentFlags().StringVar(&licenseKey, "key", "", "License key")
	err := LicenseCmd.MarkFlagRequired("key")
	if err != nil {
		panic(err)
	}

	err = viper.BindPFlag("key", LicenseCmd.PersistentFlags().Lookup("key"))
	if err != nil {
		panic(err)
	}

	LicenseCmd.PersistentFlags().StringVar(&activationID, "id", "", "License activation ID")
	err = viper.BindPFlag("id", LicenseCmd.PersistentFlags().Lookup("id"))
	if err != nil {
		panic(err)
	}

}
