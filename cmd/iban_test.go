package cmd

import (
	"regexp"
	"strings"
	"testing"
)

func TestIbanCmd(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{"electronic format", []string{"GB"}, `^GB\d{2}[A-Z]{4}\d{14}\n$`},
		{"lowercase with spaces", []string{" gb "}, `^GB\d{2}[A-Z]{4}\d{14}\n$`},
		{"paper format", []string{"GB", "--format"}, `^GB\d{2} [A-Z]{4} \d{4} \d{4} \d{4} \d{2}\n$`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, _, err := runDevtui(strings.NewReader(""), append([]string{"iban"}, tt.args...)...)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !regexp.MustCompile(tt.want).MatchString(stdout.String()) {
				t.Fatalf("expected output matching %s, got %q", tt.want, stdout.String())
			}
		})
	}
}

func TestIbanCmdInvalidCountryCode(t *testing.T) {
	tests := []struct {
		name    string
		code    string
		wantErr string
	}{
		{"too short", "G", "country code must be exactly 2 characters"},
		{"too long", "GBR", "country code must be exactly 2 characters"},
		{"unsupported", "xx", "failed to generate IBAN for country code 'XX'"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := runDevtui(strings.NewReader(""), "iban", tt.code)
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("expected error containing %q, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestIbanCmdNoInput(t *testing.T) {
	_, _, err := runDevtui(strings.NewReader(""), "iban")
	if err == nil {
		t.Fatal("expected error without a country code")
	}
}
