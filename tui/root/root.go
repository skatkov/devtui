package root

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	common.Lg = lipgloss.DefaultRenderer()
	common.Styles = ui.NewStyle(common.Lg)

	listModel := newListModel(&common)
	return RootModel{
		common:      &common,
		currentView: listModel,
		listModel:   listModel,
	}
}

// RootScreenWithRenderer creates the root model using a custom renderer.
// This is used for SSH sessions where the renderer must be session-aware
// to properly detect the client's terminal capabilities.
func RootScreenWithRenderer(lg *lipgloss.Renderer, width, height int) RootModel {
	common := ui.CommonModel{
		LastSelectedItem: 0,
		Width:            width,
		Height:           height,
	}
	common.Lg = lg
	common.Styles = ui.NewStyle(common.Lg)

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

func (m RootModel) View() string {
	return m.currentView.View()
}
