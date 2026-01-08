package yaml

import (
	"bytes"
	"encoding/json"

	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
)

// YAMLToJSON converts YAML content to JSON format.
func YAMLToJSON(yamlContent string) (string, error) {
	var data any

	if err := yaml.Unmarshal([]byte(yamlContent), &data); err != nil {
		return "", err
	}

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// JSONToYAML converts JSON content to YAML format.
func JSONToYAML(jsonContent string) (string, error) {
	var data any

	if err := json.Unmarshal([]byte(jsonContent), &data); err != nil {
		return "", err
	}

	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	if err := encoder.Encode(data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// YAMLToTOML converts YAML content to TOML format.
// Chains through JSON: YAML → JSON → TOML
func YAMLToTOML(yamlContent string) (string, error) {
	jsonStr, err := YAMLToJSON(yamlContent)
	if err != nil {
		return "", err
	}

	var data any
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return "", err
	}

	var buf bytes.Buffer
	encoder := toml.NewEncoder(&buf)
	if err := encoder.Encode(data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// TOMLToYAML converts TOML content to YAML format.
// Chains through JSON: TOML → JSON → YAML
func TOMLToYAML(tomlContent string) (string, error) {
	var data any
	if err := toml.Unmarshal([]byte(tomlContent), &data); err != nil {
		return "", err
	}

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		return "", err
	}

	jsonStr := buf.String()
	return JSONToYAML(jsonStr)
}
