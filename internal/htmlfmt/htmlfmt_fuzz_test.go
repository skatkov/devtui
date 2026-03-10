package htmlfmt

import (
	"os"
	"path/filepath"
	"testing"
)

func addHTMLSeedsFromFiles(f *testing.F, fileNames ...string) {
	for _, fileName := range fileNames {
		content, err := os.ReadFile(filepath.Join("../../testdata", fileName))
		if err != nil {
			continue
		}

		f.Add(string(content))
	}
}

func FuzzFormatHTML(f *testing.F) {
	f.Add("<html><body><h1>Hello</h1></body></html>")
	f.Add("<div><span>missing end")
	f.Add("")
	addHTMLSeedsFromFiles(f, "html-with-urls.html")

	f.Fuzz(func(t *testing.T, input string) {
		if len(input) > 4096 {
			t.Skip()
		}

		_ = Format(input)
	})
}
