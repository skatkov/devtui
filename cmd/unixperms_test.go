package cmd

import (
	"testing"
)

func TestOctalToSymbolic(t *testing.T) {
	cases := []struct {
		octal    int
		symbolic string
	}{
		{0o755, "rwxr-xr-x"},
		{0o644, "rw-r--r--"},
		{0o600, "rw-------"},
		{0o777, "rwxrwxrwx"},
		{0o000, "---------"},
		{0o444, "r--r--r--"},
	}
	for _, tc := range cases {
		got := octalToSymbolic(tc.octal)
		if got != tc.symbolic {
			t.Errorf("octalToSymbolic(%04o) = %q, want %q", tc.octal, got, tc.symbolic)
		}
	}
}

func TestSymbolicToOctal(t *testing.T) {
	cases := []struct {
		symbolic string
		octal    int
	}{
		{"rwxr-xr-x", 0o755},
		{"rw-r--r--", 0o644},
		{"rw-------", 0o600},
		{"rwxrwxrwx", 0o777},
		{"---------", 0o000},
	}
	for _, tc := range cases {
		got, err := symbolicToOctal(tc.symbolic)
		if err != nil {
			t.Errorf("symbolicToOctal(%q) unexpected error: %v", tc.symbolic, err)
			continue
		}
		if got != tc.octal {
			t.Errorf("symbolicToOctal(%q) = %04o, want %04o", tc.symbolic, got, tc.octal)
		}
	}
}

func TestSymbolicToOctalErrors(t *testing.T) {
	bad := []string{"rwx", "rwxr-xr-xx", "abc123xyz"}
	for _, s := range bad {
		if _, err := symbolicToOctal(s); err == nil {
			t.Errorf("symbolicToOctal(%q) expected error, got nil", s)
		}
	}
}

func TestRoundTrip(t *testing.T) {
	for mode := 0; mode <= 0o777; mode++ {
		sym := octalToSymbolic(mode)
		got, err := symbolicToOctal(sym)
		if err != nil {
			t.Fatalf("round-trip failed for %04o: %v", mode, err)
		}
		if got != mode {
			t.Errorf("round-trip mismatch for %04o: symbolic=%q -> %04o", mode, sym, got)
		}
	}
}
