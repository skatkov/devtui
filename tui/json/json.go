package json

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/alecthomas/chroma/quick"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mattn/go-runewidth"
	"github.com/muesli/ansi"
	"github.com/muesli/reflow/truncate"
	"github.com/skatkov/devtui/internal/editor"
	"github.com/skatkov/devtui/internal/ui"
	"github.com/tiagomelo/go-clipboard/clipboard"
)

const Title = "JSON Formatter"

var pagerHelpHeight int

type JsonModel struct {
	common *ui.CommonModel

	formatted_content string
	content           string
	viewport          viewport.Model
	showHelp          bool
	ready             bool
	state             ui.PagerState

	statusMessage      string
	statusMessageTimer *time.Timer
}

func NewJsonModel(common *ui.CommonModel) JsonModel {
	model := JsonModel{
		ready:  false,
		common: common,
		state:  ui.PagerStateBrowse,
	}

	model.setSize(common.Width, common.Height)

	return model
}

func (m JsonModel) Init() tea.Cmd {
	return nil
}

func (m JsonModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "e":
			return m, editor.OpenEditor(m.content, "json")
		case "ctrl+c", "q":
			return m, tea.Quit
		case "esc":
			return m, func() tea.Msg {
				return ui.ReturnToListMsg{
					Common: m.common,
				}
			}
		case "v":
			c := clipboard.New()
			content, err := c.PasteText()
			if err != nil {
				cmds = append(cmds, m.showErrorMessage(ui.PagerStatusMsg{Message: err.Error()}))
			} else {
				err := m.SetContent(content)

				if err != nil {
					cmds = append(cmds, m.showErrorMessage(ui.PagerStatusMsg{Message: err.Error()}))
				} else {
					cmds = append(cmds, m.showStatusMessage(ui.PagerStatusMsg{Message: "Pasted. Press 'c' to copy"}))
				}
			}
		case "c":
			c := clipboard.New()
			if err := c.CopyText(m.formatted_content); err != nil {
				cmds = append(cmds, m.showErrorMessage(ui.PagerStatusMsg{Message: err.Error()}))
			} else {
				cmds = append(cmds, m.showStatusMessage(ui.PagerStatusMsg{Message: "Copied."}))
			}
		case "?":
			m.toggleHelp()
		}
	case ui.StatusMessageTimeoutMsg:
		m.state = ui.PagerStateBrowse
	case editor.EditorFinishedMsg:
		if msg.Err != nil {
			cmds = append(cmds, m.showErrorMessage(ui.PagerStatusMsg{Message: msg.Err.Error()}))
		} else {
			err := m.SetContent(msg.Content)

			if err != nil {
				cmds = append(cmds, m.showErrorMessage(ui.PagerStatusMsg{Message: err.Error()}))
			} else {
				cmds = append(cmds, m.showStatusMessage(ui.PagerStatusMsg{Message: "Pasted. Press 'c' to copy"}))
			}
		}
	case tea.WindowSizeMsg:
		m.common.Width = msg.Width
		m.common.Height = msg.Height

		m.setSize(msg.Width, msg.Height)

		if !m.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.viewport = viewport.New(msg.Width, msg.Height-ui.StatusBarHeight)
			m.viewport.YPosition = 0
			m.viewport.SetContent(m.content)
			m.ready = true
		} else {
			m.setSize(msg.Width, msg.Height)
		}
	}
	// Handle keyboard and mouse events in the viewport
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m JsonModel) View() string {
	var b strings.Builder

	fmt.Fprint(&b, m.viewport.View()+"\n")

	// Footer
	m.statusBarView(&b)

	if m.showHelp {
		fmt.Fprint(&b, "\n"+m.helpView())
	}

	return b.String()
}

func (m *JsonModel) showErrorMessage(msg ui.PagerStatusMsg) tea.Cmd {
	m.state = ui.PagerStateErrorMessage
	m.statusMessage = msg.Message
	if m.statusMessageTimer != nil {
		m.statusMessageTimer.Stop()
	}
	m.statusMessageTimer = time.NewTimer(ui.StatusMessageTimeout)

	return ui.WaitForStatusMessageTimeout(m.statusMessageTimer)
}

func (m *JsonModel) showStatusMessage(msg ui.PagerStatusMsg) tea.Cmd {
	// Show a success message to the user
	m.state = ui.PagerStateStatusMessage
	m.statusMessage = msg.Message
	if m.statusMessageTimer != nil {
		m.statusMessageTimer.Stop()
	}
	m.statusMessageTimer = time.NewTimer(ui.StatusMessageTimeout)

	return ui.WaitForStatusMessageTimeout(m.statusMessageTimer)
}

