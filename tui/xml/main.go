package xml

import (
	"bytes"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/alecthomas/chroma/quick"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-xmlfmt/xmlfmt"
	"github.com/mattn/go-runewidth"
	"github.com/muesli/ansi"
	"github.com/muesli/reflow/truncate"
	"github.com/skatkov/devtui/internal/editor"
	"github.com/skatkov/devtui/internal/ui"
	"github.com/tiagomelo/go-clipboard/clipboard"
)

const Title = "XML Formatter"

var pagerHelpHeight int

type XMLFormatterModel struct {
	common *ui.CommonModel

	formattedContent string
	content          string
	viewport         viewport.Model
	showHelp         bool
	ready            bool
	state            ui.PagerState

	statusMessage      string
	statusMessageTimer *time.Timer
}

func NewXMLFormatterModel(common *ui.CommonModel) XMLFormatterModel {
	model := XMLFormatterModel{
		content: "",
		ready:   false,
		common:  common,
		state:   ui.PagerStateBrowse,
	}

	model.setSize(common.Width, common.Height)

	return model
}

func (m XMLFormatterModel) Init() tea.Cmd {
	return nil
}

func (m XMLFormatterModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "e":
			return m, editor.OpenEditor(m.content, "xml")
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
				panic(err)
			}
			m.SetContent(content)

			cmds = append(cmds, m.showStatusMessage(ui.PagerStatusMsg{Message: "Formatted XML"}))

		case "c":
			c := clipboard.New()
			if err := c.CopyText(m.formattedContent); err != nil {
				panic(err)
			}

			cmds = append(cmds, m.showStatusMessage(ui.PagerStatusMsg{Message: "Copied formatted XML"}))
		case "?":
			m.toggleHelp()
		}
	case ui.StatusMessageTimeoutMsg:
		m.state = ui.PagerStateBrowse
	case editor.EditorFinishedMsg:
		if msg.Err != nil {
			panic(msg.Err)
		}
		m.SetContent(msg.Content)

		cmds = append(cmds, m.showStatusMessage(ui.PagerStatusMsg{Message: "Formatted XML"}))

	case tea.WindowSizeMsg:
		m.common.Width = msg.Width
		m.common.Height = msg.Height

		m.setSize(msg.Width, msg.Height)

		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-ui.StatusBarHeight)
			m.viewport.YPosition = 0
			m.SetContent(m.content)
			m.ready = true
		} else {
			m.setSize(msg.Width, msg.Height)
		}
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m XMLFormatterModel) View() string {
	var b strings.Builder

	fmt.Fprint(&b, m.viewport.View()+"\n")

	m.statusBarView(&b)

	if m.showHelp {
		fmt.Fprint(&b, "\n"+m.helpView())
	}

	return b.String()
}

func (m *XMLFormatterModel) showStatusMessage(msg ui.PagerStatusMsg) tea.Cmd {
	m.state = ui.PagerStateStatusMessage
	m.statusMessage = msg.Message
	if m.statusMessageTimer != nil {
		m.statusMessageTimer.Stop()
	}
	m.statusMessageTimer = time.NewTimer(ui.StatusMessageTimeout)

	return ui.WaitForStatusMessageTimeout(m.statusMessageTimer)
}

func (m *XMLFormatterModel) SetContent(content string) {
	m.content = content
	m.formattedContent = xmlfmt.FormatXML(content, "\t", "  ")

	var buf bytes.Buffer
	_ = quick.Highlight(&buf, m.formattedContent, "xml", "terminal", "nord")
	m.viewport.SetContent(buf.String())
}

func (m *XMLFormatterModel) setSize(w, h int) {
	m.viewport.Width = w
	m.viewport.Height = h - ui.StatusBarHeight

	if m.showHelp {
		if pagerHelpHeight == 0 {
			pagerHelpHeight = strings.Count(m.helpView(), "\n")
		}
		m.viewport.Height -= (ui.StatusBarHeight + pagerHelpHeight)
	}
}

func (m *XMLFormatterModel) toggleHelp() {
	m.showHelp = !m.showHelp
	m.setSize(m.common.Width, m.common.Height)

	if m.viewport.PastBottom() {
		m.viewport.GotoBottom()
	}
}

func (m XMLFormatterModel) statusBarView(b *strings.Builder) {
	const (
		minPercent               float64 = 0.0
		maxPercent               float64 = 1.0
		percentToStringMagnitude float64 = 100.0
	)
	showStatusMessage := m.state == ui.PagerStateStatusMessage
	appName := ui.AppNameStyle(" " + Title + " ")

	scrollPercent := ""
	if m.content != "" {
		percent := math.Max(minPercent, math.Min(maxPercent, m.viewport.ScrollPercent()))
		scrollPercent = fmt.Sprintf(" %3.f%% ", percent*percentToStringMagnitude)
		scrollPercent = ui.StatusBarScrollPosStyle(scrollPercent)
	}
	var helpNote string
	if showStatusMessage {
		helpNote = ui.StatusBarMessageHelpStyle(" ? Help ")
	} else {
		helpNote = ui.StatusBarHelpStyle(" ? Help ")
	}

	var note string
	if showStatusMessage {
		note = m.statusMessage
	} else if m.content == "" {
		note = "Press 'v' to paste XML to format"
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

	padding := max(0,
		m.common.Width-
			ansi.PrintableRuneWidth(appName)-
			ansi.PrintableRuneWidth(note)-
			ansi.PrintableRuneWidth(scrollPercent)-
			ansi.PrintableRuneWidth(helpNote),
	)
	emptySpace := strings.Repeat(" ", padding)
	if showStatusMessage {
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

func (m XMLFormatterModel) helpView() (s string) {
	col1 := []string{
		"c              copy formatted XML",
		"e              edit XML",
		"v              paste XML to format",
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
