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

func FormatParseError(command, input string, err error) error {
	// Trim whitespace for file path detection and display
	trimmedInput := strings.TrimSpace(input)

	if !LooksLikeFilePath(trimmedInput) {
		return err
	}

	return fmt.Errorf("%w\n\nHint: '%s' looks like a file path.\n\nTo read from a file, use:\n  devtui %s < %s",
		err, trimmedInput, command, trimmedInput)
}
