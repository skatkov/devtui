package tsv2md

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
	"github.com/tiagomelo/go-clipboard/clipboard"
)

const Title = "TSV to Markdown Table Converter"

type TSV2MDModel struct {
	ui.BasePagerModel
	alignColumns bool
}

func NewTSV2MDModel(common *ui.CommonModel) TSV2MDModel {
	model := TSV2MDModel{
		BasePagerModel: ui.NewBasePagerModel(common, Title),
		alignColumns:   true,
	}

	return model
}

func (m TSV2MDModel) Init() tea.Cmd {
	return m.BasePagerModel.Init()
}

func (m TSV2MDModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			return m, editor.OpenEditor(m.Content, "tsv")
		case "v":
			c := clipboard.New()
			content, err := c.PasteText()
			if err == nil {
				err = m.setContent(content)
			}

			if err != nil {
				cmds = append(cmds, m.ShowErrorMessage(err.Error()))
			} else {
				cmds = append(cmds, m.ShowStatusMessage("Pasted. Press 'c' to copy result."))
			}
		case "a":
			m.alignColumns = !m.alignColumns
			if m.Content != "" {
				err := m.setContent(m.Content)
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
			err := m.setContent(msg.Content)

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

func (m TSV2MDModel) View() string {
	var b strings.Builder

	fmt.Fprint(&b, m.Viewport.View()+"\n")
	fmt.Fprint(&b, m.StatusBarView())

	if m.ShowHelp {
		fmt.Fprint(&b, "\n"+m.helpView())
	}

	return b.String()
}

func (m *TSV2MDModel) setContent(content string) error {
	m.Content = content

	reader := csv.NewReader(strings.NewReader(content))
	reader.Comma = '\t'

	// Read all records
	rows, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("error reading TSV: %v", err)
	}

	if len(rows) == 0 {
		return errors.New("empty result")
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

func (m TSV2MDModel) helpView() (s string) {
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
