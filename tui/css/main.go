package css

import (
	"bytes"
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/alecthomas/chroma/quick"
	"github.com/client9/csstool"

	"github.com/skatkov/devtui/internal/clipboard"
	"github.com/skatkov/devtui/internal/editor"
	"github.com/skatkov/devtui/internal/ui"
)

const Title = "CSS Formatter"

type CSSFormatterModel struct {
	ui.BasePagerModel
}

func NewCSSFormatterModel(common *ui.CommonModel) CSSFormatterModel {
	model := CSSFormatterModel{
		BasePagerModel: ui.NewBasePagerModel(common, Title),
	}

	return model
}

func (m CSSFormatterModel) Init() tea.Cmd {
	return m.BasePagerModel.Init()
}

func (m CSSFormatterModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			return m, editor.OpenEditor(m.Content, "css")
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

func (m CSSFormatterModel) View() tea.View {
	var b strings.Builder

	fmt.Fprint(&b, m.Viewport.View()+"\n")
	fmt.Fprint(&b, m.StatusBarView())

	if m.ShowHelp {
		fmt.Fprint(&b, "\n"+m.helpView())
	}

	v := m.NewView(b.String())
	v.MouseMode = tea.MouseModeCellMotion
	return v
}

func (m *CSSFormatterModel) SetContent(content string) error {
	m.Content = content

	reader := strings.NewReader(content)
	var outputBuf bytes.Buffer

	// Create CSS formatter with 2 spaces for indentation and always including semicolons
	cssformat := csstool.NewCSSFormat(2, false, nil)
	cssformat.AlwaysSemicolon = true

	// Format the CSS
	err := cssformat.Format(reader, &outputBuf)
	if err != nil {
		return err
	} else {
		m.FormattedContent = outputBuf.String()
	}

	// Syntax highlight the formatted CSS
	var highlightBuf bytes.Buffer
	err = quick.Highlight(&highlightBuf, m.FormattedContent, "css", "terminal", "nord")
	if err != nil {
		return err
	}
	m.Viewport.SetContent(highlightBuf.String())
	return nil
}

func (m *CSSFormatterModel) paste(content string, err error) tea.Cmd {
	if err == nil {
		err = m.SetContent(content)
	}

	if err != nil {
		return m.ShowErrorMessage(err.Error())
	}
	return m.ShowStatusMessage("Pasted. Press 'c' to copy result.")
}

func (m CSSFormatterModel) helpView() string {
	col1 := []string{
		"c              copy formatted CSS",
		"e              edit CSS",
		"v              paste CSS to format",
		"q/ctrl+c       quit",
	}

	return m.FormatHelpColumns(col1)
}
