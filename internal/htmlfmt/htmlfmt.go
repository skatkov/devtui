package htmlfmt

import "github.com/yosssi/gohtml"

// Format formats HTML content with consistent indentation.
func Format(content string) string {
	return gohtml.Format(content)
}
