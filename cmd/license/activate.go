package license

import (
	"errors"
	"fmt"
	"os"
	"time"

	license "github.com/skatkov/devtui/internal/license"
	"github.com/spf13/cobra"
)

var licenseKey string

var ActivateCmd = &cobra.Command{
	Use:     "activate",
	Short:   "Activate a license",
	Long:    "Activate a license",
	Example: "devtui license activate --key=YOUR_LICENSE_KEY",
	RunE: func(cmd *cobra.Command, args []string) error {
		if licenseKey == "" {
			return errors.New("license key is required")
		}

		hostname, err := os.Hostname()
		if err != nil {
			hostname = "DevTUI"
		}

		tz, _ := time.Now().Zone()
		label := fmt.Sprintf("%s-%s", hostname, tz)

		licenseData, err := license.NewLicense(licenseKey, label)
		if err != nil {
			return fmt.Errorf("failed to activate a license: %w", err)
		}

		fmt.Printf("Activation ID: %s\n", licenseData.ActivationID)
		fmt.Printf("License Key ID: %s\n", licenseData.KeyID)
		fmt.Printf("Label: %s\n", label)
		fmt.Printf("Created At: %s\n", licenseData.VerifiedAt)
		fmt.Println("\nLicense activated and stored successfully")

		return nil
	},
}

func init() {
	ActivateCmd.Flags().StringVar(&licenseKey, "key", "", "License key")
}
