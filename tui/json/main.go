package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/chroma/quick"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
	"golang.org/x/term"
)

var (
	pagerHelpHeight int

	statusBarNoteFg = lipgloss.AdaptiveColor{Light: "#656565", Dark: "#7D7D7D"}

	helpViewStyle = lipgloss.NewStyle().
			Foreground(statusBarNoteFg).
			Background(lipgloss.AdaptiveColor{Light: "#f2f2f2", Dark: "#1B1B1B"}).
			Render

	statusBarHelpStyle = lipgloss.NewStyle().
				Foreground(statusBarNoteFg).
				Background(lipgloss.AdaptiveColor{Light: "#DCDCDC", Dark: "#323232"}).
				Render
)

const (
	statusBarHeight = 1
)

type model struct {
	width             int
	height            int
	formatted_content string
	content           string
	viewport          viewport.Model
	showHelp          bool
	ready             bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		case "?":
			m.toggleHelp()
			if m.viewport.HighPerformanceRendering {
				cmds = append(cmds, viewport.Sync(m.viewport))
			}
		}
	case editorFinishedMsg:
		if msg.err != nil {
			// Handle error
			return m, nil
		}
		m.content = msg.content
		m.formatted_content = formatJSON(msg.content)
		var buf bytes.Buffer
		_ = quick.Highlight(&buf, m.formatted_content, "json", "terminal", "nord")
		m.viewport.SetContent(buf.String())
		return m, nil
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

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

func (m model) View() string {
	var b strings.Builder

	fmt.Fprint(&b, m.viewport.View()+"\n")

	// Footer
	m.statusBarView(&b)

	if m.showHelp {
		fmt.Fprint(&b, "\n"+m.helpView())
	}

	return b.String()
}

func (m *model) setSize(w, h int) {
	m.viewport.Width = w
	m.viewport.Height = h - statusBarHeight

	if m.showHelp {
		if pagerHelpHeight == 0 {
			pagerHelpHeight = strings.Count(m.helpView(), "\n")
		}
		m.viewport.Height -= (statusBarHeight + pagerHelpHeight)
	}
}

func (m *model) toggleHelp() {
	m.showHelp = !m.showHelp
	m.setSize(m.width, m.height)

	if m.viewport.PastBottom() {
		m.viewport.GotoBottom()
	}
}

func (m model) statusBarView(b *strings.Builder) {
	if m.ready {
		fmt.Fprint(b, "Ready")
	} else {
		fmt.Fprint(b, "Initializing...")
	}
}

func (m model) helpView() (s string) {
	col1 := []string{
		"g/home  go to top",
		"G/end   go to bottom",
		"c       copy contents",
		"e       edit this document",
		"r       reload this document",
		"esc     back to files",
		"q       quit",
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

	s = indent(s, 2)

	// Fill up empty cells with spaces for background coloring
	if m.width > 0 {
		lines := strings.Split(s, "\n")
		for i := 0; i < len(lines); i++ {
			l := runewidth.StringWidth(lines[i])
			n := max(m.width-l, 0)
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

func main() {
	w, h, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		fmt.Println("Error getting terminal size:", err)
		os.Exit(1)
	}
	viewport := viewport.New(w, h-statusBarHeight)
	viewport.YPosition = 0
	viewport.HighPerformanceRendering = true

	m := model{content: "", viewport: viewport}

	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
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
