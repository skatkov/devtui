package generator

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

const helpView = `package demo

func (m Model) helpView() string {
	col1 := []string{
		"c              copy text",
		"  q/ctrl+c       quit  ",
		"e              top",
		"x              up",
		"↑              scroll up",
		"↓              scroll down",
		"k              move up",
		"j              move down",
		"single spaced entry",
		"v    paste a  b",
	}
	col2 := []string{
		"w              ignored entry",
	}
	return ""
}
`

func writeFile(t *testing.T, path, content string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestExtractKeyBindings(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		content string
		want    []KeyBinding
	}{
		{"no help view", "package demo", nil},
		{"help view", helpView, []KeyBinding{
			{"c", "copy text"},
			{"q/ctrl+c", "quit"},
			{"v", "paste a  b"},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := extractKeyBindings(tt.content); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractKeyBindings() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractModuleInfo(t *testing.T) {
	t.Parallel()

	single := `const Title = "JSON Formatter"
col1 := []string{"c    copy text"}`
	tests := []struct {
		name  string
		files map[string]string
		want  func(dir string) *TUIModule
	}{
		{
			"title and bindings",
			map[string]string{"main.go": single, "README.md": `const Title = "Ignored"`},
			func(dir string) *TUIModule {
				return &TUIModule{
					Name:        "demo",
					Title:       "JSON Formatter",
					Description: "A text formatting tool that prettifies and standardizes code or data format.",
					KeyBindings: []KeyBinding{{"c", "copy text"}},
					FilePath:    filepath.Join(dir, "main.go"),
				}
			},
		},
		{
			"bindings before title",
			map[string]string{"help.go": helpView, "main.go": `const Title = "Base64 Encoder"`},
			func(dir string) *TUIModule {
				return &TUIModule{
					Name:        "demo",
					Title:       "Base64 Encoder",
					Description: "An encoding tool that transforms text into encoded format.",
					KeyBindings: []KeyBinding{{"c", "copy text"}, {"q/ctrl+c", "quit"}, {"v", "paste a  b"}},
					FilePath:    filepath.Join(dir, "help.go"),
				}
			},
		},
		{
			"bindings only",
			map[string]string{"help.go": helpView},
			func(dir string) *TUIModule {
				return &TUIModule{
					Name:        "demo",
					KeyBindings: []KeyBinding{{"c", "copy text"}, {"q/ctrl+c", "quit"}, {"v", "paste a  b"}},
					FilePath:    filepath.Join(dir, "help.go"),
				}
			},
		},
		{
			"no title or bindings",
			map[string]string{"util.go": "package demo"},
			func(string) *TUIModule { return nil },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			dir := t.TempDir()
			for name, content := range tt.files {
				writeFile(t, filepath.Join(dir, name), content)
			}

			got, err := extractModuleInfo(dir, "demo")
			if err != nil {
				t.Fatal(err)
			}
			if want := tt.want(dir); !reflect.DeepEqual(got, want) {
				t.Errorf("extractModuleInfo() = %+v, want %+v", got, want)
			}
		})
	}
}

func TestExtractModuleInfoMissingDirReturnsError(t *testing.T) {
	t.Parallel()

	if _, err := extractModuleInfo(filepath.Join(t.TempDir(), "missing"), "missing"); err == nil {
		t.Fatal("expected an error, got nil")
	}
}

func TestExtractTUIModules(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "a", "main.go"), `const Title = "Zulu Tool"`)
	writeFile(t, filepath.Join(dir, "b", "main.go"), `const Title = "Alpha Tool"`)
	writeFile(t, filepath.Join(dir, "c", "main.go"), `const Title = "Mid Tool"`)
	writeFile(t, filepath.Join(dir, "d", "README.md"), "empty module")
	writeFile(t, filepath.Join(dir, "notes.go"), `const Title = "Not A Module"`)

	modules, err := extractTUIModules(dir)
	if err != nil {
		t.Fatal(err)
	}

	got := make([]string, 0, len(modules))
	for _, module := range modules {
		got = append(got, module.Name+":"+module.Title)
	}
	want := []string{"b:Alpha Tool", "c:Mid Tool", "a:Zulu Tool"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("extractTUIModules() = %v, want %v", got, want)
	}
}

func TestExtractTUIModulesMissingDirReturnsError(t *testing.T) {
	t.Parallel()

	if _, err := extractTUIModules(filepath.Join(t.TempDir(), "missing")); err == nil {
		t.Fatal("expected an error, got nil")
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(content)
}

func TestGenerateTUIDocumentationWritesFiles(t *testing.T) {
	t.Parallel()

	sitePath := filepath.Join(t.TempDir(), "site", "tui")
	modules := []TUIModule{
		{Name: "json", Title: "JSON Formatter", KeyBindings: []KeyBinding{{"c", "copy text"}}},
		{Name: "cron", Title: "Cron Parser"},
	}
	if err := generateTUIDocumentation(modules, sitePath); err != nil {
		t.Fatal(err)
	}

	usage := "## Usage\n\n" +
		"1. Run `devtui` to open the main menu\n" +
		"2. Select \"%s\" from the list\n" +
		"3. Use the key bindings below to interact with the tool\n" +
		"4. Press `q` or `Ctrl+C` to return to the main menu\n\n" +
		"## Key Bindings\n\n"
	tests := []struct {
		file string
		want string
	}{
		{
			"json.md",
			"---\ntitle: JSON Formatter\nparent: TUI\n---\n\n# JSON Formatter\n\n" +
				strings.Replace(usage, "%s", "JSON Formatter", 1) +
				"| Key | Action |\n|-----|--------|\n| `c` | copy text |\n\n\n",
		},
		{
			"cron.md",
			"---\ntitle: Cron Parser\nparent: TUI\n---\n\n# Cron Parser\n\n" +
				strings.Replace(usage, "%s", "Cron Parser", 1) +
				"Standard key bindings apply (see main TUI documentation).\n\n\n",
		},
	}
	for _, tt := range tests {
		if got := readFile(t, filepath.Join(sitePath, tt.file)); got != tt.want {
			t.Errorf("%s = %q, want %q", tt.file, got, tt.want)
		}
	}

	if index := readFile(t, filepath.Join(sitePath, "index.md")); !strings.HasPrefix(index, "---\ntitle: TUI\nnav_order: 4\n---\n") {
		t.Errorf("expected index front matter, got %q", index)
	}
}

func TestGenerateTUIDocumentation(t *testing.T) {
	t.Parallel()

	tuiDir := t.TempDir()
	writeFile(t, filepath.Join(tuiDir, "cron", "main.go"), `const Title = "Cron Parser"`)
	sitePath := filepath.Join(t.TempDir(), "tui")

	if err := GenerateTUIDocumentation(tuiDir, sitePath); err != nil {
		t.Fatal(err)
	}
	for _, file := range []string{"index.md", "cron.md"} {
		if _, err := os.Stat(filepath.Join(sitePath, file)); err != nil {
			t.Errorf("expected %s to be generated: %v", file, err)
		}
	}
}

func TestGenerateTUIDocumentationNoModules(t *testing.T) {
	t.Parallel()

	sitePath := filepath.Join(t.TempDir(), "tui")
	if err := GenerateTUIDocumentation(t.TempDir(), sitePath); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(sitePath); !os.IsNotExist(err) {
		t.Errorf("expected no site directory, got %v", err)
	}
}

func TestGenerateTUIDocumentationMissingDirReturnsError(t *testing.T) {
	t.Parallel()

	err := GenerateTUIDocumentation(filepath.Join(t.TempDir(), "missing"), t.TempDir())
	if err == nil || !strings.Contains(err.Error(), "failed to extract TUI modules") {
		t.Fatalf("expected extract error, got %v", err)
	}
}

func TestExtractTUIModulesFindsKeyBindingsOfEveryPager(t *testing.T) {
	t.Parallel()

	modules, err := extractTUIModules(filepath.Join("..", "..", "tui"))
	if err != nil {
		t.Fatal(err)
	}

	for _, module := range modules {
		content, err := os.ReadFile(module.FilePath)
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(content), "helpView()") && len(module.KeyBindings) == 0 {
			t.Errorf("expected key bindings for %s", module.Name)
		}
	}
}
