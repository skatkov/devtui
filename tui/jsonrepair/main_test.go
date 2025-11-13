package jsonrepair

import (
	"testing"

	"github.com/skatkov/devtui/internal/ui"
)

func TestSetContent(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantErr     bool
		description string
	}{
		{
			name:        "valid JSON",
			input:       `{"key": "value"}`,
			wantErr:     false,
			description: "Should handle valid JSON",
		},
		{
			name:        "single quotes",
			input:       `{'key': 'value'}`,
			wantErr:     false,
			description: "Should repair single quotes to double quotes",
		},
		{
			name:        "unclosed array",
			input:       `[1, 2, 3, 4`,
			wantErr:     false,
			description: "Should close unclosed array",
		},
		{
			name:        "unclosed object",
			input:       `{"employees":["John", "Anna",`,
			wantErr:     false,
			description: "Should close unclosed object and array",
		},
		{
			name:        "JSON with markdown code block",
			input:       "```json\n{'key': 'value'}\n```",
			wantErr:     false,
			description: "Should strip markdown code block and repair JSON",
		},
		{
			name:        "mixed quotes",
			input:       `{'key': 'string', 'key2': false, "key3": null}`,
			wantErr:     false,
			description: "Should normalize mixed quotes",
		},
		{
			name:        "boolean uppercase",
			input:       `{"key": TRUE, "key2": FALSE, "key3": Null}`,
			wantErr:     false,
			description: "Should normalize boolean and null values",
		},
		{
			name:        "trailing comma",
			input:       `{"key":"value",}`,
			wantErr:     false,
			description: "Should handle trailing commas",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			common := &ui.CommonModel{
				Width:  80,
				Height: 24,
			}
			model := NewJSONRepairModel(common)

			err := model.SetContent(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("SetContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil && model.FormattedContent == "" {
				t.Errorf("SetContent() produced empty FormattedContent for input: %s", tt.input)
			}

			// Basic validation that output looks like JSON
			if err == nil && len(model.FormattedContent) > 0 {
				if model.FormattedContent[0] != '{' && model.FormattedContent[0] != '[' {
					t.Errorf("SetContent() produced invalid JSON format: %s", model.FormattedContent)
				}
			}
		})
	}
}

func TestNewJSONRepairModel(t *testing.T) {
	common := &ui.CommonModel{
		Width:  80,
		Height: 24,
	}

	model := NewJSONRepairModel(common)

	if model.Title != Title {
		t.Errorf("NewJSONRepairModel() title = %v, want %v", model.Title, Title)
	}

	if model.Common != common {
		t.Errorf("NewJSONRepairModel() common model not set correctly")
	}
}
