package json

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/alecthomas/chroma/quick"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
	"github.com/muesli/ansi"
	"github.com/muesli/reflow/truncate"
	"github.com/skatkov/devtui/internal/ui"
	"github.com/tiagomelo/go-clipboard/clipboard"
)

var (
	pagerHelpHeight int
	fuchsia         = lipgloss.AdaptiveColor{Light: "#EE6FF8", Dark: "#EE6FF8"}
	statusBarNoteFg = lipgloss.AdaptiveColor{Light: "#656565", Dark: "#7D7D7D"}

	helpViewStyle = lipgloss.NewStyle().
			Foreground(statusBarNoteFg).
			Background(lipgloss.AdaptiveColor{Light: "#f2f2f2", Dark: "#1B1B1B"}).
			Render

	statusBarBg = lipgloss.AdaptiveColor{Light: "#E6E6E6", Dark: "#242424"}

	statusBarNoteStyle = lipgloss.NewStyle().
				Foreground(statusBarNoteFg).
				Background(statusBarBg).Render

	appNameStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ECFD65")).
			Background(fuchsia).
			Bold(true)

	statusBarHelpStyle = lipgloss.NewStyle().
				Foreground(statusBarNoteFg).
				Background(lipgloss.AdaptiveColor{Light: "#DCDCDC", Dark: "#323232"})

	statusBarScrollPosStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#949494", Dark: "#5A5A5A"}).
				Background(statusBarBg).
				Render
)

const (
	statusBarHeight      = 1
	ellipsis             = "…"
	statusMessageTimeout = time.Second * 3
)

type JsonModel struct {
	common *ui.CommonModel

	formatted_content string
	content           string
	viewport          viewport.Model
	showHelp          bool
	ready             bool

	statusMessage      string
	statusMessageTimer *time.Timer
}

func NewJsonModel(common *ui.CommonModel) JsonModel {
	model := JsonModel{
		content: "",
		ready:   false,
		common:  common,
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
			return m, openEditor(m.content, "json")
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
				fmt.Println(err)
				os.Exit(1)
			}
			m.setContent(string(content))
		case "c":
			c := clipboard.New()
			if err := c.CopyText(m.formatted_content); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		case "?":
			m.toggleHelp()
			if m.viewport.HighPerformanceRendering {
				cmds = append(cmds, viewport.Sync(m.viewport))
			}
		}
	case editorFinishedMsg:
		if msg.err != nil {
			panic(msg.err)
		}
		m.setContent(msg.content)
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
			m.viewport = viewport.New(msg.Width, msg.Height-statusBarHeight)
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

func (m *JsonModel) setContent(content string) {
	m.content = content
	m.formatted_content = formatJSON(content)
	var buf bytes.Buffer
	_ = quick.Highlight(&buf, m.formatted_content, "json", "terminal", "nord")
	m.viewport.SetContent(buf.String())
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

func (m *JsonModel) setSize(w, h int) {
	m.viewport.Width = w
	m.viewport.Height = h - statusBarHeight

	if m.showHelp {
		if pagerHelpHeight == 0 {
			pagerHelpHeight = strings.Count(m.helpView(), "\n")
		}
		m.viewport.Height -= (statusBarHeight + pagerHelpHeight)
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

	appName := appNameStyle.Render(" JSON Formatter ")

	// Scroll percent
	scrollPercent := ""
	if m.content != "" {
		percent := math.Max(minPercent, math.Min(maxPercent, m.viewport.ScrollPercent()))
		scrollPercent = fmt.Sprintf(" %3.f%% ", percent*percentToStringMagnitude)
		scrollPercent = statusBarScrollPosStyle(scrollPercent)
	}
	helpNote := statusBarHelpStyle.Render(" ? Help ")
	var note string
	if m.content == "" {
		note = "Press 'v' to paste unformatted JSON"
	}

	note = truncate.StringWithTail(" "+note+" ", uint(max(0,
		m.common.Width-
			ansi.PrintableRuneWidth(appName)-
			ansi.PrintableRuneWidth(scrollPercent)-
			ansi.PrintableRuneWidth(helpNote),
	)), ellipsis)

	note = statusBarNoteStyle(note)

	// Empty space
	padding := max(0,
		m.common.Width-
			ansi.PrintableRuneWidth(appName)-
			ansi.PrintableRuneWidth(note)-
			ansi.PrintableRuneWidth(scrollPercent)-
			ansi.PrintableRuneWidth(helpNote),
	)
	emptySpace := strings.Repeat(" ", padding)
	emptySpace = statusBarNoteStyle(emptySpace)

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

	s = indent(s, 2)

	// Fill up empty cells with spaces for background coloring
	if m.common.Width > 0 {
		lines := strings.Split(s, "\n")
		for i := 0; i < len(lines); i++ {
			l := runewidth.StringWidth(lines[i])
			n := max(m.common.Width-l, 0)
			lines[i] += strings.Repeat(" ", n)
		}

		s = strings.Join(lines, "\n")
	}

	return helpViewStyle(s)
}

func indent(s string, n int) string {
	if n <= 0 || s == "" {
		return s
	}
	l := strings.Split(s, "\n")
	b := strings.Builder{}
	i := strings.Repeat(" ", n)
	for _, v := range l {
		fmt.Fprintf(&b, "%s%s\n", i, v)
	}
	return b.String()
}

func formatJSON(content string) string {
	var data interface{}
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
