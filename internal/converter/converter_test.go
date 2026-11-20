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

func TestConvertersKeepLargeIntegers(t *testing.T) {
	t.Parallel()

	for name, convert := range map[string]func(string) (string, error){
		"JSONToYAML": JSONToYAML,
		"JSONToXML":  JSONToXML,
	} {
		out, err := convert(`{"ts":1700000000000}`)
		if err != nil {
			t.Fatalf("%s() error = %v", name, err)
		}
		if !strings.Contains(out, "1700000000000") {
			t.Errorf("%s() = %q, want 1700000000000", name, out)
		}
	}
}

func TestConvertersRejectOutOfRangeNumbers(t *testing.T) {
	t.Parallel()

	for name, convert := range map[string]func(string) (string, error){
		"JSONToYAML": JSONToYAML,
		"JSONToXML":  JSONToXML,
	} {
		for _, input := range []string{`{"value":1e400}`, `{"a":{"b":[1,-1e400]}}`, `[[1e400]]`} {
			if out, err := convert(input); err == nil {
				t.Errorf("%s(%q) = %q, want error", name, input, out)
			}
		}
	}
}

func TestTOMLToYAMLKeepsNumericTypes(t *testing.T) {
	t.Parallel()

	out, err := TOMLToYAML("ratio = 1000000.0\nid = 9007199254740993\n")
	if err != nil {
		t.Fatal(err)
	}

	var got map[string]any
	if err := yaml.Unmarshal([]byte(out), &got); err != nil {
		t.Fatal(err)
	}
	if ratio, ok := got["ratio"].(float64); !ok || ratio != 1000000 {
		t.Errorf("ratio = %#v, want float64 1000000 (output %q)", got["ratio"], out)
	}
	if id, ok := got["id"].(int); !ok || id != 9007199254740993 {
		t.Errorf("id = %#v, want int 9007199254740993 (output %q)", got["id"], out)
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
