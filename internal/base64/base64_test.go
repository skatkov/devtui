package base64

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBase64EncodeDecode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple string",
			input:    "Hello, World!",
			expected: "SGVsbG8sIFdvcmxkIQ==",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "ASCII characters",
			input:    "abcdefghijklmnopqrstuvwxyz",
			expected: "YWJjZGVmZ2hpamtsbW5vcHFyc3R1dnd4eXo=",
		},
		{
			name:     "Numbers",
			input:    "0123456789",
			expected: "MDEyMzQ1Njc4OQ==",
		},
		{
			name:     "Special characters",
			input:    "!@#$%^&*()",
			expected: "IUAjJCVeJiooKQ==",
		},
		{
			name:     "Unicode characters",
			input:    "ñáéíóú 中文 🚀",
			expected: "w7HDocOpw63Ds8O6IOS4reaWhyDwn5qA",
		},
		{
			name:     "Multi-line text",
			input:    "Line 1\nLine 2\nLine 3",
			expected: "TGluZSAxCkxpbmUgMgpMaW5lIDM=",
		},
		{
			name:     "JSON-like structure",
			input:    `{"name": "test", "value": 123}`,
			expected: "eyJuYW1lIjogInRlc3QiLCAidmFsdWUiOiAxMjN9",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name+" encode", func(t *testing.T) {
			got := EncodeString(tt.input)
			if got != tt.expected {
				t.Errorf("EncodeString(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})

		t.Run(tt.name+" decode", func(t *testing.T) {
			if tt.expected == "" {
				return // Skip empty string decode test
			}
			got, err := DecodeToString(tt.expected)
			if err != nil {
				t.Errorf("DecodeToString(%q) error = %v", tt.expected, err)
				return
			}
			if got != tt.input {
				t.Errorf("DecodeToString(%q) = %q, want %q", tt.expected, got, tt.input)
			}
		})
	}
}

func TestBase64InvalidDecode(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{
			name:        "Invalid characters",
			input:       "SGVsbG8gV29ybGQ!invalid",
			expectError: true,
		},
		{
			name:        "Invalid padding",
			input:       "SGVsbG8===",
			expectError: true,
		},
		{
			name:        "Invalid length",
			input:       "SGVsbG8",
			expectError: true,
		},
		{
			name:        "Mixed valid/invalid",
			input:       "SGVsbG8=invalid_part",
			expectError: true,
		},
		{
			name:        "Valid base64",
			input:       "SGVsbG8=",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Decode(tt.input)
			hasError := err != nil

			if hasError != tt.expectError {
				if tt.expectError {
					t.Errorf("Expected error for input %q, but got none", tt.input)
				} else {
					t.Errorf("Unexpected error for input %q: %v", tt.input, err)
				}
			}
		})
	}
}

func TestBase64TestDataFiles(t *testing.T) {
	testDataDir := "../../testdata"

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
		t.Run(pair.textFile, func(t *testing.T) {
			textPath := filepath.Join(testDataDir, pair.textFile)
			base64Path := filepath.Join(testDataDir, pair.base64File)

			// Read text file
			textContent, err := os.ReadFile(textPath)
			if err != nil {
				t.Skipf("Could not read %s: %v", textPath, err)
			}

			// Read base64 file
			base64Content, err := os.ReadFile(base64Path)
			if err != nil {
				t.Skipf("Could not read %s: %v", base64Path, err)
			}

			// Test encoding: text file -> base64
			encoded := Encode(textContent)
			expectedBase64 := string(base64Content)

			if encoded != expectedBase64 {
				t.Errorf("Encoding mismatch for %s:\nGot:      %q\nExpected: %q", pair.textFile, encoded, expectedBase64)
			}

			// Test decoding: base64 file -> text
			if len(base64Content) > 0 {
				decoded, err := Decode(string(base64Content))
				if err != nil {
					t.Errorf("Failed to decode %s: %v", pair.base64File, err)
					return
				}

				if string(decoded) != string(textContent) {
					t.Errorf("Decoding mismatch for %s:\nGot:      %q\nExpected: %q", pair.base64File, string(decoded), string(textContent))
				}
			}
		})
	}
}

func TestBase64RoundTrip(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{
			name: "ASCII text",
			data: []byte("Hello, World! This is a test."),
		},
		{
			name: "Binary data",
			data: []byte{0x00, 0x01, 0x02, 0x03, 0xFF, 0xFE, 0xFD, 0xFC},
		},
		{
			name: "UTF-8 text",
			data: []byte("ñáéíóú 中文 🚀 ✨"),
		},
		{
			name: "Large data",
			data: make([]byte, 1024),
		},
		{
			name: "Empty data",
			data: []byte{},
		},
	}

	// Initialize large data with pattern
	for i := range tests[3].data {
		tests[3].data[i] = byte(i % 256)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encode
			encoded := Encode(tt.data)

			// Decode
			decoded, err := Decode(encoded)
			if err != nil {
				t.Errorf("Round-trip decode error: %v", err)
				return
			}

			// Compare
			if len(decoded) != len(tt.data) {
				t.Errorf("Round-trip length mismatch: got %d, want %d", len(decoded), len(tt.data))
				return
			}

			for i := range tt.data {
				if decoded[i] != tt.data[i] {
					t.Errorf("Round-trip data mismatch at byte %d: got 0x%02x, want 0x%02x", i, decoded[i], tt.data[i])
					break
				}
			}
		})
	}
}

func TestBase64WithInvalidTestData(t *testing.T) {
	testDataDir := "../../testdata"
	invalidFile := "invalid.base64"
	invalidPath := filepath.Join(testDataDir, invalidFile)

	// Check if file exists
	if _, err := os.Stat(invalidPath); os.IsNotExist(err) {
		t.Skip("invalid.base64 not found, skipping invalid data test")
	}

	// Read invalid base64 file
	invalidContent, err := os.ReadFile(invalidPath)
	if err != nil {
		t.Skipf("Could not read %s: %v", invalidPath, err)
	}

	// Try to decode - should fail
	_, err = Decode(string(invalidContent))
	if err == nil {
		t.Errorf("Expected decode error for invalid base64 data, but got none")
	}
}
