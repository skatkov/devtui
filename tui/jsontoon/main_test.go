package jsontoon

import (
	"strings"
	"testing"

	"github.com/hannes-sistemica/toon"
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

	// Check default options
	if model.indent != 2 {
		t.Errorf("NewJsonToonModel() default indent = %v, want 2", model.indent)
	}

	if model.delimiter != "," {
		t.Errorf("NewJsonToonModel() default delimiter = %v, want ,", model.delimiter)
	}

	if model.lengthMarker != "" {
		t.Errorf("NewJsonToonModel() default lengthMarker = %v, want empty string", model.lengthMarker)
	}
}

func TestConvertWithOptions(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		opts        toon.EncodeOptions
		wantContain string
		wantErr     bool
	}{
		{
			name:  "default options",
			input: `{"users":[{"id":1,"name":"Alice"},{"id":2,"name":"Bob"}]}`,
			opts: toon.EncodeOptions{
				Indent:       2,
				Delimiter:    ",",
				LengthMarker: "",
			},
			wantContain: "users[2]{id,name}:",
			wantErr:     false,
		},
		{
			name:  "with length marker",
			input: `{"users":[{"id":1,"name":"Alice"},{"id":2,"name":"Bob"}]}`,
			opts: toon.EncodeOptions{
				Indent:       2,
				Delimiter:    ",",
				LengthMarker: "#",
			},
			wantContain: "users[#2]{id,name}:",
			wantErr:     false,
		},
		{
			name:  "tab delimiter",
			input: `{"tags":["foo","bar","baz"]}`,
			opts: toon.EncodeOptions{
				Indent:       2,
				Delimiter:    "\t",
				LengthMarker: "",
			},
			wantContain: "tags[3\t]:",
			wantErr:     false,
		},
		{
			name:  "pipe delimiter",
			input: `{"tags":["foo","bar","baz"]}`,
			opts: toon.EncodeOptions{
				Indent:       2,
				Delimiter:    "|",
				LengthMarker: "",
			},
			wantContain: "tags[3|]:",
			wantErr:     false,
		},
		{
			name:  "indent 4 spaces",
			input: `{"user":{"id":1,"name":"Alice"}}`,
			opts: toon.EncodeOptions{
				Indent:       4,
				Delimiter:    ",",
				LengthMarker: "",
			},
			wantContain: "user:",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ConvertWithOptions(tt.input, tt.opts)

			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertWithOptions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil && tt.wantContain != "" {
				if !strings.Contains(result, tt.wantContain) {
					t.Errorf("ConvertWithOptions() result does not contain %q\nGot: %s", tt.wantContain, result)
				}
			}
		})
	}
}

func TestModelOptions(t *testing.T) {
	common := &ui.CommonModel{
		Width:  80,
		Height: 24,
	}

	model := NewJsonToonModel(common)
	testJSON := `{"users":[{"id":1,"name":"Alice"}]}`

	// Set initial content
	err := model.SetContent(testJSON)
	if err != nil {
		t.Fatalf("SetContent() failed: %v", err)
	}

	// Test changing indent
	model.indent = 4
	err = model.SetContent(testJSON)
	if err != nil {
		t.Errorf("SetContent() with indent=4 failed: %v", err)
	}

	// Test changing delimiter
	model.delimiter = "\t"
	err = model.SetContent(testJSON)
	if err != nil {
		t.Errorf("SetContent() with tab delimiter failed: %v", err)
	}

	model.delimiter = "|"
	err = model.SetContent(testJSON)
	if err != nil {
		t.Errorf("SetContent() with pipe delimiter failed: %v", err)
	}

	// Test changing length marker
	model.lengthMarker = "#"
	err = model.SetContent(testJSON)
	if err != nil {
		t.Errorf("SetContent() with length marker failed: %v", err)
	}

	if !strings.Contains(model.FormattedContent, "#") {
		t.Errorf("SetContent() with length marker should contain '#' in output")
	}
}
