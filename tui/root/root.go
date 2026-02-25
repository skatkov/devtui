package root

import (
	tea "charm.land/bubbletea/v2"
	"github.com/skatkov/devtui/internal/ui"
)

type RootModel struct {
	common      *ui.CommonModel
	currentView tea.Model
	listModel   *listModel
}

// RootScreen creates the root model using the default renderer.
// This is used for local terminal execution.
func RootScreen() RootModel {
	common := ui.CommonModel{LastSelectedItem: 0}
	common.Styles = ui.NewStyle()

	listModel := newListModel(&common)
	return RootModel{
		common:      &common,
		currentView: listModel,
		listModel:   listModel,
	}
}

// RootScreenWithSize creates the root model with an initial window size.
func RootScreenWithSize(width, height int) RootModel {
	common := ui.CommonModel{
		LastSelectedItem: 0,
		Width:            width,
		Height:           height,
	}
	common.Styles = ui.NewStyle()

	listModel := newListModel(&common)
	return RootModel{
		common:      &common,
		currentView: listModel,
		listModel:   listModel,
	}
}

// RootScreenWithRenderer preserves backward compatibility with old call sites.
func RootScreenWithRenderer(_ any, width, height int) RootModel {
	return RootScreenWithSize(width, height)
}

func (m RootModel) Init() tea.Cmd {
	return m.currentView.Init()
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ui.ReturnToListMsg:
		m.listModel.common = msg.Common
		m.listModel.list.SetSize(msg.Common.Width, msg.Common.Height)

		m.listModel.RefreshOrder()

		// Restore previously selected last
		m.listModel.list.Select(msg.Common.LastSelectedItem)

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

func (m RootModel) View() tea.View {
	return ui.WithAltScreen(m.currentView.View())
}
