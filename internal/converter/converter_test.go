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

func TestConvertValidInput(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		convert func(string) (string, error)
		input   string
		want    string
	}{
		{"yaml to json", YAMLToJSON, "name: x\nn: 1\n", "{\n  \"n\": 1,\n  \"name\": \"x\"\n}\n"},
		{"json to yaml", JSONToYAML, `{"name":"x","n":[1,2]}`, "\"n\":\n    - 1\n    - 2\nname: x\n"},
		{"toml to json", TOMLToJSON, "name = \"x\"\n[t]\nn = 1\n", "{\n  \"name\": \"x\",\n  \"t\": {\n    \"n\": 1\n  }\n}\n"},
		{"json to toml", JSONToTOML, `{"name":"x","t":{"n":1}}`, "name = 'x'\n\n[t]\nn = 1.0\n"},
		{"xml to json", XMLToJSON, `<root><name>x</name></root>`, "{\n  \"root\": {\n    \"name\": \"x\"\n  }\n}"},
		{"json to xml", JSONToXML, `{"root":{"name":"x","n":1}}`, "<root>\n  <n>1</n>\n  <name>x</name>\n</root>"},
		{"yaml to toml", YAMLToTOML, "name: x\nt:\n  n: 1\n", "name = 'x'\n\n[t]\nn = 1.0\n"},
		{"toml to yaml", TOMLToYAML, "name = \"x\"\n[t]\nn = 1\n", "name: x\nt:\n    \"n\": 1\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := tt.convert(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestConvertInvalidInputReturnsError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		convert func(string) (string, error)
		input   string
	}{
		{"yaml to json syntax", YAMLToJSON, "{not: valid"},
		{"yaml to json non-string key", YAMLToJSON, "1: a"},
		{"json to yaml syntax", JSONToYAML, "{"},
		{"toml to json syntax", TOMLToJSON, "name = "},
		{"toml to json nan", TOMLToJSON, "x = nan"},
		{"json to toml syntax", JSONToTOML, "{"},
		{"xml to json syntax", XMLToJSON, "<root>"},
		{"json to xml syntax", JSONToXML, "{"},
		{"json to xml without element", JSONToXML, `{"?x":{"-v":"?>"}}`},
		{"yaml to toml syntax", YAMLToTOML, "{not: valid"},
		{"toml to yaml syntax", TOMLToYAML, "name = "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := tt.convert(tt.input)
			if err == nil {
				t.Fatalf("expected error, got %q", got)
			}
			if got != "" {
				t.Fatalf("expected empty output on error, got %q", got)
			}
		})
	}
}

func TestXMLHasElement(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    bool
		wantErr bool
	}{
		{"element", "<a/>", true, false},
		{"single character text", "x", false, false},
		{"empty input", "", false, false},
		{"processing instruction only", `<?x v="?>"/>`, false, false},
		{"invalid first token", "<0/>", false, true},
		{"unclosed element", "<a>", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := xmlHasElement([]byte(tt.input))
			if (err != nil) != tt.wantErr {
				t.Fatalf("xmlHasElement(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if got != tt.want {
				t.Fatalf("xmlHasElement(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
