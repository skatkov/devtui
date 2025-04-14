package jsontoml

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
	"github.com/pelletier/go-toml/v2"
	"github.com/skatkov/devtui/internal/editor"
	"github.com/skatkov/devtui/internal/ui"
	"github.com/tiagomelo/go-clipboard/clipboard"
)

const Title = "JSON to TOML Converter"

var (
	pagerHelpHeight int
	useJsonNumber   = true // Enable JSON number handling
)

type JsonTomlModel struct {
	common *ui.CommonModel

	converted_content string
	content           string
	viewport          viewport.Model
	showHelp          bool
	ready             bool
	state             ui.PagerState

	statusMessage      string
	statusMessageTimer *time.Timer
}

func NewJsonTomlModel(common *ui.CommonModel) JsonTomlModel {
	model := JsonTomlModel{
		content: "",
		ready:   false,
		common:  common,
		state:   ui.PagerStateBrowse,
	}

	model.setSize(common.Width, common.Height)

	return model
}

func (m JsonTomlModel) Init() tea.Cmd {
	return nil
}

func (m JsonTomlModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				panic(err)
			}
			m.setContent(content)

			cmds = append(cmds, m.showStatusMessage(ui.PagerStatusMsg{Message: "Converted JSON to TOML"}))

		case "c":
			c := clipboard.New()
			if err := c.CopyText(m.converted_content); err != nil {
				panic(err)
			}

			cmds = append(cmds, m.showStatusMessage(ui.PagerStatusMsg{Message: "Copied TOML"}))
		case "?":
			m.toggleHelp()
			if m.viewport.HighPerformanceRendering {
				cmds = append(cmds, viewport.Sync(m.viewport))
			}
		}
	case ui.StatusMessageTimeoutMsg:
		m.state = ui.PagerStateBrowse
	case editor.EditorFinishedMsg:
		if msg.Err != nil {
			panic(msg.Err)
		}
		m.setContent(msg.Content)

		cmds = append(cmds, m.showStatusMessage(ui.PagerStatusMsg{Message: "Converted JSON to TOML"}))

	case tea.WindowSizeMsg:
		m.common.Width = msg.Width
		m.common.Height = msg.Height

		m.setSize(msg.Width, msg.Height)

		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-ui.StatusBarHeight)
			m.viewport.YPosition = 0
			m.viewport.HighPerformanceRendering = true
			m.viewport.SetContent(m.content)
			m.ready = true
		} else {
			m.setSize(msg.Width, msg.Height)
		}
	}
	if m.viewport.HighPerformanceRendering {
		cmds = append(cmds, viewport.Sync(m.viewport))
	}
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m JsonTomlModel) View() string {
	var b strings.Builder

	fmt.Fprint(&b, m.viewport.View()+"\n")

	m.statusBarView(&b)

	if m.showHelp {
		fmt.Fprint(&b, "\n"+m.helpView())
	}

	return b.String()
}

func (m *JsonTomlModel) showStatusMessage(msg ui.PagerStatusMsg) tea.Cmd {
	m.state = ui.PagerStateStatusMessage
	m.statusMessage = msg.Message
	if m.statusMessageTimer != nil {
		m.statusMessageTimer.Stop()
	}
	m.statusMessageTimer = time.NewTimer(ui.StatusMessageTimeout)

	return ui.WaitForStatusMessageTimeout(m.statusMessageTimer)
}

func (m *JsonTomlModel) setContent(content string) {
	m.content = content

	tomlStr, err := convert(content)

	if err != nil {
		m.converted_content = fmt.Sprintf("Error converting JSON to TOML: %v", err)
	} else {
		m.converted_content = tomlStr
	}

	var buf bytes.Buffer
	_ = quick.Highlight(&buf, m.converted_content, "TOML", "terminal", "nord")
	m.viewport.SetContent(buf.String())
}

func (m *JsonTomlModel) setSize(w, h int) {
	m.viewport.Width = w
	m.viewport.Height = h - ui.StatusBarHeight

	if m.showHelp {
		if pagerHelpHeight == 0 {
			pagerHelpHeight = strings.Count(m.helpView(), "\n")
		}
		m.viewport.Height -= (ui.StatusBarHeight + pagerHelpHeight)
	}
}

func (m *JsonTomlModel) toggleHelp() {
	m.showHelp = !m.showHelp
	m.setSize(m.common.Width, m.common.Height)

	if m.viewport.PastBottom() {
		m.viewport.GotoBottom()
	}
}

func (m JsonTomlModel) statusBarView(b *strings.Builder) {
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
		note = "Press 'v' to paste JSON"
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

func (m JsonTomlModel) helpView() (s string) {
	col1 := []string{
		"c              copy TOML",
		"e              edit JSON",
		"v              paste JSON to convert",
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

func convert(jsonContent string) (string, error) {
	var v any

	// Create a decoder that uses JSON numbers
	decoder := json.NewDecoder(strings.NewReader(jsonContent))
	if useJsonNumber {
		decoder.UseNumber()
	}

	// Decode JSON
	err := decoder.Decode(&v)
	if err != nil {
		return "", fmt.Errorf("JSON parsing error: %s", err.Error())
	}

	// Create a buffer to hold the TOML output
	var buf bytes.Buffer

	// Create a TOML encoder
	encoder := toml.NewEncoder(&buf)
	if useJsonNumber {
		encoder.SetMarshalJsonNumbers(true)
	}

	// Encode to TOML
	err = encoder.Encode(v)
	if err != nil {
		return "", fmt.Errorf("TOML encoding error: %s", err.Error())
	}

	return buf.String(), nil
}
