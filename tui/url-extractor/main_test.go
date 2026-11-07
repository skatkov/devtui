package urlextractor

import (
	"errors"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/skatkov/devtui/internal/editor"
	"github.com/skatkov/devtui/internal/ui"
)

const sampleText = "see https://a.com and example.org, then https://a.com and http://b.org"

func newTestModel() *URLExtractorModel {
	model := NewURLExtractorModel(&ui.CommonModel{Width: 80, Height: 24})
	model.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	return model
}

func keypress(text string) tea.KeyPressMsg {
	return tea.KeyPressMsg(tea.Key{Code: rune(text[0]), Text: text})
}

func TestNewURLExtractorModel(t *testing.T) {
	t.Parallel()

	model := NewURLExtractorModel(&ui.CommonModel{Width: 80, Height: 24})
	if model.Title != Title {
		t.Fatalf("expected title %q, got %q", Title, model.Title)
	}
	if model.StrictMode {
		t.Fatal("expected relaxed mode by default")
	}
	if model.Init() != nil {
		t.Fatal("expected Init to return nil")
	}
}

func TestSetContent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		strict bool
		input  string
		want   string
	}{
		{"relaxed mode", false, sampleText, "https://a.com\nexample.org\nhttp://b.org"},
		{"strict mode", true, sampleText, "https://a.com\nhttp://b.org"},
		{"no urls", false, "nothing here", ""},
		{"empty input", false, "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			model := newTestModel()
			model.StrictMode = tt.strict
			if err := model.SetContent(tt.input); err != nil {
				t.Fatalf("SetContent() error = %v", err)
			}
			if model.Content != tt.input {
				t.Fatalf("expected content %q, got %q", tt.input, model.Content)
			}
			if model.FormattedContent != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, model.FormattedContent)
			}
			if got := model.Viewport.GetContent(); got != tt.want {
				t.Fatalf("expected viewport content %q, got %q", tt.want, got)
			}
		})
	}
}

func TestSetContentEmptyClearsPreviousResult(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	if err := model.SetContent(sampleText); err != nil {
		t.Fatalf("SetContent() error = %v", err)
	}
	if err := model.SetContent(""); err != nil {
		t.Fatalf("SetContent() error = %v", err)
	}
	if model.FormattedContent != "" || model.Viewport.GetContent() != "" {
		t.Fatalf("expected empty result, got %q", model.FormattedContent)
	}
}

func TestApplyContent(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	if model.applyContent(sampleText, nil) == nil {
		t.Fatal("expected a status message command")
	}
	if model.State != ui.PagerStateStatusMessage || model.StatusMessage != "Pasted and extracted URLs. Press 'c' to copy result." {
		t.Fatalf("expected pasted status, got state %v message %q", model.State, model.StatusMessage)
	}
	if model.FormattedContent != "https://a.com\nexample.org\nhttp://b.org" {
		t.Fatalf("expected extracted URLs, got %q", model.FormattedContent)
	}
}

func TestApplyContentWithError(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	if model.applyContent(sampleText, errors.New("no clipboard")) == nil {
		t.Fatal("expected an error message command")
	}
	if model.State != ui.PagerStateErrorMessage || model.StatusMessage != "no clipboard" {
		t.Fatalf("expected clipboard error, got state %v message %q", model.State, model.StatusMessage)
	}
	if model.Content != "" {
		t.Fatalf("expected content to stay empty, got %q", model.Content)
	}
}

func TestUpdateTogglesStrictMode(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	if err := model.SetContent(sampleText); err != nil {
		t.Fatalf("SetContent() error = %v", err)
	}

	tests := []struct {
		strict      bool
		wantURLs    string
		wantMessage string
	}{
		{true, "https://a.com\nhttp://b.org", "Switched to strict mode. Press 'c' to copy result."},
		{false, "https://a.com\nexample.org\nhttp://b.org", "Switched to relaxed mode. Press 'c' to copy result."},
	}

	for _, tt := range tests {
		if _, cmd := model.Update(keypress("s")); cmd == nil {
			t.Fatal("expected a status message command")
		}
		if model.StrictMode != tt.strict {
			t.Fatalf("expected strict mode %v, got %v", tt.strict, model.StrictMode)
		}
		if model.FormattedContent != tt.wantURLs {
			t.Fatalf("expected %q, got %q", tt.wantURLs, model.FormattedContent)
		}
		if model.State != ui.PagerStateStatusMessage || model.StatusMessage != tt.wantMessage {
			t.Fatalf("expected %q, got state %v message %q", tt.wantMessage, model.State, model.StatusMessage)
		}
	}
}

func TestUpdateTogglesStrictModeWithoutContent(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	for _, want := range []string{"Switched to strict mode.", "Switched to relaxed mode."} {
		model.Update(keypress("s"))
		if model.State != ui.PagerStateStatusMessage || model.StatusMessage != want {
			t.Fatalf("expected %q, got state %v message %q", want, model.State, model.StatusMessage)
		}
	}
	if model.Content != "" {
		t.Fatalf("expected content to stay empty, got %q", model.Content)
	}
}

func TestUpdateEditorFinished(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	if _, cmd := model.Update(editor.EditorFinishedMsg{Content: sampleText}); cmd == nil {
		t.Fatal("expected a status message command")
	}
	if model.State != ui.PagerStateStatusMessage || model.StatusMessage != "URLs extracted. Press 'c' to copy result." {
		t.Fatalf("expected extracted status, got state %v message %q", model.State, model.StatusMessage)
	}
	if model.FormattedContent != "https://a.com\nexample.org\nhttp://b.org" {
		t.Fatalf("expected extracted URLs, got %q", model.FormattedContent)
	}
}

func TestUpdateEditorFinishedWithError(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	model.Update(editor.EditorFinishedMsg{Err: errors.New("editor failed"), Content: sampleText})
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
	model.Update(ui.StatusMessageTimeoutMsg{})
	if model.State != ui.PagerStateBrowse {
		t.Fatalf("expected browse state, got %v", model.State)
	}
}

func TestUpdateScrollsViewport(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	var input strings.Builder
	for i := range 30 {
		input.WriteString("https://a.com/" + strings.Repeat("x", i) + "\n")
	}
	if err := model.SetContent(input.String()); err != nil {
		t.Fatalf("SetContent() error = %v", err)
	}
	model.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyDown}))
	if offset := model.Viewport.YOffset(); offset != 1 {
		t.Fatalf("expected viewport offset 1, got %d", offset)
	}
}

func TestViewTogglesHelp(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	if strings.Contains(model.View().Content, "copy URLs") {
		t.Fatal("expected help to be hidden")
	}
	if _, cmd := model.Update(keypress("?")); cmd != nil {
		t.Fatal("expected no command when toggling help")
	}
	view := model.View()
	if !view.AltScreen {
		t.Fatal("expected alt screen view")
	}
	for _, want := range []string{Title, "copy URLs", "toggle mode (current: relaxed)"} {
		if !strings.Contains(view.Content, want) {
			t.Fatalf("expected view to contain %q, got %q", want, view.Content)
		}
	}
}

func TestHelpViewShowsStrictMode(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	model.StrictMode = true
	if !strings.Contains(model.helpView(), "toggle mode (current: strict)") {
		t.Fatalf("expected strict mode in help, got %q", model.helpView())
	}
}
