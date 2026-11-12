package base64encoder

import (
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
	for _, want := range []string{"copy base64", "edit text", "paste text to encode", "quit"} {
		if !strings.Contains(help, want) {
			t.Errorf("expected help to contain %q, got %q", want, help)
		}
	}
}

func TestSetContentEncodes(t *testing.T) {
	t.Parallel()

	m := newSizedModel(80, 24)
	if err := m.SetContent("hello"); err != nil {
		t.Fatalf("SetContent() error = %v", err)
	}
	if m.FormattedContent != "aGVsbG8=" {
		t.Errorf("expected formatted content %q, got %q", "aGVsbG8=", m.FormattedContent)
	}
	if got := m.Viewport.GetContent(); got != "aGVsbG8=" {
		t.Errorf("expected viewport content %q, got %q", "aGVsbG8=", got)
	}
}

func TestSetContentEmptyClearsOutput(t *testing.T) {
	t.Parallel()

	m := newSizedModel(80, 24)
	_ = m.SetContent("hello")
	if err := m.SetContent(""); err != nil {
		t.Fatalf("SetContent() error = %v", err)
	}
	if m.Content != "" || m.FormattedContent != "" || m.Viewport.GetContent() != "" {
		t.Errorf("expected empty model, got content %q, formatted %q", m.Content, m.FormattedContent)
	}
}

func TestSetContentWrapsToWidth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		width     int
		input     string
		wantLines []int
	}{
		{"width minus padding", 80, strings.Repeat("a", 60), []int{76, 4}},
		{"minimum width", 10, strings.Repeat("a", 30), []int{20, 20}},
		{"exact multiple", 24, strings.Repeat("a", 15), []int{20}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := newSizedModel(tt.width, 24)
			_ = m.SetContent(tt.input)
			assertLineLengths(t, m.Viewport.GetContent(), tt.wantLines)
			if strings.Contains(m.FormattedContent, "\n") {
				t.Errorf("expected unwrapped formatted content, got %q", m.FormattedContent)
			}
		})
	}
}

func TestWindowSizeMsgCreatesViewport(t *testing.T) {
	t.Parallel()

	m := newSizedModel(80, 24)
	if !m.Ready {
		t.Fatal("expected model to be ready")
	}
	if got := m.Viewport.Height(); got != 23 {
		t.Errorf("expected viewport height 23, got %d", got)
	}
	if got := m.Viewport.Width(); got != 80 {
		t.Errorf("expected viewport width 80, got %d", got)
	}
}

func TestWindowSizeMsgRewrapsContent(t *testing.T) {
	t.Parallel()

	m := newSizedModel(80, 24)
	_ = m.SetContent(strings.Repeat("a", 30))
	next, _ := m.Update(tea.WindowSizeMsg{Width: 24, Height: 24})
	got := next.(Base64Model)
	assertLineLengths(t, got.Viewport.GetContent(), []int{20, 20})
	if got.Viewport.Width() != 24 {
		t.Errorf("expected viewport width 24, got %d", got.Viewport.Width())
	}
}

func TestWindowSizeMsgKeepsScrollPosition(t *testing.T) {
	t.Parallel()

	m := newSizedModel(24, 5)
	_ = m.SetContent(strings.Repeat("a", 300))
	next, _ := m.Update(codeKeypress(tea.KeyDown))
	next, _ = next.Update(tea.WindowSizeMsg{Width: 24, Height: 5})
	got := next.(Base64Model)
	if offset := got.Viewport.YOffset(); offset != 1 {
		t.Errorf("expected viewport offset 1, got %d", offset)
	}
}

func TestPasteEncodesContent(t *testing.T) {
	t.Parallel()

	m := newSizedModel(80, 24)
	if cmd := m.paste("hello", nil); cmd == nil {
		t.Fatal("expected status message command")
	}
	assertStatus(t, m, ui.PagerStateStatusMessage, "Pasted and encoded. Press 'c' to copy result.")
	if m.FormattedContent != "aGVsbG8=" {
		t.Errorf("expected formatted content %q, got %q", "aGVsbG8=", m.FormattedContent)
	}
}

func TestPasteShowsClipboardError(t *testing.T) {
	t.Parallel()

	m := newSizedModel(80, 24)
	m.paste("hello", errors.New("no clipboard"))
	assertStatus(t, m, ui.PagerStateErrorMessage, "no clipboard")
	if m.Content != "" {
		t.Errorf("expected content to stay empty, got %q", m.Content)
	}
}

func TestEditorFinishedMsgEncodesContent(t *testing.T) {
	t.Parallel()

	m := newSizedModel(80, 24)
	next, cmd := m.Update(editor.EditorFinishedMsg{Content: "hello"})
	if cmd == nil {
		t.Fatal("expected status message command")
	}
	got := next.(Base64Model)
	assertStatus(t, &got, ui.PagerStateStatusMessage, "Content encoded. Press 'c' to copy result.")
	if got.FormattedContent != "aGVsbG8=" {
		t.Errorf("expected formatted content %q, got %q", "aGVsbG8=", got.FormattedContent)
	}
}

func TestEditorFinishedMsgShowsEditorError(t *testing.T) {
	t.Parallel()

	m := newSizedModel(80, 24)
	next, _ := m.Update(editor.EditorFinishedMsg{Err: errors.New("editor failed"), Content: "hello"})
	got := next.(Base64Model)
	assertStatus(t, &got, ui.PagerStateErrorMessage, "editor failed")
	if got.Content != "" {
		t.Errorf("expected content to stay empty, got %q", got.Content)
	}
}

func TestStatusMessageTimeoutReturnsToBrowse(t *testing.T) {
	t.Parallel()

	m := newSizedModel(80, 24)
	m.State = ui.PagerStateStatusMessage
	next, _ := m.Update(ui.StatusMessageTimeoutMsg{})
	if got := next.(Base64Model).State; got != ui.PagerStateBrowse {
		t.Errorf("expected browse state, got %v", got)
	}
}

func TestKeyPressScrollsViewport(t *testing.T) {
	t.Parallel()

	m := newSizedModel(24, 5)
	_ = m.SetContent(strings.Repeat("a", 300))
	next, _ := m.Update(codeKeypress(tea.KeyDown))
	got := next.(Base64Model)
	if offset := got.Viewport.YOffset(); offset != 1 {
		t.Errorf("expected viewport offset 1, got %d", offset)
	}
}

func TestViewShowsHelp(t *testing.T) {
	t.Parallel()

	m := newSizedModel(80, 24)
	if strings.Contains(m.View().Content, "paste text to encode") {
		t.Fatal("expected help to be hidden")
	}
	next, _ := m.Update(tea.KeyPressMsg(tea.Key{Code: '?', Text: "?"}))
	if !strings.Contains(next.(Base64Model).View().Content, "paste text to encode") {
		t.Error("expected help to be shown")
	}
}

func newSizedModel(width, height int) *Base64Model {
	m := NewBase64Model(&ui.CommonModel{Width: width, Height: height})
	next, _ := m.Update(tea.WindowSizeMsg{Width: width, Height: height})
	model := next.(Base64Model)
	return &model
}

func assertLineLengths(t *testing.T, content string, want []int) {
	t.Helper()

	lines := strings.Split(content, "\n")
	if len(lines) != len(want) {
		t.Fatalf("expected %d lines, got %d: %q", len(want), len(lines), content)
	}
	for i, line := range lines {
		if len(line) != want[i] {
			t.Errorf("expected line %d to have length %d, got %d", i, want[i], len(line))
		}
	}
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
