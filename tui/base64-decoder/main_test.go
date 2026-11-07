package base64decoder

import (
	"encoding/base64"
	"errors"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/skatkov/devtui/internal/editor"
	"github.com/skatkov/devtui/internal/ui"
)

func TestHelpViewDoesNotPanic(t *testing.T) {
	t.Parallel()

	m := NewBase64Model(&ui.CommonModel{Width: 80, Height: 24})
	if m.helpView() == "" {
		t.Fatal("expected non-empty help view")
	}
}

func TestHelpViewListsKeys(t *testing.T) {
	t.Parallel()

	m := NewBase64Model(&ui.CommonModel{Width: 80, Height: 24})
	help := m.helpView()
	for _, want := range []string{"copy text", "edit base64", "paste base64 to decode", "quit"} {
		if !strings.Contains(help, want) {
			t.Errorf("expected help to contain %q, got %q", want, help)
		}
	}
}

func TestSetContentDecodes(t *testing.T) {
	t.Parallel()

	m := newSizedModel()
	if err := m.SetContent("aGVsbG8="); err != nil {
		t.Fatalf("SetContent() error = %v", err)
	}
	if m.FormattedContent != "hello" {
		t.Errorf("expected formatted content %q, got %q", "hello", m.FormattedContent)
	}
	if got := m.Viewport.GetContent(); got != "hello" {
		t.Errorf("expected viewport content %q, got %q", "hello", got)
	}
}

func TestSetContentEmptyClearsOutput(t *testing.T) {
	t.Parallel()

	m := newSizedModel()
	_ = m.SetContent("aGVsbG8=")
	if err := m.SetContent(""); err != nil {
		t.Fatalf("SetContent() error = %v", err)
	}
	if m.Content != "" || m.FormattedContent != "" || m.Viewport.GetContent() != "" {
		t.Errorf("expected empty model, got content %q, formatted %q", m.Content, m.FormattedContent)
	}
}

func TestSetContentInvalidReturnsError(t *testing.T) {
	t.Parallel()

	m := newSizedModel()
	if err := m.SetContent("not base64!"); err == nil {
		t.Fatal("expected error for invalid base64")
	}
	if m.Content != "not base64!" {
		t.Errorf("expected content to be kept, got %q", m.Content)
	}
}

func TestPasteDecodesContent(t *testing.T) {
	t.Parallel()

	m := newSizedModel()
	if cmd := m.paste("aGVsbG8=", nil); cmd == nil {
		t.Fatal("expected status message command")
	}
	assertStatus(t, m, ui.PagerStateStatusMessage, "Pasted and decoded. Press 'c' to copy result.")
	if m.FormattedContent != "hello" {
		t.Errorf("expected formatted content %q, got %q", "hello", m.FormattedContent)
	}
}

func TestPasteShowsClipboardError(t *testing.T) {
	t.Parallel()

	m := newSizedModel()
	m.paste("aGVsbG8=", errors.New("no clipboard"))
	assertStatus(t, m, ui.PagerStateErrorMessage, "no clipboard")
	if m.Content != "" {
		t.Errorf("expected content to stay empty, got %q", m.Content)
	}
}

func TestPasteShowsDecodeError(t *testing.T) {
	t.Parallel()

	m := newSizedModel()
	m.paste("not base64!", nil)
	if m.State != ui.PagerStateErrorMessage {
		t.Errorf("expected error state, got %v", m.State)
	}
}

func TestEditorFinishedMsgDecodesContent(t *testing.T) {
	t.Parallel()

	m := newSizedModel()
	next, cmd := m.Update(editor.EditorFinishedMsg{Content: "aGVsbG8="})
	if cmd == nil {
		t.Fatal("expected status message command")
	}
	got := next.(Base64Model)
	assertStatus(t, &got, ui.PagerStateStatusMessage, "Content decoded. Press 'c' to copy result.")
	if got.FormattedContent != "hello" {
		t.Errorf("expected formatted content %q, got %q", "hello", got.FormattedContent)
	}
}

func TestEditorFinishedMsgShowsEditorError(t *testing.T) {
	t.Parallel()

	m := newSizedModel()
	next, _ := m.Update(editor.EditorFinishedMsg{Err: errors.New("editor failed"), Content: "aGVsbG8="})
	got := next.(Base64Model)
	assertStatus(t, &got, ui.PagerStateErrorMessage, "editor failed")
	if got.Content != "" {
		t.Errorf("expected content to stay empty, got %q", got.Content)
	}
}

func TestEditorFinishedMsgShowsDecodeError(t *testing.T) {
	t.Parallel()

	m := newSizedModel()
	next, _ := m.Update(editor.EditorFinishedMsg{Content: "not base64!"})
	got := next.(Base64Model)
	if got.State != ui.PagerStateErrorMessage {
		t.Errorf("expected error state, got %v", got.State)
	}
}

func TestStatusMessageTimeoutReturnsToBrowse(t *testing.T) {
	t.Parallel()

	m := newSizedModel()
	m.State = ui.PagerStateStatusMessage
	next, _ := m.Update(ui.StatusMessageTimeoutMsg{})
	if got := next.(Base64Model).State; got != ui.PagerStateBrowse {
		t.Errorf("expected browse state, got %v", got)
	}
}

func TestKeyPressScrollsViewport(t *testing.T) {
	t.Parallel()

	m := newSizedModel()
	_ = m.SetContent(base64.StdEncoding.EncodeToString([]byte(strings.Repeat("line\n", 50))))
	next, _ := m.Update(codeKeypress(tea.KeyDown))
	got := next.(Base64Model)
	if offset := got.Viewport.YOffset(); offset != 1 {
		t.Errorf("expected viewport offset 1, got %d", offset)
	}
}

func TestViewShowsHelp(t *testing.T) {
	t.Parallel()

	m := newSizedModel()
	if strings.Contains(m.View().Content, "paste base64 to decode") {
		t.Fatal("expected help to be hidden")
	}
	next, _ := m.Update(tea.KeyPressMsg(tea.Key{Code: '?', Text: "?"}))
	if !strings.Contains(next.(Base64Model).View().Content, "paste base64 to decode") {
		t.Error("expected help to be shown")
	}
}

func newSizedModel() *Base64Model {
	m := NewBase64Model(&ui.CommonModel{Width: 80, Height: 24})
	next, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	model := next.(Base64Model)
	return &model
}

func assertStatus(t *testing.T, m *Base64Model, state ui.PagerState, message string) {
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
