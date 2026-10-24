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

type passwordCharacterSet struct {
	id    string
	label string
	chars string
}

var passwordCharacterSets = []passwordCharacterSet{
	{id: "lowercase", label: "Lowercase (a-z)", chars: "abcdefghijklmnopqrstuvwxyz"},
	{id: "uppercase", label: "Uppercase (A-Z)", chars: "ABCDEFGHIJKLMNOPQRSTUVWXYZ"},
	{id: "numbers", label: "Numbers (0-9)", chars: "0123456789"},
	{id: "special", label: "Special (!@#$%^&*)", chars: "!@#$%^&*()_+-=[]{}|;:,.<>?"},
}

// runForm runs a form on the terminal; tests replace it to script the input.
var runForm = (*huh.Form).Run

func main() {
	config := &PasswordConfig{
		Length:   12, // default length
		Continue: true,
	}

	for config.Continue {
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
				Options(characterSetOptions()...).
				Value(&config.CharacterSets),
		),
	)

	// Run the form
	err := runForm(form)
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

		return runForm(result)
	}

	// Show error if no character set is selected
	errorForm := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("Error").
				Description("Please select at least one character set!"),
		),
	)

	return runForm(errorForm)
}

func characterSetOptions() []huh.Option[string] {
	options := make([]huh.Option[string], 0, len(passwordCharacterSets))
	for _, set := range passwordCharacterSets {
		options = append(options, huh.Option[string]{Key: set.label, Value: set.id})
	}

	return options
}

func selectedCharacterSets(ids []string) []passwordCharacterSet {
	sets := make([]passwordCharacterSet, 0, len(ids))
	for _, id := range ids {
		for _, set := range passwordCharacterSets {
			if set.id == id {
				sets = append(sets, set)
				break
			}
		}
	}

	return sets
}

func generateSecurePassword(config *PasswordConfig) (string, error) {
	selectedSets := selectedCharacterSets(config.CharacterSets)
	if len(selectedSets) == 0 {
		return "", errors.New("no valid character sets selected")
	}

	var charPool strings.Builder
	for _, set := range selectedSets {
		charPool.WriteString(set.chars)
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
	err := ensureAllSetsIncluded(password, selectedSets)
	if err != nil {
		return "", err
	}

	return string(password), nil
}

func ensureAllSetsIncluded(password []byte, selectedSets []passwordCharacterSet) error {
	if len(password) < len(selectedSets) {
		return errors.New("password length must be at least the number of character sets")
	}

	// Put one character of each selected set at its own random position,
	// shuffling the positions so that no set overwrites another.
	positions := make([]int, len(password))
	for i := range positions {
		positions[i] = i
	}

	for i, set := range selectedSets {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(positions)-i)))
		if err != nil {
			return err
		}
		j := i + int(n.Int64())
		positions[i], positions[j] = positions[j], positions[i]

		n, err = rand.Int(rand.Reader, big.NewInt(int64(len(set.chars))))
		if err != nil {
			return err
		}
		password[positions[i]] = set.chars[n.Int64()]
	}

	return nil
}
