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

func TestFormatReturnsErrorForInvalidFormattedYAML(t *testing.T) {
	t.Parallel()

	_, err := Format("0: \n0.:")
	if err == nil {
		t.Fatal("expected error for YAML that formats into invalid output")
	}
}
