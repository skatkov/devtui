package yamlfmt

import (
	"bytes"

	"gopkg.in/yaml.v3"
)

// Format formats YAML content with standard indentation.
func Format(content string) (string, error) {
	var data any
	if err := yaml.Unmarshal([]byte(content), &data); err != nil {
		return "", err
	}

	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	if err := encoder.Encode(data); err != nil {
		return "", err
	}
	if err := encoder.Close(); err != nil {
		return "", err
	}

	return buf.String(), nil
}
