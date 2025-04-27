package cmd

import (
	"context"
	"fmt"
	"log"

	polargo "github.com/polarsource/polar-go"
	"github.com/polarsource/polar-go/models/components"
	"github.com/spf13/cobra"
)

var deactivateCmd = &cobra.Command{
	Use:     "deactivate",
	Short:   "Deactivate a license",
	Long:    "Deactivate a license",
	Example: "devtui deactivate",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		s := polargo.New()

		res, err := s.CustomerPortal.LicenseKeys.Deactivate(ctx, components.LicenseKeyDeactivate{
			Key:            "DEVTUI-2CA57A34-E191-4290-A394-7B107BCA1036",
			OrganizationID: "afde3142-5d70-42e3-8214-71c5bbc04e6f",
			ActivationID:   "14c674ae-2e18-4997-9326-574b2c1cb280",
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

func init() {
	rootCmd.AddCommand(deactivateCmd)
}
