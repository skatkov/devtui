package license

import (
	"context"
	"fmt"
	"log"

	polargo "github.com/polarsource/polar-go"
	"github.com/polarsource/polar-go/models/components"
	parentcmd "github.com/skatkov/devtui/cmd"
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

		res, err := s.CustomerPortal.LicenseKeys.Deactivate(ctx, components.LicenseKeyDeactivate{
			Key:            parentcmd.GetLicenseKey(),
			OrganizationID: parentcmd.GetOrgID(),
			ActivationID:   parentcmd.GetActivationID(),
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
