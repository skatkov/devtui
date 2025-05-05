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
	"github.com/spf13/viper"
)

var ActivateCmd = &cobra.Command{
	Use:     "activate",
	Short:   "Activate a license",
	Long:    "Activate a license",
	Example: "devtui license activate --key=YOUR_LICENSE_KEY",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		s := polargo.New()

		hostname, err := os.Hostname()
		if err != nil {
			hostname = "DevTUI"
		}

		tz, _ := time.Now().Zone()
		label := fmt.Sprintf("%s-%s", hostname, tz)
		licenseKey := viper.GetString("key")
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
			nextCheckTime := time.Now().Add(RecheckInterval)
			hash := createLicenseHash(licenseKey, macAddress, Salt, nextCheckTime)

			err = storeLicenseData(LicenseData{
				Hash:           hash,
				LicenseKeyID:   licenseKey,
				ActivationID:   res.LicenseKeyActivationRead.ID,
				NextCheckTime:  nextCheckTime,
				LastVerifiedAt: time.Now(),
			})
			if err != nil {
				return fmt.Errorf("failed to store license data: %w", err)
			}

			fmt.Println("\nLicense activated and stored successfully")
		}

		return nil
	},
}
