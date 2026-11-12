package csv2md

import (
	"errors"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/skatkov/devtui/internal/editor"
	"github.com/skatkov/devtui/internal/ui"
)

const (
	testCSV       = "name,age\nAl,3"
	testMarkdown  = "| name | age | \n| ---- | --- | \n| Al   | 3   | "
	testUnaligned = "| name | age | \n| --- | --- | \n| Al | 3 | "
)

func TestHelpViewListsKeys(t *testing.T) {
	t.Parallel()

	m := NewCSV2MDModel(&ui.CommonModel{Width: 80, Height: 24})
	help := m.helpView()
	for _, want := range []string{"copy markdown", "edit", "paste to convert", "toggle column alignment", "quit"} {
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
	if m.FormattedContent != testMarkdown {
		t.Errorf("expected formatted content %q, got %q", testMarkdown, m.FormattedContent)
	}
	if got := m.Viewport.GetContent(); got != testMarkdown {
		t.Errorf("expected viewport content %q, got %q", testMarkdown, got)
	}
}

func TestPasteConvertsContent(t *testing.T) {
	t.Parallel()

	m := newSizedModel()
	if cmd := m.paste(testCSV, nil); cmd == nil {
		t.Fatal("expected status message command")
	}
	assertStatus(t, m, ui.PagerStateStatusMessage, "Pasted. Press 'c' to copy result.")
	if m.FormattedContent != testMarkdown {
		t.Errorf("expected formatted content %q, got %q", testMarkdown, m.FormattedContent)
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
	got := next.(CSV2MDModel)
	assertStatus(t, &got, ui.PagerStateStatusMessage, "Converted. Press 'c' to copy result.")
	if got.FormattedContent != testMarkdown {
		t.Errorf("expected formatted content %q, got %q", testMarkdown, got.FormattedContent)
	}
}

func TestEditorFinishedMsgShowsEditorError(t *testing.T) {
	t.Parallel()

	m := newSizedModel()
	next, _ := m.Update(editor.EditorFinishedMsg{Err: errors.New("editor failed"), Content: testCSV})
	got := next.(CSV2MDModel)
	assertStatus(t, &got, ui.PagerStateErrorMessage, "editor failed")
	if got.Content != "" {
		t.Errorf("expected content to stay empty, got %q", got.Content)
	}
}

func TestSetContentInvalidReturnsError(t *testing.T) {
	t.Parallel()

	m := newSizedModel()
	if err := m.SetContent("a\n\"x"); err == nil {
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
	assertStatus(t, m, ui.PagerStateErrorMessage, "empty content")
}

func TestEditorFinishedMsgShowsConvertError(t *testing.T) {
	t.Parallel()

	m := newSizedModel()
	next, _ := m.Update(editor.EditorFinishedMsg{Content: ""})
	got := next.(CSV2MDModel)
	assertStatus(t, &got, ui.PagerStateErrorMessage, "empty content")
}

func TestAlignKeyTogglesColumnAlignment(t *testing.T) {
	t.Parallel()

	m := newSizedModel()
	_ = m.SetContent(testCSV)
	next, _ := m.Update(alignKeypress())
	got := next.(CSV2MDModel)
	assertStatus(t, &got, ui.PagerStateStatusMessage, "Columns unaligned")
	if got.FormattedContent != testUnaligned {
		t.Errorf("expected formatted content %q, got %q", testUnaligned, got.FormattedContent)
	}

	next, _ = got.Update(alignKeypress())
	got = next.(CSV2MDModel)
	assertStatus(t, &got, ui.PagerStateStatusMessage, "Columns aligned")
	if got.FormattedContent != testMarkdown {
		t.Errorf("expected formatted content %q, got %q", testMarkdown, got.FormattedContent)
	}
}

func TestAlignKeyWithoutContentOnlyToggles(t *testing.T) {
	t.Parallel()

	m := newSizedModel()
	next, _ := m.Update(alignKeypress())
	got := next.(CSV2MDModel)
	if got.alignColumns {
		t.Error("expected columns to be unaligned")
	}
	if got.State != ui.PagerStateBrowse {
		t.Errorf("expected browse state, got %v", got.State)
	}
}

func TestAlignKeyShowsConvertError(t *testing.T) {
	t.Parallel()

	m := newSizedModel()
	_ = m.SetContent("a\n\"x")
	next, _ := m.Update(alignKeypress())
	got := next.(CSV2MDModel)
	if got.State != ui.PagerStateErrorMessage {
		t.Errorf("expected error state, got %v", got.State)
	}
}

func TestStatusMessageTimeoutReturnsToBrowse(t *testing.T) {
	t.Parallel()

	m := newSizedModel()
	m.State = ui.PagerStateStatusMessage
	next, _ := m.Update(ui.StatusMessageTimeoutMsg{})
	if got := next.(CSV2MDModel).State; got != ui.PagerStateBrowse {
		t.Errorf("expected browse state, got %v", got)
	}
}

func TestKeyPressScrollsViewport(t *testing.T) {
	t.Parallel()

	m := newSizedModel()
	_ = m.SetContent("n\n" + strings.Repeat("1\n", 30))
	next, _ := m.Update(codeKeypress(tea.KeyDown))
	got := next.(CSV2MDModel)
	if offset := got.Viewport.YOffset(); offset != 1 {
		t.Errorf("expected viewport offset 1, got %d", offset)
	}
}

func TestViewShowsHelp(t *testing.T) {
	t.Parallel()

	m := newSizedModel()
	if strings.Contains(m.View().Content, "toggle column alignment") {
		t.Fatal("expected help to be hidden")
	}
	next, _ := m.Update(tea.KeyPressMsg(tea.Key{Code: '?', Text: "?"}))
	if !strings.Contains(next.(CSV2MDModel).View().Content, "toggle column alignment") {
		t.Error("expected help to be shown")
	}
}

func newSizedModel() *CSV2MDModel {
	m := NewCSV2MDModel(&ui.CommonModel{Width: 80, Height: 24})
	next, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	model := next.(CSV2MDModel)
	return &model
}

func assertStatus(t *testing.T, m *CSV2MDModel, state ui.PagerState, message string) {
	t.Helper()

	if m.State != state {
		t.Errorf("expected state %v, got %v", state, m.State)
	}
	if m.StatusMessage != message {
		t.Errorf("expected status message %q, got %q", message, m.StatusMessage)
	}
}

func alignKeypress() tea.KeyPressMsg {
	return tea.KeyPressMsg(tea.Key{Code: 'a', Text: "a"})
}

func codeKeypress(code rune) tea.KeyPressMsg {
	return tea.KeyPressMsg(tea.Key{Code: code})
}
