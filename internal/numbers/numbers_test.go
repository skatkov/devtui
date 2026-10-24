package numbers

import (
	"slices"
	"testing"
)

func TestDefaultBase(t *testing.T) {
	t.Parallel()

	if got := DefaultBase(); got.Base != 10 {
		t.Fatalf("DefaultBase().Base = %d, want 10", got.Base)
	}
}

func TestParse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		base  int
		want  int64
	}{
		{"binary", "1010", 2, 10},
		{"octal", "777", 8, 511},
		{"decimal", "-42", 10, -42},
		{"hexadecimal", "ff", 16, 255},
		{"max int64", "9223372036854775807", 10, 9223372036854775807},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := Parse(tt.input, tt.base)
			if err != nil {
				t.Fatalf("Parse(%q, %d) returned error: %v", tt.input, tt.base, err)
			}
			if got != tt.want {
				t.Fatalf("Parse(%q, %d) = %d, want %d", tt.input, tt.base, got, tt.want)
			}
		})
	}
}

func TestParseUnsupportedBaseReturnsError(t *testing.T) {
	t.Parallel()

	for _, base := range []int{0, 1, 3, 9, 36} {
		got, err := Parse("1", base)
		if err == nil {
			t.Fatalf("Parse(%q, %d) expected error", "1", base)
		}
		if got != 0 {
			t.Fatalf("Parse(%q, %d) = %d, want 0", "1", base, got)
		}
	}
}

func TestParseInvalidDigitsReturnsError(t *testing.T) {
	t.Parallel()

	if _, err := Parse("2", 2); err == nil {
		t.Fatal("expected error for digit outside base")
	}
}

func TestConvert(t *testing.T) {
	t.Parallel()

	got, err := Convert("255", 10)
	if err != nil {
		t.Fatalf("Convert returned error: %v", err)
	}
	if got.Input != "255" || got.Base != 10 || got.Value != 255 {
		t.Fatalf("Convert() = %+v, want input 255, base 10, value 255", got)
	}

	want := []Conversion{
		{Label: "Base 2 (binary)", Base: 2, Value: "11111111"},
		{Label: "Base 8 (octal)", Base: 8, Value: "377"},
		{Label: "Base 10 (decimal)", Base: 10, Value: "255"},
		{Label: "Base 16 (hexadecimal)", Base: 16, Value: "ff"},
	}
	if !slices.Equal(got.Conversions, want) {
		t.Fatalf("Convert().Conversions = %+v, want %+v", got.Conversions, want)
	}
}

func TestConvertInvalidInputReturnsError(t *testing.T) {
	t.Parallel()

	if _, err := Convert("12", 2); err == nil {
		t.Fatal("expected error for invalid binary input")
	}
}
