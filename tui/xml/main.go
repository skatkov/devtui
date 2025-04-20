package xml

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/alecthomas/chroma/quick"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-xmlfmt/xmlfmt"
	"github.com/mattn/go-runewidth"
	"github.com/skatkov/devtui/internal/editor"
	"github.com/skatkov/devtui/internal/ui"
	"github.com/tiagomelo/go-clipboard/clipboard"
)

const Title = "XML Formatter"

type XMLFormatterModel struct {
	ui.BasePagerModel
}

func NewXMLFormatterModel(common *ui.CommonModel) XMLFormatterModel {
	model := XMLFormatterModel{
		BasePagerModel: ui.NewBasePagerModel(common, Title),
	}

	return model
}

func (m XMLFormatterModel) Init() tea.Cmd {
	return m.BasePagerModel.Init()
}

func (m XMLFormatterModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			return m, editor.OpenEditor(m.Content, "xml")
		case "v":
			c := clipboard.New()
			content, err := c.PasteText()
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

func (m XMLFormatterModel) View() string {
	var b strings.Builder

	fmt.Fprint(&b, m.Viewport.View()+"\n")
	fmt.Fprint(&b, m.StatusBarView())

	if m.ShowHelp {
		fmt.Fprint(&b, "\n"+m.helpView())
	}

	return b.String()
}

func (m *XMLFormatterModel) SetContent(content string) error {
	m.Content = content
	m.FormattedContent = xmlfmt.FormatXML(content, "\t", "  ")

	var buf bytes.Buffer
	err := quick.Highlight(&buf, m.FormattedContent, "xml", "terminal", "nord")
	if err != nil {
		return err
	}
	m.Viewport.SetContent(buf.String())

	return nil
}

func (m XMLFormatterModel) helpView() (s string) {
	col1 := []string{
		"c              copy formatted XML",
		"e              edit XML",
		"v              paste XML to format",
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
