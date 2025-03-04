package ui

import (
	"time"

	"github.com/charmbracelet/lipgloss"
)

type ReturnToListMsg struct {
	Common *CommonModel
}

type CommonModel struct {
	Width  int
	Height int

	LastSelectedItem int
}

var PagePaddingStyle = lipgloss.NewStyle().Padding(2)

const StatusMessageTimeout = time.Second * 3 // how long to show status messages like "stashed!"
