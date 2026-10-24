package structgen

import (
	"strings"
	"testing"
)

func TestJSONToGoStruct(t *testing.T) {
	t.Parallel()

	got, err := JSONToGoStruct(strings.NewReader(`{"name":"x","age":1,"key with spaces":true}`))
	if err != nil {
		t.Fatalf("JSONToGoStruct returned error: %v", err)
	}

	for _, want := range []string{
		"type Root struct {\n",
		"\tAge  int    `json:\"age\"`",
		"\tName string `json:\"name\"`",
		"\t// \"key with spaces\" cannot be unmarshalled into a struct field",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("expected output to contain %q, got %q", want, got)
		}
	}
}

func TestYAMLToGoStruct(t *testing.T) {
	t.Parallel()

	got, err := YAMLToGoStruct(strings.NewReader("name: x\nage: 1\n"))
	if err != nil {
		t.Fatalf("YAMLToGoStruct returned error: %v", err)
	}

	want := "package main\n\ntype Root struct {\n\tAge  int    `yaml:\"age\"`\n\tName string `yaml:\"name\"`\n}\n"
	if got != want {
		t.Fatalf("YAMLToGoStruct() = %q, want %q", got, want)
	}
}

func TestJSONToGoStructInvalidInputReturnsError(t *testing.T) {
	t.Parallel()

	if _, err := JSONToGoStruct(strings.NewReader(`{`)); err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestYAMLToGoStructInvalidInputReturnsError(t *testing.T) {
	t.Parallel()

	if _, err := YAMLToGoStruct(strings.NewReader("a: [")); err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}
