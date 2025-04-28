package license

import (
	"context"
	"fmt"
	"log"

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

		s := polargo.New()

		key, err := cmd.Flags().GetString("key")
		if err != nil {
			return err
		}

		id, err := cmd.Flags().GetString("id")
		if err != nil {
			return err
		}

		res, err := s.CustomerPortal.LicenseKeys.Deactivate(ctx, components.LicenseKeyDeactivate{
			Key:            key,
			OrganizationID: "afde3142-5d70-42e3-8214-71c5bbc04e6f",
			ActivationID:   id,
		})
		if err != nil {
			log.Fatal(err)
		}
		if res != nil {
			fmt.Println("Deactivation completed successfully")
		}
		return nil
	},
}
