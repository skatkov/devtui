package root

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/skatkov/devtui/internal/ui"
	cron "github.com/skatkov/devtui/tui/cron"
	"github.com/skatkov/devtui/tui/json"
	"github.com/skatkov/devtui/tui/markdown"
	"github.com/skatkov/devtui/tui/numbers"
	uuiddecode "github.com/skatkov/devtui/tui/uuid-decode"
	uuidgenerate "github.com/skatkov/devtui/tui/uuid-generate"
)

const listHeight = 15

type listModel struct {
	list   list.Model
	err    string
	common *ui.CommonModel
}

var (
	titleStyle        = list.DefaultStyles().Title.MarginTop(1)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(2).PaddingBottom(1)
)

type MenuOption struct {
	title string
	model func() tea.Model
}

func (i MenuOption) FilterValue() string { return i.title }
func (i MenuOption) Title() string       { return i.title }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(MenuOption)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i.title)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

func newListModel(common *ui.CommonModel) *listModel {
	items := []list.Item{
		MenuOption{
			title: "UUID Decoder",
			model: func() tea.Model { return uuiddecode.NewUUIDDecodeModel(common) },
		},
		MenuOption{
			title: "Number Base Converter",
			model: func() tea.Model { return numbers.NewNumberModel(common) },
		},
		MenuOption{
			title: "UUID Generate",
			model: func() tea.Model { return uuidgenerate.NewUUIDGenerateModel(common) },
		},
		MenuOption{
			title: "Cron Job Parser",
			model: func() tea.Model { return cron.NewCronModel(common) },
		},
		MenuOption{
			title: "JSON Formatter",
			model: func() tea.Model { return json.NewJsonModel(common) },
		},
		MenuOption{
			title: "Markdown Renderer",
			model: func() tea.Model { return markdown.NewMarkdownModel(common) },
		},
	}

	delegate := itemDelegate{}
	l := list.New(items, delegate, 20, listHeight)
	l.Title = "DevTUI"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	return &listModel{
		list:   l,
		common: common,
	}
}

func (m listModel) Init() tea.Cmd {
	return nil
}

func (m listModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.common.Width = msg.Width
		m.common.Height = msg.Height
		m.list.SetSize(msg.Width, msg.Height)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "enter":
			m.common.LastSelectedItem = m.list.Index()
			i, ok := m.list.SelectedItem().(MenuOption)
			if ok {
				newScreen := i.model()
				return newScreen, newScreen.Init()
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m listModel) View() string {
	if m.err != "" {
		return lipgloss.NewStyle().Padding(2).Render(m.err)
	}
	return "\n" + m.list.View()
}
