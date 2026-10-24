package toml2json

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

const Title = "TOML to JSON Converter"

type TomlJsonModel struct {
	ui.BasePagerModel
}

func NewTomlJsonModel(common *ui.CommonModel) TomlJsonModel {
	model := TomlJsonModel{
		BasePagerModel: ui.NewBasePagerModel(common, Title),
	}

	return model
}

func (m TomlJsonModel) Init() tea.Cmd {
	return m.BasePagerModel.Init()
}

func (m TomlJsonModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			return m, editor.OpenEditor(m.Content, "toml")
		case "v":
			content, err := clipboard.Paste()
			cmds = append(cmds, m.applyContent(content, err))
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

// applyContent shows content read from the clipboard or editor, or the error that occurred.
func (m *TomlJsonModel) applyContent(content string, err error) tea.Cmd {
	if err == nil {
		err = m.SetContent(content)
	}
	if err != nil {
		return m.ShowErrorMessage(err.Error())
	}
	return m.ShowStatusMessage("Pasted. Press 'c' to copy result.")
}

func (m TomlJsonModel) View() tea.View {
	var b strings.Builder

	fmt.Fprint(&b, m.Viewport.View()+"\n")
	fmt.Fprint(&b, m.StatusBarView())

	if m.ShowHelp {
		fmt.Fprint(&b, "\n"+m.helpView())
	}

	return m.NewView(b.String())
}

func (m *TomlJsonModel) SetContent(content string) error {
	m.Content = content
	jsonStr, err := Convert(content)

	if err != nil {
		return err
	} else {
		m.FormattedContent = jsonStr
	}

	var buf bytes.Buffer
	err = quick.Highlight(&buf, m.FormattedContent, "json", "terminal", "nord")
	if err != nil {
		return err
	}
	m.Viewport.SetContent(buf.String())
	return nil
}

func (m TomlJsonModel) helpView() string {
	col1 := []string{
		"c              copy JSON",
		"e              edit TOML",
		"v              paste TOML to convert",
		"q/ctrl+c       quit",
	}

	return m.FormatHelpColumns(col1)
}

func Convert(tomlContent string) (string, error) {
	return converter.TOMLToJSON(tomlContent)
}
