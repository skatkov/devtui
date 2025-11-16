package jsontoml

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/alecthomas/chroma/quick"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mattn/go-runewidth"
	"github.com/pelletier/go-toml/v2"

	"github.com/skatkov/devtui/internal/clipboard"
	"github.com/skatkov/devtui/internal/editor"
	"github.com/skatkov/devtui/internal/ui"
)

const Title = "JSON to TOML Converter"

var useJsonNumber = true // Enable JSON number handling

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
	case tea.KeyMsg:
		if cmd, handled := m.HandleCommonKeys(msg); handled {
			return m, cmd
		}
		switch msg.String() {
		case "e":
			return m, editor.OpenEditor(m.Content, "json")
		case "v":
			content, err := clipboard.Paste()
			if err != nil {
				cmds = append(cmds, m.ShowErrorMessage(err.Error()))
			} else {
				err = m.SetContent(content)

				if err != nil {
					cmds = append(cmds, m.ShowErrorMessage(err.Error()))
				} else {
					cmds = append(cmds, m.ShowStatusMessage("Pasted. Press 'c' to copy result."))
				}
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

func (m JsonTomlModel) View() string {
	var b strings.Builder

	fmt.Fprint(&b, m.Viewport.View()+"\n")
	fmt.Fprint(&b, m.StatusBarView())

	if m.ShowHelp {
		fmt.Fprint(&b, "\n"+m.helpView())
	}

	return b.String()
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

func (m JsonTomlModel) helpView() (s string) {
	col1 := []string{
		"c              copy TOML",
		"e              edit JSON",
		"v              paste JSON to convert",
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

func Convert(jsonContent string) (string, error) {
	var v any

	// Create a decoder that uses JSON numbers
	decoder := json.NewDecoder(strings.NewReader(jsonContent))
	if useJsonNumber {
		decoder.UseNumber()
	}

	// Decode JSON
	err := decoder.Decode(&v)
	if err != nil {
		return "", fmt.Errorf("JSON parsing error: %s", err.Error())
	}

	// Create a buffer to hold the TOML output
	var buf bytes.Buffer

	// Create a TOML encoder
	encoder := toml.NewEncoder(&buf)
	if useJsonNumber {
		encoder.SetMarshalJsonNumbers(true)
	}

	// Encode to TOML
	err = encoder.Encode(v)
	if err != nil {
		return "", fmt.Errorf("TOML encoding error: %s", err.Error())
	}

	return buf.String(), nil
}
