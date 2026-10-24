package cmd

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/pflag"
)

func TestBase64CLI(t *testing.T) {
	tests := []struct {
		name   string
		args   []string
		input  string
		want   string
		hasErr bool
	}{
		{
			name: "encode simple string",
			args: []string{"base64", "hello world"},
			want: "aGVsbG8gd29ybGQ=",
		},
		{
			name: "decode simple string",
			args: []string{"base64", "aGVsbG8gd29ybGQ=", "--decode"},
			want: "hello world",
		},
		{
			name: "decode with short flag",
			args: []string{"base64", "aGVsbG8gd29ybGQ=", "-d"},
			want: "hello world",
		},
		{
			name:  "encode from stdin",
			args:  []string{"base64"},
			input: "hello world",
			want:  "aGVsbG8gd29ybGQ=",
		},
		{
			name:  "decode from stdin",
			args:  []string{"base64", "--decode"},
			input: "aGVsbG8gd29ybGQ=",
			want:  "hello world",
		},
		{
			name:   "invalid base64 decode",
			args:   []string{"base64", "invalid_base64!", "--decode"},
			hasErr: true,
		},
		{
			name:  "empty stdin should encode to empty",
			args:  []string{"base64"},
			input: "",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr, err := runDevtui(strings.NewReader(tt.input), tt.args...)

			if tt.hasErr {
				if err == nil {
					t.Errorf("Expected error but command succeeded")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v\nStderr: %s", err, stderr.String())
				return
			}

			got := stdout.String()
			if got != tt.want {
				t.Errorf("Output mismatch:\nGot:  %q\nWant: %q", got, tt.want)
			}
		})
	}
}

func TestBase64CLIWithFiles(t *testing.T) {
	testDataDir := "../testdata"

	// Check if testdata directory exists
	if _, err := os.Stat(testDataDir); os.IsNotExist(err) {
		t.Skip("testdata directory not found, skipping file tests")
	}

	testPairs := []struct {
		textFile   string
		base64File string
	}{
		{"sample.txt", "sample.base64"},
		{"json.txt", "json.base64"},
		{"binary.txt", "binary.base64"},
	}

	for _, pair := range testPairs {
		t.Run("encode_"+pair.textFile, func(t *testing.T) {
			textPath := filepath.Join(testDataDir, pair.textFile)
			base64Path := filepath.Join(testDataDir, pair.base64File)

			// Skip if files don't exist
			if _, err := os.Stat(textPath); os.IsNotExist(err) {
				t.Skipf("File %s not found", textPath)
			}
			if _, err := os.Stat(base64Path); os.IsNotExist(err) {
				t.Skipf("File %s not found", base64Path)
			}

			// Read text file content
			textContent, err := os.ReadFile(textPath)
			if err != nil {
				t.Errorf("Failed to read %s: %v", textPath, err)
				return
			}

			// Encode using stdin
			stdout, stderr, err := runDevtui(bytes.NewReader(textContent), "base64")
			if err != nil {
				t.Errorf("Failed to encode %s: %v\nStderr: %s", textPath, err, stderr.String())
				return
			}

			// Read expected output
			expected, err := os.ReadFile(base64Path)
			if err != nil {
				t.Errorf("Failed to read %s: %v", base64Path, err)
				return
			}

			got := stdout.String()
			want := string(expected)

			if got != want {
				t.Errorf("Encoding mismatch for %s:\nGot:  %q\nWant: %q", pair.textFile, got, want)
			}
		})

		t.Run("decode_"+pair.base64File, func(t *testing.T) {
			textPath := filepath.Join(testDataDir, pair.textFile)
			base64Path := filepath.Join(testDataDir, pair.base64File)

			// Skip if files don't exist
			if _, err := os.Stat(textPath); os.IsNotExist(err) {
				t.Skipf("File %s not found", textPath)
			}
			if _, err := os.Stat(base64Path); os.IsNotExist(err) {
				t.Skipf("File %s not found", base64Path)
			}

			// Read base64 file content
			base64Content, err := os.ReadFile(base64Path)
			if err != nil {
				t.Errorf("Failed to read %s: %v", base64Path, err)
				return
			}

			// Decode using stdin
			stdout, stderr, err := runDevtui(bytes.NewReader(base64Content), "base64", "--decode")
			if err != nil {
				t.Errorf("Failed to decode %s: %v\nStderr: %s", base64Path, err, stderr.String())
				return
			}

			// Read expected output
			expected, err := os.ReadFile(textPath)
			if err != nil {
				t.Errorf("Failed to read %s: %v", textPath, err)
				return
			}

			got := stdout.String()
			want := string(expected)

			if got != want {
				t.Errorf("Decoding mismatch for %s:\nGot:  %q\nWant: %q", pair.base64File, got, want)
			}
		})
	}
}

