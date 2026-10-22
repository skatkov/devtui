package csv2json

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
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
		{
			name:  "array_scalar_then_object",
			input: "a[0],a[0].b\n1,2\n",
		},
		{
			name:  "array_object_then_scalar",
			input: "a[0].b,a[0]\n1,2\n",
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

func TestConvertArrayIndexes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "indexes_in_column_order",
			input: "a[0],a[1]\nx,y\n",
			want:  `[{"a":["x","y"]}]`,
		},
		{
			name:  "indexes_out_of_column_order",
			input: "a[1],a[0]\nx,y\n",
			want:  `[{"a":["y","x"]}]`,
		},
		{
			name:  "sparse_index_pads_with_null",
			input: "a[2]\nx\n",
			want:  `[{"a":[null,null,"x"]}]`,
		},
		{
			name:  "negative_index_is_a_plain_key",
			input: "a[-2].b\nx\n",
			want:  `[{"a[-2]":{"b":"x"}}]`,
		},
		{
			name:  "minus_one_index_is_a_plain_key",
			input: "a[-1]\nx\n",
			want:  `[{"a[-1]":"x"}]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			output, err := convertNoPanic(t, tt.input)
			if err != nil {
				t.Fatalf("Convert(%q) returned error: %v", tt.input, err)
			}

			var got, want any
			if err := json.Unmarshal([]byte(output), &got); err != nil {
				t.Fatalf("invalid JSON output %q: %v", output, err)
			}
			if err := json.Unmarshal([]byte(tt.want), &want); err != nil {
				t.Fatalf("invalid want JSON %q: %v", tt.want, err)
			}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("Convert(%q) = %s, want %s", tt.input, output, tt.want)
			}
		})
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
	f.Add("a[-2].b\nx\n")
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
