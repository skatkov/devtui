package root

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/skatkov/devtui/internal/ui"
)

type RootModel struct {
	common      *ui.CommonModel
	currentView tea.Model
	listModel   *listModel
}

func RootScreen() RootModel {
	common := ui.CommonModel{}
	listModel := newListModel(&common)
	return RootModel{
		common:      &common,
		currentView: listModel,
		listModel:   listModel,
	}
}

func (m RootModel) Init() tea.Cmd {
	return m.currentView.Init()
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ui.ReturnToListMsg:
		m.listModel.common = msg.Common
		m.listModel.list.SetSize(msg.Common.Width, msg.Common.Height)
		m.currentView = m.listModel
		return m, m.listModel.Init()
	// Window size is received when starting up and on every resize
	case tea.WindowSizeMsg:
		m.common.Width = msg.Width
		m.common.Height = msg.Height
	}
	var cmd tea.Cmd
	m.currentView, cmd = m.currentView.Update(msg)

	return m, cmd
}

func (m RootModel) View() string {
	return m.currentView.View()
}
