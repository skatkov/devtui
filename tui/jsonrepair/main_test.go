package jsonrepair

import (
	"errors"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/skatkov/devtui/internal/editor"
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

func TestRepairJSON(t *testing.T) {
	t.Parallel()

	got, err := RepairJSON(`{'a': [1, 2`)
	if err != nil {
		t.Fatalf("RepairJSON() error = %v", err)
	}
	if got != `{"a":[1,2]}` {
		t.Errorf("RepairJSON() = %q, want %q", got, `{"a":[1,2]}`)
	}
}

func TestRepairJSONTooDeepReturnsError(t *testing.T) {
	t.Parallel()

	got, err := RepairJSON(strings.Repeat("[", 1001))
	if err == nil || !strings.HasPrefix(err.Error(), "failed to repair JSON") {
		t.Errorf("expected repair error, got %q and %v", got, err)
	}
	if got != "" {
		t.Errorf("expected empty result on error, got %q", got)
	}
}

func TestEditorFinishedMsgSetsContent(t *testing.T) {
	t.Parallel()

	model := updateModel(t, newModel(), editor.EditorFinishedMsg{Content: `{'a': 1}`})

	if model.FormattedContent != `{"a":1}` {
		t.Errorf("unexpected formatted content %q", model.FormattedContent)
	}
	if !strings.Contains(model.Viewport.GetContent(), "\x1b[") {
		t.Errorf("expected highlighted viewport content, got %q", model.Viewport.GetContent())
	}
	if model.State != ui.PagerStateStatusMessage || model.StatusMessage != "Repaired. Press 'c' to copy result." {
		t.Errorf("expected status message, got state %v and %q", model.State, model.StatusMessage)
	}
}

func TestEditorFinishedMsgWithErrorShowsError(t *testing.T) {
	t.Parallel()

	model := updateModel(t, newModel(), editor.EditorFinishedMsg{Err: errors.New("editor failed"), Content: `{'a': 1}`})

	if model.Content != "" {
		t.Errorf("expected content to stay empty, got %q", model.Content)
	}
	if model.State != ui.PagerStateErrorMessage || model.StatusMessage != "editor failed" {
		t.Errorf("expected editor error, got state %v and %q", model.State, model.StatusMessage)
	}
}

func TestEditorFinishedMsgWithTooDeepContentShowsError(t *testing.T) {
	t.Parallel()

	model := updateModel(t, newModel(), editor.EditorFinishedMsg{Content: strings.Repeat("[", 1001)})

	if model.State != ui.PagerStateErrorMessage || !strings.HasPrefix(model.StatusMessage, "failed to repair JSON") {
		t.Errorf("expected repair error, got state %v and %q", model.State, model.StatusMessage)
	}
}

func TestStatusMessageTimeoutReturnsToBrowse(t *testing.T) {
	t.Parallel()

	model := updateModel(t, newModel(), editor.EditorFinishedMsg{Content: `{'a': 1}`})
	model = updateModel(t, model, ui.StatusMessageTimeoutMsg{})

	if model.State != ui.PagerStateBrowse {
		t.Errorf("expected browse state, got %v", model.State)
	}
}

func TestViewShowsHelp(t *testing.T) {
	t.Parallel()

	model := updateModel(t, newModel(), tea.WindowSizeMsg{Width: 80, Height: 24})
	if strings.Contains(model.View().Content, "copy repaired JSON") {
		t.Fatal("expected help to be hidden")
	}

	model = updateModel(t, model, textKeypress('?'))
	if !strings.Contains(model.View().Content, "copy repaired JSON") {
		t.Errorf("expected help in view, got %q", model.View().Content)
	}
}

func newModel() JSONRepairModel {
	return NewJSONRepairModel(&ui.CommonModel{Width: 80, Height: 24})
}

func updateModel(t *testing.T, model JSONRepairModel, msg tea.Msg) JSONRepairModel {
	t.Helper()

	next, _ := model.Update(msg)
	return next.(JSONRepairModel)
}

func textKeypress(r rune) tea.KeyPressMsg {
	return tea.KeyPressMsg(tea.Key{Code: r, Text: string(r)})
}
