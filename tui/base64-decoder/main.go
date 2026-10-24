// devtui/tui/base64-decoder/main.go
package base64decoder

import (
	"bytes"
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/skatkov/devtui/internal/base64"
	"github.com/skatkov/devtui/internal/clipboard"
	"github.com/skatkov/devtui/internal/editor"
	"github.com/skatkov/devtui/internal/ui"
)

const Title = "Base64 Decoder"

type Base64Model struct {
	ui.BasePagerModel
}

func NewBase64Model(common *ui.CommonModel) Base64Model {
	model := Base64Model{
		BasePagerModel: ui.NewBasePagerModel(common, Title),
	}

	return model
}

func (m Base64Model) Init() tea.Cmd {
	return m.BasePagerModel.Init()
}

func (m Base64Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			return m, editor.OpenEditor(m.Content, "base64")
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
				cmds = append(cmds, m.ShowStatusMessage("Content decoded. Press 'c' to copy result."))
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

func (m Base64Model) View() tea.View {
	var b strings.Builder

	fmt.Fprint(&b, m.Viewport.View()+"\n")
	fmt.Fprint(&b, m.StatusBarView())

	if m.ShowHelp {
		fmt.Fprint(&b, "\n"+m.helpView())
	}

	return m.NewView(b.String())
}

func (m *Base64Model) SetContent(content string) error {
	m.Content = content

	// Decode from base64
	result, err := base64.DecodeToString(content)
	if err != nil {
		return err
	}

	m.FormattedContent = result

	// Set content in viewport
	var buf bytes.Buffer
	buf.WriteString(m.FormattedContent)
	m.Viewport.SetContent(buf.String())

	return nil
}

func (m *Base64Model) paste(content string, err error) tea.Cmd {
	if err == nil {
		err = m.SetContent(content)
	}

	if err != nil {
		return m.ShowErrorMessage(err.Error())
	}
	return m.ShowStatusMessage("Pasted and decoded. Press 'c' to copy result.")
}

func (m Base64Model) helpView() string {
	col1 := []string{
		"c              copy text",
		"e              edit base64",
		"v              paste base64 to decode",
		"q/ctrl+c       quit",
	}

	return m.FormatHelpColumns(col1)
}
