// devtui/tui/csv2md/main.go
package csv2md

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mattn/go-runewidth"
	"github.com/skatkov/devtui/internal/csv2md"
	"github.com/skatkov/devtui/internal/editor"
	"github.com/skatkov/devtui/internal/ui"
	"github.com/skatkov/devtui/internal/clipboard"
)

const Title = "CSV to Markdown Table Converter"

type CSV2MDModel struct {
	ui.BasePagerModel
	alignColumns bool
}

func NewCSV2MDModel(common *ui.CommonModel) CSV2MDModel {
	model := CSV2MDModel{
		BasePagerModel: ui.NewBasePagerModel(common, Title),
		alignColumns:   true, // default to aligned columns
	}

	return model
}

func (m CSV2MDModel) Init() tea.Cmd {
	return m.BasePagerModel.Init()
}

func (m CSV2MDModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			return m, editor.OpenEditor(m.Content, "csv")
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
		case "a":
			m.alignColumns = !m.alignColumns
			if m.Content != "" {
				err := m.SetContent(m.Content)
				if err != nil {
					cmds = append(cmds, m.ShowErrorMessage(err.Error()))
				} else if m.alignColumns {
					cmds = append(cmds, m.ShowStatusMessage("Columns aligned"))
				} else {
					cmds = append(cmds, m.ShowStatusMessage("Columns unaligned"))
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
				cmds = append(cmds, m.ShowStatusMessage("Converted. Press 'c' to copy result."))
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

func (m CSV2MDModel) View() string {
	var b strings.Builder

	fmt.Fprint(&b, m.Viewport.View()+"\n")
	fmt.Fprint(&b, m.StatusBarView())

	if m.ShowHelp {
		fmt.Fprint(&b, "\n"+m.helpView())
	}

	return b.String()
}

func (m *CSV2MDModel) SetContent(content string) error {
	m.Content = content

	reader := csv.NewReader(strings.NewReader(content))

	// Read all records
	rows, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("error reading content: %v", err)
	}

	if len(rows) == 0 {
		return errors.New("empty content")
	}

	// Convert to markdown
	markdownLines := csv2md.Convert("", rows, m.alignColumns)
	m.FormattedContent = strings.Join(markdownLines, "\n")

	// Set content in viewport
	var buf bytes.Buffer
	buf.WriteString(m.FormattedContent)
	m.Viewport.SetContent(buf.String())

	return nil
}

func (m CSV2MDModel) helpView() (s string) {
	col1 := []string{
		"c              copy markdown",
		"e              edit",
		"v              paste to convert",
		"a              toggle column alignment",
		"q/ctrl+c       quit",
	}

	s += "\n"
	s += "k/↑      up                  " + col1[0] + "\n"
	s += "j/↓      down                " + col1[1] + "\n"
	s += "b/pgup   page up             " + col1[2] + "\n"
	s += "f/pgdn   page down           " + col1[3] + "\n"
	s += "u        ½ page up           " + col1[4] + "\n"
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
