package license

import (
	"context"
	"fmt"
	"os"
	"time"

	polargo "github.com/polarsource/polar-go"
	"github.com/polarsource/polar-go/models/components"
	"github.com/spf13/cobra"
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

		key, err := cmd.Flags().GetString("key")
		if err != nil {
			return err
		}

		res, err := s.CustomerPortal.LicenseKeys.Activate(ctx, components.LicenseKeyActivate{
			Key:            key,
			OrganizationID: "afde3142-5d70-42e3-8214-71c5bbc04e6f",
			Label:          label,
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
		}
		return nil
	},
}
