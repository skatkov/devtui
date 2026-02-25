package uuidgenerate

import (
	"fmt"
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

const Title = "UUID Generator"

type UUIDGenerate struct {
	common        *ui.CommonModel
	form          *huh.Form
	version       int
	namespace     string
	generatedUUID uuid.UUID
}

func NewUUIDGenerateModel(common *ui.CommonModel) *UUIDGenerate {
	m := UUIDGenerate{common: common}
	accessible, _ := strconv.ParseBool(os.Getenv("ACCESSIBLE"))

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().
				Title("UUID Version").
				Options(
					huh.NewOption("Version 1 (Time-based)", 1),
					huh.NewOption("Version 2 (Time-based)", 2),
					huh.NewOption("Version 3 (MD5 hash-based)", 3),
					huh.NewOption("Version 4 (Random)", 4),
					huh.NewOption("Version 5 (SHA1 hash-based)", 5),
					huh.NewOption("Version 6 (Time-based)", 6),
					huh.NewOption("Version 7 (Time-based)", 7),
				).
				Value(&m.version),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("Namespace").
				Value(&m.namespace),
		).WithHideFunc(func() bool { return m.hideNamespace() }),
	).WithTheme(huh.ThemeCharm()).WithAccessible(accessible).WithShowHelp(false)

	return &m
}

func (m *UUIDGenerate) hideNamespace() bool {
	switch m.version {
	case 3:
		return false
	case 5:
		return false
	default:
		return true
	}
}

func (m *UUIDGenerate) Init() tea.Cmd {
	return teacompat.Cmd(m.form.Init())
}

func (m *UUIDGenerate) View() tea.View {
	s := m.common.Styles
	switch m.form.State {
	case huh.StateCompleted:
		var rows [][]string
		rows = append(rows, []string{"Version", strconv.Itoa(m.version)})
		if m.namespace != "" {
			rows = append(rows, []string{"Namespace", m.namespace})
		}
		rows = append(rows, []string{"Generated UUID", m.generatedUUID.String()})

		tableOutput := table.New().
			Border(lipgloss.RoundedBorder()).
			Width(100).
			Rows(rows...)

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

func (m *UUIDGenerate) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Window size is received when starting up and on every resize
	case tea.WindowSizeMsg:
		m.common.Width = msg.Width
		m.common.Height = msg.Height
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "esc":
			return m, func() tea.Msg {
				return ui.ReturnToListMsg{
					Common: m.common,
				}
			}
		}
	}

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f

		if m.form.State == huh.StateCompleted {
			var err error
			m.generatedUUID, err = uuidutil.Generate(m.version, m.namespace)
			if err != nil {
				// Handle error appropriately
				fmt.Println("Error generating UUID:", err)
			}
		}
	}
	return m, teacompat.Cmd(cmd)
}
