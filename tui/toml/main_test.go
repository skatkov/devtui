package toml

import (
	"errors"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/skatkov/devtui/internal/editor"
	"github.com/skatkov/devtui/internal/ui"
)

func TestConvert(t *testing.T) {
	t.Parallel()

	got, err := Convert("b=2\na =   1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := "a = 1\nb = 2\n"; got != want {
		t.Fatalf("Convert() = %q, want %q", got, want)
	}
}

func TestConvertInvalidTOMLReturnsError(t *testing.T) {
	t.Parallel()

	_, err := Convert("a = ")
	if err == nil || !strings.Contains(err.Error(), "TOML parsing error") {
		t.Fatalf("expected parsing error, got %v", err)
	}
}

func TestNewTomlFormatModel(t *testing.T) {
	t.Parallel()

	m := NewTomlFormatModel(&ui.CommonModel{Width: 80, Height: 24})
	if m.Title != Title {
		t.Fatalf("expected title %q, got %q", Title, m.Title)
	}
	if m.Init() != nil {
		t.Fatal("expected nil init command")
	}
}

func TestSetContent(t *testing.T) {
	t.Parallel()

	m := NewTomlFormatModel(&ui.CommonModel{Width: 80, Height: 24})
	if err := m.SetContent("a =   1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.FormattedContent != "a = 1\n" {
		t.Fatalf("expected formatted content %q, got %q", "a = 1\n", m.FormattedContent)
	}
	if !strings.Contains(m.Viewport.GetContent(), "a") {
		t.Fatalf("expected highlighted content in viewport, got %q", m.Viewport.GetContent())
	}
}

func TestSetContentInvalidTOMLReturnsError(t *testing.T) {
	t.Parallel()

	m := NewTomlFormatModel(&ui.CommonModel{Width: 80, Height: 24})
	if err := m.SetContent("a = "); err == nil {
		t.Fatal("expected error for invalid TOML")
	}
	if m.Content != "a = " || m.FormattedContent != "" {
		t.Fatalf("expected raw content only, got %q and %q", m.Content, m.FormattedContent)
	}
}

func TestEditorFinishedSetsContent(t *testing.T) {
	t.Parallel()

	m := NewTomlFormatModel(&ui.CommonModel{Width: 80, Height: 24})
	next, cmd := m.Update(editor.EditorFinishedMsg{Content: "a =   1"})
	got := next.(TomlFormatModel)

	if cmd == nil {
		t.Fatal("expected a status message command")
	}
	if got.FormattedContent != "a = 1\n" {
		t.Fatalf("expected formatted content %q, got %q", "a = 1\n", got.FormattedContent)
	}
	if got.State != ui.PagerStateStatusMessage || got.StatusMessage != "Pasted. Press 'c' to copy result." {
		t.Fatalf("expected status message, got state %v and %q", got.State, got.StatusMessage)
	}
}

func TestEditorFinishedInvalidTOMLShowsError(t *testing.T) {
	t.Parallel()

	m := NewTomlFormatModel(&ui.CommonModel{Width: 80, Height: 24})
	next, _ := m.Update(editor.EditorFinishedMsg{Content: "a = "})
	got := next.(TomlFormatModel)

	if got.State != ui.PagerStateErrorMessage || !strings.Contains(got.StatusMessage, "TOML parsing error") {
		t.Fatalf("expected parsing error message, got state %v and %q", got.State, got.StatusMessage)
	}
}

func TestEditorFinishedErrorShowsMessage(t *testing.T) {
	t.Parallel()

	m := NewTomlFormatModel(&ui.CommonModel{Width: 80, Height: 24})
	next, _ := m.Update(editor.EditorFinishedMsg{Err: errors.New("editor failed"), Content: "a = 1"})
	got := next.(TomlFormatModel)

	if got.State != ui.PagerStateErrorMessage || got.StatusMessage != "editor failed" {
		t.Fatalf("expected error message, got state %v and %q", got.State, got.StatusMessage)
	}
	if got.Content != "" {
		t.Fatalf("expected content to stay empty, got %q", got.Content)
	}
}

func TestStatusMessageTimeoutReturnsToBrowse(t *testing.T) {
	t.Parallel()

	m := NewTomlFormatModel(&ui.CommonModel{Width: 80, Height: 24})
	m.State = ui.PagerStateStatusMessage
	next, _ := m.Update(ui.StatusMessageTimeoutMsg{})

	if got := next.(TomlFormatModel).State; got != ui.PagerStateBrowse {
		t.Fatalf("expected browse state, got %v", got)
	}
}

func TestWindowSizeMsgPreparesViewport(t *testing.T) {
	t.Parallel()

	m := NewTomlFormatModel(&ui.CommonModel{})
	next, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	got := next.(TomlFormatModel)

	if !got.Ready || got.Viewport.Width() != 100 {
		t.Fatalf("expected ready viewport of width 100, got ready %v width %d", got.Ready, got.Viewport.Width())
	}
}

func TestHelpKeyTogglesHelpView(t *testing.T) {
	t.Parallel()

	m := NewTomlFormatModel(&ui.CommonModel{Width: 80, Height: 24})
	if strings.Contains(m.View().Content, "copy formatted TOML") {
		t.Fatal("expected help to be hidden")
	}

	next, cmd := m.Update(tea.KeyPressMsg(tea.Key{Code: '?', Text: "?"}))
	if cmd != nil {
		t.Fatal("expected no command for help key")
	}
	if !strings.Contains(next.View().Content, "copy formatted TOML") {
		t.Fatal("expected help to be shown")
	}
}
