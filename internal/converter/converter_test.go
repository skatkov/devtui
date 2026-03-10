package converter

import "testing"

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
