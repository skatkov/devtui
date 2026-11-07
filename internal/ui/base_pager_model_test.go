package ui

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

func TestFormatHelpColumnsHandlesShortColumnList(t *testing.T) {
	t.Parallel()

	m := BasePagerModel{Common: &CommonModel{Width: 80, Height: 24}}

	if m.FormatHelpColumns([]string{"copy", "edit"}) == "" {
		t.Fatal("expected non-empty help output")
	}
}

func TestFormatHelpColumns(t *testing.T) {
	t.Parallel()

	rows := []string{
		"  ",
		"  k/↑      up                  a",
		"  j/↓      down                b",
		"  b/pgup   page up             c",
		"  f/pgdn   page down           d",
		"  u        ½ page up           e",
		"  d        ½ page down         f",
		"",
	}
	padded := make([]string, len(rows))
	for i, row := range rows {
		padded[i] = fmt.Sprintf("%-40s", row)
	}

	tests := []struct {
		name  string
		width int
		want  string
	}{
		{"pads lines to the width", 40, strings.Join(padded, "\n")},
		{"pads only shorter lines", 1, strings.Join(rows[:len(rows)-1], "\n") + "\n "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := BasePagerModel{Common: &CommonModel{Width: tt.width}}
			got := m.FormatHelpColumns([]string{"a", "b", "c", "d", "e", "f"})
			if want := HelpViewStyle(tt.want); got != want {
				t.Fatalf("FormatHelpColumns() = %q, want %q", got, want)
			}
		})
	}
}

func TestHelpColumnValueBounds(t *testing.T) {
	t.Parallel()

	columns := []string{"a", "b"}

	if got := helpColumnValue(columns, 0); got != "a" {
		t.Fatalf("expected first value to be 'a', got %q", got)
	}

	if got := helpColumnValue(columns, 1); got != "b" {
		t.Fatalf("expected second value to be 'b', got %q", got)
	}

	if got := helpColumnValue(columns, 5); got != "" {
		t.Fatalf("expected out-of-range value to be empty, got %q", got)
	}

	if got := helpColumnValue(columns, -1); got != "" {
		t.Fatalf("expected negative index value to be empty, got %q", got)
	}
}

func TestNewViewUsesAltScreen(t *testing.T) {
	t.Parallel()

	v := BasePagerModel{}.NewView("body")
	if !v.AltScreen {
		t.Fatal("expected alt screen view")
	}
}

func TestHandleCommonKeys(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		key  tea.KeyPressMsg
		want tea.Msg
	}{
		{"q quits", tea.KeyPressMsg{Code: 'q', Text: "q"}, tea.QuitMsg{}},
		{"ctrl+c quits", tea.KeyPressMsg{Code: 'c', Mod: tea.ModCtrl}, tea.QuitMsg{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := NewBasePagerModel(&CommonModel{Width: 80, Height: 24}, "Test")
			cmd, handled := m.HandleCommonKeys(tt.key)
			if !handled {
				t.Fatalf("expected %q to be handled", tt.key.String())
			}
			if got := cmd(); got != tt.want {
				t.Fatalf("expected %T, got %T", tt.want, got)
			}
		})
	}
}

func TestHandleCommonKeysEscReturnsToList(t *testing.T) {
	t.Parallel()

	common := &CommonModel{Width: 80, Height: 24}
	m := NewBasePagerModel(common, "Test")

	cmd, handled := m.HandleCommonKeys(tea.KeyPressMsg{Code: tea.KeyEscape})
	if !handled {
		t.Fatal("expected esc to be handled")
	}
	msg, ok := cmd().(ReturnToListMsg)
	if !ok || msg.Common != common {
		t.Fatalf("expected ReturnToListMsg with the common model, got %#v", msg)
	}
}

func TestHandleCommonKeysQuestionMarkTogglesHelp(t *testing.T) {
	t.Parallel()

	m := NewBasePagerModel(&CommonModel{Width: 80, Height: 24}, "Test")

	cmd, handled := m.HandleCommonKeys(tea.KeyPressMsg{Code: '?', Text: "?"})
	if !handled || cmd != nil {
		t.Fatalf("expected handled key without command, got handled=%v", handled)
	}
	if !m.ShowHelp {
		t.Fatal("expected help to be shown")
	}
}

