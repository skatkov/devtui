package uuiddecode

import (
	"os"
	"strconv"
	"strings"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
	"github.com/charmbracelet/huh"
	"github.com/google/uuid"
	"github.com/skatkov/devtui/internal/teacompat"
	"github.com/skatkov/devtui/internal/ui"
	"github.com/skatkov/devtui/internal/uuidutil"

	tea "charm.land/bubbletea/v2"
)

const Title = "UUID Decoder"

type UUIDDecode struct {
	common *ui.CommonModel
	form   *huh.Form
	uuid   string
}

func NewUUIDDecodeModel(common *ui.CommonModel) *UUIDDecode {
	m := UUIDDecode{common: common}
	accessible, _ := strconv.ParseBool(os.Getenv("ACCESSIBLE"))

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("UUID").
				Placeholder("Enter a UUID").
				Validate(func(value string) error {
					_, err := uuid.Parse(value)
					return err
				}).Value(&m.uuid),
		),
	).WithTheme(huh.ThemeCharm()).WithAccessible(accessible).WithShowHelp(false)

	return &m
}

func (m *UUIDDecode) Init() tea.Cmd {
	return teacompat.Cmd(m.form.Init())
}

func (m *UUIDDecode) View() tea.View {
	s := m.common.Styles
	switch m.form.State {
	case huh.StateCompleted:
		result, _ := uuid.Parse(m.uuid)

		fields := uuidutil.Decode(result)
		tableOutput := table.New().
			Border(lipgloss.RoundedBorder()).
			Width(100).
			Rows(uuidutil.FieldsToRows(fields)...)

		return ui.AltScreenView(s.Base.Render(tableOutput.String()))
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

func (m *UUIDDecode) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
	return m, teacompat.Cmd(cmd)
}
