package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestJSONFmtCmdOutput(t *testing.T) {
	tests := []struct {
		name  string
		input string
		args  []string
		want  string
	}{
		{
			name:  "stdin_without_newline",
			input: `{"a":1}`,
			want:  "{\n  \"a\": 1\n}\n",
		},
		{
			name:  "stdin_with_trailing_whitespace",
			input: "{\"a\":1}\r\n\n \t",
			want:  "{\n  \"a\": 1\n}\n",
		},
		{
			name: "argument",
			args: []string{`[1,2]`},
			want: "[\n  1,\n  2\n]\n",
		},
		{
			name:  "scalar",
			input: "true\n",
			want:  "true\n",
		},
		{
			name:  "invalid_input_unchanged",
			input: `{"a":`,
			want:  `{"a":`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := GetRootCmd()
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)
			cmd.SetIn(strings.NewReader(tt.input))
			cmd.SetArgs(append([]string{"jsonfmt"}, tt.args...))

			if err := cmd.Execute(); err != nil {
				t.Fatalf("jsonfmt command failed: %v", err)
			}
			if got := buf.String(); got != tt.want {
				t.Errorf("jsonfmt output = %q, want %q", got, tt.want)
			}
		})
	}
}
