package graphqlquery

import (
	"bytes"
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/alecthomas/chroma/quick"

	"github.com/skatkov/devtui/internal/clipboard"
	"github.com/skatkov/devtui/internal/editor"
	"github.com/skatkov/devtui/internal/ui"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/formatter"
	"github.com/vektah/gqlparser/v2/parser"
)

const Title = "GraphQL Query Formatter"

type GraphQLQueryModel struct {
	ui.BasePagerModel
}

func NewGraphQLQueryModel(common *ui.CommonModel) GraphQLQueryModel {
	model := GraphQLQueryModel{
		BasePagerModel: ui.NewBasePagerModel(common, Title),
	}

	return model
}

func (m GraphQLQueryModel) Init() tea.Cmd {
	return m.BasePagerModel.Init()
}

func (m GraphQLQueryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			return m, editor.OpenEditor(m.Content, "graphql")
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

func (m GraphQLQueryModel) View() tea.View {
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

func (m *GraphQLQueryModel) SetContent(content string) error {
	m.Content = content
	m.FormattedContent = formatGraphQL(content)

	// Check if GraphQL query is valid for syntax highlighting
	_, err := parser.ParseQuery(&ast.Source{Input: content})
	if err != nil {
		return err
	} else {
		var buf bytes.Buffer
		err = quick.Highlight(&buf, m.FormattedContent, "graphql", "terminal", "nord")
		if err != nil {
			return err
		} else {
			m.Viewport.SetContent(buf.String())
		}
	}

	return nil
}

func (m *GraphQLQueryModel) loadContent(content string, err error) tea.Cmd {
	if err == nil {
		err = m.SetContent(content)
	}
	if err != nil {
		return m.ShowErrorMessage(err.Error())
	}
	return m.ShowStatusMessage("Pasted. Press 'c' to copy result.")
}

func (m GraphQLQueryModel) helpView() string {
	col1 := []string{
		"c              copy formatted GraphQL",
		"e              edit unformatted GraphQL",
		"v              paste unformatted GraphQL",
		"q/ctrl+c       quit",
	}

	return m.FormatHelpColumns(col1)
}

func formatGraphQL(content string) string {
	// Parse the query
	query, err := parser.ParseQuery(&ast.Source{Input: content})
	if err != nil {
		return content // Return unmodified content if parsing fails
	}

	// Format the query
	var buf bytes.Buffer
	f := formatter.NewFormatter(&buf)
	f.FormatQueryDocument(query)

	return buf.String()
}
