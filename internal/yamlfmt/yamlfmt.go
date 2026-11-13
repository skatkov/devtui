package yamlfmt

import (
	"bytes"
	"fmt"

	"gopkg.in/yaml.v3"
)

// Format formats YAML content with standard indentation, keeping comments and key order.
func Format(content string) (string, error) {
	var node yaml.Node
	if err := yaml.Unmarshal([]byte(content), &node); err != nil {
		return "", err
	}

	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	if err := encoder.Encode(&node); err != nil {
		return "", err
	}
	if err := encoder.Close(); err != nil {
		return "", err
	}

	formatted := buf.String()
	var validation any
	if err := yaml.Unmarshal([]byte(formatted), &validation); err != nil {
		return "", fmt.Errorf("formatted YAML is invalid: %w", err)
	}

	return formatted, nil
}