func TestBase64CLIRoundTrip(t *testing.T) {
	testDataDir := "../testdata"

	// Check if testdata directory exists
	if _, err := os.Stat(testDataDir); os.IsNotExist(err) {
		t.Skip("testdata directory not found, skipping round-trip tests")
	}

	testFiles := []string{"sample.txt", "json.txt", "binary.txt"}

	for _, filename := range testFiles {
		t.Run("roundtrip_"+filename, func(t *testing.T) {
			filePath := filepath.Join(testDataDir, filename)

			// Skip if file doesn't exist
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				t.Skipf("File %s not found", filePath)
			}

			// Read original content
			original, err := os.ReadFile(filePath)
			if err != nil {
				t.Errorf("Failed to read %s: %v", filePath, err)
				return
			}

			// Step 1: Encode using stdin
			encodedOut, encodeErr, err := runDevtui(bytes.NewReader(original), "base64")
			if err != nil {
				t.Errorf("Failed to encode %s: %v\nStderr: %s", filePath, err, encodeErr.String())
				return
			}

			encoded := encodedOut.String()

			// Step 2: Decode using stdin
			decodedOut, decodeErrBuf, err := runDevtui(strings.NewReader(encoded), "base64", "--decode")
			if err != nil {
				t.Errorf("Failed to decode for %s: %v\nStderr: %s", filename, err, decodeErrBuf.String())
				return
			}

			decoded := decodedOut.String()
			originalStr := string(original)

			if decoded != originalStr {
				t.Errorf("Round-trip failed for %s:\nOriginal: %q\nDecoded:  %q", filename, originalStr, decoded)
			}
		})
	}
}

func TestBase64CLIInvalidData(t *testing.T) {
	testDataDir := "../testdata"
	invalidFile := filepath.Join(testDataDir, "invalid.base64")

	// Check if invalid file exists
	if _, err := os.Stat(invalidFile); os.IsNotExist(err) {
		t.Skip("invalid.base64 not found, skipping invalid data test")
	}

	// Read invalid base64 content
	invalidContent, err := os.ReadFile(invalidFile)
	if err != nil {
		t.Fatalf("Failed to read invalid.base64: %v", err)
	}

	// Try to decode invalid base64 using stdin - should fail
	_, stderr, err := runDevtui(bytes.NewReader(invalidContent), "base64", "--decode")
	if err == nil {
		t.Errorf("Expected error when decoding invalid base64 file, but command succeeded")
	}

	// Should have error output
	if stderr.Len() == 0 {
		t.Errorf("Expected error message in stderr, but got none")
	}
}

// runDevtui runs devtui in process with the given stdin, resetting the flags a
// previous run set.
func runDevtui(stdin io.Reader, args ...string) (*bytes.Buffer, *bytes.Buffer, error) {
	cmd := GetRootCmd()
	for _, sub := range cmd.Commands() {
		sub.Flags().VisitAll(func(flag *pflag.Flag) {
			_ = flag.Value.Set(flag.DefValue)
			flag.Changed = false
		})
	}

	var stdout, stderr bytes.Buffer
	cmd.SetIn(stdin)
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)
	cmd.SetArgs(args)

	err := cmd.Execute()
	return &stdout, &stderr, err
}
