package ui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const AppTitle = "DevTUI"

type ReturnToListMsg struct {
	Common *CommonModel
}

type CommonModel struct {
	Width  int
	Height int

	LastSelectedItem int

	Lg     *lipgloss.Renderer
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
