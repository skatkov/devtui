package license

import (
	"context"
	"errors"
	"fmt"
	"os"

	polargo "github.com/polarsource/polar-go"
	"github.com/polarsource/polar-go/models/components"
	"github.com/spf13/cobra"
)

var DeactivateCmd = &cobra.Command{
	Use:     "deactivate",
	Short:   "Deactivate a license",
	Long:    "Deactivate a license",
	Example: "devtui license deactivate",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		licenseData, err := loadLicenseData()
		if err != nil {
			return fmt.Errorf("failed to load license data: %w", err)
		}

		if licenseData == nil {
			return errors.New("no active license found")
		}
		s := polargo.New()

		res, err := s.CustomerPortal.LicenseKeys.Deactivate(ctx, components.LicenseKeyDeactivate{
			Key:            licenseData.LicenseKeyID,
			OrganizationID: OrganizationID,
			ActivationID:   licenseData.ActivationID,
		})
		if err != nil {
			return err
		}

		if res != nil {
			// Delete the license file after successful deactivation
			if err := os.Remove(LicenseFilePath); err != nil {
				fmt.Printf("Warning: License deactivated, but failed to remove license file: %v\n", err)
			}
			fmt.Println("License deactivated.")
		}
		return nil
	},
}
