package tsv2md

import (
	"errors"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/skatkov/devtui/internal/editor"
	"github.com/skatkov/devtui/internal/ui"
)

const (
	aligned   = "| id | name | \n| --- | ---- | \n| 1  | bob  | "
	unaligned = "| id | name | \n| --- | --- | \n| 1 | bob | "
)

func TestNewTSV2MDModel(t *testing.T) {
	t.Parallel()

	m := NewTSV2MDModel(&ui.CommonModel{Width: 80, Height: 24})
	if m.Title != Title {
		t.Fatalf("expected title %q, got %q", Title, m.Title)
	}
	if !m.alignColumns {
		t.Fatal("expected columns to be aligned by default")
	}
	if m.Init() != nil {
		t.Fatal("expected nil init command")
	}
}

func TestSetContent(t *testing.T) {
	t.Parallel()

	m := NewTSV2MDModel(&ui.CommonModel{Width: 80, Height: 24})
	if err := m.SetContent("id\tname\n1\tbob"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.FormattedContent != aligned {
		t.Fatalf("expected formatted content %q, got %q", aligned, m.FormattedContent)
	}
	if m.Viewport.GetContent() != aligned {
		t.Fatalf("expected table in viewport, got %q", m.Viewport.GetContent())
	}
}

func TestSetContentInvalidTSVReturnsError(t *testing.T) {
	t.Parallel()

	m := NewTSV2MDModel(&ui.CommonModel{Width: 80, Height: 24})
	if err := m.SetContent("a\"b\tc"); err == nil {
		t.Fatal("expected error for invalid TSV")
	}
	if m.Content != "a\"b\tc" || m.FormattedContent != "" {
		t.Fatalf("expected raw content only, got %q and %q", m.Content, m.FormattedContent)
	}
}

func TestEditorFinishedSetsContent(t *testing.T) {
	t.Parallel()

	m := NewTSV2MDModel(&ui.CommonModel{Width: 80, Height: 24})
	next, cmd := m.Update(editor.EditorFinishedMsg{Content: "id\tname\n1\tbob"})
	got := next.(TSV2MDModel)

	if cmd == nil {
		t.Fatal("expected a status message command")
	}
	if got.FormattedContent != aligned {
		t.Fatalf("expected formatted content %q, got %q", aligned, got.FormattedContent)
	}
	if got.State != ui.PagerStateStatusMessage || got.StatusMessage != "Pasted. Press 'c' to copy result." {
		t.Fatalf("expected status message, got state %v and %q", got.State, got.StatusMessage)
	}
}

func TestEditorFinishedInvalidTSVShowsError(t *testing.T) {
	t.Parallel()

	m := NewTSV2MDModel(&ui.CommonModel{Width: 80, Height: 24})
	next, _ := m.Update(editor.EditorFinishedMsg{Content: "a\"b\tc"})
	got := next.(TSV2MDModel)

	if got.State != ui.PagerStateErrorMessage || !strings.Contains(got.StatusMessage, "error reading TSV") {
		t.Fatalf("expected parsing error message, got state %v and %q", got.State, got.StatusMessage)
	}
}

func TestEditorFinishedErrorShowsMessage(t *testing.T) {
	t.Parallel()

	m := NewTSV2MDModel(&ui.CommonModel{Width: 80, Height: 24})
	next, _ := m.Update(editor.EditorFinishedMsg{Err: errors.New("editor failed"), Content: "id\tname\n1\tbob"})
	got := next.(TSV2MDModel)

	if got.State != ui.PagerStateErrorMessage || got.StatusMessage != "editor failed" {
		t.Fatalf("expected error message, got state %v and %q", got.State, got.StatusMessage)
	}
	if got.Content != "" {
		t.Fatalf("expected content to stay empty, got %q", got.Content)
	}
}

func TestStatusMessageTimeoutReturnsToBrowse(t *testing.T) {
	t.Parallel()

	m := NewTSV2MDModel(&ui.CommonModel{Width: 80, Height: 24})
	m.State = ui.PagerStateStatusMessage
	next, _ := m.Update(ui.StatusMessageTimeoutMsg{})

	if got := next.(TSV2MDModel).State; got != ui.PagerStateBrowse {
		t.Fatalf("expected browse state, got %v", got)
	}
}

func TestWindowSizeMsgPreparesViewport(t *testing.T) {
	t.Parallel()

	m := NewTSV2MDModel(&ui.CommonModel{})
	next, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	got := next.(TSV2MDModel)

	if !got.Ready || got.Viewport.Width() != 100 {
		t.Fatalf("expected ready viewport of width 100, got ready %v width %d", got.Ready, got.Viewport.Width())
	}
}

func TestHelpKeyTogglesHelpView(t *testing.T) {
	t.Parallel()

	m := NewTSV2MDModel(&ui.CommonModel{Width: 80, Height: 24})
	if strings.Contains(m.View().Content, "toggle column alignment") {
		t.Fatal("expected help to be hidden")
	}

	next, cmd := m.Update(tea.KeyPressMsg(tea.Key{Code: '?', Text: "?"}))
	if cmd != nil {
		t.Fatal("expected no command for help key")
	}
	if !strings.Contains(next.View().Content, "toggle column alignment") {
		t.Fatal("expected help to be shown")
	}
}

func TestSetContentEmptyReturnsError(t *testing.T) {
	t.Parallel()

	m := NewTSV2MDModel(&ui.CommonModel{Width: 80, Height: 24})
	if err := m.SetContent(""); err == nil || err.Error() != "empty result" {
		t.Fatalf("expected empty result error, got %v", err)
	}
}

func TestAlignKeyTogglesColumnAlignment(t *testing.T) {
	t.Parallel()

	m := NewTSV2MDModel(&ui.CommonModel{Width: 80, Height: 24})
	if err := m.SetContent("id\tname\n1\tbob"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	next, _ := m.Update(tea.KeyPressMsg(tea.Key{Code: 'a', Text: "a"}))
	got := next.(TSV2MDModel)
	if got.FormattedContent != unaligned || got.StatusMessage != "Columns unaligned" {
		t.Fatalf("expected unaligned table, got %q and %q", got.FormattedContent, got.StatusMessage)
	}

	next, _ = got.Update(tea.KeyPressMsg(tea.Key{Code: 'a', Text: "a"}))
	got = next.(TSV2MDModel)
	if got.FormattedContent != aligned || got.StatusMessage != "Columns aligned" {
		t.Fatalf("expected aligned table, got %q and %q", got.FormattedContent, got.StatusMessage)
	}
}

func TestAlignKeyWithoutContentOnlyToggles(t *testing.T) {
	t.Parallel()

	m := NewTSV2MDModel(&ui.CommonModel{Width: 80, Height: 24})
	next, _ := m.Update(tea.KeyPressMsg(tea.Key{Code: 'a', Text: "a"}))
	got := next.(TSV2MDModel)

	if got.alignColumns || got.State != ui.PagerStateBrowse {
		t.Fatalf("expected unaligned browse state, got align %v state %v", got.alignColumns, got.State)
	}
}

func TestAlignKeyInvalidTSVShowsError(t *testing.T) {
	t.Parallel()

	m := NewTSV2MDModel(&ui.CommonModel{Width: 80, Height: 24})
	_ = m.SetContent("a\"b\tc")
	next, _ := m.Update(tea.KeyPressMsg(tea.Key{Code: 'a', Text: "a"}))
	got := next.(TSV2MDModel)

	if got.State != ui.PagerStateErrorMessage || !strings.Contains(got.StatusMessage, "error reading TSV") {
		t.Fatalf("expected read error message, got state %v and %q", got.State, got.StatusMessage)
	}
}
