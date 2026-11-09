package json

import "testing"

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
