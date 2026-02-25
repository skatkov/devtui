package json

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/alecthomas/chroma/quick"
	"github.com/mattn/go-runewidth"
	"github.com/skatkov/devtui/internal/clipboard"
	"github.com/skatkov/devtui/internal/editor"
	"github.com/skatkov/devtui/internal/ui"
)

const Title = "JSON Formatter"

type JsonModel struct {
	ui.BasePagerModel
}

func NewJsonModel(common *ui.CommonModel) JsonModel {
	model := JsonModel{
		BasePagerModel: ui.NewBasePagerModel(common, Title),
	}

	return model
}

func (m JsonModel) Init() tea.Cmd {
	return m.BasePagerModel.Init()
}

func (m JsonModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			content, err := clipboard.Paste()
			if err != nil {
				cmds = append(cmds, m.ShowErrorMessage(err.Error()))
			} else {
				err := m.SetContent(content)

				if err != nil {
					cmds = append(cmds, m.ShowErrorMessage(err.Error()))
				} else {
					cmds = append(cmds, m.ShowStatusMessage("Pasted. Press 'c' to copy"))
				}
			}
		}
	case ui.StatusMessageTimeoutMsg:
		m.State = ui.PagerStateBrowse

	case editor.EditorFinishedMsg:
		if msg.Err != nil {
			cmds = append(cmds, m.ShowErrorMessage(msg.Err.Error()))
		} else {
			err := m.SetContent(msg.Content)

			if err != nil {
				cmds = append(cmds, m.ShowErrorMessage(err.Error()))
			} else {
				cmds = append(cmds, m.ShowStatusMessage("Pasted. Press 'c' to copy"))
			}
		}
	case tea.WindowSizeMsg:
		cmd = m.HandleWindowSizeMsg(msg)
		cmds = append(cmds, cmd)
	}
	// Handle keyboard and mouse events in the viewport
	m.Viewport, cmd = m.Viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m JsonModel) View() tea.View {
	var b strings.Builder

	fmt.Fprint(&b, m.Viewport.View()+"\n")
	fmt.Fprint(&b, m.StatusBarView())

	if m.ShowHelp {
		fmt.Fprint(&b, "\n"+m.helpView())
	}

	return m.NewView(b.String())
}

func (m *JsonModel) SetContent(content string) error {
	m.Content = content
	m.FormattedContent = FormatJSON(content)

	var buf bytes.Buffer
	err := quick.Highlight(&buf, m.FormattedContent, "json", "terminal", "nord")
	if err != nil {
		return err
	}
	m.Viewport.SetContent(buf.String())

	return nil
}

func (m JsonModel) helpView() (s string) {
	col1 := []string{
		"c              copy formatted JSON",
		"e              edit unformatted JSON",
		"v              paste unformatted JSON",
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

	// Fill up empty cells with spaces for background coloring
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

func FormatJSON(content string) string {
	var data any
	if err := json.Unmarshal([]byte(content), &data); err != nil {
		return content
	}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	if err := enc.Encode(data); err != nil {
		return content
	}
	return buf.String()
}
