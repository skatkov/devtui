// devtui/internal/ui/base_pager_model.go
package ui

import (
	"cmp"
	"fmt"
	"strings"
	"time"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"github.com/mattn/go-runewidth"
	"github.com/muesli/ansi"
	"github.com/muesli/reflow/truncate"
	"github.com/skatkov/devtui/internal/clipboard"
)

// BasePagerModel provides common functionality for pager-based TUI views. There is a lot of common
type BasePagerModel struct {
	Common             *CommonModel
	Title              string
	Content            string
	FormattedContent   string
	Viewport           viewport.Model
	ShowHelp           bool
	Ready              bool
	State              PagerState
	StatusMessage      string
	StatusMessageTimer *time.Timer
	HelpHeight         int
}

var copyToClipboard = clipboard.Copy

func NewBasePagerModel(common *CommonModel, title string) BasePagerModel {
	model := BasePagerModel{
		Common: common,
		Title:  title,
		State:  PagerStateBrowse,
	}

	model.SetSize(common.Width, common.Height)
	return model
}

func (m BasePagerModel) Init() tea.Cmd {
	return nil
}

func (m *BasePagerModel) HandleCommonKeys(msg tea.KeyPressMsg) (tea.Cmd, bool) {
	switch msg.String() {
	case "ctrl+c", "q":
		return tea.Quit, true
	case "esc":
		return func() tea.Msg {
			return ReturnToListMsg{
				Common: m.Common,
			}
		}, true
	case "?":
		m.ToggleHelp()
		return nil, true
	case "c":
		err := copyToClipboard(m.FormattedContent)
		if err != nil {
			return m.ShowErrorMessage(err.Error()), true
		}
		return m.ShowStatusMessage("Copied."), true
	}
	return nil, false
}

func (m BasePagerModel) NewView(content string) tea.View {
	v := tea.NewView(content)
	v.AltScreen = true
	return v
}

func (m *BasePagerModel) HandleWindowSizeMsg(msg tea.WindowSizeMsg) tea.Cmd {
	m.Common.Width = msg.Width
	m.Common.Height = msg.Height

	m.SetSize(msg.Width, msg.Height)

	if !m.Ready {
		m.Viewport = viewport.New(
			viewport.WithWidth(msg.Width),
			viewport.WithHeight(msg.Height-StatusBarHeight),
		)
		m.Viewport.SetContent(m.Content)
		m.Ready = true
	}
	return nil
}

func (m *BasePagerModel) ShowErrorMessage(message string) tea.Cmd {
	m.State = PagerStateErrorMessage
	m.StatusMessage = message
	if m.StatusMessageTimer != nil {
		m.StatusMessageTimer.Stop()
	}
	m.StatusMessageTimer = time.NewTimer(StatusMessageTimeout)

	return WaitForStatusMessageTimeout(m.StatusMessageTimer)
}

func (m *BasePagerModel) ShowStatusMessage(message string) tea.Cmd {
	m.State = PagerStateStatusMessage
	m.StatusMessage = message
	if m.StatusMessageTimer != nil {
		m.StatusMessageTimer.Stop()
	}
	m.StatusMessageTimer = time.NewTimer(StatusMessageTimeout)

	return WaitForStatusMessageTimeout(m.StatusMessageTimer)
}

func (m *BasePagerModel) SetSize(w, h int) {
	m.Common.Width = w
	m.Common.Height = h

	viewportHeight := m.Common.Height - StatusBarHeight

	if m.ShowHelp {
		// If help is shown, reduce viewport height; 6 lines until HelpHeight is calculated
		viewportHeight -= cmp.Or(m.HelpHeight, 6)
	}

	m.Viewport.SetWidth(m.Common.Width)
	m.Viewport.SetHeight(viewportHeight)
}

func (m *BasePagerModel) ToggleHelp() {
	m.ShowHelp = !m.ShowHelp
	m.SetSize(m.Common.Width, m.Common.Height)

	if m.Viewport.PastBottom() {
		m.Viewport.GotoBottom()
	}
}

func (m *BasePagerModel) StatusBarView() string {
	appName := AppNameStyle(" " + m.Title + " ")

	// Scroll percent
	scrollPercent := ""
	if m.Content != "" {
		scrollPercent = fmt.Sprintf(" %3.f%% ", m.Viewport.ScrollPercent()*100)
		scrollPercent = StatusBarScrollPosStyle(scrollPercent)
	}

	noteStyle, helpStyle := StatusBarNoteStyle, StatusBarHelpStyle
	var note string
	switch m.State {
	case PagerStateStatusMessage:
		noteStyle, helpStyle = StatusBarMessageStyle, StatusBarMessageHelpStyle
		note = m.StatusMessage
	case PagerStateErrorMessage:
		noteStyle, helpStyle = StatusBarErrorStyle, StatusBarErrorHelpStyle
		note = m.StatusMessage
	default:
		if m.Content == "" {
			note = "Press 'v' to paste"
		}
	}
	helpNote := helpStyle(" ? Help ")

	free := m.Common.Width -
		ansi.PrintableRuneWidth(appName) -
		ansi.PrintableRuneWidth(scrollPercent) -
		ansi.PrintableRuneWidth(helpNote)
	// Without free columns the note is left out, so the bar stays inside the window.
	if free > 0 {
		note = truncate.StringWithTail(" "+note+" ", uint(free), Ellipsis)
	} else {
		note = ""
	}
	emptySpace := strings.Repeat(" ", max(0, free-ansi.PrintableRuneWidth(note)))

	return appName + noteStyle(note) + noteStyle(emptySpace) + scrollPercent + helpNote
}

// FormatHelpColumns formats the help view with columns
func (m *BasePagerModel) FormatHelpColumns(col1 []string) string {
	s := "\n"
	s += "k/↑      up                  " + helpColumnValue(col1, 0) + "\n"
	s += "j/↓      down                " + helpColumnValue(col1, 1) + "\n"
	s += "b/pgup   page up             " + helpColumnValue(col1, 2) + "\n"
	s += "f/pgdn   page down           " + helpColumnValue(col1, 3) + "\n"
	s += "u        ½ page up           " + helpColumnValue(col1, 4) + "\n"
	s += "d        ½ page down         " + helpColumnValue(col1, 5)

	s = Indent(s, 2)

	// Fill up empty cells with spaces for background coloring
	lines := strings.Split(s, "\n")
	for i := range lines {
		lines[i] += strings.Repeat(" ", max(m.Common.Width-runewidth.StringWidth(lines[i]), 0))
	}
	s = strings.Join(lines, "\n")

	return HelpViewStyle(s)
}

func helpColumnValue(columns []string, index int) string {
	if index < 0 || index >= len(columns) {
		return ""
	}

	return columns[index]
}
