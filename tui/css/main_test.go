package css

import (
	"errors"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/skatkov/devtui/internal/editor"
	"github.com/skatkov/devtui/internal/ui"
)

const (
	testCSS          = "a{color:red}"
	testFormattedCSS = "a {\n  color: red;\n}\n"
)

func TestHelpViewListsKeys(t *testing.T) {
	t.Parallel()

	m := NewCSSFormatterModel(&ui.CommonModel{Width: 80, Height: 24})
	help := m.helpView()
	for _, want := range []string{"copy formatted CSS", "edit CSS", "paste CSS to format", "quit"} {
		if !strings.Contains(help, want) {
			t.Errorf("expected help to contain %q, got %q", want, help)
		}
	}
}

func TestSetContentFormatsCSS(t *testing.T) {
	t.Parallel()

	m := newSizedModel()
	if err := m.SetContent(testCSS); err != nil {
		t.Fatalf("SetContent() error = %v", err)
	}
	if m.Content != testCSS {
		t.Errorf("expected content %q, got %q", testCSS, m.Content)
	}
	if m.FormattedContent != testFormattedCSS {
		t.Errorf("expected formatted content %q, got %q", testFormattedCSS, m.FormattedContent)
	}
	if got := m.Viewport.GetContent(); !strings.Contains(got, "color") {
		t.Errorf("expected highlighted CSS in viewport, got %q", got)
	}
}

func TestPasteFormatsContent(t *testing.T) {
	t.Parallel()

	m := newSizedModel()
	if cmd := m.paste(testCSS, nil); cmd == nil {
		t.Fatal("expected status message command")
	}
	assertStatus(t, m, ui.PagerStateStatusMessage, "Pasted. Press 'c' to copy result.")
	if m.FormattedContent != testFormattedCSS {
		t.Errorf("expected formatted content %q, got %q", testFormattedCSS, m.FormattedContent)
	}
}

func TestPasteShowsClipboardError(t *testing.T) {
	t.Parallel()

	m := newSizedModel()
	m.paste(testCSS, errors.New("no clipboard"))
	assertStatus(t, m, ui.PagerStateErrorMessage, "no clipboard")
	if m.Content != "" {
		t.Errorf("expected content to stay empty, got %q", m.Content)
	}
}

func TestEditorFinishedMsgFormatsContent(t *testing.T) {
	t.Parallel()

	m := newSizedModel()
	next, cmd := m.Update(editor.EditorFinishedMsg{Content: testCSS})
	if cmd == nil {
		t.Fatal("expected status message command")
	}
	got := next.(CSSFormatterModel)
	assertStatus(t, &got, ui.PagerStateStatusMessage, "Pasted. Press 'c' to copy result.")
	if got.FormattedContent != testFormattedCSS {
		t.Errorf("expected formatted content %q, got %q", testFormattedCSS, got.FormattedContent)
	}
}

func TestEditorFinishedMsgShowsEditorError(t *testing.T) {
	t.Parallel()

	m := newSizedModel()
	next, _ := m.Update(editor.EditorFinishedMsg{Err: errors.New("editor failed"), Content: testCSS})
	got := next.(CSSFormatterModel)
	assertStatus(t, &got, ui.PagerStateErrorMessage, "editor failed")
	if got.Content != "" {
		t.Errorf("expected content to stay empty, got %q", got.Content)
	}
}

func TestStatusMessageTimeoutReturnsToBrowse(t *testing.T) {
	t.Parallel()

	m := newSizedModel()
	m.State = ui.PagerStateStatusMessage
	next, _ := m.Update(ui.StatusMessageTimeoutMsg{})
	if got := next.(CSSFormatterModel).State; got != ui.PagerStateBrowse {
		t.Errorf("expected browse state, got %v", got)
	}
}

func TestKeyPressScrollsViewport(t *testing.T) {
	t.Parallel()

	m := newSizedModel()
	_ = m.SetContent(strings.Repeat(testCSS, 20))
	next, _ := m.Update(codeKeypress(tea.KeyDown))
	got := next.(CSSFormatterModel)
	if offset := got.Viewport.YOffset(); offset != 1 {
		t.Errorf("expected viewport offset 1, got %d", offset)
	}
}

func TestViewShowsHelp(t *testing.T) {
	t.Parallel()

	m := newSizedModel()
	if strings.Contains(m.View().Content, "paste CSS to format") {
		t.Fatal("expected help to be hidden")
	}
	next, _ := m.Update(tea.KeyPressMsg(tea.Key{Code: '?', Text: "?"}))
	if !strings.Contains(next.(CSSFormatterModel).View().Content, "paste CSS to format") {
		t.Error("expected help to be shown")
	}
}

func TestViewEnablesMouse(t *testing.T) {
	t.Parallel()

	m := newSizedModel()
	if got := m.View().MouseMode; got != tea.MouseModeCellMotion {
		t.Errorf("expected cell motion mouse mode, got %v", got)
	}
}

func newSizedModel() *CSSFormatterModel {
	m := NewCSSFormatterModel(&ui.CommonModel{Width: 80, Height: 24})
	next, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	model := next.(CSSFormatterModel)
	return &model
}

func assertStatus(t *testing.T, m *CSSFormatterModel, state ui.PagerState, message string) {
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
