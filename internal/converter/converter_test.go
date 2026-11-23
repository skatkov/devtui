package converter

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestJSONToXMLInvalidElementName(t *testing.T) {
	t.Parallel()

	_, err := JSONToXML(`{"0":""}`)
	if err == nil {
		t.Fatal("expected error for invalid XML element name")
	}
}

func TestJSONToXMLEscapesText(t *testing.T) {
	t.Parallel()

	out, err := JSONToXML(`{"name":"A & B","title":"<t>"}`)
	if err != nil {
		t.Fatalf("JSONToXML() error = %v", err)
	}
	for _, want := range []string{"<name>A &amp; B</name>", "<title>&lt;t&gt;</title>"} {
		if !strings.Contains(out, want) {
			t.Errorf("JSONToXML() = %q, want it to contain %q", out, want)
		}
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

func TestTOMLKeepsLargeIntegers(t *testing.T) {
	t.Parallel()

	// 9007199254740993 is not exactly representable as float64.
	const want = "id = 9007199254740993\n"

	got, err := JSONToTOML(`{"id":9007199254740993}`)
	if err != nil || got != want {
		t.Errorf("JSONToTOML() = %q, %v; want %q", got, err, want)
	}

	got, err = YAMLToTOML("id: 9007199254740993\n")
	if err != nil || got != want {
		t.Errorf("YAMLToTOML() = %q, %v; want %q", got, err, want)
	}
}

func TestTOMLToYAMLKeepsScalarTypes(t *testing.T) {
	t.Parallel()

	out, err := TOMLToYAML("values = [9007199254740993, 1000000.0, 1.0, 2026-11-20, 2026-11-20T12:00:00Z]")
	if err != nil {
		t.Fatal(err)
	}

	var decoded struct {
		Values []yaml.Node
	}
	if err := yaml.Unmarshal([]byte(out), &decoded); err != nil {
		t.Fatal(err)
	}
	want := []struct{ tag, value string }{
		{"!!int", "9007199254740993"},
		{"!!float", "1.0e+06"},
		{"!!float", "1.0"},
		{"!!str", "2026-11-20"},
		{"!!str", "2026-11-20T12:00:00Z"},
	}
	if len(decoded.Values) != len(want) {
		t.Fatalf("TOMLToYAML() = %q, want %d values", out, len(want))
	}
	for i, node := range decoded.Values {
		if node.ShortTag() != want[i].tag || node.Value != want[i].value {
			t.Errorf("value %d = %s %q, want %s %q", i, node.ShortTag(), node.Value, want[i].tag, want[i].value)
		}
	}
}

func TestTOMLToYAMLRejectsNonFiniteNumbers(t *testing.T) {
	t.Parallel()

	for _, input := range []string{"value = inf", "values = [nan]"} {
		if out, err := TOMLToYAML(input); err == nil {
			t.Errorf("TOMLToYAML(%q) = %q, want error", input, out)
		}
	}
}
