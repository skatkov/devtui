package ui

import "github.com/charmbracelet/lipgloss"

type ReturnToListMsg struct {
	Common *CommonModel
}

type CommonModel struct {
	Width  int
	Height int

	LastSelectedItem int
}

var (
	PagePaddingStyle = lipgloss.NewStyle().Padding(2)
)
