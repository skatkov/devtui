package xml

import (
	"errors"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/skatkov/devtui/internal/editor"
	"github.com/skatkov/devtui/internal/ui"
)

const formattedXML = "\t\n\t<a>\n\t  <b>1</b>\n\t</a>"

func newTestModel() XMLFormatterModel {
	model := NewXMLFormatterModel(&ui.CommonModel{Width: 80, Height: 24})
	next, _ := model.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	return next.(XMLFormatterModel)
}

func TestNewXMLFormatterModel(t *testing.T) {
	t.Parallel()

	model := NewXMLFormatterModel(&ui.CommonModel{Width: 80, Height: 24})
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
	if err := model.SetContent("<a><b>1</b></a>"); err != nil {
		t.Fatalf("SetContent() error = %v", err)
	}
	if model.FormattedContent != formattedXML {
		t.Fatalf("expected formatted XML, got %q", model.FormattedContent)
	}
	if !strings.Contains(model.Viewport.GetContent(), "\x1b[") {
		t.Fatalf("expected highlighted viewport content, got %q", model.Viewport.GetContent())
	}
}

func TestApplyContent(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	if model.applyContent("<a><b>1</b></a>", nil) == nil {
		t.Fatal("expected a status message command")
	}
	if model.State != ui.PagerStateStatusMessage || model.StatusMessage != "Pasted. Press 'c' to copy result." {
		t.Fatalf("expected pasted status, got state %v message %q", model.State, model.StatusMessage)
	}
	if model.FormattedContent != formattedXML {
		t.Fatalf("expected formatted XML, got %q", model.FormattedContent)
	}
}

func TestApplyContentWithError(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	if model.applyContent("<a/>", errors.New("no clipboard")) == nil {
		t.Fatal("expected an error message command")
	}
	if model.State != ui.PagerStateErrorMessage || model.StatusMessage != "no clipboard" {
		t.Fatalf("expected clipboard error, got state %v message %q", model.State, model.StatusMessage)
	}
	if model.Content != "" {
		t.Fatalf("expected content to stay empty, got %q", model.Content)
	}
}

func TestUpdateEditorFinished(t *testing.T) {
	t.Parallel()

	next, cmd := newTestModel().Update(editor.EditorFinishedMsg{Content: "<a><b>1</b></a>"})
	model := next.(XMLFormatterModel)
	if cmd == nil {
		t.Fatal("expected a status message command")
	}
	if model.State != ui.PagerStateStatusMessage {
		t.Fatalf("expected status message state, got %v", model.State)
	}
	if model.FormattedContent != formattedXML {
		t.Fatalf("expected formatted XML, got %q", model.FormattedContent)
	}
}

func TestUpdateEditorFinishedWithError(t *testing.T) {
	t.Parallel()

	next, _ := newTestModel().Update(editor.EditorFinishedMsg{Err: errors.New("editor failed"), Content: "<a/>"})
	model := next.(XMLFormatterModel)
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
	if state := next.(XMLFormatterModel).State; state != ui.PagerStateBrowse {
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
	if err := model.SetContent("<a>" + strings.Repeat("<b>1</b>", 50) + "</a>"); err != nil {
		t.Fatalf("SetContent() error = %v", err)
	}
	next, _ := model.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyDown}))
	model = next.(XMLFormatterModel)
	if offset := model.Viewport.YOffset(); offset != 1 {
		t.Fatalf("expected viewport offset 1, got %d", offset)
	}
}

func TestViewTogglesHelp(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	if strings.Contains(model.View().Content, "copy formatted XML") {
		t.Fatal("expected help to be hidden")
	}
	next, cmd := model.Update(tea.KeyPressMsg(tea.Key{Code: '?', Text: "?"}))
	if cmd != nil {
		t.Fatal("expected no command when toggling help")
	}
	view := next.(XMLFormatterModel).View()
	if !view.AltScreen {
		t.Fatal("expected alt screen view")
	}
	for _, want := range []string{Title, "copy formatted XML", "edit XML", "paste XML to format"} {
		if !strings.Contains(view.Content, want) {
			t.Fatalf("expected view to contain %q, got %q", want, view.Content)
		}
	}
}
