package input

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestReadFromStdin(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "simple text",
			input:   "hello world",
			want:    "hello world",
			wantErr: false,
		},
		{
			name:    "empty input",
			input:   "",
			want:    "",
			wantErr: false,
		},
		{
			name:    "multiline input",
			input:   "line1\nline2\nline3",
			want:    "line1\nline2\nline3",
			wantErr: false,
		},
		{
			name:    "json input",
			input:   `{"key":"value"}`,
			want:    `{"key":"value"}`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			cmd.SetIn(strings.NewReader(tt.input))

			got, err := ReadFromStdin(cmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadFromStdin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if string(got) != tt.want {
				t.Errorf("ReadFromStdin() = %q, want %q", string(got), tt.want)
			}
		})
	}
}

func TestReadFromArgsOrStdin(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		stdin   string
		want    string
		wantErr bool
	}{
		{
			name:    "read from args",
			args:    []string{"hello world"},
			stdin:   "",
			want:    "hello world",
			wantErr: false,
		},
		{
			name:    "read from stdin when no args",
			args:    []string{},
			stdin:   "from stdin",
			want:    "from stdin",
			wantErr: false,
		},
		{
			name:    "prefer args over stdin",
			args:    []string{"from args"},
			stdin:   "from stdin",
			want:    "from args",
			wantErr: false,
		},
		{
			name:    "empty args, empty stdin",
			args:    []string{},
			stdin:   "",
			want:    "",
			wantErr: false,
		},
		{
			name:    "multiple args, use first",
			args:    []string{"first", "second", "third"},
			stdin:   "",
			want:    "first",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			cmd.SetIn(strings.NewReader(tt.stdin))

			got, err := ReadFromArgsOrStdin(cmd, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadFromArgsOrStdin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("ReadFromArgsOrStdin() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestReadBytesFromArgsOrStdin(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		stdin   string
		want    []byte
		wantErr bool
	}{
		{
			name:    "read bytes from args",
			args:    []string{"hello"},
			stdin:   "",
			want:    []byte("hello"),
			wantErr: false,
		},
		{
			name:    "read bytes from stdin",
			args:    []string{},
			stdin:   "from stdin",
			want:    []byte("from stdin"),
			wantErr: false,
		},
		{
			name:    "binary-like data from args",
			args:    []string{"test\x00data"},
			stdin:   "",
			want:    []byte("test\x00data"),
			wantErr: false,
		},
		{
			name:    "empty bytes",
			args:    []string{},
			stdin:   "",
			want:    []byte(""),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			cmd.SetIn(strings.NewReader(tt.stdin))

			got, err := ReadBytesFromArgsOrStdin(cmd, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadBytesFromArgsOrStdin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if string(got) != string(tt.want) {
				t.Errorf("ReadBytesFromArgsOrStdin() = %q, want %q", got, tt.want)
			}
		})
	}
}
