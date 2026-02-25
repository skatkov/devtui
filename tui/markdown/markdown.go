package markdown

import (
	"bytes"
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/glamour"
	"github.com/mattn/go-runewidth"
	"github.com/skatkov/devtui/internal/clipboard"
	"github.com/skatkov/devtui/internal/editor"
	"github.com/skatkov/devtui/internal/ui"
)

const Title = "Markdown Renderer"

type MarkdownModel struct {
	ui.BasePagerModel
}

func NewMarkdownModel(common *ui.CommonModel) MarkdownModel {
	return MarkdownModel{
		BasePagerModel: ui.NewBasePagerModel(common, Title),
	}
}

func (m MarkdownModel) Init() tea.Cmd {
	return m.BasePagerModel.Init()
}

func (m MarkdownModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		// First check common keys
		if cmd, handled := m.HandleCommonKeys(msg); handled {
			return m, cmd
		}

		// Then handle module-specific keys
		switch msg.String() {
		case "e":
			return m, editor.OpenEditor(m.Content, "md")
		case "v":
			content, err := clipboard.Paste()
			if err != nil {
				cmds = append(cmds, m.ShowErrorMessage(err.Error()))
			} else {
				err = m.SetContent(content)
				if err != nil {
					cmds = append(cmds, m.ShowErrorMessage(err.Error()))
				} else {
					cmds = append(cmds, m.ShowStatusMessage("Pasted contents"))
				}
			}
		}
	case ui.StatusMessageTimeoutMsg:
		m.State = ui.PagerStateBrowse

	case editor.EditorFinishedMsg:
		if msg.Err != nil {
			cmds = append(cmds, m.ShowErrorMessage(msg.Err.Error()))
		} else {
			err := m.SetContent(msg.Content)
			if err != nil {
				cmds = append(cmds, m.ShowErrorMessage(err.Error()))
			} else {
				cmds = append(cmds, m.ShowStatusMessage("Markdown content edited"))
			}
		}
	case tea.WindowSizeMsg:
		cmd = m.HandleWindowSizeMsg(msg)
		cmds = append(cmds, cmd)
	}

	m.Viewport, cmd = m.Viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m MarkdownModel) View() tea.View {
	var b strings.Builder

	fmt.Fprint(&b, m.Viewport.View()+"\n")
	fmt.Fprint(&b, m.StatusBarView())

	if m.ShowHelp {
		fmt.Fprint(&b, "\n"+m.helpView())
	}

	return m.NewView(b.String())
}

func (m *MarkdownModel) SetContent(content string) error {
	m.Content = content

	if content == "" {
		m.FormattedContent = ""
		m.Viewport.SetContent("")
		return nil
	}

	// Render markdown with glamour
	out, err := glamour.Render(content, "dark")
	if err != nil {
		// If rendering fails, just show the raw content
		out = content
	}

	m.FormattedContent = out

	// Set content in viewport
	var buf bytes.Buffer
	buf.WriteString(m.FormattedContent)
	m.Viewport.SetContent(buf.String())

	return nil
}

func (m MarkdownModel) helpView() (s string) {
	col1 := []string{
		"c              copy rendered markdown",
		"e              edit markdown",
		"v              paste markdown",
		"q/ctrl+c       quit",
		"esc            return to menu",
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

	return ui.HelpViewStyle(s)
}
