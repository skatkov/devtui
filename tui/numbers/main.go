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
	"github.com/skatkov/devtui/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

const Title = "Number Base Converter"

type NumbersModel struct {
	common *ui.CommonModel
	form   *huh.Form
	value  int64
	base   NumberBase
	input  string
}

type NumberBase struct {
	title string
	base  int
}

var (
	Base2            = NumberBase{title: "Base 2 (binary)", base: 2}
	Base8            = NumberBase{title: "Base 8 (octal)", base: 8}
	Base10           = NumberBase{title: "Base 10 (decimal)", base: 10}
	Base16           = NumberBase{title: "Base 16 (hexadecimal)", base: 16}
	ReturnedBaseList = []NumberBase{Base2, Base8, Base10, Base16}
)

func NewNumberModel(common *ui.CommonModel) NumbersModel {
	m := NumbersModel{
		common: common,
	}
	accessible, _ := strconv.ParseBool(os.Getenv("ACCESSIBLE"))

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[NumberBase]().
				Key("base").
				Options(
					huh.NewOption(Base2.title, Base2),
					huh.NewOption(Base8.title, Base8),
					huh.NewOption(Base10.title, Base10).Selected(true),
					huh.NewOption(Base16.title, Base16),
				).
				Title("Select Base").Value(&m.base),
			huh.NewInput().
				Key("input").
				Placeholder(fmt.Sprintf("Enter a %s number", m.base.title)).
				Title("Enter a number").
				Validate(func(s string) error {
					if len(s) == 0 {
						return errors.New("number cannot be empty")
					}
					_, err := strconv.ParseInt(s, m.base.base, 64)
					if err != nil {
						return fmt.Errorf("please enter a valid %s number", m.base.title)
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
			if base, ok := m.form.Get("base").(NumberBase); ok {
				if val, err := strconv.ParseInt(m.form.GetString("input"), base.base, 64); err == nil {
					m.value = val
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
		rows := make([][]string, len(ReturnedBaseList))
		for i, numberBase := range ReturnedBaseList {
			rows[i] = []string{
				numberBase.title,
				strconv.FormatInt(m.value, numberBase.base),
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
