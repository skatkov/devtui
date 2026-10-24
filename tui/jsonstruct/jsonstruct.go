package jsonstruct

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

const Title = "JSON to Go Struct Converter"

type JsonStructModel struct {
	ui.BasePagerModel
}

func NewJsonStructModel(common *ui.CommonModel) JsonStructModel {
	model := JsonStructModel{
		BasePagerModel: ui.NewBasePagerModel(common, Title),
	}

	return model
}

func (m JsonStructModel) Init() tea.Cmd {
	return m.BasePagerModel.Init()
}

func (m JsonStructModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
func (m *JsonStructModel) loadContent(content string, err error) tea.Cmd {
	if err == nil {
		err = m.SetContent(content)
	}
	if err != nil {
		return m.ShowErrorMessage(err.Error())
	}
	return m.ShowStatusMessage("Pasted. Press 'c' to copy result.")
}

func (m JsonStructModel) View() tea.View {
	var b strings.Builder

	fmt.Fprint(&b, m.Viewport.View()+"\n")
	fmt.Fprint(&b, m.StatusBarView())

	if m.ShowHelp {
		fmt.Fprint(&b, "\n"+m.helpView())
	}

	return m.NewView(b.String())
}

func (m *JsonStructModel) SetContent(content string) error {
	m.Content = content
	converted, err := structgen.JSONToGoStruct(strings.NewReader(content))
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

func (m JsonStructModel) helpView() string {
	col1 := []string{
		"c              copy Go struct",
		"e              edit JSON",
		"v              paste JSON to convert",
		"q/ctrl+c       quit",
	}

	return m.FormatHelpColumns(col1)
}
