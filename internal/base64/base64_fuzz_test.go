package base64

import (
	"bytes"
	stdbase64 "encoding/base64"
	"strings"
	"testing"
)

func FuzzBase64RoundTrip(f *testing.F) {
	f.Add([]byte(""))
	f.Add([]byte("hello world"))
	f.Add([]byte("\x00\x01\x02\xff\xfe"))
	f.Add([]byte("ñáéíóú 中文 🚀"))

	f.Fuzz(func(t *testing.T, data []byte) {
		encoded := Encode(data)

		decoded, err := Decode(encoded)
		if err != nil {
			t.Fatalf("Decode(Encode(data)) returned error: %v", err)
		}

		if !bytes.Equal(decoded, data) {
			t.Fatalf("round-trip mismatch: got %x, want %x", decoded, data)
		}

		decodedString, err := DecodeToString(encoded)
		if err != nil {
			t.Fatalf("DecodeToString(Encode(data)) returned error: %v", err)
		}

		if decodedString != string(data) {
			t.Fatalf("string round-trip mismatch: got %q, want %q", decodedString, string(data))
		}
	})
}

func FuzzDecodeMatchesStdlibBehavior(f *testing.F) {
	f.Add("")
	f.Add("SGVsbG8=")
	f.Add(" SGVsbG8= ")
	f.Add("%%%")
	f.Add("YWJj\n")

	f.Fuzz(func(t *testing.T, input string) {
		trimmed := strings.TrimSpace(input)
		want, wantErr := stdbase64.StdEncoding.DecodeString(trimmed)

		got, err := Decode(input)
		if wantErr != nil {
			if err == nil {
				t.Fatalf("expected error for input %q", input)
			}
			return
		}

		if err != nil {
			t.Fatalf("unexpected error for input %q: %v", input, err)
		}

		if !bytes.Equal(got, want) {
			t.Fatalf("decoded bytes mismatch for input %q: got %x, want %x", input, got, want)
		}
	})
}
