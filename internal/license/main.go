package license

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/adrg/xdg"
	polargo "github.com/polarsource/polar-go"
	"github.com/polarsource/polar-go/models/components"
	"github.com/skatkov/devtui/internal/macaddr"
)

const (
	organizationID  = "afde3142-5d70-42e3-8214-71c5bbc04e6f"
	salt            = "7f75d39f55414c9f9862dfc33163b598"
	recheckInterval = 168 * time.Hour // weekly
)

var (
	dataDir  = filepath.Join(xdg.DataHome, "devtui")
	filePath = filepath.Join(dataDir, "license.json")
)

type LicenseData struct {
	Hash          string    `json:"hash" validate:"required"`
	KeyID         string    `json:"key_id" validate:"required"`
	ActivationID  string    `json:"activation_id" validate:"required,uuid"`
	NextCheckTime time.Time `json:"next_check_time" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	VerifiedAt    time.Time `json:"last_verified_at" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
}

func NewLicense(licenseKey string, label string) (*LicenseData, error) {
	ctx := context.Background()
	s := polargo.New()
	macAddress := macaddr.MacUint64()

	res, err := s.CustomerPortal.LicenseKeys.Activate(ctx, components.LicenseKeyActivate{
		Key:            licenseKey,
		OrganizationID: organizationID,
		Label:          label,
		Conditions: map[string]components.LicenseKeyActivateConditions{
			"macaddr": components.CreateLicenseKeyActivateConditionsInteger(int64(macAddress)),
		},
	})
	if err != nil {
		return nil, err
	}

	data := &LicenseData{
		KeyID:        licenseKey,
		ActivationID: res.LicenseKeyActivationRead.ID,
		VerifiedAt:   time.Now(),
	}

	if err := data.store(macAddress); err != nil {
		return data, err
	}

	return data, nil
}

func LoadLicense() (*LicenseData, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, nil // Return nil if the file doesn't exist
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var licenseData LicenseData
	if err := json.Unmarshal(data, &licenseData); err != nil {
		return nil, err
	}

	return &licenseData, nil
}

func (d LicenseData) Validate() error {
	currentHash := d.buildHash(macaddr.MacUint64())
	if currentHash != d.Hash {
		return errors.New("license file is not valid")
	}

	// Skip API validation if we haven't reached the next check time
	if time.Now().Before(d.NextCheckTime) {
		return nil
	}

	s := polargo.New()
	ctx := context.Background()

	_, err := s.CustomerPortal.LicenseKeys.Validate(ctx, components.LicenseKeyValidate{
		Key:            d.KeyID,
		OrganizationID: organizationID,
		ActivationID:   polargo.String(d.ActivationID),
		Conditions: map[string]components.Conditions{
			"macaddr": components.CreateConditionsInteger(int64(macaddr.MacUint64())),
		},
	})
	if err != nil {
		// Activation ID is wrong
		// Error: {"detail":[{"loc":["body","activation_id"],"msg":"Input should be a valid UUID, invalid group length in group 0: expected 8, found 5","type":"uuid_parsing"}]}

		// MacAddress is wrong or not provided.
		// {"error":"ResourceNotFound","detail":"License key does not match required conditions"}

		return err
	}

	// Update VerifiedAt and store the updated license data
	d.VerifiedAt = time.Now()
	if err := d.store(macaddr.MacUint64()); err != nil {
		return fmt.Errorf("failed to update license data: %w", err)
	}

	return nil
}

func (d LicenseData) Deactivate() error {
	ctx := context.Background()
	s := polargo.New()

	_, err := s.CustomerPortal.LicenseKeys.Deactivate(ctx, components.LicenseKeyDeactivate{
		Key:            d.KeyID,
		OrganizationID: organizationID,
		ActivationID:   d.ActivationID,
	})
	if err != nil {
		return err
	}

	if err := os.Remove(filePath); err != nil {
		return err
	}

	return nil
}

func (d LicenseData) buildHash(macAddress uint64) string {
	// Combine the input values into a single string
	data := fmt.Sprintf("%s:%s:%d:%s:%d", d.ActivationID, d.KeyID, macAddress, salt, d.NextCheckTime.Unix())

	// Create a SHA-256 hash
	hash := sha256.Sum256([]byte(data))

	// Convert the hash to a hexadecimal string
	return hex.EncodeToString(hash[:])
}

func (d LicenseData) store(macAddress uint64) error {
	d.NextCheckTime = d.VerifiedAt.Add(recheckInterval)
	d.Hash = d.buildHash(macAddress)

	if err := os.MkdirAll(dataDir, 0o700); err != nil {
		return err
	}

	jsonData, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, jsonData, 0o600)
}
