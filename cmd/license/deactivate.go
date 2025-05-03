package license

import (
	"context"
	"fmt"

	polargo "github.com/polarsource/polar-go"
	"github.com/polarsource/polar-go/models/components"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var DeactivateCmd = &cobra.Command{
	Use:     "deactivate",
	Short:   "Deactivate a license",
	Long:    "Deactivate a license",
	Example: "devtui license deactivate",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		s := polargo.New()

		res, err := s.CustomerPortal.LicenseKeys.Deactivate(ctx, components.LicenseKeyDeactivate{
			Key:            viper.GetString("key"),
			OrganizationID: OrganizationID,
			ActivationID:   viper.GetString("id"),
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
