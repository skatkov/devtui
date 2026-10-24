// devtui/tui/url-extractor/main.go
package urlextractor

import (
	"bytes"
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/skatkov/devtui/internal/clipboard"
	"github.com/skatkov/devtui/internal/editor"
	"github.com/skatkov/devtui/internal/ui"
	"mvdan.cc/xurls/v2"
)

const (
	Title       = "URL Extractor"
	ModeStrict  = "strict"
	ModeRelaxed = "relaxed"
)

type URLExtractorModel struct {
	ui.BasePagerModel
	StrictMode bool
}

func NewURLExtractorModel(common *ui.CommonModel) *URLExtractorModel {
	model := URLExtractorModel{
		BasePagerModel: ui.NewBasePagerModel(common, Title),
		StrictMode:     false, // relaxed by default
	}

	return &model
}

func (m *URLExtractorModel) Init() tea.Cmd {
	return m.BasePagerModel.Init()
}

func (m *URLExtractorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if cmd, handled := m.HandleCommonKeys(msg); handled {
			return m, cmd
		}

		switch msg.String() {
		case "e":
			return m, editor.OpenEditor(m.Content, "txt")
		case "v":
			cmds = append(cmds, m.applyContent(clipboard.Paste()))
		case "s":
			m.StrictMode = !m.StrictMode
			if m.Content != "" {
				err := m.SetContent(m.Content)
				if err != nil {
					cmds = append(cmds, m.ShowErrorMessage(err.Error()))
				} else {
					mode := ModeRelaxed
					if m.StrictMode {
						mode = ModeStrict
					}
					cmds = append(cmds, m.ShowStatusMessage(fmt.Sprintf("Switched to %s mode. Press 'c' to copy result.", mode)))
				}
			} else {
				mode := ModeRelaxed
				if m.StrictMode {
					mode = ModeStrict
				}
				cmds = append(cmds, m.ShowStatusMessage(fmt.Sprintf("Switched to %s mode.", mode)))
			}
		}
	case ui.StatusMessageTimeoutMsg:
		m.State = ui.PagerStateBrowse

	case editor.EditorFinishedMsg:
		if msg.Err != nil {
			return m, m.ShowErrorMessage(msg.Err.Error())
		} else {
			err := m.SetContent(msg.Content)

			if err != nil {
				cmds = append(cmds, m.ShowErrorMessage(err.Error()))
			} else {
				cmds = append(cmds, m.ShowStatusMessage("URLs extracted. Press 'c' to copy result."))
			}
		}
	case tea.WindowSizeMsg:
		cmd = m.HandleWindowSizeMsg(msg)
		cmds = append(cmds, cmd)
	}

	m.Viewport, cmd = m.Viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *URLExtractorModel) View() tea.View {
	var b strings.Builder

	fmt.Fprint(&b, m.Viewport.View()+"\n")
	fmt.Fprint(&b, m.StatusBarView())

	if m.ShowHelp {
		fmt.Fprint(&b, "\n"+m.helpView())
	}

	return m.NewView(b.String())
}

func (m *URLExtractorModel) SetContent(content string) error {
	m.Content = content

	// Extract URLs based on mode
	var urls []string
	if m.StrictMode {
		urls = xurls.Strict().FindAllString(content, -1)
	} else {
		urls = xurls.Relaxed().FindAllString(content, -1)
	}

	// Remove duplicates while preserving order
	seen := make(map[string]bool)
	uniqueURLs := make([]string, 0, len(urls))
	for _, url := range urls {
		if !seen[url] {
			seen[url] = true
			uniqueURLs = append(uniqueURLs, url)
		}
	}

	result := strings.Join(uniqueURLs, "\n")
	m.FormattedContent = result

	// Set content in viewport
	var buf bytes.Buffer
	buf.WriteString(m.FormattedContent)
	m.Viewport.SetContent(buf.String())

	return nil
}

func (m *URLExtractorModel) helpView() string {
	mode := ModeRelaxed
	if m.StrictMode {
		mode = ModeStrict
	}

	col1 := []string{
		"c              copy URLs",
		"e              edit text",
		"v              paste text to extract",
		fmt.Sprintf("s              toggle mode (current: %s)", mode),
		"q/ctrl+c       quit",
	}

	return m.FormatHelpColumns(col1)
}

// applyContent sets text pasted from the clipboard and reports the outcome in
// the status bar.
func (m *URLExtractorModel) applyContent(content string, err error) tea.Cmd {
	if err == nil {
		err = m.SetContent(content)
	}
	if err != nil {
		return m.ShowErrorMessage(err.Error())
	}
	return m.ShowStatusMessage("Pasted and extracted URLs. Press 'c' to copy result.")
}
