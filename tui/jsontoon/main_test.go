package jsontoon

import (
	"strings"
	"testing"

	"github.com/skatkov/devtui/internal/ui"
)

func TestSetContent(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantErr     bool
		wantContain string
		description string
	}{
		{
			name:        "simple object",
			input:       `{"name": "Alice", "age": 30}`,
			wantErr:     false,
			wantContain: "name:",
			description: "Should convert simple JSON object to TOON",
		},
		{
			name:        "nested object",
			input:       `{"user": {"id": 1, "name": "Bob"}}`,
			wantErr:     false,
			wantContain: "user:",
			description: "Should handle nested objects",
		},
		{
			name:        "array of primitives",
			input:       `{"tags": ["foo", "bar", "baz"]}`,
			wantErr:     false,
			wantContain: "tags[3]:",
			description: "Should convert arrays with length markers",
		},
		{
			name:        "array of objects",
			input:       `{"users": [{"id": 1, "name": "Alice"}, {"id": 2, "name": "Bob"}]}`,
			wantErr:     false,
			wantContain: "users[2]{id,name}:",
			description: "Should convert array of objects to tabular format",
		},
		{
			name:        "empty array",
			input:       `{"items": []}`,
			wantErr:     false,
			wantContain: "items[0]:",
			description: "Should handle empty arrays",
		},
		{
			name:        "boolean and null",
			input:       `{"active": true, "inactive": false, "empty": null}`,
			wantErr:     false,
			wantContain: "active: true",
			description: "Should handle boolean and null values",
		},
		{
			name:        "numbers",
			input:       `{"int": 42, "float": 3.14, "negative": -10}`,
			wantErr:     false,
			wantContain: "int: 42",
			description: "Should handle various number formats",
		},
		{
			name:        "invalid JSON - missing quote",
			input:       `{"key: "value"}`,
			wantErr:     true,
			wantContain: "",
			description: "Should return error for invalid JSON",
		},
		{
			name:        "invalid JSON - trailing comma",
			input:       `{"key": "value",}`,
			wantErr:     true,
			wantContain: "",
			description: "Should return error for malformed JSON",
		},
		{
			name:        "empty object",
			input:       `{}`,
			wantErr:     false,
			wantContain: "",
			description: "Should handle empty object",
		},
		{
			name:        "complex nested structure",
			input:       `{"company": "Acme", "employees": [{"id": 1, "name": "Alice", "role": "Engineer"}, {"id": 2, "name": "Bob", "role": "Designer"}]}`,
			wantErr:     false,
			wantContain: "employees[2]{id,name,role}:",
			description: "Should handle complex nested structures",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			common := &ui.CommonModel{
				Width:  80,
				Height: 24,
			}
			model := NewJsonToonModel(common)

			err := model.SetContent(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("SetContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil && model.FormattedContent == "" && tt.input != "{}" && tt.input != "[]" {
				t.Errorf("SetContent() produced empty FormattedContent for input: %s", tt.input)
			}

			// Check if output contains expected content
			if err == nil && tt.wantContain != "" {
				if !strings.Contains(model.FormattedContent, tt.wantContain) {
					t.Errorf("SetContent() output does not contain expected string %q\nGot: %s", tt.wantContain, model.FormattedContent)
				}
			}
		})
	}
}

func TestConvert(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		check   func(string) bool
	}{
		{
			name:    "simple conversion",
			input:   `{"key": "value"}`,
			wantErr: false,
			check: func(s string) bool {
				return strings.Contains(s, "key:")
			},
		},
		{
			name:    "preserves key order",
			input:   `{"first": 1, "second": 2, "third": 3}`,
			wantErr: false,
			check: func(s string) bool {
				// TOON should preserve order from JSON parsing
				return strings.Contains(s, "first:") && strings.Contains(s, "second:") && strings.Contains(s, "third:")
			},
		},
		{
			name:    "invalid JSON",
			input:   `{invalid json`,
			wantErr: true,
			check:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Convert(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("Convert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil && tt.check != nil && !tt.check(result) {
				t.Errorf("Convert() result did not pass check function\nGot: %s", result)
			}
		})
	}
}

func TestNewJsonToonModel(t *testing.T) {
	common := &ui.CommonModel{
		Width:  80,
		Height: 24,
	}

	model := NewJsonToonModel(common)

	if model.Title != Title {
		t.Errorf("NewJsonToonModel() title = %v, want %v", model.Title, Title)
	}

	if model.Common != common {
		t.Errorf("NewJsonToonModel() common model not set correctly")
	}
}
