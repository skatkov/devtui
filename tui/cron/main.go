package cron

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/lnquy/cron"
	"github.com/skatkov/devtui/tui/messages"

	tea "github.com/charmbracelet/bubbletea"
)

type CronModel struct {
	form           *huh.Form
	cronExpression string
}

func NewCronModel() *CronModel {
	m := &CronModel{}
	accessible, _ := strconv.ParseBool(os.Getenv("ACCESSIBLE"))

	// @see https://gist.github.com/Aterfax/401875eb3d45c9c114bbef69364dd045
	// @see https://regexr.com/4jp54
	cronRegex := `^((((\d+,)+\d+|(\d+(\/|-|#)\d+)|\d+L?|\*(\/\d+)?|L(-\d+)?|\?|[A-Z]{3}(-[A-Z]{3})?) ?){5,7})|(@(annually|yearly|monthly|weekly|daily|hourly|reboot))|(@every (\d+(ns|us|µs|ms|s|m|h))+)$`

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Cron Expression").
				Placeholder("* * * * *").
				Value(&m.cronExpression).
				Validate(func(str string) error {
					// First validate with regexp
					matched, err := regexp.MatchString(cronRegex, str)
					if err != nil {
						return fmt.Errorf("Validation error: %v", err)
					}
					if !matched {
						return fmt.Errorf("invalid cron expression format")
					}

					// Then validate with cron descriptor
					expr, err := cron.NewDescriptor()
					if err != nil {
						return fmt.Errorf("invalid cron expression: %v", err)
					}
					_, err = expr.ToDescription(str, cron.Locale_en)
					if err != nil {
						return fmt.Errorf("invalid cron expression: %v", err)
					}
					return nil
				}),
		),
	).WithTheme(huh.ThemeCharm()).WithAccessible(accessible)

	return m
}

func (m *CronModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m *CronModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, func() tea.Msg {
				return messages.ReturnToListMsg{}
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}
	return m, cmd
}

func (m *CronModel) View() string {
	switch m.form.State {
	case huh.StateCompleted:
		expr, err := cron.NewDescriptor(
			cron.Use24HourTimeFormat(true),
			cron.DayOfWeekStartsAtOne(true),
		)
		if err != nil {
			return lipgloss.NewStyle().Padding(2).
				Render(fmt.Sprintf("Error parsing cron expression: %v", err))
		}

		desc, err := expr.ToDescription(m.cronExpression, cron.Locale_en)
		if err != nil {
			return lipgloss.NewStyle().Padding(2).
				Render(fmt.Sprintf("Error parsing cron expression: %v", err))
		}
		titleStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF69B4")).
			Bold(true)

		valueStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#87CEEB"))

		output := fmt.Sprintf("%s \n\n",
			titleStyle.Render(m.cronExpression)) +
			fmt.Sprintf("%s",
				valueStyle.Render(desc))

		return lipgloss.NewStyle().
			Padding(2).
			PaddingTop(1).
			Render(output)
	default:
		return lipgloss.NewStyle().Padding(2).Render(m.form.View())
	}
}
