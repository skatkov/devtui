package tomljson

import (
	"bytes"
	"encoding/json"
	"errors"
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

const Title = "TOML to JSON"

var pagerHelpHeight int

type TomlJsonModel struct {
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

func NewTomlJsonModel(common *ui.CommonModel) TomlJsonModel {
	model := TomlJsonModel{
		content: "",
		ready:   false,
		common:  common,
		state:   ui.PagerStateBrowse,
	}

	model.setSize(common.Width, common.Height)

	return model
}

func (m TomlJsonModel) Init() tea.Cmd {
	return nil
}

func (m TomlJsonModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "e":
			return m, editor.OpenEditor(m.content, "toml")
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

			cmds = append(cmds, m.showStatusMessage(ui.PagerStatusMsg{Message: "Converted TOML to JSON"}))

		case "c":
			c := clipboard.New()
			if err := c.CopyText(m.converted_content); err != nil {
				panic(err)
			}

			cmds = append(cmds, m.showStatusMessage(ui.PagerStatusMsg{Message: "Copied JSON"}))
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

		cmds = append(cmds, m.showStatusMessage(ui.PagerStatusMsg{Message: "Converted TOML to JSON"}))

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

func (m TomlJsonModel) View() string {
	var b strings.Builder

	fmt.Fprint(&b, m.viewport.View()+"\n")

	m.statusBarView(&b)

	if m.showHelp {
		fmt.Fprint(&b, "\n"+m.helpView())
	}

	return b.String()
}

func (m *TomlJsonModel) showStatusMessage(msg ui.PagerStatusMsg) tea.Cmd {
	m.state = ui.PagerStateStatusMessage
	m.statusMessage = msg.Message
	if m.statusMessageTimer != nil {
		m.statusMessageTimer.Stop()
	}
	m.statusMessageTimer = time.NewTimer(ui.StatusMessageTimeout)

	return ui.WaitForStatusMessageTimeout(m.statusMessageTimer)
}

func (m *TomlJsonModel) setContent(content string) {
	m.content = content

	jsonStr, err := convert(content)

	if err != nil {
		m.converted_content = fmt.Sprintf("Error converting TOML to JSON: %v", err)
	} else {
		m.converted_content = jsonStr
	}

	var buf bytes.Buffer
	_ = quick.Highlight(&buf, m.converted_content, "json", "terminal", "nord")
	m.viewport.SetContent(buf.String())
}

func (m *TomlJsonModel) setSize(w, h int) {
	m.viewport.Width = w
	m.viewport.Height = h - ui.StatusBarHeight

	if m.showHelp {
		if pagerHelpHeight == 0 {
			pagerHelpHeight = strings.Count(m.helpView(), "\n")
		}
		m.viewport.Height -= (ui.StatusBarHeight + pagerHelpHeight)
	}
}

func (m *TomlJsonModel) toggleHelp() {
	m.showHelp = !m.showHelp
	m.setSize(m.common.Width, m.common.Height)

	if m.viewport.PastBottom() {
		m.viewport.GotoBottom()
	}
}

func (m TomlJsonModel) statusBarView(b *strings.Builder) {
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
		note = "Press 'v' to paste TOML"
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

func (m TomlJsonModel) helpView() (s string) {
	col1 := []string{
		"c              copy JSON",
		"e              edit TOML",
		"v              paste TOML to convert",
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

func convert(tomlContent string) (string, error) {
	var v any

	err := toml.Unmarshal([]byte(tomlContent), &v)
	if err != nil {
		var derr *toml.DecodeError
		if errors.As(err, &derr) {
			return "", fmt.Errorf("TOML parsing error: %s", err.Error())
		}
		return "", err
	}

	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", fmt.Errorf("JSON encoding error: %s", err.Error())
	}

	return string(bytes), nil
}
