package yaml

import (
	"bytes"
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/alecthomas/chroma/quick"

	"github.com/skatkov/devtui/internal/clipboard"
	"github.com/skatkov/devtui/internal/editor"
	"github.com/skatkov/devtui/internal/ui"
	"github.com/skatkov/devtui/internal/yamlfmt"
)

const Title = "YAML Formatter"

type YamlModel struct {
	ui.BasePagerModel
}

func NewYamlModel(common *ui.CommonModel) YamlModel {
	model := YamlModel{
		BasePagerModel: ui.NewBasePagerModel(common, Title),
	}

	return model
}

func (m YamlModel) Init() tea.Cmd {
	return m.BasePagerModel.Init()
}

func (m YamlModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			return m, editor.OpenEditor(m.Content, "yaml")
		case "v":
			cmds = append(cmds, m.applyContent(clipboard.Paste()))
		}
	case ui.StatusMessageTimeoutMsg:
		m.State = ui.PagerStateBrowse
	case editor.EditorFinishedMsg:
		cmds = append(cmds, m.applyContent(msg.Content, msg.Err))

	case tea.WindowSizeMsg:
		cmd = m.HandleWindowSizeMsg(msg)
		cmds = append(cmds, cmd)
	}
	m.Viewport, cmd = m.Viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m YamlModel) View() tea.View {
	var b strings.Builder

	fmt.Fprint(&b, m.Viewport.View()+"\n")
	fmt.Fprint(&b, m.StatusBarView())

	if m.ShowHelp {
		fmt.Fprint(&b, "\n"+m.helpView())
	}

	return m.NewView(b.String())
}

func (m *YamlModel) SetContent(content string) error {
	m.Content = content
	formattedContent, err := yamlfmt.Format(content)
	if err != nil {
		return err
	}
	m.FormattedContent = formattedContent
	var buf bytes.Buffer

	err = quick.Highlight(&buf, m.FormattedContent, "yaml", "terminal", "nord")
	if err != nil {
		return err
	}
	m.Viewport.SetContent(buf.String())

	return nil
}

func (m YamlModel) helpView() string {
	col1 := []string{
		"c              copy formatted YAML",
		"e              edit unformatted YAML",
		"v              paste unformatted YAML",
		"q/ctrl+c       quit",
	}

	return m.FormatHelpColumns(col1)
}

// applyContent sets content that came from the clipboard or the editor and
// reports the outcome in the status bar.
func (m *YamlModel) applyContent(content string, err error) tea.Cmd {
	if err == nil {
		err = m.SetContent(content)
	}
	if err != nil {
		return m.ShowErrorMessage(err.Error())
	}
	return m.ShowStatusMessage("Pasted. Press 'c' to copy result.")
}
