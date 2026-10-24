package cron

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
	"github.com/lnquy/cron"
	"github.com/skatkov/devtui/internal/ui"

	tea "charm.land/bubbletea/v2"
)

const Title = "Cron Job Parser"

// @see https://gist.github.com/Aterfax/401875eb3d45c9c114bbef69364dd045
// @see https://regexr.com/4jp54
var cronPattern = regexp.MustCompile(`^(((\d+,)+\d+|(\d+(\/|-|#)\d+(\/\d+)?)|\d+L?|\*(\/\d+)?|L(-\d+)?|\?|[A-Z]{3}(-[A-Z]{3})?) ?){5,7}$`)

// displayOptions configure the descriptions shown to the user.
var displayOptions = []cron.Option{cron.Use24HourTimeFormat(true)}

type CronModel struct {
	common         *ui.CommonModel
	form           *huh.Form
	cronExpression string
}

func NewCronModel(common *ui.CommonModel) *CronModel {
	m := &CronModel{
		common:         common,
		cronExpression: "*/5 * * * *",
	}

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Cron Expression").
				Value(&m.cronExpression).
				Placeholder("*/5 * * * *").
				Validate(validateCronExpression).
				DescriptionFunc(m.description, &m.cronExpression),
		),
	).WithTheme(huh.ThemeFunc(huh.ThemeCharm)).WithShowHelp(false)

	return m
}

func validateCronExpression(str string) error {
	// First validate with regexp
	if !cronPattern.MatchString(str) {
		return errors.New("invalid cron expression format")
	}

	// Then validate with cron descriptor
	if _, err := describeCron(str); err != nil {
		return fmt.Errorf("invalid cron expression: %v", err)
	}
	return nil
}

func describeCron(expression string, options ...cron.Option) (string, error) {
	expr, err := cron.NewDescriptor(options...)
	if err != nil {
		return "", err
	}
	return expr.ToDescription(expression, cron.Locale_en)
}

func (m *CronModel) description() string {
	desc, err := describeCron(m.cronExpression, displayOptions...)
	if err != nil {
		return ""
	}

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#87CEEB")).PaddingLeft(20)

	return valueStyle.Render(desc)
}

func (m *CronModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m *CronModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Window size is received when starting up and on every resize
	case tea.WindowSizeMsg:
		m.common.Width = msg.Width
		m.common.Height = msg.Height
	case tea.KeyPressMsg:
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

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}
	return m, cmd
}

func (m *CronModel) View() tea.View {
	s := m.common.Styles
	switch m.form.State {
	case huh.StateCompleted:
		desc, err := describeCron(m.cronExpression, displayOptions...)
		if err != nil {
			return ui.AltScreenView(lipgloss.NewStyle().Padding(2).
				Render(fmt.Sprintf("Error parsing cron expression: %v", err)))
		}
		titleStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF69B4")).
			Bold(true)

		valueStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#87CEEB"))

		output := titleStyle.Render(m.cronExpression) + " \n\n" + valueStyle.Render(desc)

		return ui.AltScreenView(s.Base.Render(output))
	default:
		header := s.Title.Render(lipgloss.JoinHorizontal(lipgloss.Left,
			ui.AppTitle,
			" :: ",
			lipgloss.NewStyle().Bold(true).Render(Title),
		))
		v := strings.TrimSuffix(m.form.View(), "\n\n")
		form := lipgloss.NewStyle().Margin(1, 0).Render(v)
		body := lipgloss.JoinVertical(
			lipgloss.Top,
			form,
			lipgloss.PlaceVertical(
				m.common.Height-lipgloss.Height(header)-lipgloss.Height(form)-2,
				lipgloss.Bottom,
				m.form.Help().ShortHelpView(m.form.KeyBinds()),
			),
		)
		return ui.AltScreenView(s.Base.Render(header + "\n" + body))
	}
}
