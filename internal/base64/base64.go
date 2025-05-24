// devtui/internal/base64/base64.go
package base64

import (
	"encoding/base64"
	"fmt"
	"strings"
)

// Encode encodes the given data to base64 string
func Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// EncodeString encodes the given string to base64 string
func EncodeString(s string) string {
	return Encode([]byte(s))
}

// Decode decodes the given base64 string to bytes
func Decode(encoded string) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(strings.TrimSpace(encoded))
	if err != nil {
		return nil, fmt.Errorf("invalid base64 input: %v", err)
	}
	return decoded, nil
}

// DecodeToString decodes the given base64 string to string
func DecodeToString(encoded string) (string, error) {
	decoded, err := Decode(encoded)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}
