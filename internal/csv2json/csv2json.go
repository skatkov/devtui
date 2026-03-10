package csv2json

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Convert converts CSV content to JSON.
func Convert(content string) (string, error) {
	reader := csv.NewReader(strings.NewReader(content))
	rows, err := reader.ReadAll()
	if err != nil {
		return "", err
	}
	if len(rows) == 0 {
		return "", errors.New("empty CSV file")
	}

	return rowsToJSON(rows)
}

func rowsToJSON(rows [][]string) (string, error) {
	attributes := rows[0]
	entries := make([]map[string]any, 0, len(rows)-1)
	for _, row := range rows[1:] {
		entry := map[string]any{}
		for i, value := range row {
			if i >= len(attributes) {
				continue
			}

			attribute := attributes[i]
			objectSlice := strings.Split(attribute, ".")
			internal := entry
			for index, val := range objectSlice {
				key, arrayIndex := arrayContentMatch(val)
				if arrayIndex != -1 {
					if internal[key] == nil {
						internal[key] = []any{}
					}

					internalArray, ok := internal[key].([]any)
					if !ok {
						return "", fmt.Errorf("invalid CSV header %q: key %q is used as both object and array", attribute, key)
					}

					if index == len(objectSlice)-1 {
						internalArray = append(internalArray, value)
						internal[key] = internalArray
						break
					}

					for arrayIndex >= len(internalArray) {
						internalArray = append(internalArray, map[string]any{})
					}

					if internalArray[arrayIndex] == nil {
						internalArray[arrayIndex] = map[string]any{}
					}

					nextInternal, ok := internalArray[arrayIndex].(map[string]any)
					if !ok {
						return "", fmt.Errorf("invalid CSV header %q: key %q[%d] is used as both scalar and object", attribute, key, arrayIndex)
					}

					internal[key] = internalArray
					internal = nextInternal
				} else {
					if index == len(objectSlice)-1 {
						if existing, exists := internal[key]; exists {
							switch existing.(type) {
							case map[string]any, []any:
								return "", fmt.Errorf("invalid CSV header %q: key %q is used as both scalar and nested value", attribute, key)
							}
						}
						internal[key] = value
						break
					}
					if internal[key] == nil {
						internal[key] = map[string]any{}
					}

					nextInternal, ok := internal[key].(map[string]any)
					if !ok {
						return "", fmt.Errorf("invalid CSV header %q: key %q is used as both scalar and object", attribute, key)
					}

					internal = nextInternal
				}
			}
		}
		entries = append(entries, entry)
	}

	bytes, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal error: %w", err)
	}

	return string(bytes), nil
}

func arrayContentMatch(str string) (string, int) {
	i := strings.Index(str, "[")
	if i >= 0 {
		j := strings.Index(str, "]")
		if j > i {
			index, err := strconv.Atoi(str[i+1 : j])
			if err != nil {
				return str, -1
			}
			return str[0:i], index
		}
	}
	return str, -1
}
