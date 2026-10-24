package tsv2md

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/skatkov/devtui/internal/clipboard"
	"github.com/skatkov/devtui/internal/csv2md"
	"github.com/skatkov/devtui/internal/editor"
	"github.com/skatkov/devtui/internal/ui"
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
	case tea.KeyPressMsg:
		if cmd, handled := m.HandleCommonKeys(msg); handled {
			return m, cmd
		}
		switch msg.String() {
		case "e":
			return m, editor.OpenEditor(m.Content, "tsv")
		case "v":
			content, err := clipboard.Paste()
			cmds = append(cmds, m.applyContent(content, err))
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
func (m *TSV2MDModel) applyContent(content string, err error) tea.Cmd {
	if err == nil {
		err = m.SetContent(content)
	}
	if err != nil {
		return m.ShowErrorMessage(err.Error())
	}
	return m.ShowStatusMessage("Pasted. Press 'c' to copy result.")
}

func (m TSV2MDModel) View() tea.View {
	var b strings.Builder

	fmt.Fprint(&b, m.Viewport.View()+"\n")
	fmt.Fprint(&b, m.StatusBarView())

	if m.ShowHelp {
		fmt.Fprint(&b, "\n"+m.helpView())
	}

	return m.NewView(b.String())
}

func (m *TSV2MDModel) SetContent(content string) error {
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

func (m TSV2MDModel) helpView() string {
	col1 := []string{
		"c              copy markdown",
		"e              edit",
		"v              paste to convert",
		"a              toggle column alignment",
		"q/ctrl+c       quit",
	}

	return m.FormatHelpColumns(col1)
}
