package license

import (
	"context"
	"fmt"

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
			OrganizationID: OrganizationID,
			ActivationID:   id,
		})
		if err != nil {
			return err
		}
		if res != nil {
			fmt.Println("Deactivation completed successfully")
		}
		return nil
	},
}
