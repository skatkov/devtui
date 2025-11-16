package jsonrepair

import (
	"bytes"
	"fmt"
	"strings"

	jsonrepair "github.com/RealAlexandreAI/json-repair"
	"github.com/alecthomas/chroma/quick"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mattn/go-runewidth"

	"github.com/skatkov/devtui/internal/clipboard"
	"github.com/skatkov/devtui/internal/editor"
	"github.com/skatkov/devtui/internal/ui"
)

const Title = "JSON Repair"

type JSONRepairModel struct {
	ui.BasePagerModel
}

func NewJSONRepairModel(common *ui.CommonModel) JSONRepairModel {
	model := JSONRepairModel{
		BasePagerModel: ui.NewBasePagerModel(common, Title),
	}

	return model
}

func (m JSONRepairModel) Init() tea.Cmd {
	return m.BasePagerModel.Init()
}

func (m JSONRepairModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
					cmds = append(cmds, m.ShowStatusMessage("Repaired. Press 'c' to copy result."))
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
				cmds = append(cmds, m.ShowStatusMessage("Repaired. Press 'c' to copy result."))
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

func (m JSONRepairModel) View() string {
	var b strings.Builder

	fmt.Fprint(&b, m.Viewport.View()+"\n")
	fmt.Fprint(&b, m.StatusBarView())

	if m.ShowHelp {
		fmt.Fprint(&b, "\n"+m.helpView())
	}

	return b.String()
}

func (m *JSONRepairModel) SetContent(content string) error {
	m.Content = content

	// Repair the JSON
	repairedJSON, err := jsonrepair.RepairJSON(content)
	if err != nil {
		return fmt.Errorf("failed to repair JSON: %w", err)
	}

	m.FormattedContent = repairedJSON

	// Syntax highlight the repaired JSON
	var buf bytes.Buffer
	_ = quick.Highlight(&buf, m.FormattedContent, "json", "terminal", "nord")
	m.Viewport.SetContent(buf.String())
	return nil
}

func (m JSONRepairModel) helpView() (s string) {
	col1 := []string{
		"c              copy repaired JSON",
		"e              edit broken JSON",
		"v              paste broken JSON to repair",
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

// RepairJSON repairs malformed JSON string and returns the repaired version.
// This function is used by the CLI command.
func RepairJSON(content string) (string, error) {
	repairedJSON, err := jsonrepair.RepairJSON(content)
	if err != nil {
		return "", fmt.Errorf("failed to repair JSON: %w", err)
	}
	return repairedJSON, nil
}
