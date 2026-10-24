package json2toml

import (
	"bytes"
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/alecthomas/chroma/quick"

	"github.com/skatkov/devtui/internal/clipboard"
	"github.com/skatkov/devtui/internal/converter"
	"github.com/skatkov/devtui/internal/editor"
	"github.com/skatkov/devtui/internal/ui"
)

const Title = "JSON to TOML Converter"

type JsonTomlModel struct {
	ui.BasePagerModel
}

func NewJsonTomlModel(common *ui.CommonModel) JsonTomlModel {
	model := JsonTomlModel{
		BasePagerModel: ui.NewBasePagerModel(common, Title),
	}

	return model
}

func (m JsonTomlModel) Init() tea.Cmd {
	return m.BasePagerModel.Init()
}

func (m JsonTomlModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
	m.Viewport, cmd = m.Viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// loadContent shows content from the clipboard or editor, or the error that came with it.
func (m *JsonTomlModel) loadContent(content string, err error) tea.Cmd {
	if err == nil {
		err = m.SetContent(content)
	}
	if err != nil {
		return m.ShowErrorMessage(err.Error())
	}
	return m.ShowStatusMessage("Pasted. Press 'c' to copy result.")
}

func (m JsonTomlModel) View() tea.View {
	var b strings.Builder

	fmt.Fprint(&b, m.Viewport.View()+"\n")
	fmt.Fprint(&b, m.StatusBarView())

	if m.ShowHelp {
		fmt.Fprint(&b, "\n"+m.helpView())
	}

	return m.NewView(b.String())
}

func (m *JsonTomlModel) SetContent(content string) error {
	m.Content = content

	tomlStr, err := Convert(content)

	if err != nil {
		return fmt.Errorf("error converting JSON to TOML: %v", err)
	} else {
		m.FormattedContent = tomlStr
	}

	var buf bytes.Buffer
	_ = quick.Highlight(&buf, m.FormattedContent, "TOML", "terminal", "nord")
	m.Viewport.SetContent(buf.String())
	return nil
}

func (m JsonTomlModel) helpView() string {
	col1 := []string{
		"c              copy TOML",
		"e              edit JSON",
		"v              paste JSON to convert",
		"q/ctrl+c       quit",
	}

	return m.FormatHelpColumns(col1)
}

func Convert(jsonContent string) (string, error) {
	return converter.JSONToTOML(jsonContent)
}
