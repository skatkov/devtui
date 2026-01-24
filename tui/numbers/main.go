package numbers

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/skatkov/devtui/internal/numbers"
	"github.com/skatkov/devtui/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

const Title = "Number Base Converter"

type NumbersModel struct {
	common *ui.CommonModel
	form   *huh.Form
	base   numbers.Base
	input  string
	result numbers.Result
}

func NewNumberModel(common *ui.CommonModel) NumbersModel {
	m := NumbersModel{
		common: common,
		base:   numbers.DefaultBase(),
	}
	accessible, _ := strconv.ParseBool(os.Getenv("ACCESSIBLE"))

	options := make([]huh.Option[numbers.Base], 0, len(numbers.Bases))
	defaultBase := numbers.DefaultBase()
	for _, base := range numbers.Bases {
		option := huh.NewOption(base.Label, base)
		if base.Base == defaultBase.Base {
			option = option.Selected(true)
		}
		options = append(options, option)
	}

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[numbers.Base]().
				Key("base").
				Options(options...).
				Title("Select Base").Value(&m.base),
			huh.NewInput().
				Key("input").
				Placeholder(fmt.Sprintf("Enter a %s number", m.base.Label)).
				Title("Enter a number").
				Validate(func(s string) error {
					if len(s) == 0 {
						return errors.New("number cannot be empty")
					}
					if _, err := numbers.Parse(s, m.base.Base); err != nil {
						return fmt.Errorf("please enter a valid %s number", m.base.Label)
					}
					return nil
				}).Value(&m.input),
		),
	).WithTheme(huh.ThemeCharm()).WithAccessible(accessible).WithShowHelp(false)

	return m
}

func (m NumbersModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m NumbersModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Window size is received when starting up and on every resize
	case tea.WindowSizeMsg:
		m.common.Width = msg.Width
		m.common.Height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, func() tea.Msg {
				return ui.ReturnToListMsg{
					Common: m.common,
				}
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	var cmds []tea.Cmd

	// Process the form
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f

		// If the form is completed, parse the input value
		if m.form.State == huh.StateCompleted {
			if base, ok := m.form.Get("base").(numbers.Base); ok {
				result, err := numbers.Convert(m.form.GetString("input"), base.Base)
				if err == nil {
					m.result = result
				}
			}
		}

		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m NumbersModel) View() string {
	s := m.common.Styles
	switch m.form.State {
	case huh.StateCompleted:
		rows := make([][]string, len(m.result.Conversions))
		for i, conversion := range m.result.Conversions {
			rows[i] = []string{
				conversion.Label,
				conversion.Value,
			}
		}
		t := table.New().
			Border(lipgloss.RoundedBorder()).
			Width(100).
			Headers("Base", "Value").
			Rows(rows...)
		return s.Base.Render(t.String())
	default:
		header := s.Title.Render(lipgloss.JoinHorizontal(lipgloss.Left,
			ui.AppTitle,
			" :: ",
			lipgloss.NewStyle().Bold(true).Render(Title),
		))
		v := strings.TrimSuffix(m.form.View(), "\n\n")
		form := m.common.Lg.NewStyle().Margin(1, 0).Render(v)
		body := lipgloss.JoinVertical(
			lipgloss.Top,
			form,
			lipgloss.PlaceVertical(
				m.common.Height-lipgloss.Height(header)-lipgloss.Height(form)-2,
				lipgloss.Bottom,
				m.form.Help().ShortHelpView(m.form.KeyBinds()),
			),
		)
		return s.Base.Render(header + "\n" + body)
	}
}
