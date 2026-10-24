package csv2json

import (
	"errors"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/skatkov/devtui/internal/editor"
	"github.com/skatkov/devtui/internal/ui"
)

const (
	testCSV  = "a,b\n1,2"
	testJSON = "[\n  {\n    \"a\": \"1\",\n    \"b\": \"2\"\n  }\n]"
)

func TestHelpViewListsKeys(t *testing.T) {
	t.Parallel()

	m := NewCSVJsonModel(&ui.CommonModel{Width: 80, Height: 24})
	help := m.helpView()
	for _, want := range []string{"copy JSON", "edit CSV", "paste CSV to convert", "quit"} {
		if !strings.Contains(help, want) {
			t.Errorf("expected help to contain %q, got %q", want, help)
		}
	}
}

func TestSetContentConvertsCSV(t *testing.T) {
	t.Parallel()

	m := newSizedModel()
	if err := m.SetContent(testCSV); err != nil {
		t.Fatalf("SetContent() error = %v", err)
	}
	if m.Content != testCSV {
		t.Errorf("expected content %q, got %q", testCSV, m.Content)
	}
	if m.FormattedContent != testJSON {
		t.Errorf("expected formatted content %q, got %q", testJSON, m.FormattedContent)
	}
	if got := m.Viewport.GetContent(); !strings.Contains(got, "\"b\"") {
		t.Errorf("expected highlighted JSON in viewport, got %q", got)
	}
}

func TestPasteConvertsContent(t *testing.T) {
	t.Parallel()

	m := newSizedModel()
	if cmd := m.paste(testCSV, nil); cmd == nil {
		t.Fatal("expected status message command")
	}
	assertStatus(t, m, ui.PagerStateStatusMessage, "Pasted. Press 'c' to copy result.")
	if m.FormattedContent != testJSON {
		t.Errorf("expected formatted content %q, got %q", testJSON, m.FormattedContent)
	}
}

func TestPasteShowsClipboardError(t *testing.T) {
	t.Parallel()

	m := newSizedModel()
	m.paste(testCSV, errors.New("no clipboard"))
	assertStatus(t, m, ui.PagerStateErrorMessage, "no clipboard")
	if m.Content != "" {
		t.Errorf("expected content to stay empty, got %q", m.Content)
	}
}

func TestEditorFinishedMsgConvertsContent(t *testing.T) {
	t.Parallel()

	m := newSizedModel()
	next, cmd := m.Update(editor.EditorFinishedMsg{Content: testCSV})
	if cmd == nil {
		t.Fatal("expected status message command")
	}
	got := next.(CSVJsonModel)
	assertStatus(t, &got, ui.PagerStateStatusMessage, "Pasted. Press 'c' to copy result.")
	if got.FormattedContent != testJSON {
		t.Errorf("expected formatted content %q, got %q", testJSON, got.FormattedContent)
	}
}

func TestEditorFinishedMsgShowsEditorError(t *testing.T) {
	t.Parallel()

	m := newSizedModel()
	next, _ := m.Update(editor.EditorFinishedMsg{Err: errors.New("editor failed"), Content: testCSV})
	got := next.(CSVJsonModel)
	assertStatus(t, &got, ui.PagerStateErrorMessage, "editor failed")
	if got.Content != "" {
		t.Errorf("expected content to stay empty, got %q", got.Content)
	}
}

func TestSetContentInvalidReturnsError(t *testing.T) {
	t.Parallel()

	m := newSizedModel()
	if err := m.SetContent("a,b\n1,2,3"); err == nil {
		t.Fatal("expected error for invalid CSV")
	}
	if m.FormattedContent != "" {
		t.Errorf("expected empty formatted content, got %q", m.FormattedContent)
	}
}

func TestPasteShowsConvertError(t *testing.T) {
	t.Parallel()

	m := newSizedModel()
	m.paste("", nil)
	assertStatus(t, m, ui.PagerStateErrorMessage, "empty CSV file")
}

func TestEditorFinishedMsgShowsConvertError(t *testing.T) {
	t.Parallel()

	m := newSizedModel()
	next, _ := m.Update(editor.EditorFinishedMsg{Content: ""})
	got := next.(CSVJsonModel)
	assertStatus(t, &got, ui.PagerStateErrorMessage, "empty CSV file")
}

func TestStatusMessageTimeoutReturnsToBrowse(t *testing.T) {
	t.Parallel()

	m := newSizedModel()
	m.State = ui.PagerStateStatusMessage
	next, _ := m.Update(ui.StatusMessageTimeoutMsg{})
	if got := next.(CSVJsonModel).State; got != ui.PagerStateBrowse {
		t.Errorf("expected browse state, got %v", got)
	}
}

func TestKeyPressScrollsViewport(t *testing.T) {
	t.Parallel()

	m := newSizedModel()
	_ = m.SetContent("n\n" + strings.Repeat("1\n", 30))
	next, _ := m.Update(codeKeypress(tea.KeyDown))
	got := next.(CSVJsonModel)
	if offset := got.Viewport.YOffset(); offset != 1 {
		t.Errorf("expected viewport offset 1, got %d", offset)
	}
}

func TestViewShowsHelp(t *testing.T) {
	t.Parallel()

	m := newSizedModel()
	if strings.Contains(m.View().Content, "paste CSV to convert") {
		t.Fatal("expected help to be hidden")
	}
	next, _ := m.Update(tea.KeyPressMsg(tea.Key{Code: '?', Text: "?"}))
	if !strings.Contains(next.(CSVJsonModel).View().Content, "paste CSV to convert") {
		t.Error("expected help to be shown")
	}
}

func newSizedModel() *CSVJsonModel {
	m := NewCSVJsonModel(&ui.CommonModel{Width: 80, Height: 24})
	next, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	model := next.(CSVJsonModel)
	return &model
}

func assertStatus(t *testing.T, m *CSVJsonModel, state ui.PagerState, message string) {
	t.Helper()

	if m.State != state {
		t.Errorf("expected state %v, got %v", state, m.State)
	}
	if m.StatusMessage != message {
		t.Errorf("expected status message %q, got %q", message, m.StatusMessage)
	}
}

func codeKeypress(code rune) tea.KeyPressMsg {
	return tea.KeyPressMsg(tea.Key{Code: code})
}
