package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestToml2jsonCmd(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		args        []string
		wantContain string
		wantErr     bool
		description string
	}{
		{
			name: "simple key-value conversion",
			input: `name = "Alice"
age = 30`,
			args:        []string{},
			wantContain: `"name": "Alice"`,
			wantErr:     false,
			description: "Should convert simple TOML to JSON",
		},
		{
			name: "nested tables conversion",
			input: `[user]
name = "Bob"
email = "bob@example.com"`,
			args:        []string{},
			wantContain: `"user"`,
			wantErr:     false,
			description: "Should convert nested TOML tables to JSON object",
		},
		{
			name: "arrays of tables conversion",
			input: `[[items]]
id = 1
name = "first"

[[items]]
id = 2
name = "second"`,
			args:        []string{},
			wantContain: `"items"`,
			wantErr:     false,
			description: "Should convert TOML array of tables to JSON array",
		},
		{
			name:        "invalid TOML input",
			input:       `{invalid toml}`,
			args:        []string{},
			wantContain: "",
			wantErr:     true,
			description: "Should error on invalid TOML",
		},
		{
			name: "boolean and number values",
			input: `enabled = true
count = 42
ratio = 3.14`,
			args:        []string{},
			wantContain: `"enabled": true`,
			wantErr:     false,
			description: "Should preserve boolean and number types",
		},
		{
			name:        "array conversion",
			input:       `tags = ["foo", "bar", "baz"]`,
			args:        []string{},
			wantContain: `"tags"`,
			wantErr:     false,
			description: "Should convert TOML arrays to JSON arrays",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := GetRootCmd()
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)
			cmd.SetIn(strings.NewReader(tt.input))

			args := []string{"toml2json"}
			args = append(args, tt.args...)
			cmd.SetArgs(args)

			err := cmd.Execute()

			if (err != nil) != tt.wantErr {
				t.Errorf("toml2json command error = %v, wantErr %v\nDescription: %s", err, tt.wantErr, tt.description)
				return
			}

			if !tt.wantErr && tt.wantContain != "" {
				output := buf.String()
				if !strings.Contains(output, tt.wantContain) {
					t.Errorf("toml2json output does not contain expected string %q\nGot: %s\nDescription: %s",
						tt.wantContain, output, tt.description)
				}
			}
		})
	}
}

func TestToml2jsonCmdHelp(t *testing.T) {
	cmd := GetRootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"toml2json", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("help command failed: %v", err)
	}

	output := buf.String()

	expectedStrings := []string{
		"Convert TOML",
		"--tui",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("help output missing expected string %q", expected)
		}
	}
}
