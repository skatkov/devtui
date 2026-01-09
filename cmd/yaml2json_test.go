package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestYaml2jsonCmd(t *testing.T) {
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
			input:       "name: Alice\nage: 30",
			args:        []string{},
			wantContain: `"name": "Alice"`,
			wantErr:     false,
			description: "Should convert simple YAML to JSON",
		},
		{
			name:        "nested objects conversion",
			input:       "user:\n  name: Bob\n  email: bob@example.com",
			args:        []string{},
			wantContain: `"user"`,
			wantErr:     false,
			description: "Should convert nested YAML objects to JSON",
		},
		{
			name:        "arrays conversion",
			input:       "items:\n  - apple\n  - banana\n  - cherry",
			args:        []string{},
			wantContain: `"items"`,
			wantErr:     false,
			description: "Should convert YAML arrays to JSON arrays",
		},
		{
			name:        "mixed types conversion",
			input:       "string: text\nnumber: 42\nfloat: 3.14\nboolean: true",
			args:        []string{},
			wantContain: `"string": "text"`,
			wantErr:     false,
			description: "Should handle mixed YAML types",
		},
		{
			name:        "invalid YAML input",
			input:       "not: valid: yaml: here",
			args:        []string{},
			wantErr:     true,
			description: "Should error on invalid YAML",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := GetRootCmd()
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)
			cmd.SetIn(strings.NewReader(tt.input))

			args := []string{"yaml2json"}
			args = append(args, tt.args...)
			cmd.SetArgs(args)

			err := cmd.Execute()

			if (err != nil) != tt.wantErr {
				t.Errorf("yaml2json command error = %v, wantErr %v\nDescription: %s", err, tt.wantErr, tt.description)
				return
			}

			if !tt.wantErr && tt.wantContain != "" {
				output := buf.String()
				if !strings.Contains(output, tt.wantContain) {
					t.Errorf("yaml2json output does not contain expected string %q\nGot: %s\nDescription: %s",
						tt.wantContain, output, tt.description)
				}
			}
		})
	}
}

func TestYaml2jsonCmdNoInput(t *testing.T) {
	cmd := GetRootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetIn(strings.NewReader(""))
	cmd.SetArgs([]string{"yaml2json"})

	err := cmd.Execute()

	if err == nil {
		t.Error("yaml2json command should return error when no input provided")
	}

	if !strings.Contains(err.Error(), "no input provided") {
		t.Errorf("yaml2json command error message should mention 'no input provided', got: %v", err)
	}
}

func TestYaml2jsonCmdArgumentInput(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		input       string
		wantContain string
		description string
	}{
		{
			name:        "argument input",
			args:        []string{"yaml2json", "name: myapp\nversion: 1.0.0"},
			input:       "",
			wantContain: `"name": "myapp"`,
			description: "Should handle YAML string argument",
		},
		{
			name:        "stdin input",
			args:        []string{"yaml2json"},
			input:       "name: myapp\nversion: 1.0.0",
			wantContain: `"name": "myapp"`,
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
				t.Fatalf("yaml2json command failed: %v", err)
			}

			output := buf.String()
			if !strings.Contains(output, tt.wantContain) {
				t.Errorf("yaml2json output does not contain expected string %q\nGot: %s\nDescription: %s",
					tt.wantContain, output, tt.description)
			}
		})
	}
}
