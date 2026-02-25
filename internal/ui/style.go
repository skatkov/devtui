package ui

import (
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/compat"
)

var (
	green           = lipgloss.Color("#04B575")
	mintGreen       = compat.AdaptiveColor{Light: lipgloss.Color("#89F0CB"), Dark: lipgloss.Color("#89F0CB")}
	darkGreen       = compat.AdaptiveColor{Light: lipgloss.Color("#1C8760"), Dark: lipgloss.Color("#1C8760")}
	statusBarNoteFg = compat.AdaptiveColor{Light: lipgloss.Color("#656565"), Dark: lipgloss.Color("#7D7D7D")}
	statusBarBg     = compat.AdaptiveColor{Light: lipgloss.Color("#E6E6E6"), Dark: lipgloss.Color("#242424")}
	lineNumberFg    = compat.AdaptiveColor{Light: lipgloss.Color("#656565"), Dark: lipgloss.Color("#7D7D7D")}
	lightRed        = compat.AdaptiveColor{Light: lipgloss.Color("#FFAAAA"), Dark: lipgloss.Color("#FFAAAA")}
	darkRed         = compat.AdaptiveColor{Light: lipgloss.Color("#CC0000"), Dark: lipgloss.Color("#AA0000")}

	HelpViewStyle = lipgloss.NewStyle().
			Foreground(statusBarNoteFg).
			Background(compat.AdaptiveColor{Light: lipgloss.Color("#f2f2f2"), Dark: lipgloss.Color("#1B1B1B")}).
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
				Background(compat.AdaptiveColor{Light: lipgloss.Color("#DCDCDC"), Dark: lipgloss.Color("#323232")}).Render

	StatusBarScrollPosStyle = lipgloss.NewStyle().
				Foreground(compat.AdaptiveColor{Light: lipgloss.Color("#949494"), Dark: lipgloss.Color("#5A5A5A")}).
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

func NewStyle() *Styles {
	s := Styles{
		Base: lipgloss.NewStyle().Padding(1, 4, 1, 2),
		Title: lipgloss.NewStyle().
			Background(lipgloss.Color("62")).
			Foreground(lipgloss.Color("230")).
			Padding(0, 1),
		Help: lipgloss.NewStyle().Foreground(lipgloss.Color("240")),
	}
	return &s
}
