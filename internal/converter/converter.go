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

// mxj writes values as-is unless escaping is on, so "A & B" gave invalid XML.
// Set once here: the setter writes a package global that XmlIndent reads.
func init() {
	mxj.XMLEscapeChars(true)
}

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
	data, err := decodeJSONNumbers(jsonContent)
	if err != nil {
		return "", fmt.Errorf("JSON parsing error: %w", err)
	}

	return encodeYAML(data)
}

func encodeYAML(data any) (string, error) {
	data, err := mapNumbers(data, yamlNumber)
	if err != nil {
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

// decodeJSONNumbers decodes JSON with numbers kept as json.Number, so integers
// stay integers (encoding/json would make every number a float64).
func decodeJSONNumbers(content string) (any, error) {
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
	return v, nil
}

// mapNumbers replaces every json.Number in v, nested ones included, with fn's result.
func mapNumbers(v any, fn func(json.Number) (any, error)) (any, error) {
	switch value := v.(type) {
	case json.Number:
		return fn(value)
	case map[string]any:
		for key, item := range value {
			n, err := mapNumbers(item, fn)
			if err != nil {
				return nil, err
			}
			value[key] = n
		}
	case []any:
		for i, item := range value {
			n, err := mapNumbers(item, fn)
			if err != nil {
				return nil, err
			}
			value[i] = n
		}
	}
	return v, nil
}

// yamlNumber keeps integers as int64, so YAML gets 1700000000000 rather than 1.7e+12.
func yamlNumber(n json.Number) (any, error) {
	if i, err := n.Int64(); err == nil {
		return i, nil
	}
	return n.Float64()
}

// xmlNumber keeps the number as written, but rejects one out of float64 range
// (such as 1e400) as encoding/json does.
func xmlNumber(n json.Number) (any, error) {
	if _, err := n.Float64(); err != nil {
		return nil, err
	}
	return n, nil
}

// keepTOMLFloats turns numbers back into floats where TOML had a float,
// since the JSON round trip writes 1000000.0 as 1000000.
func keepTOMLFloats(v, typed any) any {
	switch value := v.(type) {
	case json.Number:
		if f, ok := typed.(float64); ok {
			return f
		}
	case map[string]any:
		typedMap, _ := typed.(map[string]any)
		for key, item := range value {
			value[key] = keepTOMLFloats(item, typedMap[key])
		}
	case []any:
		typedSlice, _ := typed.([]any)
		for i, item := range value {
			if i < len(typedSlice) {
				value[i] = keepTOMLFloats(item, typedSlice[i])
			}
		}
	}
	return v
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
	v, err := decodeJSONNumbers(jsonContent)
	if err != nil {
		return "", fmt.Errorf("JSON parsing error: %w", err)
	}

	var buf bytes.Buffer
	encoder := toml.NewEncoder(&buf)
	encoder.SetMarshalJSONNumbers(true)
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
	// Keep numbers as written; as float64, 1700000000000 is written as 1.7e+12.
	data, err := decodeJSONNumbers(jsonContent)
	if err != nil {
		return "", fmt.Errorf("JSON parsing error: %w", err)
	}
	data, err = mapNumbers(data, xmlNumber)
	if err != nil {
		return "", fmt.Errorf("JSON parsing error: %w", err)
	}

	var mv mxj.Map
	switch value := data.(type) {
	case map[string]any:
		mv = value
	case []any:
		// Same root key mxj.NewMapJson gives a top-level array.
		mv = mxj.Map{"object": value}
	default:
		return "", errors.New("JSON parsing error: expected an object or array")
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
	var typed any
	if err := toml.Unmarshal([]byte(tomlContent), &typed); err != nil {
		return "", fmt.Errorf("TOML parsing error: %w", err)
	}

	jsonStr, err := TOMLToJSON(tomlContent)
	if err != nil {
		return "", err
	}

	data, err := decodeJSONNumbers(jsonStr)
	if err != nil {
		return "", fmt.Errorf("JSON parsing error: %w", err)
	}

	return encodeYAML(keepTOMLFloats(data, typed))
}
