package yamlfmt

import (
	"strings"
	"testing"
)

func TestFormatValidYAML(t *testing.T) {
	t.Parallel()

	formatted, err := Format("name: Alice\nage: 30\n")
	if err != nil {
		t.Fatalf("Format returned error for valid YAML: %v", err)
	}

	if !strings.Contains(formatted, "name: Alice") {
		t.Fatalf("unexpected formatted YAML output: %q", formatted)
	}
}

func TestFormatKeepsDistinctNumericKeys(t *testing.T) {
	t.Parallel()

	formatted, err := Format("0: \n0.:")
	if err != nil {
		t.Fatalf("Format returned error: %v", err)
	}

	want := "0:\n0.:\n"
	if formatted != want {
		t.Fatalf("Format() = %q, want %q", formatted, want)
	}
}

func TestFormatReturnsErrorForDuplicateKeys(t *testing.T) {
	t.Parallel()

	if _, err := Format("a: 1\na: 2\n"); err == nil {
		t.Fatal("expected error for duplicate mapping keys")
	}
}

func TestFormatKeepsCommentsAndKeyOrder(t *testing.T) {
	t.Parallel()

	formatted, err := Format("# keep\nz:   1 # why\na: 2\n")
	if err != nil {
		t.Fatalf("Format returned error: %v", err)
	}

	want := "# keep\nz: 1 # why\na: 2\n"
	if formatted != want {
		t.Fatalf("Format() = %q, want %q", formatted, want)
	}
}
