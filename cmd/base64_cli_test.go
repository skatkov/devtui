package cmd

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestBase64CLI(t *testing.T) {
	// Build the binary first
	buildCmd := exec.Command("go", "build", "-o", "devtui_test")
	buildCmd.Dir = ".."
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build devtui: %v", err)
	}
	defer os.Remove("../devtui_test")

	binary := "../devtui_test"

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
			name:   "no input provided",
			args:   []string{"base64"},
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)

			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			if tt.input != "" {
				cmd.Stdin = strings.NewReader(tt.input)
			}

			err := cmd.Run()

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
	// Build the binary first
	buildCmd := exec.Command("go", "build", "-o", "devtui_test")
	buildCmd.Dir = ".."
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build devtui: %v", err)
	}
	defer os.Remove("../devtui_test")

	binary := "../devtui_test"
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

			// Encode the text file
			cmd := exec.Command(binary, "base64", textPath)
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
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

			// Decode the base64 file
			cmd := exec.Command(binary, "base64", base64Path, "--decode")
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
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
	// Build the binary first
	buildCmd := exec.Command("go", "build", "-o", "devtui_test")
	buildCmd.Dir = ".."
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build devtui: %v", err)
	}
	defer os.Remove("../devtui_test")

	binary := "../devtui_test"
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

			// Step 1: Encode
			encodeCmd := exec.Command(binary, "base64", filePath)
			var encodedOut, encodeErr bytes.Buffer
			encodeCmd.Stdout = &encodedOut
			encodeCmd.Stderr = &encodeErr

			err = encodeCmd.Run()
			if err != nil {
				t.Errorf("Failed to encode %s: %v\nStderr: %s", filePath, err, encodeErr.String())
				return
			}

			encoded := encodedOut.String()

			// Step 2: Decode
			decodeCmd := exec.Command(binary, "base64", "--decode")
			decodeCmd.Stdin = strings.NewReader(encoded)
			var decodedOut, decodeErrBuf bytes.Buffer
			decodeCmd.Stdout = &decodedOut
			decodeCmd.Stderr = &decodeErrBuf

			err = decodeCmd.Run()
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
	// Build the binary first
	buildCmd := exec.Command("go", "build", "-o", "devtui_test")
	buildCmd.Dir = ".."
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build devtui: %v", err)
	}
	defer os.Remove("../devtui_test")

	binary := "../devtui_test"
	testDataDir := "../testdata"
	invalidFile := filepath.Join(testDataDir, "invalid.base64")

	// Check if invalid file exists
	if _, err := os.Stat(invalidFile); os.IsNotExist(err) {
		t.Skip("invalid.base64 not found, skipping invalid data test")
	}

	// Try to decode invalid base64 file - should fail
	cmd := exec.Command(binary, "base64", invalidFile, "--decode")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err == nil {
		t.Errorf("Expected error when decoding invalid base64 file, but command succeeded")
	}

	// Should have error output
	if stderr.Len() == 0 {
		t.Errorf("Expected error message in stderr, but got none")
	}
}
