package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestToml2yamlCmd(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		checkOutput func(string) bool
		wantErr     bool
		description string
	}{
		{
			name:  "simple key-value conversion",
			input: "name = \"myapp\"\nversion = \"1.0.0\"",
			checkOutput: func(output string) bool {
				return strings.Contains(output, "name:") &&
					strings.Contains(output, "myapp") &&
					strings.Contains(output, "version:") &&
					strings.Contains(output, "1.0.0")
			},
			description: "Should convert simple TOML to YAML",
		},
		{
			name:  "nested tables conversion",
			input: "[user]\nname = \"Bob\"\nemail = \"bob@example.com\"",
			checkOutput: func(output string) bool {
				return strings.Contains(output, "user:") &&
					strings.Contains(output, "name:") &&
					strings.Contains(output, "Bob") &&
					strings.Contains(output, "email:") &&
					strings.Contains(output, "bob@example.com")
			},
			description: "Should convert nested TOML tables to YAML",
		},
		{
			name:  "arrays conversion",
			input: "items = [\"apple\", \"banana\", \"cherry\"]",
			checkOutput: func(output string) bool {
				return strings.Contains(output, "items:") &&
					strings.Contains(output, "- apple") &&
					strings.Contains(output, "- banana") &&
					strings.Contains(output, "- cherry")
			},
			description: "Should convert TOML arrays to YAML",
		},
		{
			name:  "boolean and number values",
			input: "enabled = true\ncount = 42\nratio = 3.14",
			checkOutput: func(output string) bool {
				return strings.Contains(output, "enabled:") &&
					strings.Contains(output, "true") &&
					strings.Contains(output, "count:") &&
					strings.Contains(output, "42") &&
					strings.Contains(output, "ratio:") &&
					strings.Contains(output, "3.14")
			},
			description: "Should preserve boolean and number types",
		},
		{
			name:  "array of tables",
			input: "[[items]]\nid = 1\nname = \"first\"\n\n[[items]]\nid = 2\nname = \"second\"",
			checkOutput: func(output string) bool {
				return strings.Contains(output, "items:") &&
					strings.Contains(output, "id:") &&
					strings.Contains(output, "name:")
			},
			description: "Should convert TOML array of tables to YAML",
		},
		{
			name:        "invalid TOML input",
			input:       "{invalid toml}",
			checkOutput: nil,
			wantErr:     true,
			description: "Should error on invalid TOML",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := GetRootCmd()
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)
			cmd.SetIn(strings.NewReader(tt.input))
			cmd.SetArgs([]string{"toml2yaml"})

			err := cmd.Execute()

			if (err != nil) != tt.wantErr {
				t.Errorf("toml2yaml command error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.checkOutput != nil {
				output := buf.String()
				if !tt.checkOutput(output) {
					t.Errorf("toml2yaml command output check failed.\nInput: %s\nOutput: %s\nDescription: %s",
						tt.input, output, tt.description)
				}
			}
		})
	}
}

func TestToml2yamlCmdNoInput(t *testing.T) {
	cmd := GetRootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetIn(strings.NewReader(""))
	cmd.SetArgs([]string{"toml2yaml"})

	err := cmd.Execute()

	if err == nil {
		t.Error("toml2yaml command should return error when no input provided")
	}

	if !strings.Contains(err.Error(), "no input provided") {
		t.Errorf("toml2yaml command error message should mention 'no input provided', got: %v", err)
	}
}

func TestToml2yamlCmdArgumentInput(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		input       string
		checkOutput func(string) bool
		description string
	}{
		{
			name:  "argument input",
			args:  []string{"toml2yaml", "name = \"myapp\""},
			input: "",
			checkOutput: func(output string) bool {
				return strings.Contains(output, "name:") && strings.Contains(output, "myapp")
			},
			description: "Should handle TOML string argument",
		},
		{
			name:  "stdin input",
			args:  []string{"toml2yaml"},
			input: "name = \"myapp\"",
			checkOutput: func(output string) bool {
				return strings.Contains(output, "name:") && strings.Contains(output, "myapp")
			},
			description: "Should handle TOML from stdin",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := GetRootCmd()
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)
			cmd.SetIn(strings.NewReader(tt.input))
			cmd.SetArgs(tt.args)

			err := cmd.Execute()
			if err != nil {
				t.Fatalf("toml2yaml command failed: %v", err)
			}

			output := buf.String()
			if !tt.checkOutput(output) {
				t.Errorf("toml2yaml command output check failed.\nOutput: %s\nDescription: %s",
					output, tt.description)
			}
		})
	}
}
