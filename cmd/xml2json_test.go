package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestXml2jsonCmd(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantContain string
		wantErr     bool
		description string
	}{
		{
			name:        "simple element conversion",
			input:       `<root><item>value</item></root>`,
			wantContain: `"item": "value"`,
			wantErr:     false,
			description: "Should convert simple XML to JSON",
		},
		{
			name:        "nested elements conversion",
			input:       `<root><child><value>data</value></child></root>`,
			wantContain: `"child"`,
			wantErr:     false,
			description: "Should convert nested XML elements to JSON objects",
		},
		{
			name:        "attributes to elements",
			input:       `<root><item id="1">value</item></root>`,
			wantContain: `"item"`,
			wantErr:     false,
			description: "Should convert XML with attributes to JSON",
		},
		{
			name:        "multiple elements",
			input:       `<root><item>first</item><item>second</item></root>`,
			wantContain: `"item"`,
			wantErr:     false,
			description: "Should convert XML with multiple elements to JSON arrays",
		},
		{
			name:        "invalid XML input",
			input:       `<invalid>xml</unclosed`,
			wantContain: "",
			wantErr:     true,
			description: "Should error on invalid XML",
		},
		{
			name:        "empty root",
			input:       `<root></root>`,
			wantContain: `"root"`,
			wantErr:     false,
			description: "Should handle empty XML element",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := GetRootCmd()
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)
			cmd.SetIn(strings.NewReader(tt.input))

			args := []string{"xml2json"}
			cmd.SetArgs(args)

			err := cmd.Execute()

			if (err != nil) != tt.wantErr {
				t.Errorf("xml2json command error = %v, wantErr %v\nDescription: %s", err, tt.wantErr, tt.description)
				return
			}

			if !tt.wantErr && tt.wantContain != "" {
				output := buf.String()
				if !strings.Contains(output, tt.wantContain) {
					t.Errorf("xml2json output does not contain expected string %q\nGot: %s\nDescription: %s",
						tt.wantContain, output, tt.description)
				}
			}
		})
	}
}

func TestXml2jsonCmdNoInput(t *testing.T) {
	cmd := GetRootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetIn(strings.NewReader(""))
	cmd.SetArgs([]string{"xml2json"})

	err := cmd.Execute()

	if err == nil {
		t.Error("xml2json command should return error when no input provided")
	}

	if !strings.Contains(err.Error(), "no input provided") {
		t.Errorf("xml2json command error message should mention 'no input provided', got: %v", err)
	}
}

func TestXml2jsonCmdArgumentInput(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		input       string
		wantContain string
		description string
	}{
		{
			name:        "argument input",
			args:        []string{"xml2json", `<root><item>value</item></root>`},
			input:       "",
			wantContain: `"item": "value"`,
			description: "Should handle XML string argument",
		},
		{
			name:        "stdin input",
			args:        []string{"xml2json"},
			input:       `<root><item>value</item></root>`,
			wantContain: `"item": "value"`,
			description: "Should handle XML from stdin",
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
				t.Fatalf("xml2json command failed: %v", err)
			}

			output := buf.String()
			if !strings.Contains(output, tt.wantContain) {
				t.Errorf("xml2json output does not contain expected string %q\nGot: %s\nDescription: %s",
					tt.wantContain, output, tt.description)
			}
		})
	}
}
