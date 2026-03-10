package yamlfmt

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

func addYAMLSeedsFromFiles(f *testing.F, fileNames ...string) {
	for _, fileName := range fileNames {
		content, err := os.ReadFile(filepath.Join("../../testdata", fileName))
		if err != nil {
			continue
		}

		f.Add(string(content))
	}
}

func FuzzFormatYAML(f *testing.F) {
	f.Add("name: Alice\nage: 30\n")
	f.Add("items:\n  - one\n  - two\n")
	f.Add("{not: valid")
	f.Add("")
	addYAMLSeedsFromFiles(f, "example.yaml", "nested.yaml")

	f.Fuzz(func(t *testing.T, input string) {
		if len(input) > 4096 {
			t.Skip()
		}

		output, err := Format(input)
		if err != nil {
			return
		}

		var data any
		if err := yaml.Unmarshal([]byte(output), &data); err != nil {
			t.Fatalf("Format returned invalid YAML: %v", err)
		}
	})
}
