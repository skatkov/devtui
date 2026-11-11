package converter

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"math"
	"time"

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

// decodeJSONForTOML keeps integers as integers: TOML has separate integer and float types.
func decodeJSONForTOML(content string) (any, error) {
	// Unmarshal checks the whole input, so trailing data is still rejected.
	var raw json.RawMessage
	if err := json.Unmarshal([]byte(content), &raw); err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()

	var v any
	if err := decoder.Decode(&v); err != nil {
		return nil, err
	}
	return convertJSONNumbers(v)
}

func convertJSONNumbers(v any) (any, error) {
	switch value := v.(type) {
	case json.Number:
		if i, err := value.Int64(); err == nil {
			return i, nil
		}
		return value.Float64()
	case map[string]any:
		for key, item := range value {
			converted, err := convertJSONNumbers(item)
			if err != nil {
				return nil, err
			}
			value[key] = converted
		}
	case []any:
		for i, item := range value {
			converted, err := convertJSONNumbers(item)
			if err != nil {
				return nil, err
			}
			value[i] = converted
		}
	}
	return v, nil
}

// yamlValuesForTOML keeps the numeric types YAML decoded and turns the rest
// into what the JSON round trip used to give.
func yamlValuesForTOML(v any) any {
	switch value := v.(type) {
	case time.Time:
		return value.Format(time.RFC3339Nano)
	case uint64:
		if value > math.MaxInt64 {
			return float64(value)
		}
	case map[string]any:
		for key, item := range value {
			value[key] = yamlValuesForTOML(item)
		}
	case []any:
		for i, item := range value {
			value[i] = yamlValuesForTOML(item)
		}
	}
	return v
}

func JSONToTOML(jsonContent string) (string, error) {
	v, err := decodeJSONForTOML(jsonContent)
	if err != nil {
		return "", fmt.Errorf("JSON parsing error: %w", err)
	}

	var buf bytes.Buffer
	encoder := toml.NewEncoder(&buf)
	err = encoder.Encode(v)
	if err != nil {
		return "", fmt.Errorf("TOML encoding error: %w", err)
	}

	tomlContent := buf.String()
	var validation any
	if err := toml.Unmarshal([]byte(tomlContent), &validation); err != nil {
		return "", fmt.Errorf("TOML encoding error: %w", err)
	}

	return tomlContent, nil
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

	if len(bytes.TrimSpace(xmlBytes)) == 0 {
		return "", errors.New("XML encoding error: empty XML output")
	}

	decoder := xml.NewDecoder(bytes.NewReader(xmlBytes))
	seenStart := false
	for {
		tok, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", fmt.Errorf("XML encoding error: %w", err)
		}

		if _, ok := tok.(xml.StartElement); ok {
			seenStart = true
		}
	}

	if !seenStart {
		return "", errors.New("XML encoding error: empty XML output")
	}

	var buf bytes.Buffer
	buf.Write(xmlBytes)
	return buf.String(), nil
}

func YAMLToTOML(yamlContent string) (string, error) {
	var data any
	if err := yaml.Unmarshal([]byte(yamlContent), &data); err != nil {
		return "", fmt.Errorf("YAML parsing error: %w", err)
	}

	// Reject what YAMLToJSON rejects (non-string keys, NaN, Inf).
	if _, err := json.Marshal(data); err != nil {
		return "", fmt.Errorf("JSON encoding error: %w", err)
	}

	var buf bytes.Buffer
	encoder := toml.NewEncoder(&buf)
	if err := encoder.Encode(yamlValuesForTOML(data)); err != nil {
		return "", fmt.Errorf("TOML encoding error: %w", err)
	}

	tomlContent := buf.String()
	var validation any
	if err := toml.Unmarshal([]byte(tomlContent), &validation); err != nil {
		return "", fmt.Errorf("TOML encoding error: %w", err)
	}

	return tomlContent, nil
}

func TOMLToYAML(tomlContent string) (string, error) {
	jsonStr, err := TOMLToJSON(tomlContent)
	if err != nil {
		return "", err
	}

	return JSONToYAML(jsonStr)
}
