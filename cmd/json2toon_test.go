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
		description string
	}{
		{
			name:        "simple object with defaults",
			input:       `{"name":"Alice","age":30}`,
			args:        []string{},
			wantContain: "name: Alice",
			wantErr:     false,
			description: "Should convert simple JSON to TOON",
		},
		{
			name:        "array with length marker",
			input:       `{"tags":["foo","bar"]}`,
			args:        []string{"-l", "#"},
			wantContain: "tags[#2]:",
			wantErr:     false,
			description: "Should add length marker prefix",
		},
		{
			name:        "with tab delimiter",
			input:       `{"tags":["foo","bar","baz"]}`,
			args:        []string{"-d", "tab"},
			wantContain: "foo\tbar\tbaz",
			wantErr:     false,
			description: "Should use tab delimiter",
		},
		{
			name:        "with pipe delimiter",
			input:       `{"tags":["foo","bar","baz"]}`,
			args:        []string{"-d", "pipe"},
			wantContain: "foo|bar|baz",
			wantErr:     false,
			description: "Should use pipe delimiter",
		},
		{
			name:        "with 4-space indent",
			input:       `{"user":{"id":1}}`,
			args:        []string{"-i", "4"},
			wantContain: "user:",
			wantErr:     false,
			description: "Should use 4-space indentation",
		},
		{
			name:        "combined options",
			input:       `{"items":[{"id":1},{"id":2}]}`,
			args:        []string{"-l", "#", "-d", "pipe", "-i", "4"},
			wantContain: "items[#2|]{id}:",
			wantErr:     false,
			description: "Should combine all options",
		},
		{
			name:        "invalid delimiter",
			input:       `{"test":"value"}`,
			args:        []string{"-d", "invalid"},
			wantContain: "",
			wantErr:     true,
			description: "Should error on invalid delimiter",
		},
		{
			name:        "invalid JSON",
			input:       `{invalid json}`,
			args:        []string{},
			wantContain: "",
			wantErr:     true,
			description: "Should error on invalid JSON",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset global flags
			json2toonIndent = 2
			json2toonDelimiter = "comma"
			json2toonLengthMarker = ""

			// Create fresh command and root for each test
			cmd := GetRootCmd()
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)
			cmd.SetIn(strings.NewReader(tt.input))

			// Build args with json2toon command
			args := []string{"json2toon"}
			args = append(args, tt.args...)
			cmd.SetArgs(args)

			// Execute command
			err := cmd.Execute()

			// Check error expectation
			if (err != nil) != tt.wantErr {
				t.Errorf("json2toon command error = %v, wantErr %v\nDescription: %s", err, tt.wantErr, tt.description)
				return
			}

			// Check output content if no error expected
			if !tt.wantErr && tt.wantContain != "" {
				output := buf.String()
				if !strings.Contains(output, tt.wantContain) {
					t.Errorf("json2toon output does not contain expected string %q\nGot: %s\nDescription: %s",
						tt.wantContain, output, tt.description)
				}
			}
		})
	}
}

func TestJson2toonCmdHelp(t *testing.T) {
	cmd := GetRootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"json2toon", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("help command failed: %v", err)
	}

	output := buf.String()

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
