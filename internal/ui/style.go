package ui

import "github.com/charmbracelet/lipgloss"

var (
	green           = lipgloss.Color("#04B575")
	mintGreen       = lipgloss.AdaptiveColor{Light: "#89F0CB", Dark: "#89F0CB"}
	darkGreen       = lipgloss.AdaptiveColor{Light: "#1C8760", Dark: "#1C8760"}
	statusBarNoteFg = lipgloss.AdaptiveColor{Light: "#656565", Dark: "#7D7D7D"}
	statusBarBg     = lipgloss.AdaptiveColor{Light: "#E6E6E6", Dark: "#242424"}
	lineNumberFg    = lipgloss.AdaptiveColor{Light: "#656565", Dark: "#7D7D7D"}
	lightRed        = lipgloss.AdaptiveColor{Light: "#FFAAAA", Dark: "#FFAAAA"}
	darkRed         = lipgloss.AdaptiveColor{Light: "#CC0000", Dark: "#AA0000"}

	HelpViewStyle = lipgloss.NewStyle().
			Foreground(statusBarNoteFg).
			Background(lipgloss.AdaptiveColor{Light: "#f2f2f2", Dark: "#1B1B1B"}).
			Render

	StatusBarNoteStyle = lipgloss.NewStyle().
				Foreground(statusBarNoteFg).
				Background(statusBarBg).Render

	AppNameStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("230")).
			Background(lipgloss.Color("62")).
			Bold(true).Render

	StatusBarErrorHelpStyle = lipgloss.NewStyle().
				Foreground(lightRed).
				Background(darkRed).Render

	StatusBarMessageHelpStyle = lipgloss.NewStyle().
					Foreground(lipgloss.Color("#B6FFE4")).
					Background(green).Render

	StatusBarHelpStyle = lipgloss.NewStyle().
				Foreground(statusBarNoteFg).
				Background(lipgloss.AdaptiveColor{Light: "#DCDCDC", Dark: "#323232"}).Render

	StatusBarScrollPosStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#949494", Dark: "#5A5A5A"}).
				Background(statusBarBg).
				Render
	StatusBarMessageStyle = lipgloss.NewStyle().
				Foreground(mintGreen).
				Background(darkGreen).Render
	StatusBarErrorStyle = lipgloss.NewStyle().
				Foreground(lightRed).
				Background(darkRed).Render

	LineNumberStyle = lipgloss.NewStyle().
			Foreground(lineNumberFg).
			Render
)

// Use huh/examples/dynamic-bubbletea as an example how to style an application with a huh form in it.
type Styles struct {
	Base,
	Title,
	Subtitle,
	Help lipgloss.Style
}

func NewStyle(lg *lipgloss.Renderer) *Styles {
	s := Styles{
		Base: lg.NewStyle().Padding(1, 4, 1, 2),
		Title: lg.NewStyle().
			Background(lipgloss.Color("62")).
			Foreground(lipgloss.Color("230")).
			Padding(0, 1),
		Help: lg.NewStyle().Foreground(lipgloss.Color("240")),
	}
	return &s
}
