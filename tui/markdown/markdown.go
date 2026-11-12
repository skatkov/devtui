package markdown

import (
	"bytes"
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/glamour"
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
			cmds = append(cmds, m.applyContent(content, err, "Pasted contents"))
		}
	case ui.StatusMessageTimeoutMsg:
		m.State = ui.PagerStateBrowse

	case editor.EditorFinishedMsg:
		cmds = append(cmds, m.applyContent(msg.Content, msg.Err, "Markdown content edited"))
	case tea.WindowSizeMsg:
		cmd = m.HandleWindowSizeMsg(msg)
		cmds = append(cmds, cmd)
	}

	m.Viewport, cmd = m.Viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// applyContent shows content read from the clipboard or editor, or the error that occurred.
func (m *MarkdownModel) applyContent(content string, err error, status string) tea.Cmd {
	if err == nil {
		err = m.SetContent(content)
	}
	if err != nil {
		return m.ShowErrorMessage(err.Error())
	}
	return m.ShowStatusMessage(status)
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

func (m MarkdownModel) helpView() string {
	col1 := []string{
		"c              copy rendered markdown",
		"e              edit markdown",
		"v              paste markdown",
		"q/ctrl+c       quit",
		"esc            return to menu",
	}

	return m.FormatHelpColumns(col1)
}
