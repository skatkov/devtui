package htmlfmt

import "testing"

func FuzzFormatHTML(f *testing.F) {
	f.Add("<html><body><h1>Hello</h1></body></html>")
	f.Add("<div><span>missing end")
	f.Add("")

	f.Fuzz(func(t *testing.T, input string) {
		if len(input) > 4096 {
			t.Skip()
		}

		_ = Format(input)
	})
}
