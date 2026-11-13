package converter

import (
	"strings"
	"testing"
)

func TestJSONToXMLInvalidElementName(t *testing.T) {
	t.Parallel()

	_, err := JSONToXML(`{"0":""}`)
	if err == nil {
		t.Fatal("expected error for invalid XML element name")
	}
}

func TestJSONToTOMLTopLevelArrayReturnsError(t *testing.T) {
	t.Parallel()

	_, err := JSONToTOML(`[]`)
	if err == nil {
		t.Fatal("expected error for top-level array TOML encoding")
	}
}

func TestYAMLToTOMLScalarReturnsError(t *testing.T) {
	t.Parallel()

	_, err := YAMLToTOML("00")
	if err == nil {
		t.Fatal("expected error for scalar TOML encoding")
	}
}

func TestJSONToTOMLNumberOverflowReturnsError(t *testing.T) {
	t.Parallel()

	if out, err := JSONToTOML(`{"value":1e400}`); err == nil {
		t.Fatalf("expected error for out-of-range number, got %q", out)
	}
}

func TestJSONToTOMLRejectsTrailingData(t *testing.T) {
	t.Parallel()

	for _, input := range []string{`{"port":8080} {"port":9090}`, `{"port":8080} trailing garbage`} {
		if out, err := JSONToTOML(input); err == nil {
			t.Errorf("JSONToTOML(%q) = %q, want error", input, out)
		}
	}
}

func TestYAMLToTOMLKeepsNumericTypes(t *testing.T) {
	t.Parallel()

	got, err := YAMLToTOML("int: 1\nfloat: 1.0\n")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got, "int = 1\n") || !strings.Contains(got, "float = 1.0\n") {
		t.Errorf("YAMLToTOML() = %q, want int = 1 and float = 1.0", got)
	}
}
