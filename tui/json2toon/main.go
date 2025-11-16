package json2toon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hannes-sistemica/toon"
	"github.com/mattn/go-runewidth"
	"github.com/skatkov/devtui/internal/clipboard"
	"github.com/skatkov/devtui/internal/editor"
	"github.com/skatkov/devtui/internal/ui"
)

const Title = "JSON to TOON Converter"

type JsonToonModel struct {
	ui.BasePagerModel
	indent       int
	lengthMarker string
}

func NewJsonToonModel(common *ui.CommonModel) JsonToonModel {
	model := JsonToonModel{
		BasePagerModel: ui.NewBasePagerModel(common, Title),
		indent:         2,
		lengthMarker:   "",
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
		case "i":
			// Cycle through indent options: 2, 4
			if m.indent == 2 {
				m.indent = 4
			} else {
				m.indent = 2
			}
			if m.Content != "" {
				err := m.SetContent(m.Content)
				if err != nil {
					cmds = append(cmds, m.ShowErrorMessage(err.Error()))
				} else {
					cmds = append(cmds, m.ShowStatusMessage(fmt.Sprintf("Indent: %d spaces", m.indent)))
				}
			}
		case "l":
			// Toggle length marker: "" or "#"
			if m.lengthMarker == "" {
				m.lengthMarker = "#"
			} else {
				m.lengthMarker = ""
			}
			if m.Content != "" {
				err := m.SetContent(m.Content)
				if err != nil {
					cmds = append(cmds, m.ShowErrorMessage(err.Error()))
				} else {
					status := "off"
					if m.lengthMarker != "" {
						status = "on"
					}
					cmds = append(cmds, m.ShowStatusMessage("Length marker: "+status))
				}
			}
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

	opts := toon.EncodeOptions{
		Indent:       m.indent,
		Delimiter:    ",",
		LengthMarker: m.lengthMarker,
	}

	toonStr, err := ConvertWithOptions(content, opts)

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
	lengthStatus := "off"
	if m.lengthMarker != "" {
		lengthStatus = "on"
	}

	col1 := []string{
		"c              copy TOON",
		"e              edit JSON",
		"v              paste JSON to convert",
		fmt.Sprintf("i              toggle indent (current: %d)", m.indent),
		fmt.Sprintf("l              toggle length marker (current: %s)", lengthStatus),
		"q/ctrl+c       quit",
	}

	s += "\n"
	s += "k/↑      up                  " + col1[0] + "\n"
	s += "j/↓      down                " + col1[1] + "\n"
	s += "b/pgup   page up             " + col1[2] + "\n"
	s += "f/pgdn   page down           " + col1[3] + "\n"
	s += "u        ½ page up           " + col1[4] + "\n"
	s += "d        ½ page down         " + col1[5]

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
	opts := toon.EncodeOptions{
		Indent:       2,
		Delimiter:    ",",
		LengthMarker: "",
	}
	return ConvertWithOptions(jsonContent, opts)
}

func ConvertWithOptions(jsonContent string, opts toon.EncodeOptions) (string, error) {
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

	// Encode to TOON format with specified options
	toonStr, err := toon.EncodeWithOptions(data, opts)
	if err != nil {
		return "", fmt.Errorf("TOON encoding error: %s", err.Error())
	}

	return toonStr, nil
}
