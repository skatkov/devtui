package yamlfmt

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func FuzzFormatYAML(f *testing.F) {
	f.Add("name: Alice\nage: 30\n")
	f.Add("items:\n  - one\n  - two\n")
	f.Add("{not: valid")
	f.Add("")

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
