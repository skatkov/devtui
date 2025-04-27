package license

import (
	"context"
	"fmt"
	"log"

	polargo "github.com/polarsource/polar-go"
	"github.com/polarsource/polar-go/models/components"
	"github.com/spf13/cobra"
)

var ValidateCmd = &cobra.Command{
	Use:     "validate",
	Short:   "Validate a license",
	Long:    "Validate a license",
	Example: "devtui validate",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		s := polargo.New()

		res, err := s.CustomerPortal.LicenseKeys.Validate(ctx, components.LicenseKeyValidate{
			Key:            "DEVTUI-2CA57A34-E191-4290-A394-7B107BCA1036",
			OrganizationID: "afde3142-5d70-42e3-8214-71c5bbc04e6f",
			ActivationID:   polargo.String("802c9aa5-9156-4b48-9a6e-3e716c335955"),
		})
		if err != nil {
			log.Fatal(err)
		}
		if res.ValidatedLicenseKey != nil {
			fmt.Printf("ID: %s\n", res.ValidatedLicenseKey.ID)
			fmt.Printf("Organization ID: %s\n", res.ValidatedLicenseKey.OrganizationID)
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
