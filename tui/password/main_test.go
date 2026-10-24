package main

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"testing/iotest"

	"charm.land/huh/v2"
)

const lowercase = "abcdefghijklmnopqrstuvwxyz"

var errStop = errors.New("stop")

func TestShowInterfaceGeneratesPassword(t *testing.T) {
	out, _ := scriptForms(t, "8\n1\n0\n", "y\n")
	config := &PasswordConfig{Length: 12}

	if err := showInterface(config); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if config.Length != 8 || len(config.GeneratedPass) != 8 {
		t.Fatalf("expected 8 character password, got length %d and %q", config.Length, config.GeneratedPass)
	}
	if strings.Trim(config.GeneratedPass, lowercase) != "" {
		t.Fatalf("expected only lowercase letters, got %q", config.GeneratedPass)
	}
	if !config.Continue {
		t.Fatal("expected continue to be confirmed")
	}
	if !strings.Contains(out.String(), config.GeneratedPass) {
		t.Fatalf("expected password in output, got %q", out.String())
	}
}

func TestShowInterfaceWithoutCharacterSetsShowsError(t *testing.T) {
	out, _ := scriptForms(t, "8\n0\n", "")
	config := &PasswordConfig{Length: 12}

	if err := showInterface(config); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if config.GeneratedPass != "" {
		t.Fatalf("expected no password, got %q", config.GeneratedPass)
	}
	if !strings.Contains(out.String(), "Please select at least one character set!") {
		t.Fatalf("expected error note in output, got %q", out.String())
	}
}

func TestShowInterfaceInvalidLengthReturnsError(t *testing.T) {
	scriptForms(t, "abc\n1\n0\n")

	err := showInterface(&PasswordConfig{Length: 12})
	if err == nil || !strings.Contains(err.Error(), "invalid length") {
		t.Fatalf("expected invalid length error, got %v", err)
	}
}

func TestShowInterfaceGenerationErrorIsReturned(t *testing.T) {
	scriptForms(t, "0\n1\n0\n")

	err := showInterface(&PasswordConfig{Length: 12})
	if err == nil || !strings.Contains(err.Error(), "failed to generate password") {
		t.Fatalf("expected generation error, got %v", err)
	}
}

func TestShowInterfaceFormErrorIsReturned(t *testing.T) {
	scriptForms(t)
	config := &PasswordConfig{Length: 12}

	if err := showInterface(config); !errors.Is(err, errStop) {
		t.Fatalf("expected form error, got %v", err)
	}
	if config.Length != 12 {
		t.Fatalf("expected length to stay 12, got %d", config.Length)
	}
}

func TestShowInterfaceResultFormErrorIsReturned(t *testing.T) {
	scriptForms(t, "8\n1\n0\n")

	if err := showInterface(&PasswordConfig{Length: 12}); !errors.Is(err, errStop) {
		t.Fatalf("expected result form error, got %v", err)
	}
}

func TestMainStopsWhenNotContinuing(t *testing.T) {
	_, calls := scriptForms(t, "\n1\n0\n", "n\n")

	main()

	if *calls != 2 {
		t.Fatalf("expected main to stop after answering no, got %d forms", *calls)
	}
}

func TestMainUsesDefaultLengthAndStopsOnError(t *testing.T) {
	out, calls := scriptForms(t, "\n1\n0\n", "y\n")

	main()

	if *calls != 3 {
		t.Fatalf("expected main to stop at the failing third form, got %d forms", *calls)
	}

	for line := range strings.SplitSeq(out.String(), "\n") {
		if len(line) == 12 && strings.Trim(line, lowercase) == "" {
			return
		}
	}
	t.Fatalf("expected a 12 character password in output, got %q", out.String())
}

func TestCharacterSetOptions(t *testing.T) {
	t.Parallel()

	options := characterSetOptions()
	if len(options) != 4 {
		t.Fatalf("expected 4 options, got %d", len(options))
	}
	if options[0].Key != "Lowercase (a-z)" || options[0].Value != "lowercase" {
		t.Fatalf("expected lowercase option first, got %+v", options[0])
	}
}

func TestSelectedCharacterSets(t *testing.T) {
	t.Parallel()

	sets := selectedCharacterSets([]string{"numbers", "unknown", "lowercase"})
	if len(sets) != 2 || sets[0].id != "numbers" || sets[1].id != "lowercase" {
		t.Fatalf("expected numbers and lowercase, got %+v", sets)
	}
}

func TestGenerateSecurePassword(t *testing.T) {
	t.Parallel()

	got, err := generateSecurePassword(&PasswordConfig{Length: 16, CharacterSets: []string{"lowercase", "numbers"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 16 || strings.Trim(got, lowercase+"0123456789") != "" {
		t.Fatalf("expected 16 lowercase letters and digits, got %q", got)
	}
}

func TestGenerateSecurePasswordUnknownSetReturnsError(t *testing.T) {
	t.Parallel()

	if _, err := generateSecurePassword(&PasswordConfig{Length: 8, CharacterSets: []string{"unknown"}}); err == nil {
		t.Fatal("expected error for unknown character set")
	}
}

func TestGenerateSecurePasswordIncludesEverySet(t *testing.T) {
	t.Parallel()

	ids := []string{"lowercase", "uppercase", "numbers", "special"}
	sets := selectedCharacterSets(ids)
	for range 200 {
		got, err := generateSecurePassword(&PasswordConfig{Length: len(ids), CharacterSets: ids})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		for _, set := range sets {
			if !strings.ContainsAny(got, set.chars) {
				t.Fatalf("expected %q to contain a %s character", got, set.id)
			}
		}
	}
}

func TestEnsureAllSetsIncludedRequiresLength(t *testing.T) {
	t.Parallel()

	sets := selectedCharacterSets([]string{"lowercase", "numbers"})
	tests := []struct {
		name    string
		length  int
		wantErr bool
	}{
		{"shorter than sets", 1, true},
		{"as long as sets", 2, false},
		{"longer than sets", 3, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			password := make([]byte, tt.length)
			err := ensureAllSetsIncluded(password, sets)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ensureAllSetsIncluded(%d) error = %v, wantErr %v", tt.length, err, tt.wantErr)
			}
			if !tt.wantErr && bytes.Count(password, []byte{0}) == tt.length {
				t.Fatalf("expected characters to be placed, got %q", password)
			}
		})
	}
}

// scriptForms answers each form in accessible mode with the next input and
// fails the next form with errStop.
func scriptForms(t *testing.T, inputs ...string) (*bytes.Buffer, *int) {
	t.Helper()

	var out bytes.Buffer
	calls := 0
	original := runForm
	t.Cleanup(func() { runForm = original })

	runForm = func(f *huh.Form) error {
		calls++
		if calls > len(inputs)+1 {
			t.Fatalf("unexpected form %d", calls)
		}
		if calls > len(inputs) {
			return errStop
		}
		input := iotest.OneByteReader(strings.NewReader(inputs[calls-1]))
		return f.WithAccessible(true).WithInput(input).WithOutput(&out).Run()
	}

	return &out, &calls
}
