package jsontoon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hannes-sistemica/toon"
	"github.com/mattn/go-runewidth"
	"github.com/skatkov/devtui/internal/editor"
	"github.com/skatkov/devtui/internal/ui"
	"github.com/tiagomelo/go-clipboard/clipboard"
)

const Title = "JSON to TOON Converter"

type JsonToonModel struct {
	ui.BasePagerModel
}

func NewJsonToonModel(common *ui.CommonModel) JsonToonModel {
	model := JsonToonModel{
		BasePagerModel: ui.NewBasePagerModel(common, Title),
	}

	return model
}

func (m JsonToonModel) Init() tea.Cmd {
	return m.BasePagerModel.Init()
}

func (m JsonToonModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			c := clipboard.New()
			content, err := c.PasteText()
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

func (m JsonToonModel) View() string {
	var b strings.Builder

	fmt.Fprint(&b, m.Viewport.View()+"\n")
	fmt.Fprint(&b, m.StatusBarView())

	if m.ShowHelp {
		fmt.Fprint(&b, "\n"+m.helpView())
	}

	return b.String()
}

func (m *JsonToonModel) SetContent(content string) error {
	m.Content = content

	toonStr, err := Convert(content)

	if err != nil {
		return fmt.Errorf("error converting JSON to TOON: %v", err)
	} else {
		m.FormattedContent = toonStr
	}

	var buf bytes.Buffer
	buf.WriteString(m.FormattedContent)
	m.Viewport.SetContent(buf.String())
	return nil
}

func (m JsonToonModel) helpView() (s string) {
	col1 := []string{
		"c              copy TOON",
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
	// Parse JSON using toon's ParseJSON which preserves key order
	data, err := toon.ParseJSON(strings.NewReader(jsonContent))
	if err != nil {
		// Try to provide more helpful error message
		var syntaxErr *json.SyntaxError
		if e, ok := err.(*json.SyntaxError); ok {
			syntaxErr = e
			return "", fmt.Errorf("JSON syntax error at offset %d: %v", syntaxErr.Offset, err)
		}
		return "", fmt.Errorf("JSON parsing error: %s", err.Error())
	}

	// Encode to TOON format with default options
	toonStr, err := toon.Encode(data)
	if err != nil {
		return "", fmt.Errorf("TOON encoding error: %s", err.Error())
	}

	return toonStr, nil
}
