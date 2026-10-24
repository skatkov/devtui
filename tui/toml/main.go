package toml

import (
	"bytes"
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/alecthomas/chroma/quick"
	"github.com/pelletier/go-toml/v2"

	"github.com/skatkov/devtui/internal/clipboard"
	"github.com/skatkov/devtui/internal/editor"
	"github.com/skatkov/devtui/internal/ui"
)

const Title = "TOML Formatter"

type TomlFormatModel struct {
	ui.BasePagerModel
}

func NewTomlFormatModel(common *ui.CommonModel) TomlFormatModel {
	model := TomlFormatModel{
		BasePagerModel: ui.NewBasePagerModel(common, Title),
	}

	return model
}

func (m TomlFormatModel) Init() tea.Cmd {
	return m.BasePagerModel.Init()
}

func (m TomlFormatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
func (m *TomlFormatModel) applyContent(content string, err error) tea.Cmd {
	if err == nil {
		err = m.SetContent(content)
	}
	if err != nil {
		return m.ShowErrorMessage(err.Error())
	}
	return m.ShowStatusMessage("Pasted. Press 'c' to copy result.")
}

func (m TomlFormatModel) View() tea.View {
	var b strings.Builder

	fmt.Fprint(&b, m.Viewport.View()+"\n")
	fmt.Fprint(&b, m.StatusBarView())

	if m.ShowHelp {
		fmt.Fprint(&b, "\n"+m.helpView())
	}

	return m.NewView(b.String())
}

func (m *TomlFormatModel) SetContent(content string) error {
	m.Content = content
	formattedStr, err := Convert(content)

	if err != nil {
		return err
	} else {
		m.FormattedContent = formattedStr
	}

	var buf bytes.Buffer
	err = quick.Highlight(&buf, m.FormattedContent, "TOML", "terminal", "nord")
	if err != nil {
		return err
	}
	m.Viewport.SetContent(buf.String())
	return nil
}

func (m TomlFormatModel) helpView() string {
	col1 := []string{
		"c              copy formatted TOML",
		"e              edit TOML",
		"v              paste TOML to format",
		"q/ctrl+c       quit",
	}

	return m.FormatHelpColumns(col1)
}

func Convert(tomlContent string) (string, error) {
	var v any

	// Parse the TOML content
	err := toml.Unmarshal([]byte(tomlContent), &v)
	if err != nil {
		return "", fmt.Errorf("TOML parsing error: %s", err.Error())
	}

	// Create a buffer to hold the formatted TOML output
	var buf bytes.Buffer

	// Create a TOML encoder
	encoder := toml.NewEncoder(&buf)

	// Encode back to TOML (formatting happens during this process)
	err = encoder.Encode(v)
	if err != nil {
		return "", fmt.Errorf("TOML encoding error: %s", err.Error())
	}

	return buf.String(), nil
}
