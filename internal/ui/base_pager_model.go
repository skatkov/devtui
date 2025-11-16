// devtui/internal/ui/base_pager_model.go
package ui

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
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

func NewBasePagerModel(common *CommonModel, title string) BasePagerModel {
	model := BasePagerModel{
		Common: common,
		Title:  title,
		Ready:  false,
		State:  PagerStateBrowse,
	}

	model.SetSize(common.Width, common.Height)
	return model
}

func (m BasePagerModel) Init() tea.Cmd {
	return nil
}

func (m *BasePagerModel) HandleCommonKeys(msg tea.KeyMsg) (tea.Cmd, bool) {
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
		err := clipboard.Copy(m.FormattedContent)
		if err != nil {
			return m.ShowErrorMessage(err.Error()), true
		}
		return m.ShowStatusMessage("Copied."), true
	}
	return nil, false
}

func (m *BasePagerModel) HandleWindowSizeMsg(msg tea.WindowSizeMsg) tea.Cmd {
	m.Common.Width = msg.Width
	m.Common.Height = msg.Height

	m.SetSize(msg.Width, msg.Height)

	if !m.Ready {
		m.Viewport = viewport.New(msg.Width, msg.Height-StatusBarHeight)
		m.Viewport.YPosition = 0
		m.Viewport.SetContent(m.Content)
		m.Ready = true
	} else {
		m.SetSize(msg.Width, msg.Height)
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
		// If help is shown, reduce viewport height
		// Make sure HelpHeight is at least a minimum value (like 6 lines)
		helpHeight := m.HelpHeight
		if helpHeight == 0 {
			helpHeight = 6 // Default if not calculated yet
		}
		viewportHeight -= helpHeight
	}

	m.Viewport.Width = m.Common.Width
	m.Viewport.Height = viewportHeight
}

func (m *BasePagerModel) ToggleHelp() {
	m.ShowHelp = !m.ShowHelp
	m.SetSize(m.Common.Width, m.Common.Height)

	if m.Viewport.PastBottom() {
		m.Viewport.GotoBottom()
	}
}

func (m *BasePagerModel) StatusBarView() string {
	var b strings.Builder

	const (
		minPercent               float64 = 0.0
		maxPercent               float64 = 1.0
		percentToStringMagnitude float64 = 100.0
	)
	showStatusMessage := m.State == PagerStateStatusMessage
	showErrorMessage := m.State == PagerStateErrorMessage
	appName := AppNameStyle(" " + m.Title + " ")

	// Scroll percent
	scrollPercent := ""
	if m.Content != "" {
		percent := math.Max(minPercent, math.Min(maxPercent, m.Viewport.ScrollPercent()))
		scrollPercent = fmt.Sprintf(" %3.f%% ", percent*percentToStringMagnitude)
		scrollPercent = StatusBarScrollPosStyle(scrollPercent)
	}

	var helpNote string
	if showErrorMessage {
		helpNote = StatusBarErrorHelpStyle(" ? Help ")
	} else if showStatusMessage {
		helpNote = StatusBarMessageHelpStyle(" ? Help ")
	} else {
		helpNote = StatusBarHelpStyle(" ? Help ")
	}

	var note string
	if showStatusMessage || showErrorMessage {
		note = m.StatusMessage
	} else if m.Content == "" {
		note = "Press 'v' to paste"
	}

	note = truncate.StringWithTail(" "+note+" ", uint(max(0,
		m.Common.Width-
			ansi.PrintableRuneWidth(appName)-
			ansi.PrintableRuneWidth(scrollPercent)-
			ansi.PrintableRuneWidth(helpNote),
	)), Ellipsis)

	if showErrorMessage {
		note = StatusBarErrorStyle(note)
	} else if showStatusMessage {
		note = StatusBarMessageStyle(note)
	} else {
		note = StatusBarNoteStyle(note)
	}

	// Empty space
	padding := max(0,
		m.Common.Width-
			ansi.PrintableRuneWidth(appName)-
			ansi.PrintableRuneWidth(note)-
			ansi.PrintableRuneWidth(scrollPercent)-
			ansi.PrintableRuneWidth(helpNote),
	)
	emptySpace := strings.Repeat(" ", padding)
	if showErrorMessage {
		emptySpace = StatusBarErrorStyle(emptySpace)
	} else if showStatusMessage {
		emptySpace = StatusBarMessageStyle(emptySpace)
	} else {
		emptySpace = StatusBarNoteStyle(emptySpace)
	}

	fmt.Fprintf(&b, "%s%s%s%s%s",
		appName,
		note,
		emptySpace,
		scrollPercent,
		helpNote,
	)

	return b.String()
}

// FormatHelpColumns formats the help view with columns
func (m *BasePagerModel) FormatHelpColumns(col1 []string) string {
	s := "\n"
	s += "k/↑      up                  " + col1[0] + "\n"
	s += "j/↓      down                " + col1[1] + "\n"
	s += "b/pgup   page up             " + col1[2] + "\n"
	s += "f/pgdn   page down           " + col1[3] + "\n"
	s += "u        ½ page up           " + col1[4] + "\n"
	s += "d        ½ page down         "

	if len(col1) > 5 {
		s += col1[5]
	}

	s = Indent(s, 2)

	// Fill up empty cells with spaces for background coloring
	if m.Common.Width > 0 {
		lines := strings.Split(s, "\n")
		for i := range lines {
			l := runewidth.StringWidth(lines[i])
			n := max(m.Common.Width-l, 0)
			lines[i] += strings.Repeat(" ", n)
		}

		s = strings.Join(lines, "\n")
	}

	return HelpViewStyle(s)
}
