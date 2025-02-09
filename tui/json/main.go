package main

import (
	"bytes"
	"encoding/json"
	"log"

	"github.com/charmbracelet/huh"
)

func main() {
	var json string
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewText().
				Title("Unformatted JSON").
				EditorExtension("json").
				Placeholder("Paste JSON here or open editor with Ctrl+E").
				Value(&json),
		),

		huh.NewGroup(huh.NewNote().Height(20).Title("Formatted JSON").
			DescriptionFunc(func() string {
				fmd := formatJSON(json)
				return fmd
			}, &json)),
	).WithLayout(huh.LayoutColumns(2)).
		Run()

	if err != nil {
		log.Fatal(err)
	}
}

func formatJSON(input string) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(input), "", "  "); err != nil {
		return input // Return original if not valid JSON
	}
	return prettyJSON.String()
}
