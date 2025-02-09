package root

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type rootScreenModel struct {
	model tea.Model
}

func RootScreen() rootScreenModel {
	items := []list.Item{
		item("UUID Decode"),
		item("Number Base Converter"),
		item("UUID Generate"),
	}

	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "Choose your weapon!"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	m := listModel{list: l}
	return rootScreenModel{model: m}

}

func (m rootScreenModel) Init() tea.Cmd {
	return m.model.Init()
}

func (m rootScreenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.model.Update(msg)
}

func (m rootScreenModel) View() string {
	return m.model.View()
}
