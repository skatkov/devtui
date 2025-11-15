package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestJson2toonCmd(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		args        []string
		wantContain string
		wantErr     bool
	}{
		{
			name:        "simple object with defaults",
			input:       `{"name":"Alice","age":30}`,
			args:        []string{},
			wantContain: "name: Alice",
			wantErr:     false,
		},
		{
			name:        "array with length marker",
			input:       `{"tags":["foo","bar"]}`,
			args:        []string{"-l", "#"},
			wantContain: "tags[#2]:",
			wantErr:     false,
		},
		{
			name:        "with tab delimiter",
			input:       `{"tags":["foo","bar","baz"]}`,
			args:        []string{"-d", "tab"},
			wantContain: "tags[3\t]:",
			wantErr:     false,
		},
		{
			name:        "with pipe delimiter",
			input:       `{"tags":["foo","bar","baz"]}`,
			args:        []string{"-d", "pipe"},
			wantContain: "tags[3|]:",
			wantErr:     false,
		},
		{
			name:        "with 4-space indent",
			input:       `{"user":{"id":1}}`,
			args:        []string{"-i", "4"},
			wantContain: "user:\n    id: 1",
			wantErr:     false,
		},
		{
			name:        "combined options",
			input:       `{"items":[{"id":1},{"id":2}]}`,
			args:        []string{"-l", "#", "-d", "pipe", "-i", "4"},
			wantContain: "items[#2|]{id}:",
			wantErr:     false,
		},
		{
			name:        "invalid delimiter",
			input:       `{"test":"value"}`,
			args:        []string{"-d", "invalid"},
			wantContain: "",
			wantErr:     true,
		},
		{
			name:        "invalid JSON",
			input:       `{invalid json}`,
			args:        []string{},
			wantContain: "",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flags to defaults before each test
			json2toonIndent = 2
			json2toonDelimiter = "comma"
			json2toonLengthMarker = ""

			// Create command and set input
			cmd := json2toonCmd
			cmd.SetArgs(tt.args)

			// Capture stdin
			stdin := strings.NewReader(tt.input)
			cmd.SetIn(stdin)

			// Capture output
			var outBuf bytes.Buffer
			cmd.SetOut(&outBuf)
			cmd.SetErr(&outBuf)

			// Execute command
			err := cmd.Execute()

			// Check error expectation
			if (err != nil) != tt.wantErr {
				t.Errorf("json2toon command error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Check output content if no error expected
			if !tt.wantErr && tt.wantContain != "" {
				output := outBuf.String()
				if !strings.Contains(output, tt.wantContain) {
					t.Errorf("json2toon output does not contain expected string %q\nGot: %s", tt.wantContain, output)
				}
			}
		})
	}
}

func TestJson2toonCmdHelp(t *testing.T) {
	cmd := json2toonCmd
	cmd.SetArgs([]string{"--help"})

	var outBuf bytes.Buffer
	cmd.SetOut(&outBuf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("help command failed: %v", err)
	}

	output := outBuf.String()

	expectedStrings := []string{
		"Convert JSON to TOON",
		"--indent",
		"--delimiter",
		"--length-marker",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("help output missing expected string %q", expected)
		}
	}
}
