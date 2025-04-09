package textanalyzer

import (
	"strings"
	"unicode"
)

// TextStats represents the statistical information extracted from a text
type TextStats struct {
	Characters int
	Words      int
	Spaces     int
}

// Analyze processes the given text and returns statistical information
func Analyze(text string) (TextStats, error) {
	stats := TextStats{}

	stats.Characters = len([]rune(text))

	// Count spaces
	for _, r := range text {
		if unicode.IsSpace(r) {
			stats.Spaces++
		}
	}

	// Count words
	fields := strings.Fields(text)
	stats.Words = len(fields)

	return stats, nil
}
