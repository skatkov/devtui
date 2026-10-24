package jsonstruct

import (
	"errors"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/skatkov/devtui/internal/editor"
	"github.com/skatkov/devtui/internal/ui"
)

func TestEditorFinishedMsgSetsContent(t *testing.T) {
	t.Parallel()

	model := updateModel(t, newModel(), editor.EditorFinishedMsg{Content: `{"id":1}`})

	if !strings.Contains(model.FormattedContent, "type Root struct {\n\tID int") {
		t.Errorf("unexpected formatted content %q", model.FormattedContent)
	}
	if !strings.Contains(model.Viewport.GetContent(), "\x1b[") {
		t.Errorf("expected highlighted viewport content, got %q", model.Viewport.GetContent())
	}
	if model.State != ui.PagerStateStatusMessage || model.StatusMessage != "Pasted. Press 'c' to copy result." {
		t.Errorf("expected status message, got state %v and %q", model.State, model.StatusMessage)
	}
}

func TestEditorFinishedMsgWithErrorShowsError(t *testing.T) {
	t.Parallel()

	model := updateModel(t, newModel(), editor.EditorFinishedMsg{Err: errors.New("editor failed"), Content: `{"id":1}`})

	if model.Content != "" {
		t.Errorf("expected content to stay empty, got %q", model.Content)
	}
	if model.State != ui.PagerStateErrorMessage || model.StatusMessage != "editor failed" {
		t.Errorf("expected editor error, got state %v and %q", model.State, model.StatusMessage)
	}
}

func TestEditorFinishedMsgWithInvalidContentShowsError(t *testing.T) {
	t.Parallel()

	model := updateModel(t, newModel(), editor.EditorFinishedMsg{Content: `{`})

	if model.State != ui.PagerStateErrorMessage || !strings.HasPrefix(model.StatusMessage, "unexpected EOF") {
		t.Errorf("expected conversion error, got state %v and %q", model.State, model.StatusMessage)
	}
}

func TestStatusMessageTimeoutReturnsToBrowse(t *testing.T) {
	t.Parallel()

	model := updateModel(t, newModel(), editor.EditorFinishedMsg{Content: `{"id":1}`})
	model = updateModel(t, model, ui.StatusMessageTimeoutMsg{})

	if model.State != ui.PagerStateBrowse {
		t.Errorf("expected browse state, got %v", model.State)
	}
}

func TestViewShowsHelp(t *testing.T) {
	t.Parallel()

	model := updateModel(t, newModel(), tea.WindowSizeMsg{Width: 80, Height: 24})
	if strings.Contains(model.View().Content, "copy Go struct") {
		t.Fatal("expected help to be hidden")
	}

	model = updateModel(t, model, textKeypress('?'))
	if !strings.Contains(model.View().Content, "copy Go struct") {
		t.Errorf("expected help in view, got %q", model.View().Content)
	}
}

func newModel() JsonStructModel {
	return NewJsonStructModel(&ui.CommonModel{Width: 80, Height: 24})
}

func updateModel(t *testing.T, model JsonStructModel, msg tea.Msg) JsonStructModel {
	t.Helper()

	next, _ := model.Update(msg)
	return next.(JsonStructModel)
}

func textKeypress(r rune) tea.KeyPressMsg {
	return tea.KeyPressMsg(tea.Key{Code: r, Text: string(r)})
}
