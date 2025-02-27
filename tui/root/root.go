package root

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/skatkov/devtui/internal/ui"
)

type RootModel struct {
	currentView tea.Model
	listModel   *listModel
}

func RootScreen() RootModel {
	listModel := newListModel()
	return RootModel{
		currentView: listModel,
		listModel:   listModel,
	}
}

func (m RootModel) Init() tea.Cmd {
	return m.currentView.Init()
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case ui.ReturnToListMsg:
		m.currentView = m.listModel
		return m, m.listModel.Init()
	}
	var cmd tea.Cmd
	m.currentView, cmd = m.currentView.Update(msg)

	return m, cmd
}

func (m RootModel) View() string {
	return m.currentView.View()
}