func TestHandleCommonKeysIgnoresOtherKeys(t *testing.T) {
	t.Parallel()

	m := NewBasePagerModel(&CommonModel{Width: 80, Height: 24}, "Test")

	cmd, handled := m.HandleCommonKeys(tea.KeyPressMsg{Code: 'x', Text: "x"})
	if handled || cmd != nil {
		t.Fatalf("expected unhandled key, got handled=%v", handled)
	}
}

func TestHandleCommonKeysCopy(t *testing.T) {
	tests := []struct {
		name        string
		copyErr     error
		wantState   PagerState
		wantMessage string
	}{
		{"copies formatted content", nil, PagerStateStatusMessage, "Copied."},
		{"shows copy error", errors.New("no clipboard"), PagerStateErrorMessage, "no clipboard"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			original := copyToClipboard
			t.Cleanup(func() { copyToClipboard = original })
			var copied string
			copyToClipboard = func(text string) error {
				copied = text
				return tt.copyErr
			}

			m := NewBasePagerModel(&CommonModel{Width: 80, Height: 24}, "Test")
			m.FormattedContent = "formatted"

			cmd, handled := m.HandleCommonKeys(tea.KeyPressMsg{Code: 'c', Text: "c"})
			m.StatusMessageTimer.Stop()
			if !handled || cmd == nil {
				t.Fatalf("expected handled key with command, got handled=%v", handled)
			}
			if copied != "formatted" {
				t.Fatalf("expected formatted content to be copied, got %q", copied)
			}
			if m.State != tt.wantState || m.StatusMessage != tt.wantMessage {
				t.Fatalf("got state %v message %q, want %v %q", m.State, m.StatusMessage, tt.wantState, tt.wantMessage)
			}
		})
	}
}

func TestHandleWindowSizeMsg(t *testing.T) {
	t.Parallel()

	common := &CommonModel{}
	m := NewBasePagerModel(common, "Test")
	m.Content = "hello"

	m.HandleWindowSizeMsg(tea.WindowSizeMsg{Width: 30, Height: 10})
	if !m.Ready {
		t.Fatal("expected model to be ready")
	}
	if common.Width != 30 || common.Height != 10 {
		t.Fatalf("expected common size 30x10, got %dx%d", common.Width, common.Height)
	}
	if m.Viewport.Width() != 30 || m.Viewport.Height() != 9 {
		t.Fatalf("expected viewport 30x9, got %dx%d", m.Viewport.Width(), m.Viewport.Height())
	}
	if !strings.Contains(m.Viewport.View(), "hello") {
		t.Fatalf("expected viewport to show content, got %q", m.Viewport.View())
	}

	m.Content = "other"
	m.HandleWindowSizeMsg(tea.WindowSizeMsg{Width: 40, Height: 20})
	if m.Viewport.Width() != 40 || m.Viewport.Height() != 19 {
		t.Fatalf("expected viewport 40x19, got %dx%d", m.Viewport.Width(), m.Viewport.Height())
	}
	if !strings.Contains(m.Viewport.View(), "hello") {
		t.Fatalf("expected resize to keep the viewport, got %q", m.Viewport.View())
	}
}

func TestSetSize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		showHelp   bool
		helpHeight int
		want       int
	}{
		{"without help", false, 4, 23},
		{"with default help height", true, 0, 17},
		{"with measured help height", true, 4, 19},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := BasePagerModel{Common: &CommonModel{}, ShowHelp: tt.showHelp, HelpHeight: tt.helpHeight}
			m.SetSize(80, 24)
			if m.Viewport.Width() != 80 || m.Viewport.Height() != tt.want {
				t.Fatalf("expected viewport 80x%d, got %dx%d", tt.want, m.Viewport.Width(), m.Viewport.Height())
			}
		})
	}
}

func TestToggleHelp(t *testing.T) {
	t.Parallel()

	m := NewBasePagerModel(&CommonModel{Width: 80, Height: 24}, "Test")

	m.ToggleHelp()
	if !m.ShowHelp || m.Viewport.Height() != 17 {
		t.Fatalf("expected help shown with height 17, got %v %d", m.ShowHelp, m.Viewport.Height())
	}

	m.ToggleHelp()
	if m.ShowHelp || m.Viewport.Height() != 23 {
		t.Fatalf("expected help hidden with height 23, got %v %d", m.ShowHelp, m.Viewport.Height())
	}
}

