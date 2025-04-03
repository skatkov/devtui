package textanalyzer

import (
	"fmt"
	"strings"
	"unicode"
)

// TextStats represents the statistical information extracted from a text
type TextStats struct {
	Text       string
	Characters int
	Words      int
	Spaces     int
}

// Analyze processes the given text and returns statistical information
func Analyze(text string) (TextStats, error) {
	stats := TextStats{Text: text}

	// Count characters (runes)
	runes := []rune(text)
	stats.Characters = len(runes)

	// Count spaces
	for _, r := range runes {
		if unicode.IsSpace(r) {
			stats.Spaces++
		}
	}

	// Count words
	fields := strings.Fields(text)
	fmt.Println(fields)
	stats.Words = len(fields)

	return stats, nil
}
