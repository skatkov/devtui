package tomljson

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/alecthomas/chroma/quick"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mattn/go-runewidth"
	"github.com/pelletier/go-toml/v2"
	"github.com/skatkov/devtui/internal/editor"
	"github.com/skatkov/devtui/internal/ui"
	"github.com/skatkov/devtui/internal/clipboard"
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
	case tea.KeyMsg:
		if cmd, handled := m.HandleCommonKeys(msg); handled {
			return m, cmd
		}
		switch msg.String() {
		case "e":
			return m, editor.OpenEditor(m.Content, "toml")
		case "v":
			content, err := clipboard.Paste()
			if err == nil {
				err = m.SetContent(content)
			}

			if err != nil {
				cmds = append(cmds, m.ShowErrorMessage(err.Error()))
			} else {
				cmds = append(cmds, m.ShowStatusMessage("Pasted. Press 'c' to copy result."))
			}
		}
	case ui.StatusMessageTimeoutMsg:
		m.State = ui.PagerStateBrowse
	case editor.EditorFinishedMsg:
		if msg.Err != nil {
			return m, m.ShowErrorMessage(msg.Err.Error())
		} else {
			err := m.SetContent(msg.Content)

			if err != nil {
				cmds = append(cmds, m.ShowErrorMessage(err.Error()))
			} else {
				cmds = append(cmds, m.ShowStatusMessage("Pasted. Press 'c' to copy result."))
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

func (m TomlJsonModel) View() string {
	var b strings.Builder

	fmt.Fprint(&b, m.Viewport.View()+"\n")
	fmt.Fprint(&b, m.StatusBarView())

	if m.ShowHelp {
		fmt.Fprint(&b, "\n"+m.helpView())
	}

	return b.String()
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

func (m TomlJsonModel) helpView() (s string) {
	col1 := []string{
		"c              copy JSON",
		"e              edit TOML",
		"v              paste TOML to convert",
		"q/ctrl+c       quit",
	}

	s += "\n"
	s += "k/↑      up                  " + col1[0] + "\n"
	s += "j/↓      down                " + col1[1] + "\n"
	s += "b/pgup   page up             " + col1[2] + "\n"
	s += "f/pgdn   page down           " + col1[3] + "\n"
	s += "u        ½ page up           " + "\n"
	s += "d        ½ page down         "

	if len(col1) > 5 {
		s += col1[5]
	}

	s = ui.Indent(s, 2)

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

func Convert(tomlContent string) (string, error) {
	var v any

	err := toml.Unmarshal([]byte(tomlContent), &v)
	if err != nil {
		var derr *toml.DecodeError
		if errors.As(err, &derr) {
			return "", fmt.Errorf("TOML parsing error: %s", err.Error())
		}
		return "", err
	}

	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", fmt.Errorf("JSON encoding error: %s", err.Error())
	}

	return string(bytes), nil
}
