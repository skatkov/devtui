package numbers

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"

	tea "github.com/charmbracelet/bubbletea"
)

type NumbersModel struct {
	form  *huh.Form
	value int64
	base  NumberBase
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

func NewNumberModel() NumbersModel {
	m := NumbersModel{}
	accessible, _ := strconv.ParseBool(os.Getenv("ACCESSIBLE"))

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[NumberBase]().
				Options(
					huh.NewOption(Base2.title, Base2),
					huh.NewOption(Base8.title, Base8),
					huh.NewOption(Base10.title, Base10).Selected(true),
					huh.NewOption(Base16.title, Base16),
				).
				Title("Select Base").Value(&m.base),
			huh.NewInput().
				Placeholder(fmt.Sprintf("Enter a %s number", m.base.title)).
				Title("Enter a number").
				Validate(func(s string) error {
					if len(s) == 0 {
						return errors.New("number cannot be empty")
					}
					val, err := strconv.ParseInt(s, m.base.base, 64)
					if err != nil {
						return errors.New(fmt.Sprintf("please enter a valid %s number", m.base.title))
					}
					m.value = val
					return nil
				}).Value(new(string)),
		),
	).WithTheme(huh.ThemeCharm()).WithAccessible(accessible)

	return m
}

func (m NumbersModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m NumbersModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Interrupt
		case "esc", "q":
			return m, tea.Quit
		}
	}

	var cmds []tea.Cmd

	// Process the form
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m NumbersModel) View() string {
	switch m.form.State {
	case huh.StateCompleted:
		rows := make([][]string, len(ReturnedBaseList))
		for i, numberBase := range ReturnedBaseList {
			rows[i] = []string{
				numberBase.title,
				strconv.FormatInt(m.number.Value, numberBase.base),
			}
		}
		t := table.New().
			Border(lipgloss.RoundedBorder()).
			Width(100).
			Headers("Base", "Value").
			Rows(rows...)
		return t.String()
	default:
		return m.form.View()
	}
}
