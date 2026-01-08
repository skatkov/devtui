package toml

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/pelletier/go-toml/v2"
	"github.com/skatkov/devtui/internal/yaml"
)

// TOMLToJSON converts TOML content to JSON format.
func TOMLToJSON(tomlContent string) (string, error) {
	var v any

	err := toml.Unmarshal([]byte(tomlContent), &v)
	if err != nil {
		return "", fmt.Errorf("TOML parsing error: %w", err)
	}

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(v)
	if err != nil {
		return "", fmt.Errorf("JSON encoding error: %w", err)
	}

	return buf.String(), nil
}

// JSONToTOML converts JSON content to TOML format.
func JSONToTOML(jsonContent string) (string, error) {
	var v any

	err := json.Unmarshal([]byte(jsonContent), &v)
	if err != nil {
		return "", fmt.Errorf("JSON parsing error: %w", err)
	}

	var buf bytes.Buffer
	encoder := toml.NewEncoder(&buf)
	err = encoder.Encode(v)
	if err != nil {
		return "", fmt.Errorf("TOML encoding error: %w", err)
	}

	return buf.String(), nil
}

// TOMLToYAML converts TOML content to YAML format.
// Chains through JSON: TOML → JSON → YAML
func TOMLToYAML(tomlContent string) (string, error) {
	jsonStr, err := TOMLToJSON(tomlContent)
	if err != nil {
		return "", err
	}

	return yaml.JSONToYAML(jsonStr)
}
