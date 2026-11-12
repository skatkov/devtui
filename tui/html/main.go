package html

import (
	"bytes"
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/alecthomas/chroma/quick"

	"github.com/skatkov/devtui/internal/clipboard"
	"github.com/skatkov/devtui/internal/editor"
	"github.com/skatkov/devtui/internal/htmlfmt"
	"github.com/skatkov/devtui/internal/ui"
)

const Title = "HTML Formatter"

type HTMLFormatterModel struct {
	ui.BasePagerModel
}

func NewHTMLFormatterModel(common *ui.CommonModel) HTMLFormatterModel {
	model := HTMLFormatterModel{
		BasePagerModel: ui.NewBasePagerModel(common, Title),
	}

	return model
}

func (m HTMLFormatterModel) Init() tea.Cmd {
	return m.BasePagerModel.Init()
}

func (m HTMLFormatterModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			return m, editor.OpenEditor(m.Content, "html")
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
func (m *HTMLFormatterModel) loadContent(content string, err error) tea.Cmd {
	if err == nil {
		err = m.SetContent(content)
	}
	if err != nil {
		return m.ShowErrorMessage(err.Error())
	}
	return m.ShowStatusMessage("Pasted. Press 'c' to copy result.")
}

func (m HTMLFormatterModel) View() tea.View {
	var b strings.Builder

	fmt.Fprint(&b, m.Viewport.View()+"\n")
	fmt.Fprint(&b, m.StatusBarView())

	if m.ShowHelp {
		fmt.Fprint(&b, "\n"+m.helpView())
	}

	return m.NewView(b.String())
}

func (m *HTMLFormatterModel) SetContent(content string) error {
	m.Content = content
	m.FormattedContent = htmlfmt.Format(content)

	var buf bytes.Buffer
	err := quick.Highlight(&buf, m.FormattedContent, "html", "terminal", "nord")
	if err != nil {
		return err
	}
	m.Viewport.SetContent(buf.String())
	return nil
}

func (m HTMLFormatterModel) helpView() string {
	col1 := []string{
		"c              copy formatted HTML",
		"e              edit HTML",
		"v              paste HTML to format",
		"q/ctrl+c       quit",
	}

	return m.FormatHelpColumns(col1)
}
