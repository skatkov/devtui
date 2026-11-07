package json

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/alecthomas/chroma/quick"
	"github.com/skatkov/devtui/internal/clipboard"
	"github.com/skatkov/devtui/internal/editor"
	"github.com/skatkov/devtui/internal/ui"
)

const Title = "JSON Formatter"

type JsonModel struct {
	ui.BasePagerModel
}

func NewJsonModel(common *ui.CommonModel) JsonModel {
	model := JsonModel{
		BasePagerModel: ui.NewBasePagerModel(common, Title),
	}

	return model
}

func (m JsonModel) Init() tea.Cmd {
	return m.BasePagerModel.Init()
}

func (m JsonModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if cmd, handled := m.HandleCommonKeys(msg); handled {
			return m, cmd
		}
		switch msg.String() {
		case "e":
			return m, editor.OpenEditor(m.Content, "json")
		case "v":
			cmds = append(cmds, m.loadContent(clipboard.Paste()))
		}
	case ui.StatusMessageTimeoutMsg:
		m.State = ui.PagerStateBrowse

	case editor.EditorFinishedMsg:
		cmds = append(cmds, m.loadContent(msg.Content, msg.Err))
	case tea.WindowSizeMsg:
		cmd = m.HandleWindowSizeMsg(msg)
		cmds = append(cmds, cmd)
	}
	// Handle keyboard and mouse events in the viewport
	m.Viewport, cmd = m.Viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// loadContent shows content from the clipboard or editor, or the error that came with it.
func (m *JsonModel) loadContent(content string, err error) tea.Cmd {
	if err == nil {
		err = m.SetContent(content)
	}
	if err != nil {
		return m.ShowErrorMessage(err.Error())
	}
	return m.ShowStatusMessage("Pasted. Press 'c' to copy")
}

func (m JsonModel) View() tea.View {
	var b strings.Builder

	fmt.Fprint(&b, m.Viewport.View()+"\n")
	fmt.Fprint(&b, m.StatusBarView())

	if m.ShowHelp {
		fmt.Fprint(&b, "\n"+m.helpView())
	}

	return m.NewView(b.String())
}

func (m *JsonModel) SetContent(content string) error {
	m.Content = content
	m.FormattedContent = FormatJSON(content)

	var buf bytes.Buffer
	err := quick.Highlight(&buf, m.FormattedContent, "json", "terminal", "nord")
	if err != nil {
		return err
	}
	m.Viewport.SetContent(buf.String())

	return nil
}

func (m JsonModel) helpView() string {
	col1 := []string{
		"c              copy formatted JSON",
		"e              edit unformatted JSON",
		"v              paste unformatted JSON",
		"q/ctrl+c       quit",
	}

	return m.FormatHelpColumns(col1)
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
