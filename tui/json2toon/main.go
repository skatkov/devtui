package json2toon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/hannes-sistemica/toon"
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
	case tea.KeyPressMsg:
		if cmd, handled := m.HandleCommonKeys(msg); handled {
			return m, cmd
		}
		switch msg.String() {
		case "e":
			return m, editor.OpenEditor(m.Content, "json")
		case "i":
			// Toggle indent between 2 and 4
			m.indent = 6 - m.indent
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
func (m *JsonToonModel) loadContent(content string, err error) tea.Cmd {
	if err == nil {
		err = m.SetContent(content)
	}
	if err != nil {
		return m.ShowErrorMessage(err.Error())
	}
	return m.ShowStatusMessage("Pasted. Press 'c' to copy result.")
}

func (m JsonToonModel) View() tea.View {
	var b strings.Builder

	fmt.Fprint(&b, m.Viewport.View()+"\n")
	fmt.Fprint(&b, m.StatusBarView())

	if m.ShowHelp {
		fmt.Fprint(&b, "\n"+m.helpView())
	}

	return m.NewView(b.String())
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

func (m JsonToonModel) helpView() string {
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

	return m.FormatHelpColumns(col1)
}

func Convert(jsonContent string) (string, error) {
	// toon fills in its defaults: two-space indent and comma delimiter
	return ConvertWithOptions(jsonContent, toon.EncodeOptions{})
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
