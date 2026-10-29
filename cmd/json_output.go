package cmd

import (
	"encoding/json"
	"io"
)

var outputJSON bool

func writeJSONValue(out io.Writer, value any) error {
	encoder := json.NewEncoder(out)
	encoder.SetEscapeHTML(false)
	return encoder.Encode(value)
}
