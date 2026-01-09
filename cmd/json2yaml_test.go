package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestJson2yamlCmd(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		args        []string
		wantContain string
		wantErr     bool
		description string
	}{
		{
			name:        "simple object conversion",
			input:       `{"name": "myapp", "version": "1.0.0"}`,
			args:        []string{},
			wantContain: "name:",
			wantErr:     false,
			description: "Should convert simple JSON object to YAML",
		},
		{
			name:        "nested objects conversion",
			input:       `{"database": {"host": "localhost", "port": 5432}}`,
			args:        []string{},
			wantContain: "database:",
			wantErr:     false,
			description: "Should convert nested JSON objects to YAML",
		},
		{
			name:        "arrays conversion",
			input:       `{"fruits": ["apple", "banana", "cherry"]}`,
			args:        []string{},
			wantContain: "- apple",
			wantErr:     false,
			description: "Should convert JSON arrays to YAML",
		},
		{
			name:        "mixed types conversion",
			input:       `{"string": "text", "number": 123, "float": 3.14, "boolean": true}`,
			args:        []string{},
			wantContain: "string:",
			wantErr:     false,
			description: "Should handle mixed JSON types",
		},
		{
			name:        "empty object",
			input:       `{}`,
			args:        []string{},
			wantContain: "{}",
			wantErr:     false,
			description: "Should handle empty JSON object",
		},
		{
			name:        "invalid JSON input",
			input:       `{invalid json}`,
			args:        []string{},
			wantContain: "",
			wantErr:     true,
			description: "Should return error for invalid JSON",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := GetRootCmd()
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)
			cmd.SetIn(strings.NewReader(tt.input))

			args := []string{"json2yaml"}
			args = append(args, tt.args...)
			cmd.SetArgs(args)

			err := cmd.Execute()

			if (err != nil) != tt.wantErr {
				t.Errorf("json2yaml command error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.wantContain != "" {
				output := buf.String()
				if !strings.Contains(output, tt.wantContain) {
					t.Errorf("json2yaml command output check failed.\nInput: %s\nOutput: %s\nDescription: %s",
						tt.input, output, tt.description)
				}
			}
		})
	}
}

func TestJson2yamlCmdNoInput(t *testing.T) {
	cmd := GetRootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetIn(strings.NewReader(""))
	cmd.SetArgs([]string{"json2yaml"})

	err := cmd.Execute()

	if err == nil {
		t.Error("json2yaml command should return error when no input provided")
	}

	if !strings.Contains(err.Error(), "no input provided") {
		t.Errorf("json2yaml command error message should mention 'no input provided', got: %v", err)
	}
}

func TestJson2yamlCmdArgumentInput(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		input       string
		wantContain string
		description string
	}{
		{
			name:        "argument input",
			args:        []string{"json2yaml", `{"key": "value"}`},
			input:       "",
			wantContain: "key:",
			description: "Should handle JSON string argument",
		},
		{
			name:        "stdin input",
			args:        []string{"json2yaml"},
			input:       `{"key": "value"}`,
			wantContain: "key:",
			description: "Should handle JSON from stdin",
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
				t.Fatalf("json2yaml command failed: %v", err)
			}

			output := buf.String()
			if !strings.Contains(output, tt.wantContain) {
				t.Errorf("json2yaml command output check failed.\nOutput: %s\nDescription: %s",
					output, tt.description)
			}
		})
	}
}
