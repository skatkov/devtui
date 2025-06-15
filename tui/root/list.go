package root

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/adrg/xdg"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/skatkov/devtui/internal/ui"
	base64decoder "github.com/skatkov/devtui/tui/base64-decoder"
	base64encoder "github.com/skatkov/devtui/tui/base64-encoder"
	cron "github.com/skatkov/devtui/tui/cron"
	"github.com/skatkov/devtui/tui/css"
	"github.com/skatkov/devtui/tui/csv2md"
	"github.com/skatkov/devtui/tui/csvjson"
	graphqlquery "github.com/skatkov/devtui/tui/graphql-query"
	"github.com/skatkov/devtui/tui/html"
	"github.com/skatkov/devtui/tui/iban"
	js "github.com/skatkov/devtui/tui/json"
	"github.com/skatkov/devtui/tui/jsonstruct"
	"github.com/skatkov/devtui/tui/jsontoml"
	"github.com/skatkov/devtui/tui/markdown"
	"github.com/skatkov/devtui/tui/numbers"
	"github.com/skatkov/devtui/tui/toml"
	"github.com/skatkov/devtui/tui/tomljson"
	"github.com/skatkov/devtui/tui/tsv2md"
	urlextractor "github.com/skatkov/devtui/tui/url-extractor"
	uuiddecode "github.com/skatkov/devtui/tui/uuid-decode"
	uuidgenerate "github.com/skatkov/devtui/tui/uuid-generate"
	"github.com/skatkov/devtui/tui/xml"
	"github.com/skatkov/devtui/tui/yaml"
	"github.com/skatkov/devtui/tui/yamlstruct"
)

const listHeight = 15

type listModel struct {
	list   list.Model
	err    string
	common *ui.CommonModel
	items  []MenuOption
}

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(2).PaddingBottom(1)
)

type MenuOption struct {
	id         string
	title      string
	model      func() tea.Model
	usageCount int
}

func (i MenuOption) FilterValue() string { return i.title }
func (i MenuOption) Title() string       { return i.title }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(MenuOption)
	if !ok {
		return
	}

	str := fmt.Sprintf("%s %s", "•", i.title)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	_, err := fmt.Fprint(w, fn(str))
	if err != nil {
		fmt.Println("Error rendering item:", err)
	}
}

func getMenuOptions(common *ui.CommonModel) []MenuOption {
	return []MenuOption{
		{
			id:    "base64-encoder",
			title: base64encoder.Title,
			model: func() tea.Model { return base64encoder.NewBase64Model(common) },
		},
		{
			id:    "base64-decoder",
			title: base64decoder.Title,
			model: func() tea.Model { return base64decoder.NewBase64Model(common) },
		},
		{
			id:    "uuiddecode",
			title: uuiddecode.Title,
			model: func() tea.Model { return uuiddecode.NewUUIDDecodeModel(common) },
		},
		{
			id:    "numbers",
			title: numbers.Title,
			model: func() tea.Model { return numbers.NewNumberModel(common) },
		},
		{
			id:    "uuidgenerate",
			title: uuidgenerate.Title,
			model: func() tea.Model { return uuidgenerate.NewUUIDGenerateModel(common) },
		},
		{
			id:    "iban",
			title: iban.Title,
			model: func() tea.Model { return iban.NewIBANGenerateModel(common) },
		},
		{
			id:    "cron",
			title: cron.Title,
			model: func() tea.Model { return cron.NewCronModel(common) },
		},
		{
			id:    "json",
			title: js.Title,
			model: func() tea.Model { return js.NewJsonModel(common) },
		},
		{
			id:    "yaml",
			title: yaml.Title,
			model: func() tea.Model { return yaml.NewYamlModel(common) },
		},
		{
			id:    "markdown",
			title: markdown.Title,
			model: func() tea.Model { return markdown.NewMarkdownModel(common) },
		},
		{
			id:    "jsonstruct",
			title: jsonstruct.Title,
			model: func() tea.Model { return jsonstruct.NewJsonStructModel(common) },
		},
		{
			id:    "yamlstruct",
			title: yamlstruct.Title,
			model: func() tea.Model { return yamlstruct.NewYamlStructModel(common) },
		},
		{
			id:    "csvjson",
			title: csvjson.Title,
			model: func() tea.Model { return csvjson.NewCSVJsonModel(common) },
		},
		{
			id:    "tomljson",
			title: tomljson.Title,
			model: func() tea.Model { return tomljson.NewTomlJsonModel(common) },
		},
		{
			id:    "jsontoml",
			title: jsontoml.Title,
			model: func() tea.Model { return jsontoml.NewJsonTomlModel(common) },
		},
		{
			id:    "toml",
			title: toml.Title,
			model: func() tea.Model { return toml.NewTomlFormatModel(common) },
		},
		{
			id:    "html",
			title: html.Title,
			model: func() tea.Model { return html.NewHTMLFormatterModel(common) },
		},
		{
			id:    "xml",
			title: xml.Title,
			model: func() tea.Model { return xml.NewXMLFormatterModel(common) },
		},
		{
			id:    "css",
			title: css.Title,
			model: func() tea.Model { return css.NewCSSFormatterModel(common) },
		},
		{
			id:    "graphql-query",
			title: graphqlquery.Title,
			model: func() tea.Model { return graphqlquery.NewGraphQLQueryModel(common) },
		},
		{
			id:    "csv2md",
			title: csv2md.Title,
			model: func() tea.Model { return csv2md.NewCSV2MDModel(common) },
		},
		{
			id:    "tsv2md",
			title: tsv2md.Title,
			model: func() tea.Model { return tsv2md.NewTSV2MDModel(common) },
		},
		{
			id:    "url-extractor",
			title: urlextractor.Title,
			model: func() tea.Model { return urlextractor.NewURLExtractorModel(common) },
		},
	}
}

