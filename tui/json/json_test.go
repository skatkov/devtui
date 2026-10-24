package json

import (
	"errors"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/skatkov/devtui/internal/editor"
	"github.com/skatkov/devtui/internal/ui"
)

func TestFormatJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"object", `{"b":1,"a":[true,null]}`, "{\n  \"a\": [\n    true,\n    null\n  ],\n  \"b\": 1\n}\n"},
		{"invalid json", `{"a":`, `{"a":`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := FormatJSON(tt.input); got != tt.want {
				t.Errorf("FormatJSON(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestSetContent(t *testing.T) {
	t.Parallel()

	model := newModel()
	if err := model.SetContent(`{"a":1}`); err != nil {
		t.Fatalf("SetContent() error = %v", err)
	}

	if model.Content != `{"a":1}` {
		t.Errorf("expected raw content to be kept, got %q", model.Content)
	}
	if model.FormattedContent != "{\n  \"a\": 1\n}\n" {
		t.Errorf("expected formatted JSON, got %q", model.FormattedContent)
	}
	if !strings.Contains(model.Viewport.GetContent(), "\x1b[") {
		t.Errorf("expected highlighted viewport content, got %q", model.Viewport.GetContent())
	}
}

func TestEditorFinishedMsgSetsContent(t *testing.T) {
	t.Parallel()

	model := updateModel(t, newModel(), editor.EditorFinishedMsg{Content: `[1]`})

	if model.FormattedContent != "[\n  1\n]\n" {
		t.Errorf("expected formatted JSON, got %q", model.FormattedContent)
	}
	if model.State != ui.PagerStateStatusMessage {
		t.Errorf("expected status message state, got %v", model.State)
	}
	if model.StatusMessage != "Pasted. Press 'c' to copy" {
		t.Errorf("expected paste status, got %q", model.StatusMessage)
	}
}

func TestEditorFinishedMsgWithErrorShowsError(t *testing.T) {
	t.Parallel()

	model := updateModel(t, newModel(), editor.EditorFinishedMsg{Err: errors.New("editor failed"), Content: `[1]`})

	if model.Content != "" {
		t.Errorf("expected content to stay empty, got %q", model.Content)
	}
	if model.State != ui.PagerStateErrorMessage {
		t.Errorf("expected error message state, got %v", model.State)
	}
	if model.StatusMessage != "editor failed" {
		t.Errorf("expected editor error, got %q", model.StatusMessage)
	}
}

func TestStatusMessageTimeoutReturnsToBrowse(t *testing.T) {
	t.Parallel()

	model := updateModel(t, newModel(), editor.EditorFinishedMsg{Content: `[1]`})
	model = updateModel(t, model, ui.StatusMessageTimeoutMsg{})

	if model.State != ui.PagerStateBrowse {
		t.Errorf("expected browse state, got %v", model.State)
	}
}

func TestViewShowsHelp(t *testing.T) {
	t.Parallel()

	model := updateModel(t, newModel(), tea.WindowSizeMsg{Width: 80, Height: 24})
	if strings.Contains(model.View().Content, "copy formatted JSON") {
		t.Fatal("expected help to be hidden")
	}

	model = updateModel(t, model, textKeypress('?'))
	if !strings.Contains(model.View().Content, "copy formatted JSON") {
		t.Errorf("expected help in view, got %q", model.View().Content)
	}
}

func newModel() JsonModel {
	return NewJsonModel(&ui.CommonModel{Width: 80, Height: 24})
}

func updateModel(t *testing.T, model JsonModel, msg tea.Msg) JsonModel {
	t.Helper()

	next, _ := model.Update(msg)
	return next.(JsonModel)
}

func textKeypress(r rune) tea.KeyPressMsg {
	return tea.KeyPressMsg(tea.Key{Code: r, Text: string(r)})
}
