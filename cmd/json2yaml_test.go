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
		checkOutput func(string) bool
		wantErr     bool
		description string
	}{
		{
			name:  "simple object conversion",
			input: `{"name": "myapp", "version": "1.0.0"}`,
			checkOutput: func(output string) bool {
				return strings.Contains(output, "name:") &&
					strings.Contains(output, "myapp") &&
					strings.Contains(output, "version:") &&
					strings.Contains(output, "1.0.0")
			},
			description: "Should convert simple JSON object to YAML",
		},
		{
			name:  "nested objects conversion",
			input: `{"database": {"host": "localhost", "port": 5432}}`,
			checkOutput: func(output string) bool {
				return strings.Contains(output, "database:") &&
					strings.Contains(output, "host:") &&
					strings.Contains(output, "localhost") &&
					strings.Contains(output, "port:") &&
					strings.Contains(output, "5432")
			},
			description: "Should convert nested JSON objects to YAML",
		},
		{
			name:  "arrays conversion",
			input: `{"fruits": ["apple", "banana", "cherry"]}`,
			checkOutput: func(output string) bool {
				return strings.Contains(output, "fruits:") &&
					strings.Contains(output, "- apple") &&
					strings.Contains(output, "- banana") &&
					strings.Contains(output, "- cherry")
			},
			description: "Should convert JSON arrays to YAML",
		},
		{
			name:  "mixed types conversion",
			input: `{"string": "text", "number": 123, "float": 3.14, "boolean": true}`,
			checkOutput: func(output string) bool {
				return strings.Contains(output, "string:") &&
					strings.Contains(output, "number:") &&
					strings.Contains(output, "float:") &&
					strings.Contains(output, "boolean:")
			},
			description: "Should handle mixed JSON types",
		},
		{
			name:  "empty object",
			input: `{}`,
			checkOutput: func(output string) bool {
				return strings.TrimSpace(output) == "{}" || strings.TrimSpace(output) == ""
			},
			description: "Should handle empty JSON object",
		},
		{
			name:  "nested arrays",
			input: `{"matrix": [[1, 2], [3, 4]]}`,
			checkOutput: func(output string) bool {
				return strings.Contains(output, "matrix:")
			},
			description: "Should convert nested arrays correctly",
		},
		{
			name:        "invalid JSON input",
			input:       `{invalid json}`,
			checkOutput: nil,
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
			cmd.SetArgs([]string{"json2yaml"})

			err := cmd.Execute()

			if (err != nil) != tt.wantErr {
				t.Errorf("json2yaml command error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.checkOutput != nil {
				output := buf.String()
				if !tt.checkOutput(output) {
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
		checkOutput func(string) bool
		description string
	}{
		{
			name:        "argument input",
			args:        []string{"json2yaml", `{"key": "value"}`},
			input:       "",
			checkOutput: func(output string) bool { return strings.Contains(output, "key:") && strings.Contains(output, "value") },
			description: "Should handle JSON string argument",
		},
		{
			name:        "stdin input",
			args:        []string{"json2yaml"},
			input:       `{"key": "value"}`,
			checkOutput: func(output string) bool { return strings.Contains(output, "key:") && strings.Contains(output, "value") },
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
			if !tt.checkOutput(output) {
				t.Errorf("json2yaml command output check failed.\nOutput: %s\nDescription: %s",
					output, tt.description)
			}
		})
	}
}
