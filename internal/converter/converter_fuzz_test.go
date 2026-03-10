package converter

import (
	"encoding/json"
	"testing"

	"github.com/clbanning/mxj/v2"
	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
)

func FuzzYAMLToJSON(f *testing.F) {
	f.Add("name: Alice\nage: 30\n")
	f.Add("items:\n  - one\n  - two\n")
	f.Add("{not: valid")
	f.Add("")

	f.Fuzz(func(t *testing.T, input string) {
		if len(input) > 4096 {
			t.Skip()
		}

		output, err := YAMLToJSON(input)
		if err != nil {
			return
		}

		if !json.Valid([]byte(output)) {
			t.Fatalf("YAMLToJSON returned invalid JSON: %q", output)
		}
	})
}

func FuzzJSONToYAML(f *testing.F) {
	f.Add(`{"name":"Alice","age":30}`)
	f.Add(`[1,2,3]`)
	f.Add(`{invalid`)
	f.Add("")

	f.Fuzz(func(t *testing.T, input string) {
		if len(input) > 4096 {
			t.Skip()
		}

		output, err := JSONToYAML(input)
		if err != nil {
			return
		}

		var data any
		if err := yaml.Unmarshal([]byte(output), &data); err != nil {
			t.Fatalf("JSONToYAML returned invalid YAML: %v", err)
		}
	})
}

func FuzzTOMLToJSON(f *testing.F) {
	f.Add("name = \"Alice\"\nage = 30\n")
	f.Add("[user]\nname = \"Bob\"\n")
	f.Add("[[items]]\nname = \"one\"\n")
	f.Add("")

	f.Fuzz(func(t *testing.T, input string) {
		if len(input) > 4096 {
			t.Skip()
		}

		output, err := TOMLToJSON(input)
		if err != nil {
			return
		}

		if !json.Valid([]byte(output)) {
			t.Fatalf("TOMLToJSON returned invalid JSON: %q", output)
		}
	})
}

func FuzzJSONToTOML(f *testing.F) {
	f.Add(`{"name":"Alice","age":30}`)
	f.Add(`{"items":[{"name":"one"}]}`)
	f.Add(`[]`)
	f.Add("")

	f.Fuzz(func(t *testing.T, input string) {
		if len(input) > 4096 {
			t.Skip()
		}

		output, err := JSONToTOML(input)
		if err != nil {
			return
		}

		var data any
		if err := toml.Unmarshal([]byte(output), &data); err != nil {
			t.Fatalf("JSONToTOML returned invalid TOML: %v", err)
		}
	})
}

func FuzzXMLToJSON(f *testing.F) {
	f.Add("<root><name>Alice</name></root>")
	f.Add("<root attr=\"x\"><item>1</item></root>")
	f.Add("<root>")
	f.Add("")

	f.Fuzz(func(t *testing.T, input string) {
		if len(input) > 4096 {
			t.Skip()
		}

		output, err := XMLToJSON(input)
		if err != nil {
			return
		}

		if !json.Valid([]byte(output)) {
			t.Fatalf("XMLToJSON returned invalid JSON: %q", output)
		}
	})
}

func FuzzJSONToXML(f *testing.F) {
	f.Add(`{"root":{"name":"Alice"}}`)
	f.Add(`{"root":{"items":[1,2,3]}}`)
	f.Add(`[]`)
	f.Add("")

	f.Fuzz(func(t *testing.T, input string) {
		if len(input) > 4096 {
			t.Skip()
		}

		output, err := JSONToXML(input)
		if err != nil {
			return
		}

		if _, err := mxj.NewMapXml([]byte(output)); err != nil {
			t.Fatalf("JSONToXML returned invalid XML: %v", err)
		}
	})
}

func FuzzYAMLToTOML(f *testing.F) {
	f.Add("name: Alice\nage: 30\n")
	f.Add("items:\n  - one\n  - two\n")
	f.Add("{not: valid")
	f.Add("")

	f.Fuzz(func(t *testing.T, input string) {
		if len(input) > 4096 {
			t.Skip()
		}

		output, err := YAMLToTOML(input)
		if err != nil {
			return
		}

		var data any
		if err := toml.Unmarshal([]byte(output), &data); err != nil {
			t.Fatalf("YAMLToTOML returned invalid TOML: %v", err)
		}
	})
}

func FuzzTOMLToYAML(f *testing.F) {
	f.Add("name = \"Alice\"\nage = 30\n")
	f.Add("[user]\nname = \"Bob\"\n")
	f.Add("[[items]]\nname = \"one\"\n")
	f.Add("")

	f.Fuzz(func(t *testing.T, input string) {
		if len(input) > 4096 {
			t.Skip()
		}

		output, err := TOMLToYAML(input)
		if err != nil {
			return
		}

		var data any
		if err := yaml.Unmarshal([]byte(output), &data); err != nil {
			t.Fatalf("TOMLToYAML returned invalid YAML: %v", err)
		}
	})
}
