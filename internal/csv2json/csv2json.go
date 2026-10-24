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
	entries := make([]map[string]any, 0, len(rows))
	for _, row := range rows[1:] {
		entry := map[string]any{}
		for i, value := range row {
			if err := setValue(entry, attributes[i], value); err != nil {
				return "", err
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

// setValue stores value in entry under the dotted, optionally indexed path in attribute.
func setValue(entry map[string]any, attribute, value string) error {
	parts := strings.Split(attribute, ".")
	internal := entry
	for _, part := range parts[:len(parts)-1] {
		next, err := descend(internal, part, attribute)
		if err != nil {
			return err
		}
		internal = next
	}

	return assign(internal, parts[len(parts)-1], value, attribute)
}

// descend returns the object stored under part, creating it when missing.
func descend(internal map[string]any, part, attribute string) (map[string]any, error) {
	key, arrayIndex, isArray := arrayContentMatch(part)
	if !isArray {
		if internal[key] == nil {
			internal[key] = map[string]any{}
		}
		next, ok := internal[key].(map[string]any)
		if !ok {
			return nil, fmt.Errorf("invalid CSV header %q: key %q is used as both scalar and object", attribute, key)
		}
		return next, nil
	}

	internalArray, err := arrayAt(internal, key, attribute)
	if err != nil {
		return nil, err
	}
	internalArray = padArray(internalArray, arrayIndex, func() any { return map[string]any{} })
	if internalArray[arrayIndex] == nil {
		internalArray[arrayIndex] = map[string]any{}
	}
	next, ok := internalArray[arrayIndex].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid CSV header %q: key %q[%d] is used as both scalar and object", attribute, key, arrayIndex)
	}
	internal[key] = internalArray
	return next, nil
}

// assign stores value under the last path part, which may be indexed.
func assign(internal map[string]any, part, value, attribute string) error {
	key, arrayIndex, isArray := arrayContentMatch(part)
	if !isArray {
		switch internal[key].(type) {
		case map[string]any, []any:
			return fmt.Errorf("invalid CSV header %q: key %q is used as both scalar and nested value", attribute, key)
		}
		internal[key] = value
		return nil
	}

	internalArray, err := arrayAt(internal, key, attribute)
	if err != nil {
		return err
	}
	internalArray = padArray(internalArray, arrayIndex, func() any { return nil })
	if nested, ok := internalArray[arrayIndex].(map[string]any); ok && len(nested) > 0 {
		return fmt.Errorf("invalid CSV header %q: key %q[%d] is used as both object and scalar", attribute, key, arrayIndex)
	}
	internalArray[arrayIndex] = value
	internal[key] = internalArray
	return nil
}

// arrayAt returns the array stored under key, or nil when the key is unset.
func arrayAt(internal map[string]any, key, attribute string) ([]any, error) {
	existing, exists := internal[key]
	if !exists {
		return nil, nil
	}
	internalArray, ok := existing.([]any)
	if !ok {
		return nil, fmt.Errorf("invalid CSV header %q: key %q is used as both object and array", attribute, key)
	}
	return internalArray, nil
}

// padArray grows arr until index is in range, filling new slots with fill.
func padArray(arr []any, index int, fill func() any) []any {
	for range index + 1 - len(arr) {
		arr = append(arr, fill())
	}
	return arr
}

// arrayContentMatch splits "key[index]" into key and index; ok is false for plain keys.
func arrayContentMatch(str string) (key string, index int, ok bool) {
	key, rest, opened := strings.Cut(str, "[")
	digits, _, closed := strings.Cut(rest, "]")
	if !opened || !closed || strings.Contains(key, "]") {
		return str, index, false
	}

	index, err := strconv.Atoi(digits)
	if err != nil || index < 0 {
		return str, index, false
	}
	return key, index, true
}
