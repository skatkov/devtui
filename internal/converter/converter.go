package converter

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/clbanning/mxj/v2"
	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
)

func YAMLToJSON(yamlContent string) (string, error) {
	var data any

	if err := yaml.Unmarshal([]byte(yamlContent), &data); err != nil {
		return "", fmt.Errorf("YAML parsing error: %w", err)
	}

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		return "", fmt.Errorf("JSON encoding error: %w", err)
	}

	return buf.String(), nil
}

func JSONToYAML(jsonContent string) (string, error) {
	var data any

	if err := json.Unmarshal([]byte(jsonContent), &data); err != nil {
		return "", fmt.Errorf("JSON parsing error: %w", err)
	}

	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	if err := encoder.Encode(data); err != nil {
		return "", fmt.Errorf("YAML encoding error: %w", err)
	}

	return buf.String(), nil
}

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

func XMLToJSON(xmlContent string) (string, error) {
	mv, err := mxj.NewMapXml([]byte(xmlContent))
	if err != nil {
		return "", fmt.Errorf("XML parsing error: %w", err)
	}

	jsonBytes, err := json.MarshalIndent(mv, "", "  ")
	if err != nil {
		return "", fmt.Errorf("JSON encoding error: %w", err)
	}

	return string(jsonBytes), nil
}

func JSONToXML(jsonContent string) (string, error) {
	mv, err := mxj.NewMapJson([]byte(jsonContent))
	if err != nil {
		return "", fmt.Errorf("JSON parsing error: %w", err)
	}

	xmlBytes, err := mv.XmlIndent("", "  ")
	if err != nil {
		return "", fmt.Errorf("XML encoding error: %w", err)
	}

	var buf bytes.Buffer
	buf.Write(xmlBytes)
	return buf.String(), nil
}

func YAMLToTOML(yamlContent string) (string, error) {
	jsonStr, err := YAMLToJSON(yamlContent)
	if err != nil {
		return "", err
	}

	var data any
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return "", fmt.Errorf("JSON parsing error: %w", err)
	}

	var buf bytes.Buffer
	encoder := toml.NewEncoder(&buf)
	if err := encoder.Encode(data); err != nil {
		return "", fmt.Errorf("TOML encoding error: %w", err)
	}

	return buf.String(), nil
}

func TOMLToYAML(tomlContent string) (string, error) {
	jsonStr, err := TOMLToJSON(tomlContent)
	if err != nil {
		return "", err
	}

	return JSONToYAML(jsonStr)
}
