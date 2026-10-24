package root

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"github.com/adrg/xdg"
	"github.com/skatkov/devtui/internal/ui"
)

var ansiPattern = regexp.MustCompile(`\x1b\[[0-9;]*m`)

type otherItem struct{}

func (otherItem) FilterValue() string { return "other" }

func TestItemDelegateRender(t *testing.T) {
	t.Parallel()

	items := []list.Item{MenuOption{title: "A"}, MenuOption{title: "B"}, MenuOption{title: "C"}}
	l := list.New(items, itemDelegate{}, 20, 10)
	l.Select(1)

	tests := []struct {
		name  string
		index int
		want  string
	}{
		{"before selection", 0, "    • A"},
		{"selection", 1, "  > • B"},
		{"after selection", 2, "    • C"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var b strings.Builder
			itemDelegate{}.Render(&b, l, tt.index, items[tt.index])
			if got := stripANSI(b.String()); got != tt.want {
				t.Fatalf("Render(%d) = %q, want %q", tt.index, got, tt.want)
			}
		})
	}
}

func TestItemDelegateRenderIgnoresOtherItems(t *testing.T) {
	t.Parallel()

	var b strings.Builder
	itemDelegate{}.Render(&b, list.New(nil, itemDelegate{}, 20, 10), 0, otherItem{})
	if b.String() != "" {
		t.Fatalf("expected nothing rendered, got %q", b.String())
	}
}

func TestItemDelegateSize(t *testing.T) {
	t.Parallel()

	d := itemDelegate{}
	if d.Height() != 1 || d.Spacing() != 0 {
		t.Fatalf("expected height 1 and spacing 0, got %d and %d", d.Height(), d.Spacing())
	}
}

func TestNewListModelLayout(t *testing.T) {
	useTempConfig(t)

	lines := strings.Split(stripANSI(newListModel(&ui.CommonModel{}).View().Content), "\n")
	if help := "  ↑/k up • ↓/j down • / filter • v select • ? more"; lines[15] != help {
		t.Fatalf("expected help line %q without trailing padding, got %q", help, lines[15])
	}
	for i := range lines {
		lines[i] = strings.TrimRight(lines[i], " ")
	}

	want := []string{
		"",
		"   DevTUI",
		"",
		"  > • Base64 Encoder",
		"    • Base64 Decoder",
		"    • UUID Decoder",
		"    • Number Base Converter",
		"    • UUID Generator",
		"    • IBAN Generator",
		"    • Cron Job Parser",
		"    • JSON Formatter",
		"    • YAML Formatter",
		"",
		"    •••",
		"",
		"  ↑/k up • ↓/j down • / filter • v select • ? more",
		"",
	}
	if strings.Join(lines, "\n") != strings.Join(want, "\n") {
		t.Fatalf("unexpected menu layout:\n%s", strings.Join(lines, "\n"))
	}
}

func TestNewListModelKeepsMenuOrderWithoutUsage(t *testing.T) {
	useTempConfig(t)

	m := newListModel(&ui.CommonModel{})
	options := getMenuOptions(&ui.CommonModel{})
	for i, item := range m.list.Items() {
		if got := item.(MenuOption).id; got != options[i].id {
			t.Fatalf("expected %q at position %d, got %q", options[i].id, i, got)
		}
	}
	if !m.list.FilteringEnabled() {
		t.Fatal("expected filtering to be enabled")
	}
}

func TestNewListModelSortsByUsage(t *testing.T) {
	writeStats(t, useTempConfig(t), `{"yaml": 3, "json": 5}`)

	m := newListModel(&ui.CommonModel{})
	first := m.list.Items()[0].(MenuOption)
	second := m.list.Items()[1].(MenuOption)
	if first.id != "json" || first.usageCount != 5 || second.id != "yaml" || second.usageCount != 3 {
		t.Fatalf("expected json then yaml, got %+v and %+v", first, second)
	}
}