func (m *JsonModel) SetContent(content string) error {
	m.content = content
	m.formatted_content = FormatJSON(content)

	var buf bytes.Buffer
	err := quick.Highlight(&buf, m.formatted_content, "json", "terminal", "nord")
	if err != nil {
		return err
	}
	m.viewport.SetContent(buf.String())

	return nil
}

func (m *JsonModel) setSize(w, h int) {
	m.viewport.Width = w
	m.viewport.Height = h - ui.StatusBarHeight

	if m.showHelp {
		if pagerHelpHeight == 0 {
			pagerHelpHeight = strings.Count(m.helpView(), "\n")
		}
		m.viewport.Height -= (ui.StatusBarHeight + pagerHelpHeight)
	}
}

func (m *JsonModel) toggleHelp() {
	m.showHelp = !m.showHelp
	m.setSize(m.common.Width, m.common.Height)

	if m.viewport.PastBottom() {
		m.viewport.GotoBottom()
	}
}

func (m JsonModel) statusBarView(b *strings.Builder) {
	const (
		minPercent               float64 = 0.0
		maxPercent               float64 = 1.0
		percentToStringMagnitude float64 = 100.0
	)
	showStatusMessage := m.state == ui.PagerStateStatusMessage
	showErrorMessage := m.state == ui.PagerStateErrorMessage
	appName := ui.AppNameStyle(" " + Title + " ")

	// Scroll percent
	scrollPercent := ""
	if m.content != "" {
		percent := math.Max(minPercent, math.Min(maxPercent, m.viewport.ScrollPercent()))
		scrollPercent = fmt.Sprintf(" %3.f%% ", percent*percentToStringMagnitude)
		scrollPercent = ui.StatusBarScrollPosStyle(scrollPercent)
	}
	var helpNote string

	if showErrorMessage {
		helpNote = ui.StatusBarErrorHelpStyle(" ? Help ")
	} else if showStatusMessage {
		helpNote = ui.StatusBarMessageHelpStyle(" ? Help ")
	} else {
		helpNote = ui.StatusBarHelpStyle(" ? Help ")
	}

	var note string
	if showStatusMessage || showErrorMessage {
		note = m.statusMessage
	} else if m.content == "" {
		note = "Press 'v' to paste"
	}

	note = truncate.StringWithTail(" "+note+" ", uint(max(0,
		m.common.Width-
			ansi.PrintableRuneWidth(appName)-
			ansi.PrintableRuneWidth(scrollPercent)-
			ansi.PrintableRuneWidth(helpNote),
	)), ui.Ellipsis)

	if showStatusMessage {
		note = ui.StatusBarMessageStyle(note)
	} else {
		note = ui.StatusBarNoteStyle(note)
	}

	// Empty space
	padding := max(0,
		m.common.Width-
			ansi.PrintableRuneWidth(appName)-
			ansi.PrintableRuneWidth(note)-
			ansi.PrintableRuneWidth(scrollPercent)-
			ansi.PrintableRuneWidth(helpNote),
	)
	emptySpace := strings.Repeat(" ", padding)
	if showErrorMessage {
		emptySpace = ui.StatusBarErrorStyle(emptySpace)
	} else if showStatusMessage {
		emptySpace = ui.StatusBarMessageStyle(emptySpace)
	} else {
		emptySpace = ui.StatusBarNoteStyle(emptySpace)
	}

	fmt.Fprintf(b, "%s%s%s%s%s",
		appName,
		note,
		emptySpace,
		scrollPercent,
		helpNote,
	)
}

func (m JsonModel) helpView() (s string) {
	col1 := []string{
		"c              copy formatted JSON",
		"e              edit unformatted JSON",
		"v              paste unformatted JSON",
		"q/ctrl+c       quit",
	}

	s += "\n"
	s += "k/↑      up                  " + col1[0] + "\n"
	s += "j/↓      down                " + col1[1] + "\n"
	s += "b/pgup   page up             " + col1[2] + "\n"
	s += "f/pgdn   page down           " + col1[3] + "\n"
	s += "u        ½ page up           " + "\n"
	s += "d        ½ page down         "

	if len(col1) > 5 {
		s += col1[5]
	}

	s = ui.Indent(s, 2)

	// Fill up empty cells with spaces for background coloring
	if m.common.Width > 0 {
		lines := strings.Split(s, "\n")
		for i := range lines {
			l := runewidth.StringWidth(lines[i])
			n := max(m.common.Width-l, 0)
			lines[i] += strings.Repeat(" ", n)
		}

		s = strings.Join(lines, "\n")
	}

	return ui.HelpViewStyle(s)
}

func FormatJSON(content string) string {
	var data any
	if err := json.Unmarshal([]byte(content), &data); err != nil {
		return content
	}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	if err := enc.Encode(data); err != nil {
		return content
	}
	return buf.String()
}
