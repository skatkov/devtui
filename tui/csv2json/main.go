package csv2json

import (
	"bytes"
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/alecthomas/chroma/quick"

	"github.com/skatkov/devtui/internal/clipboard"
	csv2jsonconv "github.com/skatkov/devtui/internal/csv2json"
	"github.com/skatkov/devtui/internal/editor"
	"github.com/skatkov/devtui/internal/ui"
)

const Title = "CSV to JSON Converter"

type CSVJsonModel struct {
	ui.BasePagerModel
}

func NewCSVJsonModel(common *ui.CommonModel) CSVJsonModel {
	model := CSVJsonModel{
		BasePagerModel: ui.NewBasePagerModel(common, Title),
	}

	return model
}

func (m CSVJsonModel) Init() tea.Cmd {
	return m.BasePagerModel.Init()
}

func (m CSVJsonModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			return m, editor.OpenEditor(m.Content, "csv")
		case "v":
			content, err := clipboard.Paste()
			cmds = append(cmds, m.paste(content, err))
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

func (m CSVJsonModel) View() tea.View {
	var b strings.Builder

	fmt.Fprint(&b, m.Viewport.View()+"\n")
	fmt.Fprint(&b, m.StatusBarView())

	if m.ShowHelp {
		fmt.Fprint(&b, "\n"+m.helpView())
	}

	return m.NewView(b.String())
}

func (m *CSVJsonModel) SetContent(content string) error {
	m.Content = content

	jsonStr, err := csv2jsonconv.Convert(content)
	if err != nil {
		return err
	}
	m.FormattedContent = jsonStr

	var buf bytes.Buffer
	_ = quick.Highlight(&buf, m.FormattedContent, "json", "terminal", "nord")
	m.Viewport.SetContent(buf.String())
	return nil
}

func (m *CSVJsonModel) paste(content string, err error) tea.Cmd {
	if err == nil {
		err = m.SetContent(content)
	}

	if err != nil {
		return m.ShowErrorMessage(err.Error())
	}
	return m.ShowStatusMessage("Pasted. Press 'c' to copy result.")
}

func (m CSVJsonModel) helpView() string {
	col1 := []string{
		"c              copy JSON",
		"e              edit CSV",
		"v              paste CSV to convert",
		"q/ctrl+c       quit",
	}

	return m.FormatHelpColumns(col1)
}