func TestNewListModelKeepsMenuOrderAmongUnusedTools(t *testing.T) {
	writeStats(t, useTempConfig(t), `{"json": 2}`)

	m := newListModel(&ui.CommonModel{})
	var want []string
	for _, option := range getMenuOptions(&ui.CommonModel{}) {
		if option.id != "json" {
			want = append(want, option.id)
		}
	}
	for i, item := range m.list.Items()[1:] {
		if got := item.(MenuOption).id; got != want[i] {
			t.Fatalf("expected %q at position %d, got %q", want[i], i+1, got)
		}
	}
}

func TestNewListModelReportsBrokenStats(t *testing.T) {
	dir := useTempConfig(t)

	if out := captureStderr(t, func() { newListModel(&ui.CommonModel{}) }); out != "" {
		t.Fatalf("expected no warning without stats, got %q", out)
	}

	writeStats(t, dir, "not json")
	out := captureStderr(t, func() { newListModel(&ui.CommonModel{}) })
	if !strings.Contains(out, "Failed to load usage stats") {
		t.Fatalf("expected load warning, got %q", out)
	}
}

func TestLoadUsageStats(t *testing.T) {
	dir := useTempConfig(t)

	stats, err := loadUsageStats()
	if err != nil || stats == nil || len(stats) != 0 {
		t.Fatalf("expected empty stats without a file, got %v and %v", stats, err)
	}

	writeStats(t, dir, `{"json": 2}`)
	stats, err = loadUsageStats()
	if err != nil || stats["json"] != 2 {
		t.Fatalf("expected json count 2, got %v and %v", stats, err)
	}

	writeStats(t, dir, "not json")
	if _, err := loadUsageStats(); err == nil {
		t.Fatal("expected error for invalid stats")
	}
}

func TestLoadUsageStatsUnreadableFileReturnsError(t *testing.T) {
	dir := useTempConfig(t)
	if err := os.MkdirAll(filepath.Join(dir, "devtui", "usage_stats.json"), 0o750); err != nil {
		t.Fatal(err)
	}

	if _, err := loadUsageStats(); err == nil {
		t.Fatal("expected error when stats path is a directory")
	}
}

func TestSaveUsageStats(t *testing.T) {
	dir := useTempConfig(t)

	if err := saveUsageStats([]MenuOption{{id: "json", usageCount: 2}, {id: "yaml"}}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := readStats(t, dir); len(got) != 2 || got["json"] != 2 || got["yaml"] != 0 {
		t.Fatalf("expected json 2 and yaml 0, got %v", got)
	}
	assertPerm(t, filepath.Join(dir, "devtui"), 0o750)
	assertPerm(t, filepath.Join(dir, "devtui", "usage_stats.json"), 0o600)
}

func TestSaveUsageStatsReturnsError(t *testing.T) {
	dir := useTempConfig(t)
	if err := os.WriteFile(filepath.Join(dir, "devtui"), nil, 0o600); err != nil {
		t.Fatal(err)
	}

	if err := saveUsageStats(nil); err == nil {
		t.Fatal("expected error when config path is a file")
	}
}

func TestEnterOpensToolAndCountsUsage(t *testing.T) {
	dir := useTempConfig(t)
	common := &ui.CommonModel{Width: 80, Height: 24}
	m := newListModel(common)

	for range 3 {
		m.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyDown}))
	}
	next, cmd := m.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyEnter}))

	if _, ok := next.(*listModel); ok {
		t.Fatal("expected the selected tool to open")
	}
	if !sendsWindowSize(cmd, tea.WindowSizeMsg{Width: 80, Height: 24}) {
		t.Fatal("expected the tool to receive the window size")
	}
	if common.LastSelectedItem != 3 {
		t.Fatalf("expected last selected item 3, got %d", common.LastSelectedItem)
	}

	for id, count := range readStats(t, dir) {
		want := 0
		if id == "numbers" {
			want = 1
		}
		if count != want {
			t.Fatalf("expected %s used %d times, got %d", id, want, count)
		}
	}
}

