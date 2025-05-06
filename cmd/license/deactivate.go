package license

import (
	"errors"
	"fmt"

	license "github.com/skatkov/devtui/internal/license"
	"github.com/spf13/cobra"
)

var DeactivateCmd = &cobra.Command{
	Use:     "deactivate",
	Short:   "Deactivate a license",
	Long:    "Deactivate a license",
	Example: "devtui license deactivate",
	RunE: func(cmd *cobra.Command, args []string) error {
		licenseData, err := license.LoadLicense()
		if err != nil {
			return fmt.Errorf("failed to load license: %w", err)
		}

		if licenseData == nil {
			return errors.New("no active license found")
		}

		err = licenseData.Deactivate()
		if err != nil {
			return fmt.Errorf("failed to deactivate license: %w", err)
		}

		fmt.Println("License deactivated")

		return nil
	},
}
