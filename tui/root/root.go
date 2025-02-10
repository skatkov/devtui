package root

import (
	tea "github.com/charmbracelet/bubbletea"
)

type rootScreenModel struct {
	model tea.Model
}

func RootScreen() rootScreenModel {
	m := newListModel()
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
