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
		args        []string
		wantContain string
		wantErr     bool
		description string
	}{
		{
			name:        "simple object conversion",
			input:       `{"item": "value"}`,
			args:        []string{},
			wantContain: "<item>",
			wantErr:     false,
			description: "Should convert simple JSON object to XML",
		},
		{
			name:        "nested objects conversion",
			input:       `{"root": {"child": {"value": "data"}}}`,
			args:        []string{},
			wantContain: "<root>",
			wantErr:     false,
			description: "Should convert nested JSON objects to XML hierarchy",
		},
		{
			name:        "arrays conversion",
			input:       `{"items": ["apple", "banana"]}`,
			args:        []string{},
			wantContain: "<items>",
			wantErr:     false,
			description: "Should convert JSON arrays to XML",
		},
		{
			name:        "mixed types conversion",
			input:       `{"string": "text", "number": 123, "boolean": true}`,
			args:        []string{},
			wantContain: "<string>",
			wantErr:     false,
			description: "Should handle mixed JSON types",
		},
		{
			name:        "empty object",
			input:       `{}`,
			args:        []string{},
			wantContain: "<doc/>",
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

			args := []string{"json2xml"}
			args = append(args, tt.args...)
			cmd.SetArgs(args)

			err := cmd.Execute()

			if (err != nil) != tt.wantErr {
				t.Errorf("json2xml command error = %v, wantErr %v\nDescription: %s", err, tt.wantErr, tt.description)
				return
			}

			if !tt.wantErr && tt.wantContain != "" {
				output := buf.String()
				if !strings.Contains(output, tt.wantContain) {
					t.Errorf("json2xml output does not contain expected string %q\nGot: %s\nDescription: %s",
						tt.wantContain, output, tt.description)
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
		wantContain string
		description string
	}{
		{
			name:        "argument input",
			args:        []string{"json2xml", `{"key": "value"}`},
			input:       "",
			wantContain: "<key>",
			description: "Should handle JSON string argument",
		},
		{
			name:        "stdin input",
			args:        []string{"json2xml"},
			input:       `{"key": "value"}`,
			wantContain: "<key>",
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
			if !strings.Contains(output, tt.wantContain) {
				t.Errorf("json2xml output does not contain expected string %q\nGot: %s\nDescription: %s",
					tt.wantContain, output, tt.description)
			}
		})
	}
}
