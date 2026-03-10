package json2toon

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hannes-sistemica/toon"
)

func addJSONSeedsFromFiles(f *testing.F, fileNames ...string) {
	for _, fileName := range fileNames {
		content, err := os.ReadFile(filepath.Join("../../testdata", fileName))
		if err != nil {
			continue
		}

		f.Add(string(content))
	}
}

func addJSONOptionSeedsFromFiles(f *testing.F, fileNames ...string) {
	for _, fileName := range fileNames {
		content, err := os.ReadFile(filepath.Join("../../testdata", fileName))
		if err != nil {
			continue
		}

		f.Add(string(content), uint8(2), "")
	}
}

func FuzzConvertNoPanic(f *testing.F) {
	f.Add(`{"name":"Alice"}`)
	f.Add(`{"items":[1,2,3]}`)
	f.Add(`{invalid`)
	f.Add("")
	addJSONSeedsFromFiles(f, "example.json", "nested.json", "json-with-urls.json")

	f.Fuzz(func(t *testing.T, input string) {
		if len(input) > 4096 {
			t.Skip()
		}

		_, _ = Convert(input)
	})
}

func FuzzConvertWithOptionsNoPanic(f *testing.F) {
	f.Add(`{"name":"Alice"}`, uint8(2), "")
	f.Add(`{"items":[{"id":1}]}`, uint8(4), "#")
	f.Add(`{invalid`, uint8(0), "!")
	addJSONOptionSeedsFromFiles(f, "example.json", "nested.json", "json-with-urls.json")

	f.Fuzz(func(t *testing.T, input string, indent uint8, marker string) {
		if len(input) > 4096 || len(marker) > 8 {
			t.Skip()
		}

		opts := toon.EncodeOptions{
			Indent:       int(indent % 8),
			Delimiter:    ",",
			LengthMarker: marker,
		}

		_, _ = ConvertWithOptions(input, opts)
	})
}
