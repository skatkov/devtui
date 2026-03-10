package jsonrepair

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func addJSONSeedsFromFiles(f *testing.F, fileNames ...string) {
	for _, fileName := range fileNames {
		content, err := os.ReadFile(filepath.Join("../../testdata", fileName))
		if err != nil {
			continue
		}

		f.Add(string(content))
	}
}

func FuzzRepairJSONProducesValidJSON(f *testing.F) {
	f.Add(`{"name":"Alice"}`)
	f.Add(`{'name':'Alice'}`)
	f.Add("```json\n{'name': 'Alice'}\n```")
	f.Add("")
	addJSONSeedsFromFiles(f, "example.json", "nested.json", "json-with-urls.json")

	f.Fuzz(func(t *testing.T, input string) {
		if len(input) > 4096 {
			t.Skip()
		}

		output, err := RepairJSON(input)
		if err != nil {
			return
		}

		if !json.Valid([]byte(output)) {
			t.Fatalf("RepairJSON returned invalid JSON: %q", output)
		}
	})
}
