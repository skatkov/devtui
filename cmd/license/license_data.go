package license

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/adrg/xdg"
)

const (
	OrganizationID  = "afde3142-5d70-42e3-8214-71c5bbc04e6f"
	Salt            = "7f75d39f55414c9f9862dfc33163b598"
	RecheckInterval = 168 * time.Hour // weekly
)

var (
	DataDir         = filepath.Join(xdg.DataHome, "devtui")
	LicenseFilePath = filepath.Join(DataDir, "license.json")
)

type LicenseData struct {
	Hash          string    `json:"hash" validate:"required"`
	LicenseKeyID  string    `json:"license_key_id" validate:"required"`
	ActivationID  string    `json:"activation_id" validate:"required,uuid"`
	NextCheckTime time.Time `json:"next_check_time" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	VerifiedAt    time.Time `json:"last_verified_at" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
}

func createLicenseHash(activationID string, licenseKey string, macAddress uint64, nextCheckTime time.Time) string {
	// Combine the input values into a single string
	data := fmt.Sprintf("%s:%s:%d:%s:%d", activationID, licenseKey, macAddress, Salt, nextCheckTime.Unix())

	// Create a SHA-256 hash
	hash := sha256.Sum256([]byte(data))

	// Convert the hash to a hexadecimal string
	return hex.EncodeToString(hash[:])
}

// storeLicenseData stores the license data in a JSON file in the XDG data directory
func storeLicense(data LicenseData, macAddress uint64) error {
	nextCheckTime := time.Now().Add(RecheckInterval)
	hash := createLicenseHash(data.ActivationID, data.LicenseKeyID, macAddress, nextCheckTime)

	data.NextCheckTime = nextCheckTime
	data.Hash = hash

	if err := os.MkdirAll(DataDir, 0o700); err != nil {
		return err
	}

	// Convert license data to JSON
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(LicenseFilePath, jsonData, 0o600)
}

// loadLicenseData loads the license data from the license.json file
func loadLicenseData() (*LicenseData, error) {
	if _, err := os.Stat(LicenseFilePath); os.IsNotExist(err) {
		return nil, nil // Return nil if the file doesn't exist
	}

	// Read the file
	data, err := os.ReadFile(LicenseFilePath)
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
