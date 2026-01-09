package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestJson2xmlCmd(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		checkOutput func(string) bool
		wantErr     bool
		description string
	}{
		{
			name:  "simple object conversion",
			input: `{"item": "value"}`,
			checkOutput: func(output string) bool {
				return strings.Contains(output, "<item>") &&
					strings.Contains(output, "value") &&
					strings.Contains(output, "</item>")
			},
			description: "Should convert simple JSON object to XML",
		},
		{
			name:  "nested objects conversion",
			input: `{"root": {"child": {"value": "data"}}}`,
			checkOutput: func(output string) bool {
				return strings.Contains(output, "<root>") &&
					strings.Contains(output, "<child>") &&
					strings.Contains(output, "<value>") &&
					strings.Contains(output, "data")
			},
			description: "Should convert nested JSON objects to XML hierarchy",
		},
		{
			name:  "arrays conversion",
			input: `{"items": ["apple", "banana"]}`,
			checkOutput: func(output string) bool {
				return strings.Contains(output, "<items>") &&
					strings.Contains(output, "apple") &&
					strings.Contains(output, "banana")
			},
			description: "Should convert JSON arrays to XML",
		},
		{
			name:  "mixed types conversion",
			input: `{"string": "text", "number": 123, "boolean": true}`,
			checkOutput: func(output string) bool {
				return strings.Contains(output, "<string>") &&
					strings.Contains(output, "<number>") &&
					strings.Contains(output, "<boolean>")
			},
			description: "Should handle mixed JSON types",
		},
		{
			name:  "empty object",
			input: `{}`,
			checkOutput: func(output string) bool {
				return strings.Contains(output, "<doc/>") || strings.TrimSpace(output) == ""
			},
			description: "Should handle empty JSON object",
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
			cmd.SetArgs([]string{"json2xml"})

			err := cmd.Execute()

			if (err != nil) != tt.wantErr {
				t.Errorf("json2xml command error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.checkOutput != nil {
				output := buf.String()
				if !tt.checkOutput(output) {
					t.Errorf("json2xml command output check failed.\nInput: %s\nOutput: %s\nDescription: %s",
						tt.input, output, tt.description)
				}
			}
		})
	}
}

func TestJson2xmlCmdNoInput(t *testing.T) {
	cmd := GetRootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetIn(strings.NewReader(""))
	cmd.SetArgs([]string{"json2xml"})

	err := cmd.Execute()

	if err == nil {
		t.Error("json2xml command should return error when no input provided")
	}

	if !strings.Contains(err.Error(), "no input provided") {
		t.Errorf("json2xml command error message should mention 'no input provided', got: %v", err)
	}
}

func TestJson2xmlCmdArgumentInput(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		input       string
		checkOutput func(string) bool
		description string
	}{
		{
			name:  "argument input",
			args:  []string{"json2xml", `{"key": "value"}`},
			input: "",
			checkOutput: func(output string) bool {
				return strings.Contains(output, "<key>") && strings.Contains(output, "value")
			},
			description: "Should handle JSON string argument",
		},
		{
			name:  "stdin input",
			args:  []string{"json2xml"},
			input: `{"key": "value"}`,
			checkOutput: func(output string) bool {
				return strings.Contains(output, "<key>") && strings.Contains(output, "value")
			},
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
				t.Fatalf("json2xml command failed: %v", err)
			}

			output := buf.String()
			if !tt.checkOutput(output) {
				t.Errorf("json2xml command output check failed.\nOutput: %s\nDescription: %s",
					output, tt.description)
			}
		})
	}
}