func TestEnterShowsSaveError(t *testing.T) {
	dir := useTempConfig(t)
	m := newListModel(&ui.CommonModel{Width: 80, Height: 24})
	if err := os.WriteFile(filepath.Join(dir, "devtui"), nil, 0o600); err != nil {
		t.Fatal(err)
	}

	m.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyEnter}))

	lines := strings.Split(m.View().Content, "\n")
	if len(lines) != 5 || strings.TrimSpace(lines[1]) != "" || !strings.HasPrefix(lines[2], "  Failed to save usage stats:") {
		t.Fatalf("expected padded save error, got %q", m.View().Content)
	}
}

func TestListWindowSizeMsgResizesList(t *testing.T) {
	useTempConfig(t)
	common := &ui.CommonModel{}
	m := newListModel(common)

	_, cmd := m.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	if cmd != nil {
		t.Fatal("expected no command")
	}
	if common.Width != 100 || common.Height != 40 || m.list.Width() != 100 || m.list.Height() != 40 {
		t.Fatalf("expected 100x40, got common %dx%d list %dx%d", common.Width, common.Height, m.list.Width(), m.list.Height())
	}
}

func TestListQuitKeys(t *testing.T) {
	useTempConfig(t)
	m := newListModel(&ui.CommonModel{})

	for _, key := range []tea.KeyPressMsg{
		tea.KeyPressMsg(tea.Key{Code: 'q', Text: "q"}),
		tea.KeyPressMsg(tea.Key{Code: 'c', Mod: tea.ModCtrl}),
	} {
		_, cmd := m.Update(key)
		if _, ok := cmd().(tea.QuitMsg); !ok {
			t.Fatalf("expected quit for %q", key.String())
		}
	}
}

func TestRefreshOrder(t *testing.T) {
	useTempConfig(t)
	m := newListModel(&ui.CommonModel{})

	m.RefreshOrder()
	options := getMenuOptions(&ui.CommonModel{})
	for i, item := range m.list.Items() {
		if got := item.(MenuOption).id; got != options[i].id {
			t.Fatalf("expected %q at position %d, got %q", options[i].id, i, got)
		}
	}

	m.items[0].usageCount = 2
	m.items[24].usageCount = 1
	m.RefreshOrder()
	first := m.list.Items()[0].(MenuOption).id
	second := m.list.Items()[1].(MenuOption).id
	if first != options[0].id || second != options[24].id {
		t.Fatalf("expected %q then %q, got %q then %q", options[0].id, options[24].id, first, second)
	}
}

// useTempConfig points the xdg config home at a temporary directory, so tests
// never touch the real usage stats.
func useTempConfig(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()
	original := xdg.ConfigHome
	xdg.ConfigHome = dir
	t.Cleanup(func() { xdg.ConfigHome = original })

	return dir
}

func writeStats(t *testing.T, dir, content string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Join(dir, "devtui"), 0o750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "devtui", "usage_stats.json"), []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func readStats(t *testing.T, dir string) map[string]int {
	t.Helper()

	data, err := os.ReadFile(filepath.Join(dir, "devtui", "usage_stats.json"))
	if err != nil {
		t.Fatal(err)
	}

	var stats map[string]int
	if err := json.Unmarshal(data, &stats); err != nil {
		t.Fatal(err)
	}
	return stats
}

func assertPerm(t *testing.T, path string, want os.FileMode) {
	t.Helper()

	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if got := info.Mode().Perm(); got != want {
		t.Fatalf("expected %s to have mode %v, got %v", path, want, got)
	}
}

func captureStderr(t *testing.T, fn func()) string {
	t.Helper()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	original := os.Stderr
	os.Stderr = w
	fn()
	os.Stderr = original
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}

	out, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	return string(out)
}

func stripANSI(s string) string {
	return ansiPattern.ReplaceAllString(s, "")
}

func sendsWindowSize(cmd tea.Cmd, want tea.WindowSizeMsg) bool {
	batch, ok := cmd().(tea.BatchMsg)
	if !ok {
		return false
	}
	for _, c := range batch {
		done := make(chan tea.Msg, 1)
		go func() { done <- c() }()
		select {
		case msg := <-done:
			if msg == want {
				return true
			}
		case <-time.After(200 * time.Millisecond):
		}
	}
	return false
}
