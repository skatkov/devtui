package cmderror

import (
	"fmt"
	"path/filepath"
	"strings"
)

var knownExtensions = []string{
	".json", ".toml", ".yaml", ".yml", ".xml", ".csv", ".tsv",
	".gql", ".graphql", ".css", ".html", ".htm", ".md", ".txt",
}

func LooksLikeFilePath(input string) bool {
	if len(input) > 500 {
		return false
	}

	trimmed := strings.TrimSpace(input)
	if len(trimmed) > 0 {
		first := trimmed[0]
		if first == '{' || first == '[' || first == '<' || first == '#' {
			return false
		}
	}

	if strings.Contains(input, "/") || strings.Contains(input, "\\") {
		return true
	}

	ext := strings.ToLower(filepath.Ext(trimmed))
	for _, knownExt := range knownExtensions {
		if ext == knownExt {
			return true
		}
	}

	return false
}

func FormatParseError(format, command, input string, err error) error {
	if !LooksLikeFilePath(input) {
		return fmt.Errorf("%s parsing error: %w", format, err)
	}

	return fmt.Errorf("%s parsing error: %w\n\nHint: '%s' looks like a file path. To read from a file, use:\n  devtui %s < %s",
		format, err, input, command, input)
}
