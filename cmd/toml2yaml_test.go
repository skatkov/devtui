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
		args        []string
		wantContain string
		wantErr     bool
		description string
	}{
		{
			name:        "simple key-value conversion",
			input:       "name = \"myapp\"\nversion = \"1.0.0\"",
			args:        []string{},
			wantContain: "name:",
			wantErr:     false,
			description: "Should convert simple TOML to YAML",
		},
		{
			name:        "nested tables conversion",
			input:       "[user]\nname = \"Bob\"\nemail = \"bob@example.com\"",
			args:        []string{},
			wantContain: "user:",
			wantErr:     false,
			description: "Should convert nested TOML tables to YAML",
		},
		{
			name:        "arrays conversion",
			input:       "items = [\"apple\", \"banana\", \"cherry\"]",
			args:        []string{},
			wantContain: "- apple",
			wantErr:     false,
			description: "Should convert TOML arrays to YAML",
		},
		{
			name:        "boolean and number values",
			input:       "enabled = true\ncount = 42\nratio = 3.14",
			args:        []string{},
			wantContain: "enabled:",
			wantErr:     false,
			description: "Should preserve boolean and number types",
		},
		{
			name:        "array of tables",
			input:       "[[items]]\nid = 1\nname = \"first\"\n\n[[items]]\nid = 2\nname = \"second\"",
			args:        []string{},
			wantContain: "items:",
			wantErr:     false,
			description: "Should convert TOML array of tables to YAML",
		},
		{
			name:        "invalid TOML input",
			input:       "{invalid toml}",
			args:        []string{},
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

			args := []string{"toml2yaml"}
			args = append(args, tt.args...)
			cmd.SetArgs(args)

			err := cmd.Execute()

			if (err != nil) != tt.wantErr {
				t.Errorf("toml2yaml command error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.wantContain != "" {
				output := buf.String()
				if !strings.Contains(output, tt.wantContain) {
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
		wantContain string
		description string
	}{
		{
			name:        "argument input",
			args:        []string{"toml2yaml", "name = \"myapp\""},
			input:       "",
			wantContain: "name:",
			description: "Should handle TOML string argument",
		},
		{
			name:        "stdin input",
			args:        []string{"toml2yaml"},
			input:       "name = \"myapp\"",
			wantContain: "name:",
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
			if !strings.Contains(output, tt.wantContain) {
				t.Errorf("toml2yaml command output check failed.\nOutput: %s\nDescription: %s",
					output, tt.description)
			}
		})
	}
}
