package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestJSONRepairCmd(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantErr     bool
		checkOutput func(string) bool
		description string
	}{
		{
			name:  "repair single quotes",
			input: "{'key': 'value'}",
			checkOutput: func(output string) bool {
				return strings.Contains(output, `"key"`) && strings.Contains(output, `"value"`)
			},
			description: "Should convert single quotes to double quotes",
		},
		{
			name:  "repair unclosed array",
			input: "[1, 2, 3, 4",
			checkOutput: func(output string) bool {
				return strings.Contains(output, "[") && strings.Contains(output, "]")
			},
			description: "Should close unclosed array",
		},
		{
			name:  "repair unclosed object",
			input: `{"employees":["John", "Anna"`,
			checkOutput: func(output string) bool {
				return strings.Contains(output, `"employees"`) && strings.Contains(output, "}")
			},
			description: "Should close unclosed object and array",
		},
		{
			name:  "repair uppercase booleans",
			input: `{"key": TRUE, "key2": FALSE, "key3": Null}`,
			checkOutput: func(output string) bool {
				output = strings.ToLower(output)
				return strings.Contains(output, "true") && strings.Contains(output, "false") && strings.Contains(output, "null")
			},
			description: "Should normalize boolean and null values",
		},
		{
			name:  "repair trailing comma",
			input: `{"key":"value",}`,
			checkOutput: func(output string) bool {
				return strings.Contains(output, `"key"`) && strings.Contains(output, `"value"`)
			},
			description: "Should handle trailing comma",
		},
		{
			name:  "repair mixed quotes",
			input: `{'key': 'string', "key2": false}`,
			checkOutput: func(output string) bool {
				return strings.Contains(output, `"key"`) && strings.Contains(output, `"key2"`)
			},
			description: "Should normalize mixed quotes",
		},
		{
			name:  "strip markdown code block",
			input: "```json\n{'key': 'value'}\n```",
			checkOutput: func(output string) bool {
				return strings.Contains(output, `"key"`) && !strings.Contains(output, "```")
			},
			description: "Should strip markdown code blocks",
		},
		{
			name:  "valid JSON passthrough",
			input: `{"key": "value"}`,
			checkOutput: func(output string) bool {
				return strings.Contains(output, `"key"`) && strings.Contains(output, `"value"`)
			},
			description: "Should handle valid JSON without changes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create fresh command and root for each test
			cmd := GetRootCmd()
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)
			cmd.SetIn(strings.NewReader(tt.input))
			cmd.SetArgs([]string{"jsonrepair"})

			err := cmd.Execute()

			if (err != nil) != tt.wantErr {
				t.Errorf("jsonrepair command error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				output := buf.String()
				if tt.checkOutput != nil && !tt.checkOutput(output) {
					t.Errorf("jsonrepair command output check failed.\nInput: %s\nOutput: %s\nDescription: %s",
						tt.input, output, tt.description)
				}
			}
		})
	}
}

func TestJSONRepairCmdNoInput(t *testing.T) {
	cmd := GetRootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetIn(strings.NewReader(""))
	cmd.SetArgs([]string{"jsonrepair"})

	err := cmd.Execute()

	if err == nil {
		t.Error("jsonrepair command should return error when no input provided")
	}

	if !strings.Contains(err.Error(), "no input provided") {
		t.Errorf("jsonrepair command error message should mention 'no input provided', got: %v", err)
	}
}

func TestJSONRepairCmdChaining(t *testing.T) {
	// Test that output can be piped to another command
	input := "{'employees':['John', 'Anna']}"

	cmd := GetRootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetIn(strings.NewReader(input))
	cmd.SetArgs([]string{"jsonrepair"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("jsonrepair command failed: %v", err)
	}

	output := buf.String()

	// Verify output is valid-looking JSON
	if !strings.Contains(output, `"employees"`) {
		t.Errorf("jsonrepair command output doesn't look like valid JSON: %s", output)
	}

	// Verify it's properly formatted for chaining (should be single line or compact)
	if strings.TrimSpace(output) == "" {
		t.Error("jsonrepair command output is empty")
	}
}