func newListModel(common *ui.CommonModel) *listModel {
	menuOptions := getMenuOptions(common)

	// Load usage stats
	stats, err := loadUsageStats()
	if err != nil {
		// Just log the error and continue with zero counts
		fmt.Fprintf(os.Stderr, "Failed to load usage stats: %v\n", err)
	}

	// Apply usage counts to the items
	for i := range menuOptions {
		menuOptions[i].usageCount = stats[menuOptions[i].id]
	}

	// Sort items by usage count (descending)
	sort.Slice(menuOptions, func(i, j int) bool {
		return menuOptions[i].usageCount > menuOptions[j].usageCount
	})

	// Convert to list.Item interface for bubbles/list
	listItems := make([]list.Item, len(menuOptions))
	for i, item := range menuOptions {
		listItems[i] = item
	}

	delegate := itemDelegate{}
	l := list.New(listItems, delegate, 20, listHeight)
	l.Title = ui.AppTitle
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = l.Styles.Title.MarginTop(1)
	l.FilterInput.PromptStyle = l.FilterInput.PromptStyle.MarginTop(1)

	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	return &listModel{
		list:   l,
		common: common,
		items:  menuOptions,
	}
}

func (m listModel) Init() tea.Cmd {
	return nil
}

func (m *listModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.common.Width = msg.Width
		m.common.Height = msg.Height
		m.list.SetSize(msg.Width, msg.Height)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "enter":
			m.common.LastSelectedItem = m.list.Index()
			i, ok := m.list.SelectedItem().(MenuOption)
			if ok {
				// Find the selected item in our items slice and increment usage
				for idx := range m.items {
					if m.items[idx].id == i.id {
						m.items[idx].usageCount++
						break
					}
				}

				// Save updated usage stats
				if err := saveUsageStats(m.items); err != nil {
					m.err = fmt.Sprintf("Failed to save usage stats: %v", err)
				}
				newScreen := i.model()
				return newScreen, newScreen.Init()
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m listModel) View() string {
	if m.err != "" {
		return lipgloss.NewStyle().Padding(2).Render(m.err)
	}
	return m.list.View()
}

func (m *listModel) RefreshOrder() {
	// Re-sort items by usage count
	sort.Slice(m.items, func(i, j int) bool {
		return m.items[i].usageCount > m.items[j].usageCount
	})

	// Convert MenuOptions to list.Items
	items := make([]list.Item, len(m.items))
	for i, opt := range m.items {
		items[i] = opt
	}

	// Update the list with the new sorted items
	m.list.SetItems(items)
}

// Update saveUsageStats to use xdg.
func saveUsageStats(items []MenuOption) error {
	// Get the proper config path using xdg
	configPath := filepath.Join(xdg.ConfigHome, "devtui")

	if err := os.MkdirAll(configPath, 0o750); err != nil {
		return err
	}

	statsFile := filepath.Join(configPath, "usage_stats.json")

	// Create a map of title -> count for serialization
	stats := make(map[string]int)
	for _, item := range items {
		stats[item.id] = item.usageCount
	}

	data, err := json.Marshal(stats)
	if err != nil {
		return err
	}

	return os.WriteFile(statsFile, data, 0o600)
}

// Update loadUsageStats to use xdg.
func loadUsageStats() (map[string]int, error) {
	statsFile := filepath.Join(xdg.ConfigHome, "devtui", "usage_stats.json")

	data, err := os.ReadFile(statsFile)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]int), nil
		}
		return nil, err
	}

	var stats map[string]int
	if err := json.Unmarshal(data, &stats); err != nil {
		return nil, err
	}

	return stats, nil
}
