package license

import (
	"errors"
	"fmt"

	license "github.com/skatkov/devtui/internal/license"
	"github.com/spf13/cobra"
)

var ValidateCmd = &cobra.Command{
	Use:     "validate",
	Short:   "Validate a license",
	Long:    "Reads a license and validates it",
	Example: "devtui license validate",
	RunE: func(cmd *cobra.Command, args []string) error {
		licenseData, err := license.LoadLicense()
		if err != nil {
			return fmt.Errorf("failed to load license: %w", err)
		}

		if licenseData == nil {
			return errors.New("no active license found")
		}

		if err := licenseData.Validate(); err != nil {
			return errors.New("license file is not valid")
		}

		return nil
	},
}
