package yaml

import (
	"errors"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/skatkov/devtui/internal/editor"
	"github.com/skatkov/devtui/internal/ui"
)

func newTestModel() YamlModel {
	model := NewYamlModel(&ui.CommonModel{Width: 80, Height: 24})
	next, _ := model.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	return next.(YamlModel)
}

func TestNewYamlModel(t *testing.T) {
	t.Parallel()

	model := NewYamlModel(&ui.CommonModel{Width: 80, Height: 24})
	if model.Title != Title {
		t.Fatalf("expected title %q, got %q", Title, model.Title)
	}
	if model.Init() != nil {
		t.Fatal("expected Init to return nil")
	}
}

func TestSetContent(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	if err := model.SetContent("name:   Alice\nage: 30"); err != nil {
		t.Fatalf("SetContent() error = %v", err)
	}
	if model.FormattedContent != "age: 30\nname: Alice\n" {
		t.Fatalf("expected formatted YAML, got %q", model.FormattedContent)
	}
	if !strings.Contains(model.Viewport.GetContent(), "\x1b[") {
		t.Fatalf("expected highlighted viewport content, got %q", model.Viewport.GetContent())
	}
}

func TestSetContentInvalidYAMLReturnsError(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	if err := model.SetContent("{not: valid"); err == nil {
		t.Fatal("expected error for invalid YAML")
	}
	if model.Viewport.GetContent() != "" {
		t.Fatalf("expected empty viewport, got %q", model.Viewport.GetContent())
	}
}

func TestApplyContent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		content     string
		err         error
		wantState   ui.PagerState
		wantMessage string
	}{
		{"valid yaml", "a: 1", nil, ui.PagerStateStatusMessage, "Pasted. Press 'c' to copy result."},
		{"invalid yaml", "{not: valid", nil, ui.PagerStateErrorMessage, ""},
		{"source error", "a: 1", errors.New("no clipboard"), ui.PagerStateErrorMessage, "no clipboard"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			model := newTestModel()
			if model.applyContent(tt.content, tt.err) == nil {
				t.Fatal("expected a status message command")
			}
			if model.State != tt.wantState {
				t.Fatalf("expected state %v, got %v", tt.wantState, model.State)
			}
			if tt.wantMessage != "" && model.StatusMessage != tt.wantMessage {
				t.Fatalf("expected message %q, got %q", tt.wantMessage, model.StatusMessage)
			}
			if tt.err != nil && model.Content != "" {
				t.Fatalf("expected content to stay empty, got %q", model.Content)
			}
		})
	}
}

func TestUpdateEditorFinished(t *testing.T) {
	t.Parallel()

	next, cmd := newTestModel().Update(editor.EditorFinishedMsg{Content: "a:   1"})
	model := next.(YamlModel)
	if cmd == nil {
		t.Fatal("expected a status message command")
	}
	if model.State != ui.PagerStateStatusMessage {
		t.Fatalf("expected status message state, got %v", model.State)
	}
	if model.FormattedContent != "a: 1\n" {
		t.Fatalf("expected formatted YAML, got %q", model.FormattedContent)
	}
}

func TestUpdateEditorFinishedWithError(t *testing.T) {
	t.Parallel()

	next, _ := newTestModel().Update(editor.EditorFinishedMsg{Err: errors.New("editor failed"), Content: "a: 1"})
	model := next.(YamlModel)
	if model.State != ui.PagerStateErrorMessage || model.StatusMessage != "editor failed" {
		t.Fatalf("expected editor error, got state %v message %q", model.State, model.StatusMessage)
	}
	if model.Content != "" {
		t.Fatalf("expected content to stay empty, got %q", model.Content)
	}
}

func TestUpdateStatusMessageTimeout(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	model.State = ui.PagerStateStatusMessage
	next, _ := model.Update(ui.StatusMessageTimeoutMsg{})
	if state := next.(YamlModel).State; state != ui.PagerStateBrowse {
		t.Fatalf("expected browse state, got %v", state)
	}
}

func TestUpdateWindowSize(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	if !model.Ready {
		t.Fatal("expected model to be ready after a window size message")
	}
	if model.Viewport.Height() != 23 {
		t.Fatalf("expected viewport height 23, got %d", model.Viewport.Height())
	}
}

func TestUpdateScrollsViewport(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	if err := model.SetContent(strings.Repeat("- item\n", 50)); err != nil {
		t.Fatalf("SetContent() error = %v", err)
	}
	next, _ := model.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyDown}))
	model = next.(YamlModel)
	if offset := model.Viewport.YOffset(); offset != 1 {
		t.Fatalf("expected viewport offset 1, got %d", offset)
	}
}

func TestViewTogglesHelp(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	if strings.Contains(model.View().Content, "copy formatted YAML") {
		t.Fatal("expected help to be hidden")
	}
	next, cmd := model.Update(tea.KeyPressMsg(tea.Key{Code: '?', Text: "?"}))
	if cmd != nil {
		t.Fatal("expected no command when toggling help")
	}
	view := next.(YamlModel).View()
	if !view.AltScreen {
		t.Fatal("expected alt screen view")
	}
	for _, want := range []string{Title, "copy formatted YAML", "edit unformatted YAML", "paste unformatted YAML"} {
		if !strings.Contains(view.Content, want) {
			t.Fatalf("expected view to contain %q, got %q", want, view.Content)
		}
	}
}
