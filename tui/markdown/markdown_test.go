package markdown

import (
	"errors"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/skatkov/devtui/internal/editor"
	"github.com/skatkov/devtui/internal/ui"
)

func TestNewMarkdownModel(t *testing.T) {
	t.Parallel()

	m := NewMarkdownModel(&ui.CommonModel{Width: 80, Height: 24})
	if m.Title != Title {
		t.Fatalf("expected title %q, got %q", Title, m.Title)
	}
	if m.Init() != nil {
		t.Fatal("expected nil init command")
	}
}

func TestSetContentRendersMarkdown(t *testing.T) {
	t.Parallel()

	m := NewMarkdownModel(&ui.CommonModel{Width: 80, Height: 24})
	if err := m.SetContent("# Hello"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Content != "# Hello" {
		t.Fatalf("expected content %q, got %q", "# Hello", m.Content)
	}
	if !strings.Contains(m.FormattedContent, "Hello") || strings.Contains(m.FormattedContent, "# Hello") {
		t.Fatalf("expected rendered heading, got %q", m.FormattedContent)
	}
}

func TestSetContentEmptyClearsOutput(t *testing.T) {
	t.Parallel()

	m := NewMarkdownModel(&ui.CommonModel{Width: 80, Height: 24})
	if err := m.SetContent("# Hello"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := m.SetContent(""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.FormattedContent != "" {
		t.Fatalf("expected empty output, got %q", m.FormattedContent)
	}
}

func TestEditorFinishedSetsContent(t *testing.T) {
	t.Parallel()

	m := NewMarkdownModel(&ui.CommonModel{Width: 80, Height: 24})
	next, cmd := m.Update(editor.EditorFinishedMsg{Content: "# Hello"})
	got := next.(MarkdownModel)

	if cmd == nil {
		t.Fatal("expected a status message command")
	}
	if got.Content != "# Hello" {
		t.Fatalf("expected content %q, got %q", "# Hello", got.Content)
	}
	if got.State != ui.PagerStateStatusMessage || got.StatusMessage != "Markdown content edited" {
		t.Fatalf("expected status message, got state %v and %q", got.State, got.StatusMessage)
	}
}

func TestEditorFinishedErrorShowsMessage(t *testing.T) {
	t.Parallel()

	m := NewMarkdownModel(&ui.CommonModel{Width: 80, Height: 24})
	next, _ := m.Update(editor.EditorFinishedMsg{Err: errors.New("editor failed"), Content: "# Hello"})
	got := next.(MarkdownModel)

	if got.State != ui.PagerStateErrorMessage || got.StatusMessage != "editor failed" {
		t.Fatalf("expected error message, got state %v and %q", got.State, got.StatusMessage)
	}
	if got.Content != "" {
		t.Fatalf("expected content to stay empty, got %q", got.Content)
	}
}

func TestStatusMessageTimeoutReturnsToBrowse(t *testing.T) {
	t.Parallel()

	m := NewMarkdownModel(&ui.CommonModel{Width: 80, Height: 24})
	m.State = ui.PagerStateStatusMessage
	next, _ := m.Update(ui.StatusMessageTimeoutMsg{})

	if got := next.(MarkdownModel).State; got != ui.PagerStateBrowse {
		t.Fatalf("expected browse state, got %v", got)
	}
}

func TestWindowSizeMsgPreparesViewport(t *testing.T) {
	t.Parallel()

	m := NewMarkdownModel(&ui.CommonModel{})
	next, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	got := next.(MarkdownModel)

	if !got.Ready || got.Viewport.Width() != 100 {
		t.Fatalf("expected ready viewport of width 100, got ready %v width %d", got.Ready, got.Viewport.Width())
	}
}

func TestHelpKeyTogglesHelpView(t *testing.T) {
	t.Parallel()

	m := NewMarkdownModel(&ui.CommonModel{Width: 80, Height: 24})
	if strings.Contains(m.View().Content, "copy rendered markdown") {
		t.Fatal("expected help to be hidden")
	}

	next, cmd := m.Update(tea.KeyPressMsg(tea.Key{Code: '?', Text: "?"}))
	if cmd != nil {
		t.Fatal("expected no command for help key")
	}
	if !strings.Contains(next.View().Content, "copy rendered markdown") {
		t.Fatal("expected help to be shown")
	}
}
