package uuiddecode

import (
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/google/uuid"
	"github.com/skatkov/devtui/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
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
	return m.form.Init()
}

func (m *UUIDDecode) View() string {
	s := m.common.Styles
	switch m.form.State {
	case huh.StateCompleted:
		result, _ := uuid.Parse(m.uuid)

		tableOutput := table.New().
			Border(lipgloss.RoundedBorder()).
			Width(100).
			Rows(extractUUIDData(result)...)

		return s.Base.Render(tableOutput.String())
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

func (m *UUIDDecode) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}
	return m, cmd
}

func extractUUIDData(id uuid.UUID) [][]string {
	if id == uuid.Nil {
		return [][]string{}
	}
	var i big.Int
	i.SetString(strings.Replace(id.String(), "-", "", 4), 16)

	var result [][]string

	result = [][]string{
		{"Standard String Format", id.String()},
		{"Single Integer Value", i.String()},
		{"Version", fmt.Sprintf("%d", id.Version())},
		{"Variant", mapVariant(id.Variant())},
	}
	switch id.Version() {
	case uuid.Version(1):
		t := id.Time()
		sec, nsec := t.UnixTime()
		timeStamp := time.Unix(sec, nsec)
		node := id.NodeID()
		clockSeq := id.ClockSequence()

		result = append(result, []string{"Contents - Time", timeStamp.UTC().Format("2006-01-02 15:04:05.999999999 UTC")})
		result = append(result, []string{"Contents - Clock", strconv.Itoa(clockSeq)})
		result = append(result, []string{"Contents - Node", fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", node[0], node[1], node[2], node[3], node[4], node[5])})

	default:
		formatted := strings.ToUpper(strings.ReplaceAll(id.String(), "-", ""))
		var pairs []string
		for i := 0; i < len(formatted); i += 2 {
			if i < len(formatted)-1 {
				pairs = append(pairs, formatted[i:i+2])
			}
		}
		result = append(result, []string{"Contents", strings.Join(pairs, ":")})
	}

	return result
}

func mapVariant(v uuid.Variant) string {
	switch v {
	case uuid.Invalid:
		return "Invalid UUID"
	case uuid.RFC4122:
		return "DCE 1.1, ISO/IEC 11578:1996"
	case uuid.Reserved:
		return "Reserved (NCS backward compatibility)"
	case uuid.Microsoft:
		return "Reserved (Microsoft GUID)"
	case uuid.Future:
		return "Reserved (future use)"
	default:
		return "Unknown"
	}
}
