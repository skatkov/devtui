package license

import (
	"context"
	"errors"
	"fmt"

	polargo "github.com/polarsource/polar-go"
	"github.com/polarsource/polar-go/models/components"
	"github.com/skatkov/devtui/internal/macaddr"
	"github.com/spf13/cobra"
)

var ValidateCmd = &cobra.Command{
	Use:     "validate",
	Short:   "Validate a license",
	Long:    "Reads a license and validates it",
	Example: "devtui license validate",
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

		res, err := s.CustomerPortal.LicenseKeys.Validate(ctx, components.LicenseKeyValidate{
			Key:            licenseData.LicenseKeyID,
			OrganizationID: OrganizationID,
			ActivationID:   polargo.String(licenseData.ActivationID),
			Conditions: map[string]components.Conditions{
				"macaddr": components.CreateConditionsInteger(int64(macaddr.MacUint64())),
			},
		})
		if err != nil {
			// Activation ID is wrong
			// Error: {"detail":[{"loc":["body","activation_id"],"msg":"Input should be a valid UUID, invalid group length in group 0: expected 8, found 5","type":"uuid_parsing"}]}

			// MacAddress is wrong or not provided.
			// {"error":"ResourceNotFound","detail":"License key does not match required conditions"}

			return err
		}
		if res.ValidatedLicenseKey != nil {
			fmt.Println("\nLicense valid.")
		}
		return nil
	},
}
