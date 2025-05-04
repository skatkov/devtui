package license

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/adrg/xdg"
	polargo "github.com/polarsource/polar-go"
	"github.com/polarsource/polar-go/models/components"
	"github.com/skatkov/devtui/internal/macaddr"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	OrganizationID  = "afde3142-5d70-42e3-8214-71c5bbc04e6f"
	Salt            = "7f75d39f55414c9f9862dfc33163b598"
	RecheckInterval = 168 * time.Hour // weekly
)

type LicenseData struct {
	Hash           string    `json:"hash"`
	LicenseKeyID   string    `json:"license_key_id"`
	ActivationID   string    `json:"activation_id"`
	NextCheckTime  time.Time `json:"next_check_time"`
	LastVerifiedAt time.Time `json:"last_verified_at"`
}
var ActivateCmd = &cobra.Command{
	Use:     "activate",
	Short:   "Activate a license",
	Long:    "Activate a license",
	Example: "devtui license activate --key=YOUR_LICENSE_KEY",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		s := polargo.New()

		hostname, err := os.Hostname()
		if err != nil {
			hostname = "DevTUI"
		}

		tz, _ := time.Now().Zone()
		label := fmt.Sprintf("%s-%s", hostname, tz)
		licenseKey := viper.GetString("key")
		macAddress := macaddr.MacUint64()

		res, err := s.CustomerPortal.LicenseKeys.Activate(ctx, components.LicenseKeyActivate{
			Key:            licenseKey,
			OrganizationID: OrganizationID,
			Label:          label,
			Conditions: map[string]components.LicenseKeyActivateConditions{
				"macaddr": components.CreateLicenseKeyActivateConditionsInteger(int64(macAddress)),
			},
		})
		if err != nil {
			return err
		}
		if res.LicenseKeyActivationRead != nil {
			nextCheckTime := time.Now().Add(RecheckInterval)
			hash := createLicenseHash(licenseKey, macAddress, Salt, nextCheckTime)

			err = storeLicenseData(LicenseData{
				Hash:           hash,
				LicenseKeyID:   res.LicenseKeyActivationRead.LicenseKeyID,
				ActivationID:   res.LicenseKeyActivationRead.ID,
				NextCheckTime:  nextCheckTime,
				LastVerifiedAt: time.Now(),
			})
			if err != nil {
				return fmt.Errorf("failed to store license data: %w", err)
			}

			fmt.Println("License activated and stored successfully")
		}

		return nil
	},
}

// createLicenseHash creates a SHA-256 hash from the license key, MAC address, salt, and next check time
func createLicenseHash(licenseKey string, macAddress uint64, salt string, nextCheckTime time.Time) string {
	// Combine the input values into a single string
	data := fmt.Sprintf("%s:%d:%s:%d", licenseKey, macAddress, salt, nextCheckTime.Unix())

	// Create a SHA-256 hash
	hash := sha256.Sum256([]byte(data))

	// Convert the hash to a hexadecimal string
	return hex.EncodeToString(hash[:])
}

// storeLicenseData stores the license data in a JSON file in the XDG data directory
func storeLicenseData(data LicenseData) error {
	// Create directories if they don't exist
	dataDir := filepath.Join(xdg.DataHome, "devtui")
	if err := os.MkdirAll(dataDir, 0700); err != nil {
		return err
	}

	// Convert license data to JSON
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	// Write to file
	licenseFilePath := filepath.Join(dataDir, "license.json")
	return os.WriteFile(licenseFilePath, jsonData, 0600)
}
