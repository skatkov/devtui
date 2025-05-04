package license

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
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

		licenseData, err := loadLicenseData()
		if err != nil {
			return fmt.Errorf("failed to load license data: %w", err)
		}

		if licenseData == nil {
			return fmt.Errorf("no active license found")
		}
    fmt.Println("Deactivating license...")
		s := polargo.New()

		res, err := s.CustomerPortal.LicenseKeys.Deactivate(ctx, components.LicenseKeyDeactivate{
			Key:            licenseData.LicenseKeyID,
			OrganizationID: OrganizationID,
			ActivationID:   licenseData.ActivationID,
		})
		if err != nil {
			return err
		}

		fmt.Println("License deactivated successfully")

		if res != nil {
			// Delete the license file after successful deactivation
			if err := removeLicenseFile(); err != nil {
				fmt.Printf("Warning: License deactivated but failed to remove license file: %v\n", err)
			}
			fmt.Println("Deactivation completed successfully")
		}
		return nil
	},
}

// loadLicenseData loads the license data from the license.json file
func loadLicenseData() (*LicenseData, error) {
	dataDir := filepath.Join(xdg.DataHome, "devtui")
	licenseFilePath := filepath.Join(dataDir, "license.json")

	// Check if the file exists
	if _, err := os.Stat(licenseFilePath); os.IsNotExist(err) {
		return nil, nil // Return nil if the file doesn't exist
	}

	// Read the file
	data, err := os.ReadFile(licenseFilePath)
	if err != nil {
		return nil, err
	}

	// Parse the JSON data
	var licenseData LicenseData
	if err := json.Unmarshal(data, &licenseData); err != nil {
		return nil, err
	}

	return &licenseData, nil
}

// removeLicenseFile deletes the license.json file
func removeLicenseFile() error {
	dataDir := filepath.Join(xdg.DataHome, "devtui")
	licenseFilePath := filepath.Join(dataDir, "license.json")
	return os.Remove(licenseFilePath)
}
