package main

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"charm.land/huh/v2"
)

// PasswordConfig holds the configuration for password generation.
type PasswordConfig struct {
	Length        int
	CharacterSets []string
	GeneratedPass string
	Continue      bool
}

const (
	lowercase = "abcdefghijklmnopqrstuvwxyz"
	uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers   = "0123456789"
	special   = "!@#$%^&*()_+-=[]{}|;:,.<>?"

	charsetLowercase = "Lowercase (a-z)"
	charsetUppercase = "Uppercase (A-Z)"
	charsetNumbers   = "Numbers (0-9)"
	charsetSpecial   = "Special (!@#$%^&*)"
)

var characterSetOptions = []huh.Option[string]{
	{Key: charsetLowercase, Value: charsetLowercase},
	{Key: charsetUppercase, Value: charsetUppercase},
	{Key: charsetNumbers, Value: charsetNumbers},
	{Key: charsetSpecial, Value: charsetSpecial},
}

func main() {
	config := &PasswordConfig{
		Length: 12, // default length
	}

	for {
		err := showInterface(config)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	}
}

func showInterface(config *PasswordConfig) error {
	lengthStr := strconv.Itoa(config.Length)

	// Create the form
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Password Length").
				Value(&lengthStr),
			huh.NewMultiSelect[string]().
				Title("Select character sets to include").
				Options(characterSetOptions...).
				Value(&config.CharacterSets),
		),
	)

	// Run the form
	err := form.Run()
	if err != nil {
		return err
	}

	// Convert length back to int
	config.Length, err = strconv.Atoi(lengthStr)
	if err != nil {
		return fmt.Errorf("invalid length: %w", err)
	}

	// Generate password if at least one character set is selected
	if len(config.CharacterSets) > 0 {
		config.GeneratedPass, err = generateSecurePassword(config)
		if err != nil {
			return fmt.Errorf("failed to generate password: %w", err)
		}

		// Show the generated password
		result := huh.NewForm(
			huh.NewGroup(
				huh.NewNote().
					Title("Generated Password").
					Description(config.GeneratedPass),
				huh.NewConfirm().
					Title("Generate another password?").
					Value(&config.Continue),
			),
		)

		err = result.Run()
		if err != nil {
			return err
		}

		return nil
	}

	// Show error if no character set is selected
	errorForm := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("Error").
				Description("Please select at least one character set!"),
		),
	)

	return errorForm.Run()
}

func generateSecurePassword(config *PasswordConfig) (string, error) {
	// Create the character pool based on selected options
	var charPool strings.Builder

	for _, set := range config.CharacterSets {
		switch set {
		case charsetLowercase:
			charPool.WriteString(lowercase)
		case charsetUppercase:
			charPool.WriteString(uppercase)
		case charsetNumbers:
			charPool.WriteString(numbers)
		case charsetSpecial:
			charPool.WriteString(special)
		}
	}

	pool := charPool.String()
	poolLength := big.NewInt(int64(len(pool)))
	password := make([]byte, config.Length)

	for i := range config.Length {
		// Generate cryptographically secure random number
		n, err := rand.Int(rand.Reader, poolLength)
		if err != nil {
			return "", fmt.Errorf("failed to generate random number: %w", err)
		}
		password[i] = pool[n.Int64()]
	}

	// Ensure at least one character from each selected set
	err := ensureAllSetsIncluded(password, config)
	if err != nil {
		return "", err
	}

	return string(password), nil
}

func ensureAllSetsIncluded(password []byte, config *PasswordConfig) error {
	if len(password) < len(config.CharacterSets) {
		return errors.New("password length must be at least the number of character sets")
	}

	// For each selected character set, ensure at least one character is included
	positions := make([]int, len(config.CharacterSets))
	for i := range positions {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(password))))
		if err != nil {
			return err
		}
		positions[i] = int(n.Int64())
	}

	for i, set := range config.CharacterSets {
		var chars string
		switch set {
		case charsetLowercase:
			chars = lowercase
		case charsetUppercase:
			chars = uppercase
		case charsetNumbers:
			chars = numbers
		case charsetSpecial:
			chars = special
		}

		if len(chars) > 0 {
			n, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
			if err != nil {
				return err
			}
			password[positions[i]] = chars[n.Int64()]
		}
	}

	return nil
}
