package yamlstruct

import (
	"bytes"
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/alecthomas/chroma/quick"

	"github.com/skatkov/devtui/internal/clipboard"
	"github.com/skatkov/devtui/internal/editor"
	structgen "github.com/skatkov/devtui/internal/structgen"
	"github.com/skatkov/devtui/internal/ui"
)

const Title = "YAML to Go Struct Converter"

type YamlStructModel struct {
	ui.BasePagerModel
}

func NewYamlStructModel(common *ui.CommonModel) YamlStructModel {
	model := YamlStructModel{
		BasePagerModel: ui.NewBasePagerModel(common, Title),
	}

	return model
}

func (m YamlStructModel) Init() tea.Cmd {
	return m.BasePagerModel.Init()
}

func (m YamlStructModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m YamlStructModel) View() tea.View {
	var b strings.Builder

	fmt.Fprint(&b, m.Viewport.View()+"\n")
	fmt.Fprint(&b, m.StatusBarView())

	if m.ShowHelp {
		fmt.Fprint(&b, "\n"+m.helpView())
	}

	return m.NewView(b.String())
}

func (m *YamlStructModel) SetContent(content string) error {
	m.Content = content
	converted, err := structgen.YAMLToGoStruct(strings.NewReader(content))
	if err != nil {
		return err
	}
	m.FormattedContent = converted
	var buf bytes.Buffer

	err = quick.Highlight(&buf, m.FormattedContent, "go", "terminal", "nord")
	if err != nil {
		return err
	}
	m.Viewport.SetContent(buf.String())
	return nil
}

func (m YamlStructModel) helpView() string {
	col1 := []string{
		"c              copy Go struct",
		"e              edit YAML",
		"v              paste YAML to convert",
		"q/ctrl+c       quit",
	}

	return m.FormatHelpColumns(col1)
}

// applyContent sets content that came from the clipboard or the editor and
// reports the outcome in the status bar.
func (m *YamlStructModel) applyContent(content string, err error) tea.Cmd {
	if err == nil {
		err = m.SetContent(content)
	}
	if err != nil {
		return m.ShowErrorMessage(err.Error())
	}
	return m.ShowStatusMessage("Pasted. Press 'c' to copy result.")
}
