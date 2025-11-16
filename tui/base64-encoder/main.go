// devtui/tui/base64-encoder/main.go
package base64encoder

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mattn/go-runewidth"

	"github.com/skatkov/devtui/internal/base64"
	"github.com/skatkov/devtui/internal/editor"
	"github.com/skatkov/devtui/internal/ui"
	"github.com/skatkov/devtui/internal/clipboard"
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
	case tea.KeyMsg:
		if cmd, handled := m.HandleCommonKeys(msg); handled {
			return m, cmd
		}

		switch msg.String() {
		case "e":
			return m, editor.OpenEditor(m.Content, "txt")
		case "v":
			content, err := clipboard.Paste()
			if err == nil {
				err = m.SetContent(content)
			}

			if err != nil {
				cmds = append(cmds, m.ShowErrorMessage(err.Error()))
			} else {
				cmds = append(cmds, m.ShowStatusMessage("Pasted and encoded. Press 'c' to copy result."))
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
				cmds = append(cmds, m.ShowStatusMessage("Content encoded. Press 'c' to copy result."))
			}
		}
	case tea.WindowSizeMsg:
		// Handle window size change with custom logic for re-wrapping
		m.Common.Width = msg.Width
		m.Common.Height = msg.Height

		m.SetSize(msg.Width, msg.Height)

		if !m.Ready {
			m.Viewport = viewport.New(msg.Width, msg.Height-ui.StatusBarHeight)
			m.Viewport.YPosition = 0
			m.Ready = true
		} else {
			m.SetSize(msg.Width, msg.Height)
		}

		// Re-wrap content with new width if we have content
		if m.Content != "" {
			err := m.SetContent(m.Content)
			if err != nil {
				cmds = append(cmds, m.ShowErrorMessage(err.Error()))
			}
		}
	}

	m.Viewport, cmd = m.Viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Base64Model) View() string {
	var b strings.Builder

	fmt.Fprint(&b, m.Viewport.View()+"\n")
	fmt.Fprint(&b, m.StatusBarView())

	if m.ShowHelp {
		fmt.Fprint(&b, "\n"+m.helpView())
	}

	return b.String()
}

func (m *Base64Model) SetContent(content string) error {
	m.Content = content

	if content == "" {
		m.FormattedContent = ""
		m.Viewport.SetContent("")
		return nil
	}

	// Encode to base64
	result := base64.EncodeString(content)

	// Wrap the base64 output to fit the screen width
	// Use common width minus some padding for readability
	wrapWidth := m.Common.Width - 4
	if wrapWidth < 20 {
		wrapWidth = 20 // Minimum width to ensure readability on very narrow screens
	}

	// Hard wrap the base64 string at character boundaries
	var wrappedResult strings.Builder
	for i, char := range result {
		if i > 0 && i%wrapWidth == 0 {
			wrappedResult.WriteRune('\n')
		}
		wrappedResult.WriteRune(char)
	}

	m.FormattedContent = result // Keep original unwrapped for copying

	// Set content in viewport
	var buf bytes.Buffer
	buf.WriteString(wrappedResult.String())
	m.Viewport.SetContent(buf.String())

	return nil
}

func (m Base64Model) helpView() (s string) {
	col1 := []string{
		"c              copy base64",
		"e              edit text",
		"v              paste text to encode",
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
