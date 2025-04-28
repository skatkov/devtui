package license

import (
	"context"
	"fmt"

	polargo "github.com/polarsource/polar-go"
	"github.com/polarsource/polar-go/models/components"
	"github.com/spf13/cobra"
)

var ValidateCmd = &cobra.Command{
	Use:     "validate",
	Short:   "Validate a license",
	Long:    "Validate a license",
	Example: "devtui license validate",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		s := polargo.New()

		key, err := cmd.Flags().GetString("key")
		if err != nil {
			return err
		}

		id, err := cmd.Flags().GetString("id")
		if err != nil {
			return err
		}

		res, err := s.CustomerPortal.LicenseKeys.Validate(ctx, components.LicenseKeyValidate{
			Key:            key,
			OrganizationID: OrganizationID,
			ActivationID:   polargo.String(id),
		})
		if err != nil {
			return err
		}
		if res.ValidatedLicenseKey != nil {
			fmt.Printf("ID: %s\n", res.ValidatedLicenseKey.ID)
			// fmt.Printf("Organization ID: %s\n", res.ValidatedLicenseKey.OrganizationID)
			fmt.Printf("Customer ID: %s\n", res.ValidatedLicenseKey.CustomerID)
			fmt.Printf("Customer: %+v\n", res.ValidatedLicenseKey.Customer)
			fmt.Printf("Benefit ID: %s\n", res.ValidatedLicenseKey.BenefitID)
			fmt.Printf("Key: %s\n", res.ValidatedLicenseKey.Key)
			fmt.Printf("Display Key: %s\n", res.ValidatedLicenseKey.DisplayKey)
			fmt.Printf("Status: %s\n", res.ValidatedLicenseKey.Status)
			fmt.Printf("Limit Activations: %v\n", res.ValidatedLicenseKey.LimitActivations)
			fmt.Printf("Usage: %d\n", res.ValidatedLicenseKey.Usage)
			fmt.Printf("Limit Usage: %v\n", res.ValidatedLicenseKey.LimitUsage)
			fmt.Printf("Validations: %d\n", res.ValidatedLicenseKey.Validations)
			fmt.Printf("Last Validated At: %v\n", res.ValidatedLicenseKey.LastValidatedAt)
			fmt.Printf("Expires At: %v\n", res.ValidatedLicenseKey.ExpiresAt)
			fmt.Printf("Activation: %+v\n", res.ValidatedLicenseKey.Activation)
		}
		return nil
	},
}
