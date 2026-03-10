package csv2json

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func addCSVSeedsFromFiles(f *testing.F, fileNames ...string) {
	for _, fileName := range fileNames {
		content, err := os.ReadFile(filepath.Join("../../testdata", fileName))
		if err != nil {
			continue
		}

		f.Add(string(content))
	}
}

func convertNoPanic(t *testing.T, input string) (result string, err error) {
	t.Helper()

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Convert panicked: %v", r)
		}
	}()

	return Convert(input)
}

func TestConvertConflictingHeadersReturnError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "scalar_then_object",
			input: "a,a.b\n1,2\n",
		},
		{
			name:  "object_then_scalar",
			input: "a.b,a\n1,2\n",
		},
		{
			name:  "array_then_object",
			input: "a[0],a.b\n1,2\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := convertNoPanic(t, tt.input)
			if err == nil {
				t.Fatal("expected error for conflicting CSV headers")
			}
		})
	}
}

func TestConvertHighArrayIndexDoesNotPanic(t *testing.T) {
	t.Parallel()

	output, err := convertNoPanic(t, "items[2].name\nhello\n")
	if err != nil {
		t.Fatalf("Convert should not fail for sparse array indexes: %v", err)
	}

	if !json.Valid([]byte(output)) {
		t.Fatalf("expected valid JSON output, got: %s", output)
	}

	if !strings.Contains(output, "hello") {
		t.Fatalf("expected output to include converted value, got: %s", output)
	}
}

func TestConvertMalformedArrayHeaderDoesNotPanic(t *testing.T) {
	t.Parallel()

	output, err := convertNoPanic(t, "][\n0\n")
	if err != nil {
		return
	}

	if !json.Valid([]byte(output)) {
		t.Fatalf("expected valid JSON output, got: %s", output)
	}
}

func FuzzConvertDoesNotPanic(f *testing.F) {
	f.Add("name,age\nAlice,30\n")
	f.Add("a,a.b\n1,2\n")
	f.Add("items[2].name\nhello\n")
	f.Add("")
	addCSVSeedsFromFiles(f, "example.csv")

	f.Fuzz(func(t *testing.T, input string) {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("Convert panicked for input %q: %v", input, r)
			}
		}()

		_, _ = Convert(input)
	})
}
