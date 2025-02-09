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
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			screenTwo := screenTwo{}
			m.model = screenTwo
			return screenTwo, screenTwo.Init()
		}

	}
	return m.model.Update(msg)
}

func (m rootScreenModel) View() string {
	return m.model.View()
}

type screenTwo struct {
}

func (m screenTwo) Init() tea.Cmd {
	return nil
}

func (m screenTwo) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m screenTwo) View() string {
	return "Second view"
}
