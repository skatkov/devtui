package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestYaml2tomlCmd(t *testing.T) {
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
			input:       "name: myapp\nversion: \"1.0.0\"",
			args:        []string{},
			wantContain: "name",
			wantErr:     false,
			description: "Should convert simple YAML to TOML",
		},
		{
			name:        "nested tables conversion",
			input:       "user:\n  name: Bob\n  email: bob@example.com",
			args:        []string{},
			wantContain: "[user]",
			wantErr:     false,
			description: "Should convert nested YAML to TOML tables",
		},
		{
			name:        "arrays conversion",
			input:       "items:\n  - apple\n  - banana",
			args:        []string{},
			wantContain: "items",
			wantErr:     false,
			description: "Should convert YAML arrays to TOML arrays",
		},
		{
			name:        "mixed types conversion",
			input:       "string: text\nnumber: 42\nfloat: 3.14\nboolean: true",
			args:        []string{},
			wantContain: "string",
			wantErr:     false,
			description: "Should handle mixed YAML types",
		},
		{
			name:        "invalid YAML input",
			input:       "not: valid: yaml: here",
			args:        []string{},
			wantErr:     true,
			description: "Should return error for invalid YAML",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := GetRootCmd()
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)
			cmd.SetIn(strings.NewReader(tt.input))

			args := []string{"yaml2toml"}
			args = append(args, tt.args...)
			cmd.SetArgs(args)

			err := cmd.Execute()

			if (err != nil) != tt.wantErr {
				t.Errorf("yaml2toml command error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.wantContain != "" {
				output := buf.String()
				if !strings.Contains(output, tt.wantContain) {
					t.Errorf("yaml2toml command output check failed.\nInput: %s\nOutput: %s\nDescription: %s",
						tt.input, output, tt.description)
				}
			}
		})
	}
}

func TestYaml2tomlCmdNoInput(t *testing.T) {
	cmd := GetRootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetIn(strings.NewReader(""))
	cmd.SetArgs([]string{"yaml2toml"})

	err := cmd.Execute()

	if err == nil {
		t.Error("yaml2toml command should return error when no input provided")
	}

	if !strings.Contains(err.Error(), "no input provided") {
		t.Errorf("yaml2toml command error message should mention 'no input provided', got: %v", err)
	}
}

func TestYaml2tomlCmdArgumentInput(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		input       string
		wantContain string
		description string
	}{
		{
			name:        "argument input",
			args:        []string{"yaml2toml", "name: myapp"},
			input:       "",
			wantContain: "name",
			description: "Should handle YAML string argument",
		},
		{
			name:        "stdin input",
			args:        []string{"yaml2toml"},
			input:       "name: myapp",
			wantContain: "name",
			description: "Should handle YAML from stdin",
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
				t.Fatalf("yaml2toml command failed: %v", err)
			}

			output := buf.String()
			if !strings.Contains(output, tt.wantContain) {
				t.Errorf("yaml2toml command output check failed.\nOutput: %s\nDescription: %s",
					output, tt.description)
			}
		})
	}
}
