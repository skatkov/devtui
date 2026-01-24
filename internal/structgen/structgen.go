package structgen

import (
	"io"

	"github.com/twpayne/go-jsonstruct/v3"
)

// JSONToGoStruct converts JSON input into a Go struct definition.
func JSONToGoStruct(input io.Reader) (string, error) {
	generator := newGenerator()
	if err := generator.ObserveJSONReader(input); err != nil {
		return "", err
	}

	bytes, err := generator.Generate()
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// YAMLToGoStruct converts YAML input into a Go struct definition.
func YAMLToGoStruct(input io.Reader) (string, error) {
	generator := newGenerator()
	if err := generator.ObserveYAMLReader(input); err != nil {
		return "", err
	}

	bytes, err := generator.Generate()
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func newGenerator() *jsonstruct.Generator {
	options := []jsonstruct.GeneratorOption{
		jsonstruct.WithSkipUnparsableProperties(true),
		jsonstruct.WithStructTagName("yaml"),
		jsonstruct.WithGoFormat(true),
		jsonstruct.WithOmitEmptyTags(jsonstruct.OmitEmptyTagsAuto),
		jsonstruct.WithTypeName("Root"),
	}

	return jsonstruct.NewGenerator(options...)
}
