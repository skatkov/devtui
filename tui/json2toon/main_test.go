package json2toon

import (
	"errors"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/hannes-sistemica/toon"
	"github.com/skatkov/devtui/internal/editor"
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

func TestConvertUsesTwoSpaceIndent(t *testing.T) {
	t.Parallel()

	got, err := Convert(`{"user":{"id":1}}`)
	if err != nil {
		t.Fatalf("Convert() error = %v", err)
	}
	if got != "user:\n  id: 1" {
		t.Errorf("expected two-space indent, got %q", got)
	}
}

func TestConvertWithOptionsErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr string
	}{
		{"syntax error", `{"a" 1}`, "JSON syntax error at offset"},
		{"empty input", ``, "JSON parsing error: EOF"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := ConvertWithOptions(tt.input, toon.EncodeOptions{})
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("ConvertWithOptions(%q) error = %v, want %q", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestIndentKeyTogglesIndent(t *testing.T) {
	t.Parallel()

	model := updateModel(t, newModel(), editor.EditorFinishedMsg{Content: `{"user":{"id":1}}`})

	model = updateModel(t, model, textKeypress('i'))
	if model.FormattedContent != "user:\n    id: 1" {
		t.Errorf("expected four-space indent, got %q", model.FormattedContent)
	}
	if model.StatusMessage != "Indent: 4 spaces" {
		t.Errorf("expected indent status, got %q", model.StatusMessage)
	}

	model = updateModel(t, model, textKeypress('i'))
	if model.indent != 2 {
		t.Errorf("expected indent 2 after second toggle, got %d", model.indent)
	}
}

func TestLengthMarkerKeyTogglesMarker(t *testing.T) {
	t.Parallel()

	model := updateModel(t, newModel(), editor.EditorFinishedMsg{Content: `{"tags":["a"]}`})

	model = updateModel(t, model, textKeypress('l'))
	if model.FormattedContent != "tags[#1]: a" {
		t.Errorf("expected length marker, got %q", model.FormattedContent)
	}
	if model.StatusMessage != "Length marker: on" {
		t.Errorf("expected length marker status, got %q", model.StatusMessage)
	}

	model = updateModel(t, model, textKeypress('l'))
	if model.FormattedContent != "tags[1]: a" {
		t.Errorf("expected no length marker, got %q", model.FormattedContent)
	}
	if model.StatusMessage != "Length marker: off" {
		t.Errorf("expected length marker status, got %q", model.StatusMessage)
	}
}

func TestOptionKeysWithoutContentShowNoStatus(t *testing.T) {
	t.Parallel()

	for _, key := range []rune{'i', 'l'} {
		model := updateModel(t, newModel(), textKeypress(key))
		if model.State != ui.PagerStateBrowse {
			t.Errorf("expected browse state after %q, got %v", key, model.State)
		}
	}
}

func TestOptionKeysWithInvalidContentShowError(t *testing.T) {
	t.Parallel()

	for _, key := range []rune{'i', 'l'} {
		model := updateModel(t, newModel(), editor.EditorFinishedMsg{Content: `{`})
		model = updateModel(t, model, ui.StatusMessageTimeoutMsg{})
		if model.State != ui.PagerStateBrowse {
			t.Fatalf("expected browse state after timeout, got %v", model.State)
		}

		model = updateModel(t, model, textKeypress(key))
		if model.State != ui.PagerStateErrorMessage {
			t.Errorf("expected error state after %q, got %v", key, model.State)
		}
		if !strings.HasPrefix(model.StatusMessage, "error converting JSON to TOON") {
			t.Errorf("expected conversion error after %q, got %q", key, model.StatusMessage)
		}
	}
}

func TestEditorFinishedMsgSetsContent(t *testing.T) {
	t.Parallel()

	model := updateModel(t, newModel(), editor.EditorFinishedMsg{Content: `{"a":1}`})

	if model.Viewport.GetContent() != "a: 1" {
		t.Errorf("expected TOON in viewport, got %q", model.Viewport.GetContent())
	}
	if model.StatusMessage != "Pasted. Press 'c' to copy result." {
		t.Errorf("expected paste status, got %q", model.StatusMessage)
	}
}

func TestEditorFinishedMsgWithErrorShowsError(t *testing.T) {
	t.Parallel()

	model := updateModel(t, newModel(), editor.EditorFinishedMsg{Err: errors.New("editor failed"), Content: `{"a":1}`})

	if model.Content != "" {
		t.Errorf("expected content to stay empty, got %q", model.Content)
	}
	if model.State != ui.PagerStateErrorMessage || model.StatusMessage != "editor failed" {
		t.Errorf("expected editor error, got state %v and %q", model.State, model.StatusMessage)
	}
}

func TestHelpViewShowsCurrentOptions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		indent       int
		lengthMarker string
		want         []string
	}{
		{"defaults", 2, "", []string{"toggle indent (current: 2)", "toggle length marker (current: off)"}},
		{"toggled", 4, "#", []string{"toggle indent (current: 4)", "toggle length marker (current: on)"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			model := newModel()
			model.indent = tt.indent
			model.lengthMarker = tt.lengthMarker

			help := model.helpView()
			for _, want := range tt.want {
				if !strings.Contains(help, want) {
					t.Errorf("expected help to contain %q, got %q", want, help)
				}
			}
		})
	}
}

func TestViewShowsHelp(t *testing.T) {
	t.Parallel()

	model := updateModel(t, newModel(), tea.WindowSizeMsg{Width: 80, Height: 24})
	if strings.Contains(model.View().Content, "copy TOON") {
		t.Fatal("expected help to be hidden")
	}

	model = updateModel(t, model, textKeypress('?'))
	if !strings.Contains(model.View().Content, "copy TOON") {
		t.Errorf("expected help in view, got %q", model.View().Content)
	}
}

func newModel() JsonToonModel {
	return NewJsonToonModel(&ui.CommonModel{Width: 80, Height: 24})
}

func updateModel(t *testing.T, model JsonToonModel, msg tea.Msg) JsonToonModel {
	t.Helper()

	next, _ := model.Update(msg)
	return next.(JsonToonModel)
}

func textKeypress(r rune) tea.KeyPressMsg {
	return tea.KeyPressMsg(tea.Key{Code: r, Text: string(r)})
}
