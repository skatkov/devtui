// devtui/tui/base64-encoder/main.go
package base64encoder

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"

	"github.com/skatkov/devtui/internal/base64"
	"github.com/skatkov/devtui/internal/clipboard"
	"github.com/skatkov/devtui/internal/editor"
	"github.com/skatkov/devtui/internal/ui"
)

const Title = "Base64 Encoder"

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
			return m, editor.OpenEditor(m.Content, "txt")
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
				cmds = append(cmds, m.ShowStatusMessage("Content encoded. Press 'c' to copy result."))
			}
		}
	case tea.WindowSizeMsg:
		// Handle window size change with custom logic for re-wrapping
		m.Common.Width = msg.Width
		m.Common.Height = msg.Height

		m.SetSize(msg.Width, msg.Height)

		if !m.Ready {
			m.Viewport = viewport.New(
				viewport.WithWidth(msg.Width),
				viewport.WithHeight(msg.Height-ui.StatusBarHeight),
			)
			m.Ready = true
		} else {
			m.SetSize(msg.Width, msg.Height)
		}

		// Re-wrap content with new width
		err := m.SetContent(m.Content)
		if err != nil {
			cmds = append(cmds, m.ShowErrorMessage(err.Error()))
		}
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

	// Encode to base64
	result := base64.EncodeString(content)

	// Wrap the base64 output to fit the screen width
	// Use common width minus some padding for readability
	// Minimum width of 20 to ensure readability on very narrow screens
	wrapWidth := max(m.Common.Width-4, 20)

	// Hard wrap the base64 string at character boundaries
	var lines []string
	rest := result
	for len(rest) > wrapWidth {
		lines = append(lines, rest[:wrapWidth])
		rest = rest[wrapWidth:]
	}
	lines = append(lines, rest)

	m.FormattedContent = result // Keep original unwrapped for copying

	m.Viewport.SetContent(strings.Join(lines, "\n"))

	return nil
}

func (m *Base64Model) paste(content string, err error) tea.Cmd {
	if err == nil {
		err = m.SetContent(content)
	}

	if err != nil {
		return m.ShowErrorMessage(err.Error())
	}
	return m.ShowStatusMessage("Pasted and encoded. Press 'c' to copy result.")
}

func (m Base64Model) helpView() string {
	col1 := []string{
		"c              copy base64",
		"e              edit text",
		"v              paste text to encode",
		"q/ctrl+c       quit",
	}

	return m.FormatHelpColumns(col1)
}
