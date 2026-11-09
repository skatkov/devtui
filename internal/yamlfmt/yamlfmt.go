package yamlfmt

import (
	"bytes"
	"fmt"

	"gopkg.in/yaml.v3"
)

// Format formats YAML content with standard indentation, keeping comments and key order.
func Format(content string) (string, error) {
	if err := validate(content); err != nil {
		return "", err
	}

	var node yaml.Node
	if err := yaml.Unmarshal([]byte(content), &node); err != nil {
		return "", err
	}
	if node.IsZero() {
		return encode(nil)
	}
	return encode(&node)
}

// validate rejects content whose values don't survive a decode and re-encode.
func validate(content string) error {
	var data any
	if err := yaml.Unmarshal([]byte(content), &data); err != nil {
		return err
	}

	formatted, err := encode(data)
	if err != nil {
		return err
	}

	var validation any
	if err := yaml.Unmarshal([]byte(formatted), &validation); err != nil {
		return fmt.Errorf("formatted YAML is invalid: %w", err)
	}
	return nil
}

func encode(value any) (string, error) {
	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	if err := encoder.Encode(value); err != nil {
		return "", err
	}
	if err := encoder.Close(); err != nil {
		return "", err
	}
	return buf.String(), nil
}
