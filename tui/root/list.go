package root

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/skatkov/devtui/tui/numbers"
	uuiddecode "github.com/skatkov/devtui/tui/uuid-decode"
	uuidgenerate "github.com/skatkov/devtui/tui/uuid-generate"
)

const listHeight = 15

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type item string

func (i item) FilterValue() string { return string(i) }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

func newListModel() *listModel {
	items := []list.Item{
		item("UUID Decode"),
		item("Number Base Converter"),
		item("UUID Generate"),
	}

	delegate := itemDelegate{}
	l := list.New(items, delegate, 20, listHeight)
	l.Title = "Choose your weapon!"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	return &listModel{
		list: l,
	}
}

type listModel struct {
	list list.Model
	err  string
}

func (m listModel) Init() tea.Cmd {
	return nil
}

func (m listModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				switch string(i) {
				case "Number Base Converter":
					newScreen := numbers.NewNumberModel()
					return newScreen, newScreen.Init()
				case "UUID Decode":
					newScreen := uuiddecode.NewUUIDDecodeModel()
					return newScreen, newScreen.Init()
				case "UUID Generate":
					newScreen := uuidgenerate.NewUUIDGenerateModel()
					return newScreen, newScreen.Init()
				default:
					m.err = fmt.Sprintf("%s app is not available", string(i))
				}
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
