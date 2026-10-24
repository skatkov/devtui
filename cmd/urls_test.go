package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestUrlsCmd(t *testing.T) {
	tests := []struct {
		name  string
		input string
		args  []string
		want  string
	}{
		{
			name:  "relaxed mode keeps first occurrence of each url",
			input: "Visit google.com and https://example.com, then google.com again",
			want:  "google.com\nhttps://example.com\n",
		},
		{
			name:  "strict mode requires a scheme",
			input: "Visit google.com and https://example.com and http://example.org",
			args:  []string{"--strict"},
			want:  "https://example.com\nhttp://example.org\n",
		},
		{
			name: "argument input",
			args: []string{"see https://example.com"},
			want: "https://example.com\n",
		},
		{
			name:  "no urls prints nothing",
			input: "nothing to see here",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strictFlag := urlsCmd.Flags().Lookup("strict")
			if err := strictFlag.Value.Set(strictFlag.DefValue); err != nil {
				t.Fatalf("reset strict flag: %v", err)
			}

			cmd := GetRootCmd()
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)
			cmd.SetIn(strings.NewReader(tt.input))
			cmd.SetArgs(append([]string{"urls"}, tt.args...))

			if err := cmd.Execute(); err != nil {
				t.Fatalf("urls command failed: %v", err)
			}
			if got := buf.String(); got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}
