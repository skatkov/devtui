package json

import (
	"strings"
	"testing"
)

func TestFormatJSONOnlyChangesWhitespace(t *testing.T) {
	input := `{"id":1234567890123456789,"zeta":1,"alpha":2.50,"u":"a&b<c>"}`
	want := `{
  "id": 1234567890123456789,
  "zeta": 1,
  "alpha": 2.50,
  "u": "a&b<c>"
}
`
	if got := FormatJSON(input); got != want {
		t.Errorf("FormatJSON() = %q, want %q", got, want)
	}
}

func TestFormatJSONReturnsInvalidInputUnchanged(t *testing.T) {
	input := `{"a":`
	if got := FormatJSON(input); got != input {
		t.Errorf("FormatJSON() = %q, want %q", got, input)
	}
}

func TestFormatJSONNewlineTerminatedInputIsIdempotent(t *testing.T) {
	for _, input := range []string{"{}", "{}\n", "{}\n\n", "{\"a\":1}\r\n", "  [1,2] \t\n"} {
		once := FormatJSON(input)
		if strings.HasSuffix(once, "\n\n") || !strings.HasSuffix(once, "\n") {
			t.Errorf("FormatJSON(%q) = %q, want exactly one trailing newline", input, once)
		}
		if twice := FormatJSON(once); twice != once {
			t.Errorf("FormatJSON(FormatJSON(%q)) = %q, want %q", input, twice, once)
		}
	}
	if got, want := FormatJSON("{}\n"), "{}\n"; got != want {
		t.Errorf("FormatJSON(%q) = %q, want %q", "{}\n", got, want)
	}
}
