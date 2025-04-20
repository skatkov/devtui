// devtui/tui/csv2md/main.go
package csv2md

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mattn/go-runewidth"
	"github.com/muesli/ansi"
	"github.com/muesli/reflow/truncate"
	"github.com/skatkov/devtui/internal/csv2md"
	"github.com/skatkov/devtui/internal/editor"
	"github.com/skatkov/devtui/internal/ui"
	"github.com/tiagomelo/go-clipboard/clipboard"
)

const Title = "CSV to Markdown Table Converter"

var pagerHelpHeight int

type CSV2MDModel struct {
	common *ui.CommonModel

	converted_content string
	content           string
	viewport          viewport.Model
	showHelp          bool
	ready             bool
	state             ui.PagerState
	alignColumns      bool

	statusMessage      string
	statusMessageTimer *time.Timer
}

func NewCSV2MDModel(common *ui.CommonModel) CSV2MDModel {
	model := CSV2MDModel{
		content:      "",
		ready:        false,
		common:       common,
		state:        ui.PagerStateBrowse,
		alignColumns: true, // default to aligned columns
	}

	model.setSize(common.Width, common.Height)

	return model
}

func (m CSV2MDModel) Init() tea.Cmd {
	return nil
}

func (m CSV2MDModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "e":
			return m, editor.OpenEditor(m.content, "csv")
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
			if err == nil {
				err = m.SetContent(content)
			}

			if err != nil {
				cmds = append(cmds, m.showErrorMessage(ui.PagerStatusMsg{Message: err.Error()}))
			} else {
				cmds = append(cmds, m.showStatusMessage(ui.PagerStatusMsg{Message: "Converted " + Title + ". Press 'c' to copy result."}))
			}
		case "c":
			c := clipboard.New()
			err := c.CopyText(m.converted_content)

			if err != nil {
				cmds = append(cmds, m.showErrorMessage(ui.PagerStatusMsg{Message: err.Error()}))
			} else {
				cmds = append(cmds, m.showStatusMessage(ui.PagerStatusMsg{Message: "Copied"}))
			}
		case "a":
			m.alignColumns = !m.alignColumns
			if m.content != "" {
				err := m.SetContent(m.content)
				if err != nil {
					cmds = append(cmds, m.showErrorMessage(ui.PagerStatusMsg{Message: err.Error()}))
				} else if m.alignColumns {
					cmds = append(cmds, m.showStatusMessage(ui.PagerStatusMsg{Message: "Columns aligned"}))
				} else {
					cmds = append(cmds, m.showStatusMessage(ui.PagerStatusMsg{Message: "Columns unaligned"}))
				}
			}
		case "?":
			m.toggleHelp()
		}
	case ui.StatusMessageTimeoutMsg:
		m.state = ui.PagerStateBrowse
	case editor.EditorFinishedMsg:
		if msg.Err != nil {
			panic(msg.Err)
		}
		err := m.SetContent(msg.Content)

		if err != nil {
			cmds = append(cmds, m.showErrorMessage(ui.PagerStatusMsg{Message: err.Error()}))
		} else {
			cmds = append(cmds, m.showStatusMessage(ui.PagerStatusMsg{Message: "Converted. Press 'c' to copy result."}))
		}
	case tea.WindowSizeMsg:
		m.common.Width = msg.Width
		m.common.Height = msg.Height

		m.setSize(msg.Width, msg.Height)

		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-ui.StatusBarHeight)
			m.viewport.YPosition = 0
			m.viewport.SetContent(m.content)
			m.ready = true
		} else {
			m.setSize(msg.Width, msg.Height)
		}
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m CSV2MDModel) View() string {
	var b strings.Builder

	fmt.Fprint(&b, m.viewport.View()+"\n")

	m.statusBarView(&b)

	if m.showHelp {
		fmt.Fprint(&b, "\n"+m.helpView())
	}

	return b.String()
}

func (m *CSV2MDModel) showErrorMessage(msg ui.PagerStatusMsg) tea.Cmd {
	m.state = ui.PagerStateErrorMessage
	m.statusMessage = msg.Message
	if m.statusMessageTimer != nil {
		m.statusMessageTimer.Stop()
	}
	m.statusMessageTimer = time.NewTimer(ui.StatusMessageTimeout)

	return ui.WaitForStatusMessageTimeout(m.statusMessageTimer)
}

func (m *CSV2MDModel) showStatusMessage(msg ui.PagerStatusMsg) tea.Cmd {
	m.state = ui.PagerStateStatusMessage
	m.statusMessage = msg.Message
	if m.statusMessageTimer != nil {
		m.statusMessageTimer.Stop()
	}
	m.statusMessageTimer = time.NewTimer(ui.StatusMessageTimeout)

	return ui.WaitForStatusMessageTimeout(m.statusMessageTimer)
}

func (m *CSV2MDModel) SetContent(content string) error {
	m.content = content

	reader := csv.NewReader(strings.NewReader(content))

	// Read all records
	rows, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("error reading content: %v", err)
	}

	if len(rows) == 0 {
		return fmt.Errorf("empty content")
	}

	// Convert to markdown
	markdownLines := csv2md.Convert("", rows, m.alignColumns)
	m.converted_content = strings.Join(markdownLines, "\n")

	// Set content in viewport
	var buf bytes.Buffer
	buf.WriteString(m.converted_content)
	m.viewport.SetContent(buf.String())

	return nil
}

func (m *CSV2MDModel) setSize(w, h int) {
	m.viewport.Width = w
	m.viewport.Height = h - ui.StatusBarHeight

	if m.showHelp {
		if pagerHelpHeight == 0 {
			pagerHelpHeight = strings.Count(m.helpView(), "\n")
		}
		m.viewport.Height -= (ui.StatusBarHeight + pagerHelpHeight)
	}
}

func (m *CSV2MDModel) toggleHelp() {
	m.showHelp = !m.showHelp
	m.setSize(m.common.Width, m.common.Height)

	if m.viewport.PastBottom() {
		m.viewport.GotoBottom()
	}
}

func (m CSV2MDModel) statusBarView(b *strings.Builder) {
	const (
		minPercent               float64 = 0.0
		maxPercent               float64 = 1.0
		percentToStringMagnitude float64 = 100.0
	)
	showStatusMessage := m.state == ui.PagerStateStatusMessage
	showErrorMessage := m.state == ui.PagerStateErrorMessage
	appName := ui.AppNameStyle(" " + Title + " ")

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

	if showErrorMessage {
		note = ui.StatusBarErrorStyle(note)
	} else if showStatusMessage {
		note = ui.StatusBarMessageStyle(note)
	} else {
		note = ui.StatusBarNoteStyle(note)
	}

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

func (m CSV2MDModel) helpView() (s string) {
	col1 := []string{
		"c              copy markdown",
		"e              edit",
		"v              paste to convert",
		"a              toggle column alignment",
		"q/ctrl+c       quit",
	}

	s += "\n"
	s += "k/↑      up                  " + col1[0] + "\n"
	s += "j/↓      down                " + col1[1] + "\n"
	s += "b/pgup   page up             " + col1[2] + "\n"
	s += "f/pgdn   page down           " + col1[3] + "\n"
	s += "u        ½ page up           " + col1[4] + "\n"
	s += "d        ½ page down         "

	if len(col1) > 5 {
		s += col1[5]
	}

	s = ui.Indent(s, 2)

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
