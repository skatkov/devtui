package generator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/skatkov/devtui/cmd"
)

func TestGenerate(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "tui", "json", "main.go"), `const Title = "JSON Formatter"`)
	if err := os.MkdirAll(filepath.Join(dir, "site", "cli"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Chdir(dir)

	Generate(cmd.GetRootCmd())

	for _, file := range []string{"site/cli/base64.md", "site/tui/index.md", "site/tui/json.md"} {
		if _, err := os.Stat(file); err != nil {
			t.Errorf("expected %s to be generated: %v", file, err)
		}
	}
}

func TestProjectPath(t *testing.T) {
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, "site"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Chdir(dir)

	tests := []struct {
		name string
		want string
	}{
		{"site", "site"},
		{"tui", filepath.Join("..", "tui")},
	}
	for _, tt := range tests {
		if got := projectPath(tt.name); got != tt.want {
			t.Errorf("projectPath(%q) = %q, want %q", tt.name, got, tt.want)
		}
	}
}
