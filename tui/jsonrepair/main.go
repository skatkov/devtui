package jsonrepair

import (
	"bytes"
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	jsonrepair "github.com/RealAlexandreAI/json-repair"
	"github.com/alecthomas/chroma/quick"

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
	case tea.KeyPressMsg:
		if cmd, handled := m.HandleCommonKeys(msg); handled {
			return m, cmd
		}
		switch msg.String() {
		case "e":
			return m, editor.OpenEditor(m.Content, "json")
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
func (m *JSONRepairModel) loadContent(content string, err error) tea.Cmd {
	if err == nil {
		err = m.SetContent(content)
	}
	if err != nil {
		return m.ShowErrorMessage(err.Error())
	}
	return m.ShowStatusMessage("Repaired. Press 'c' to copy result.")
}

func (m JSONRepairModel) View() tea.View {
	var b strings.Builder

	fmt.Fprint(&b, m.Viewport.View()+"\n")
	fmt.Fprint(&b, m.StatusBarView())

	if m.ShowHelp {
		fmt.Fprint(&b, "\n"+m.helpView())
	}

	return m.NewView(b.String())
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

func (m JSONRepairModel) helpView() string {
	col1 := []string{
		"c              copy repaired JSON",
		"e              edit broken JSON",
		"v              paste broken JSON to repair",
		"q/ctrl+c       quit",
	}

	return m.FormatHelpColumns(col1)
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