func TestShowMessages(t *testing.T) {
	t.Parallel()

	m := NewBasePagerModel(&CommonModel{Width: 80, Height: 24}, "Test")

	if m.ShowStatusMessage("saved") == nil {
		t.Fatal("expected a timeout command")
	}
	if m.State != PagerStateStatusMessage || m.StatusMessage != "saved" {
		t.Fatalf("got state %v message %q", m.State, m.StatusMessage)
	}
	first := m.StatusMessageTimer

	if m.ShowErrorMessage("failed") == nil {
		t.Fatal("expected a timeout command")
	}
	if m.State != PagerStateErrorMessage || m.StatusMessage != "failed" {
		t.Fatalf("got state %v message %q", m.State, m.StatusMessage)
	}
	second := m.StatusMessageTimer
	if first.Stop() {
		t.Fatal("expected the previous timer to be stopped")
	}

	m.ShowStatusMessage("again")
	if second.Stop() {
		t.Fatal("expected the previous timer to be stopped")
	}
	m.StatusMessageTimer.Stop()
}

func TestStatusBarView(t *testing.T) {
	t.Parallel()

	long := strings.Repeat("x", 40)

	tests := []struct {
		name    string
		state   PagerState
		message string
		content string
		width   int
		want    string
	}{
		{
			"empty content asks to paste", PagerStateBrowse, "", "", 40,
			AppNameStyle(" T ") + StatusBarNoteStyle(" Press 'v' to paste ") + StatusBarNoteStyle(strings.Repeat(" ", 9)) +
				StatusBarHelpStyle(" ? Help "),
		},
		{
			"content shows scroll percent", PagerStateBrowse, "", "hello", 40,
			AppNameStyle(" T ") + StatusBarNoteStyle("  ") + StatusBarNoteStyle(strings.Repeat(" ", 21)) +
				StatusBarScrollPosStyle(" 100% ") + StatusBarHelpStyle(" ? Help "),
		},
		{
			"status message", PagerStateStatusMessage, "Copied.", "hello", 40,
			AppNameStyle(" T ") + StatusBarMessageStyle(" Copied. ") + StatusBarMessageStyle(strings.Repeat(" ", 14)) +
				StatusBarScrollPosStyle(" 100% ") + StatusBarMessageHelpStyle(" ? Help "),
		},
		{
			"long error message is truncated", PagerStateErrorMessage, long, "hello", 40,
			AppNameStyle(" T ") + StatusBarErrorStyle(" "+long[:21]+Ellipsis) + StatusBarErrorStyle("") +
				StatusBarScrollPosStyle(" 100% ") + StatusBarErrorHelpStyle(" ? Help "),
		},
		{
			"one free column leaves only an ellipsis", PagerStateBrowse, "", "", 12,
			AppNameStyle(" T ") + StatusBarNoteStyle(Ellipsis) + StatusBarNoteStyle("") + StatusBarHelpStyle(" ? Help "),
		},
		{
			"no free columns leaves the note out", PagerStateBrowse, "", "", 11,
			AppNameStyle(" T ") + StatusBarNoteStyle("") + StatusBarNoteStyle("") + StatusBarHelpStyle(" ? Help "),
		},
		{
			"narrower than the bar leaves the note out", PagerStateBrowse, "", "", 5,
			AppNameStyle(" T ") + StatusBarNoteStyle("") + StatusBarNoteStyle("") + StatusBarHelpStyle(" ? Help "),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := NewBasePagerModel(&CommonModel{}, "T")
			m.Content = tt.content
			m.HandleWindowSizeMsg(tea.WindowSizeMsg{Width: tt.width, Height: 10})
			m.State = tt.state
			m.StatusMessage = tt.message

			if got := m.StatusBarView(); got != tt.want {
				t.Fatalf("StatusBarView() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestStatusBarViewScrollPercent(t *testing.T) {
	t.Parallel()

	m := NewBasePagerModel(&CommonModel{}, "T")
	m.Content = strings.Repeat("line\n", 30)
	m.HandleWindowSizeMsg(tea.WindowSizeMsg{Width: 40, Height: 10})

	if got := m.StatusBarView(); !strings.Contains(got, "   0% ") {
		t.Fatalf("expected 0%% at the top, got %q", got)
	}

	m.Viewport.GotoBottom()
	if got := m.StatusBarView(); !strings.Contains(got, " 100% ") {
		t.Fatalf("expected 100%% at the bottom, got %q", got)
	}
}
