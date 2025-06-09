package iban

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/jacoelho/banking/iban"
	"github.com/skatkov/devtui/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

const Title = "IBAN Generator"

type IBANGenerate struct {
	common            *ui.CommonModel
	form              *huh.Form
	countryCode       string
	generatedIBAN     string
	formattedIBAN     string
	error             string
}

type CountryOption struct {
	Code string
	Name string
}

func getCountryOptions() []huh.Option[string] {
	countries := []CountryOption{
		{"AD", "Andorra"},
		{"AE", "United Arab Emirates"},
		{"AL", "Albania"},
		{"AT", "Austria"},
		{"AZ", "Azerbaijan"},
		{"BA", "Bosnia and Herzegovina"},
		{"BE", "Belgium"},
		{"BG", "Bulgaria"},
		{"BH", "Bahrain"},
		{"BI", "Burundi"},
		{"BR", "Brazil"},
		{"BY", "Belarus"},
		{"CH", "Switzerland"},
		{"CR", "Costa Rica"},
		{"CY", "Cyprus"},
		{"CZ", "Czech Republic"},
		{"DE", "Germany"},
		{"DJ", "Djibouti"},
		{"DK", "Denmark"},
		{"DO", "Dominican Republic"},
		{"EE", "Estonia"},
		{"EG", "Egypt"},
		{"ES", "Spain"},
		{"FI", "Finland"},
		{"FK", "Falkland Islands"},
		{"FO", "Faroe Islands"},
		{"FR", "France"},
		{"GB", "United Kingdom"},
		{"GE", "Georgia"},
		{"GI", "Gibraltar"},
		{"GL", "Greenland"},
		{"GR", "Greece"},
		{"GT", "Guatemala"},
		{"HN", "Honduras"},
		{"HR", "Croatia"},
		{"HU", "Hungary"},
		{"IE", "Ireland"},
		{"IL", "Israel"},
		{"IQ", "Iraq"},
		{"IS", "Iceland"},
		{"IT", "Italy"},
		{"JO", "Jordan"},
		{"KW", "Kuwait"},
		{"KZ", "Kazakhstan"},
		{"LB", "Lebanon"},
		{"LC", "Saint Lucia"},
		{"LI", "Liechtenstein"},
		{"LT", "Lithuania"},
		{"LU", "Luxembourg"},
		{"LV", "Latvia"},
		{"LY", "Libya"},
		{"MC", "Monaco"},
		{"MD", "Moldova"},
		{"ME", "Montenegro"},
		{"MK", "Macedonia"},
		{"MN", "Mongolia"},
		{"MR", "Mauritania"},
		{"MT", "Malta"},
		{"MU", "Mauritius"},
		{"NI", "Nicaragua"},
		{"NL", "Netherlands"},
		{"NO", "Norway"},
		{"OM", "Oman"},
		{"PK", "Pakistan"},
		{"PL", "Poland"},
		{"PS", "Palestine"},
		{"PT", "Portugal"},
		{"QA", "Qatar"},
		{"RO", "Romania"},
		{"RS", "Serbia"},
		{"RU", "Russia"},
		{"SA", "Saudi Arabia"},
		{"SC", "Seychelles"},
		{"SD", "Sudan"},
		{"SE", "Sweden"},
		{"SI", "Slovenia"},
		{"SK", "Slovakia"},
		{"SM", "San Marino"},
		{"SO", "Somalia"},
		{"ST", "Sao Tome and Principe"},
		{"SV", "El Salvador"},
		{"TL", "Timor-Leste"},
		{"TN", "Tunisia"},
		{"TR", "Turkey"},
		{"UA", "Ukraine"},
		{"VA", "Vatican City"},
		{"VG", "British Virgin Islands"},
		{"XK", "Kosovo"},
		{"GF", "French Guyana"},
		{"GP", "Guadeloupe"},
		{"MQ", "Martinique"},
		{"RE", "Reunion"},
		{"FP", "French Polynesia"},
		{"TF", "French Southern Territories"},
		{"YT", "Mayotte"},
		{"NC", "New Caledonia"},
		{"BL", "Saint Barthelemy"},
		{"MF", "Saint Martin"},
		{"PM", "Saint Pierre et Miquelon"},
		{"WF", "Wallis and Futuna Islands"},
		{"YE", "Yemen"},
	}

	options := make([]huh.Option[string], len(countries))
	for i, country := range countries {
		options[i] = huh.NewOption(fmt.Sprintf("%s - %s", country.Code, country.Name), country.Code)
	}

	return options
}

func NewIBANGenerateModel(common *ui.CommonModel) *IBANGenerate {
	m := IBANGenerate{common: common}
	accessible, _ := strconv.ParseBool(os.Getenv("ACCESSIBLE"))

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select Country").
				Options(getCountryOptions()...).
				Height(10).
				Value(&m.countryCode),
		),
	).WithTheme(huh.ThemeCharm()).WithAccessible(accessible).WithShowHelp(false)

	return &m
}

func (m *IBANGenerate) Init() tea.Cmd {
	return m.form.Init()
}

func (m *IBANGenerate) View() string {
	s := m.common.Styles
	switch m.form.State {
	case huh.StateCompleted:
		var rows [][]string
		
		// Find country name for display
		countryName := ""
		for _, option := range getCountryOptions() {
			if option.Value == m.countryCode {
				countryName = strings.TrimPrefix(option.Key, m.countryCode+" - ")
				break
			}
		}
		
		rows = append(rows, []string{"Country", fmt.Sprintf("%s (%s)", countryName, m.countryCode)})
		
		if m.error != "" {
			rows = append(rows, []string{"Error", m.error})
		} else {
			rows = append(rows, []string{"IBAN", m.generatedIBAN})
			rows = append(rows, []string{"IBAN (Formatted)", m.formattedIBAN})
		}

		tableOutput := table.New().
			Border(lipgloss.RoundedBorder()).
			Width(100).
			Rows(rows...)

		header := s.Title.Render(lipgloss.JoinHorizontal(lipgloss.Left,
			m.common.AppTitle(),
			" :: ",
			lipgloss.NewStyle().Bold(true).Render(Title),
		))
		
		results := m.common.Lg.NewStyle().Margin(1, 0).Render(tableOutput.String())
		body := lipgloss.JoinVertical(lipgloss.Top, results)

		return s.Base.Render(header + "\n" + body)
	default:
		header := s.Title.Render(lipgloss.JoinHorizontal(lipgloss.Left,
			m.common.AppTitle(),
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

func (m *IBANGenerate) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Window size is received when starting up and on every resize
	case tea.WindowSizeMsg:
		m.common.Width = msg.Width
		m.common.Height = msg.Height
	case tea.KeyMsg:
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
			m.generateIBAN()
		}
	}
	return m, cmd
}

func (m *IBANGenerate) generateIBAN() {
	if m.countryCode == "" {
		m.error = "Country code is required"
		return
	}

	generatedIban, err := iban.Generate(m.countryCode)
	if err != nil {
		m.error = fmt.Sprintf("Failed to generate IBAN: %v", err)
		return
	}

	m.error = ""
	m.generatedIBAN = generatedIban
	m.formattedIBAN = iban.PaperFormat(generatedIban)
}