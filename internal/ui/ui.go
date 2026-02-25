package ui

import (
	"fmt"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

const AppTitle = "DevTUI"

type ReturnToListMsg struct {
	Common *CommonModel
}

type CommonModel struct {
	Width  int
	Height int

	LastSelectedItem int

	Styles *Styles
}

var PagePaddingStyle = lipgloss.NewStyle().Padding(2)

const StatusMessageTimeout = time.Second * 3 // how long to show status messages like "stashed!"

const (
	StatusBarHeight = 1
	Ellipsis        = "…"
)

type PagerState int

const (
	PagerStateBrowse PagerState = iota
	PagerStateStatusMessage
	PagerStateErrorMessage
)

type PagerStatusMsg struct {
	Message string
}

type StatusMessageTimeoutMsg struct{}

func WaitForStatusMessageTimeout(t *time.Timer) tea.Cmd {
	return func() tea.Msg {
		<-t.C
		return StatusMessageTimeoutMsg{}
	}
}

func AltScreenView(content string) tea.View {
	v := tea.NewView(content)
	v.AltScreen = true
	return v
}

func WithAltScreen(v tea.View) tea.View {
	v.AltScreen = true
	return v
}

func Indent(s string, n int) string {
	if n <= 0 || s == "" {
		return s
	}
	l := strings.Split(s, "\n")
	b := strings.Builder{}
	i := strings.Repeat(" ", n)
	for _, v := range l {
		fmt.Fprintf(&b, "%s%s\n", i, v)
	}
	return b.String()
}
