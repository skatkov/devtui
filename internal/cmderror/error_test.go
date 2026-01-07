package cmderror

import (
	"errors"
	"strings"
	"testing"
)

func TestLooksLikeFilePath(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"json extension", "config.json", true},
		{"toml extension", "settings.toml", true},
		{"yaml extension", "data.yaml", true},
		{"yml extension", "config.yml", true},
		{"xml extension", "feed.xml", true},
		{"csv extension", "data.csv", true},
		{"tsv extension", "data.tsv", true},
		{"gql extension", "query.gql", true},
		{"graphql extension", "query.graphql", true},
		{"css extension", "style.css", true},
		{"html extension", "page.html", true},
		{"htm extension", "page.htm", true},
		{"md extension", "readme.md", true},
		{"txt extension", "file.txt", true},
		{"path with directory", "path/to/file.json", true},
		{"relative path", "./config.json", true},
		{"parent directory", "../config.json", true},
		{"absolute unix path", "/etc/config.toml", true},

		{"json object", `{"key": "value"}`, false},
		{"json array", `["item1", "item2"]`, false},
		{"toml content", `name = "value"`, false},
		{"xml content", `<root><item/></root>`, false},
		{"markdown heading", `# Title`, false},
		{"very long string", strings.Repeat("a", 600), false},
		{"empty string", "", false},
		{"whitespace only", "   ", false},
		{"no extension simple word", "hello", false},
		{"json-like but no braces", `key: value`, false},
		{"mixed case extension", "CONFIG.JSON", true},
		{"extension with spaces", " config.json ", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LooksLikeFilePath(tt.input)
			if got != tt.want {
				t.Errorf("LooksLikeFilePath(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestFormatParseError(t *testing.T) {
	baseErr := errors.New("invalid character 'c' in literal")

	tests := []struct {
		name        string
		command     string
		input       string
		wantContain []string
		wantExclude []string
	}{
		{
			name:    "file path shows hint",
			command: "json2toml",
			input:   "config.json",
			wantContain: []string{
				"invalid character",
				"looks like a file path",
				"devtui json2toml < config.json",
			},
			wantExclude: []string{},
		},
		{
			name:    "actual content no hint",
			command: "json2toml",
			input:   `{"invalid": }`,
			wantContain: []string{
				"invalid character",
			},
			wantExclude: []string{
				"looks like a file path",
				"Hint:",
			},
		},
		{
			name:    "toml2json hint",
			command: "toml2json",
			input:   "config.toml",
			wantContain: []string{
				"devtui toml2json < config.toml",
			},
			wantExclude: []string{},
		},
		{
			name:    "gqlquery hint",
			command: "gqlquery",
			input:   "query.gql",
			wantContain: []string{
				"devtui gqlquery < query.gql",
			},
			wantExclude: []string{},
		},
		{
			name:    "cssfmt hint",
			command: "cssfmt",
			input:   "style.css",
			wantContain: []string{
				"devtui cssfmt < style.css",
			},
			wantExclude: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := FormatParseError(tt.command, tt.input, baseErr)
			errStr := err.Error()

			for _, want := range tt.wantContain {
				if !strings.Contains(errStr, want) {
					t.Errorf("error should contain %q, got: %s", want, errStr)
				}
			}

			for _, exclude := range tt.wantExclude {
				if strings.Contains(errStr, exclude) {
					t.Errorf("error should not contain %q, got: %s", exclude, errStr)
				}
			}
		})
	}
}
