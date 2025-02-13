package uuidgenerate

import (
	"fmt"
	"os"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/google/uuid"

	tea "github.com/charmbracelet/bubbletea"
)

type UUIDGenerate struct {
	form          *huh.Form
	version       int
	namespace     string
	generatedUUID uuid.UUID
}

func NewUUIDGenerateModel() *UUIDGenerate {
	m := UUIDGenerate{}
	accessible, _ := strconv.ParseBool(os.Getenv("ACCESSIBLE"))

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().
				Title("UUID Version").
				Options(
					huh.NewOption("Version 1 (Time-based)", 1),
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
		).WithHide(m.version == 3 || m.version == 5),
	).WithTheme(huh.ThemeCharm()).WithAccessible(accessible)

	return &m
}

func (m *UUIDGenerate) Init() tea.Cmd {
	return m.form.Init()
}

func (m *UUIDGenerate) View() string {
	switch m.form.State {
	case huh.StateCompleted:
		tableOutput := table.New().
			Border(lipgloss.RoundedBorder()).
			Width(100).
			Rows(
				[]string{"Version", fmt.Sprintf("%d", m.version)},
				[]string{"Generated UUID", m.generatedUUID.String()},
			)

		return lipgloss.NewStyle().Padding(2).PaddingTop(1).Render(tableOutput.String())
	default:
		return lipgloss.NewStyle().Padding(2).Render(m.form.View())
	}
}

func (m *UUIDGenerate) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Interrupt
		case "q", "esc":
			return m, tea.Quit
		}
	}

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f

		if m.form.State == huh.StateCompleted {
			var err error
			m.generatedUUID, err = m.generateUUID()
			if err != nil {
				// Handle error appropriately
				fmt.Println("Error generating UUID:", err)
			}
		}
	}
	return m, cmd
}

func (m *UUIDGenerate) generateUUID() (uuid.UUID, error) {
	switch m.version {
	case 1:
		return uuid.NewUUID()
	case 3:
		return uuid.NewMD5(uuid.NameSpaceURL, []byte(m.namespace)), nil
	case 4:
		return uuid.NewRandom()
	case 5:
		return uuid.NewSHA1(uuid.NameSpaceURL, []byte(m.namespace)), nil
	case 6:
		return uuid.NewV6()
	case 7:
		return uuid.NewV7()
	default:
		return uuid.Nil, fmt.Errorf("unsupported UUID version: %s", m.version)
	}
}
