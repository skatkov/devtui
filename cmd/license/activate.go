package license

import (
	"context"
	"fmt"
	"os"
	"time"

	polargo "github.com/polarsource/polar-go"
	"github.com/polarsource/polar-go/models/components"
	"github.com/skatkov/devtui/internal/macaddr"
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
			return fmt.Errorf("license key is required")
		}

		ctx := context.Background()

		s := polargo.New()

		hostname, err := os.Hostname()
		if err != nil {
			hostname = "DevTUI"
		}

		tz, _ := time.Now().Zone()
		label := fmt.Sprintf("%s-%s", hostname, tz)
		macAddress := macaddr.MacUint64()

		res, err := s.CustomerPortal.LicenseKeys.Activate(ctx, components.LicenseKeyActivate{
			Key:            licenseKey,
			OrganizationID: OrganizationID,
			Label:          label,
			Conditions: map[string]components.LicenseKeyActivateConditions{
				"macaddr": components.CreateLicenseKeyActivateConditionsInteger(int64(macAddress)),
			},
		})
		if err != nil {
			return err
		}
		if res.LicenseKeyActivationRead != nil {
			fmt.Printf("ID: %s\n", res.LicenseKeyActivationRead.ID)
			fmt.Printf("License Key ID: %s\n", res.LicenseKeyActivationRead.LicenseKeyID)
			fmt.Printf("Label: %s\n", res.LicenseKeyActivationRead.Label)
			fmt.Printf("Created At: %s\n", res.LicenseKeyActivationRead.CreatedAt)
			if res.LicenseKeyActivationRead.ModifiedAt != nil {
				fmt.Printf("Modified At: %s\n", *res.LicenseKeyActivationRead.ModifiedAt)
			}

			err = storeLicense(LicenseData{
				LicenseKeyID: licenseKey,
				ActivationID: res.LicenseKeyActivationRead.ID,
				VerifiedAt:   time.Now(),
			}, macAddress)
			if err != nil {
				return fmt.Errorf("failed to store license data: %w", err)
			}

			fmt.Println("\nLicense activated and stored successfully")
		}

		return nil
	},
}

func init() {
	ActivateCmd.Flags().StringVar(&licenseKey, "key", "", "License key")
}
