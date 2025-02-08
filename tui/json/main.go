package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	initialInputs = 2
	maxInputs     = 2
	minInputs     = 1
	helpHeight    = 5
)

var (
	cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))

	cursorLineStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("57")).
			Foreground(lipgloss.Color("230"))

	placeholderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("238"))

	endOfBufferStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("235"))

	focusedPlaceholderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("99"))

	focusedBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("238"))

	blurredBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.HiddenBorder())

	jsonKeyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("203"))

	jsonValueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86"))

	jsonStringStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("83"))

	jsonNumberStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("141"))
)

type keymap = struct {
	next, prev, quit, up, down key.Binding
}

func newTextInput() textinput.Model {
	t := textinput.New()
	t.Prompt = ""
	t.Placeholder = "Paste JSON here"
	t.PlaceholderStyle = placeholderStyle
	t.CharLimit = -1
	return t
}

type model struct {
	width     int
	height    int
	keymap    keymap
	help      help.Model
	textinput textinput.Model
	viewport  viewport.Model
	focus     int
}

func newModel() model {
	ti := newTextInput()
	vp := viewport.New(0, 0)
	vp.Style = blurredBorderStyle

	m := model{
		help:      help.New(),
		textinput: ti,
		viewport:  vp,
		keymap: keymap{
			next: key.NewBinding(
				key.WithKeys("tab"),
				key.WithHelp("tab", "next"),
			),
			prev: key.NewBinding(
				key.WithKeys("shift+tab"),
				key.WithHelp("shift+tab", "prev"),
			),
			quit: key.NewBinding(
				key.WithKeys("esc", "ctrl+c"),
				key.WithHelp("esc", "quit"),
			),
			up: key.NewBinding(
				key.WithKeys("up"),
				key.WithHelp("↑", "scroll up"),
			),
			down: key.NewBinding(
				key.WithKeys("down"),
				key.WithHelp("↓", "scroll down"),
			),
		},
	}

	m.textinput.Focus()
	return m
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func formatJSON(input string) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(input), "", "  "); err != nil {
		return input // Return original if not valid JSON
	}
	return prettyJSON.String()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.quit):
			return m, tea.Quit
		case key.Matches(msg, m.keymap.next):
			if m.focus == 0 {
				m.focus = 1
			} else {
				m.focus = 0
				cmds = append(cmds, m.textinput.Focus())
			}
		case key.Matches(msg, m.keymap.prev):
			if m.focus == 1 {
				m.focus = 0
				cmds = append(cmds, m.textinput.Focus())
			} else {
				m.focus = 1
			}
		case key.Matches(msg, m.keymap.up):
			if m.focus == 1 {
				m.viewport.LineUp(1)
			}
		case key.Matches(msg, m.keymap.down):
			if m.focus == 1 {
				m.viewport.LineDown(1)
			}
		}
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		m.sizeComponents()
	}

	// Handle JSON formatting
	if m.focus == 0 {
		formattedJSON := formatJSON(m.textinput.Value())
		m.viewport.SetContent(formattedJSON)
	}

	// Update active component
	if m.focus == 0 {
		newTextinput, cmd := m.textinput.Update(msg)
		m.textinput = newTextinput
		cmds = append(cmds, cmd)
	} else {
		newViewport, cmd := m.viewport.Update(msg)
		m.viewport = newViewport
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *model) sizeComponents() {
	m.textinput.Width = m.width / 2
	m.viewport.Width = m.width / 2
	m.viewport.Height = m.height - helpHeight
}

func (m model) View() string {
	help := m.help.ShortHelpView([]key.Binding{
		m.keymap.next,
		m.keymap.prev,
		m.keymap.quit,
		m.keymap.up,
		m.keymap.down,
	})

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.textinput.View(),
		m.viewport.View(),
	) + "\n\n" + help
}

func main() {
	if _, err := tea.NewProgram(newModel(), tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error while running program:", err)
		os.Exit(1)
	}
}
