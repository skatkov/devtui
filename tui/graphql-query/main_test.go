package graphqlquery

import (
	"errors"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/skatkov/devtui/internal/editor"
	"github.com/skatkov/devtui/internal/ui"
)

const (
	testQuery          = "query{a}"
	testFormattedQuery = "query {\n\ta\n}\n"
	testInvalidQuery   = "query {"
)

func TestHelpViewListsKeys(t *testing.T) {
	t.Parallel()

	m := NewGraphQLQueryModel(&ui.CommonModel{Width: 80, Height: 24})
	help := m.helpView()
	for _, want := range []string{"copy formatted GraphQL", "edit unformatted GraphQL", "paste unformatted GraphQL", "quit"} {
		if !strings.Contains(help, want) {
			t.Errorf("expected help to contain %q, got %q", want, help)
		}
	}
}

func TestSetContentFormatsQuery(t *testing.T) {
	t.Parallel()

	m := newSizedModel()
	if err := m.SetContent(testQuery); err != nil {
		t.Fatalf("SetContent() error = %v", err)
	}
	if m.Content != testQuery {
		t.Errorf("expected content %q, got %q", testQuery, m.Content)
	}
	if m.FormattedContent != testFormattedQuery {
		t.Errorf("expected formatted content %q, got %q", testFormattedQuery, m.FormattedContent)
	}
	if got := m.Viewport.GetContent(); !strings.Contains(got, "query") {
		t.Errorf("expected highlighted query in viewport, got %q", got)
	}
}

func TestLoadContentFormatsContent(t *testing.T) {
	t.Parallel()

	m := newSizedModel()
	if cmd := m.loadContent(testQuery, nil); cmd == nil {
		t.Fatal("expected status message command")
	}
	assertStatus(t, m, ui.PagerStateStatusMessage, "Pasted. Press 'c' to copy result.")
	if m.FormattedContent != testFormattedQuery {
		t.Errorf("expected formatted content %q, got %q", testFormattedQuery, m.FormattedContent)
	}
}

func TestLoadContentShowsClipboardError(t *testing.T) {
	t.Parallel()

	m := newSizedModel()
	m.loadContent("", errors.New("no clipboard"))
	assertStatus(t, m, ui.PagerStateErrorMessage, "no clipboard")
	if m.Content != "" {
		t.Errorf("expected content to stay empty, got %q", m.Content)
	}
}

func TestLoadContentShowsParseError(t *testing.T) {
	t.Parallel()

	m := newSizedModel()
	m.loadContent(testInvalidQuery, nil)
	assertStatus(t, m, ui.PagerStateErrorMessage, "input:1:8: Expected Name, found <EOF>")
}

func TestEditorFinishedMsgFormatsContent(t *testing.T) {
	t.Parallel()

	m := newSizedModel()
	next, cmd := m.Update(editor.EditorFinishedMsg{Content: testQuery})
	if cmd == nil {
		t.Fatal("expected status message command")
	}
	got := next.(GraphQLQueryModel)
	assertStatus(t, &got, ui.PagerStateStatusMessage, "Pasted. Press 'c' to copy result.")
	if got.FormattedContent != testFormattedQuery {
		t.Errorf("expected formatted content %q, got %q", testFormattedQuery, got.FormattedContent)
	}
}

func TestEditorFinishedMsgShowsEditorError(t *testing.T) {
	t.Parallel()

	m := newSizedModel()
	next, _ := m.Update(editor.EditorFinishedMsg{Err: errors.New("editor failed"), Content: testQuery})
	got := next.(GraphQLQueryModel)
	assertStatus(t, &got, ui.PagerStateErrorMessage, "editor failed")
	if got.Content != "" {
		t.Errorf("expected content to stay empty, got %q", got.Content)
	}
}

func TestSetContentInvalidReturnsError(t *testing.T) {
	t.Parallel()

	m := newSizedModel()
	if err := m.SetContent(testInvalidQuery); err == nil {
		t.Fatal("expected error for invalid query")
	}
	if m.FormattedContent != testInvalidQuery {
		t.Errorf("expected unformatted content %q, got %q", testInvalidQuery, m.FormattedContent)
	}
	if got := m.Viewport.GetContent(); got != "" {
		t.Errorf("expected empty viewport, got %q", got)
	}
}

func TestFormatGraphQL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"valid query", testQuery, testFormattedQuery},
		{"invalid query", testInvalidQuery, testInvalidQuery},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := formatGraphQL(tt.input); got != tt.want {
				t.Errorf("formatGraphQL(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestEditorFinishedMsgShowsParseError(t *testing.T) {
	t.Parallel()

	m := newSizedModel()
	next, _ := m.Update(editor.EditorFinishedMsg{Content: testInvalidQuery})
	got := next.(GraphQLQueryModel)
	assertStatus(t, &got, ui.PagerStateErrorMessage, "input:1:8: Expected Name, found <EOF>")
}

func TestStatusMessageTimeoutReturnsToBrowse(t *testing.T) {
	t.Parallel()

	m := newSizedModel()
	m.State = ui.PagerStateStatusMessage
	next, _ := m.Update(ui.StatusMessageTimeoutMsg{})
	if got := next.(GraphQLQueryModel).State; got != ui.PagerStateBrowse {
		t.Errorf("expected browse state, got %v", got)
	}
}

func TestKeyPressScrollsViewport(t *testing.T) {
	t.Parallel()

	m := newSizedModel()
	_ = m.SetContent("query {" + strings.Repeat(" a", 40) + " }")
	next, _ := m.Update(codeKeypress(tea.KeyDown))
	got := next.(GraphQLQueryModel)
	if offset := got.Viewport.YOffset(); offset != 1 {
		t.Errorf("expected viewport offset 1, got %d", offset)
	}
}

func TestViewShowsHelp(t *testing.T) {
	t.Parallel()

	m := newSizedModel()
	if strings.Contains(m.View().Content, "paste unformatted GraphQL") {
		t.Fatal("expected help to be hidden")
	}
	next, _ := m.Update(tea.KeyPressMsg(tea.Key{Code: '?', Text: "?"}))
	if !strings.Contains(next.(GraphQLQueryModel).View().Content, "paste unformatted GraphQL") {
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

func newSizedModel() *GraphQLQueryModel {
	m := NewGraphQLQueryModel(&ui.CommonModel{Width: 80, Height: 24})
	next, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	model := next.(GraphQLQueryModel)
	return &model
}

func assertStatus(t *testing.T, m *GraphQLQueryModel, state ui.PagerState, message string) {
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
