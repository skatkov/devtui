package textanalyzer

import (
	"reflect"
	"testing"
)

func TestTextAnalyzer_AnalyzeString(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected TextStats
	}{
		{
			name: "Empty text",
			text: "",
			expected: TextStats{
				Text:       "",
				Characters: 0,
				Words:      0,
				Spaces:     0,
			},
		},
		{
			name: "Simple sentence",
			text: "Hello, world!",
			expected: TextStats{
				Text:       "Hello, world!",
				Characters: 13,
				Words:      2,
				Spaces:     1,
			},
		},
		{
			name: "Multiple sentences",
			text: "This is a test. It has multiple sentences!",
			expected: TextStats{
				Text:       "This is a test. It has multiple sentences!",
				Characters: 42,
				Words:      8,
				Spaces:     7,
			},
		},
		{
			name: "Multiple paragraphs",
			text: `This is the first paragraph.\n\nThis is the second paragraph.`,
			expected: TextStats{
				Text:       `This is the first paragraph.\n\nThis is the second paragraph.`,
				Characters: 61,
				Words:      10,
				Spaces:     11, // Including newlines
			},
		},
		{
			name: "Text without sentence terminators",
			text: "This text has no sentence terminators",
			expected: TextStats{
				Text:       "This text has no sentence terminators",
				Characters: 37,
				Words:      6,
				Spaces:     5,
			},
		},
		{
			name: "Text with multiple spaces",
			text: "This   has   extra   spaces",
			expected: TextStats{
				Text:       "This   has   extra   spaces",
				Characters: 28,
				Words:      4,
				Spaces:     12,
			},
		},
		{
			name: "Text with multiple consecutive newlines",
			text: `Paragraph one.\n\n\n\nParagraph two.`,
			expected: TextStats{
				Text:       `Paragraph one.\n\n\n\nParagraph two.`,
				Characters: 32,
				Words:      4,
				Spaces:     5, // Including newlines
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Analyze(tt.text)

			if err != nil {
				t.Errorf("TextAnalyzer.AnalyzeString() error = %v", err)
				return
			}

			// Compare field by field for better error messages
			if got.Text != tt.expected.Text {
				t.Errorf("Text = %v, want %v", got.Text, tt.expected.Text)
			}
			if got.Characters != tt.expected.Characters {
				t.Errorf("Characters = %v, want %v", got.Characters, tt.expected.Characters)
			}
			if got.Words != tt.expected.Words {
				t.Errorf("Words = %v, want %v", got.Words, tt.expected.Words)
			}
			if got.Spaces != tt.expected.Spaces {
				t.Errorf("Spaces = %v, want %v", got.Spaces, tt.expected.Spaces)
			}
		})
	}
}

func TestTextAnalyzer_Analyze(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected TextStats
	}{
		{
			name: "ASCII text as bytes",
			text: "Hello, world!",
			expected: TextStats{
				Text:       "Hello, world!",
				Characters: 13,
				Words:      2,
				Spaces:     1,
			},
		},
		{
			name: "Text with Japanese characters (UTF-8)",
			text: "こんにちは世界！",
			expected: TextStats{
				Text:       "こんにちは世界！",
				Characters: 8, // 7 characters + 1 full-width exclamation
				Words:      1, // Fields will count this as 1 word since no spaces
				Spaces:     0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Analyze(tt.text)

			if err != nil {
				t.Errorf("TextAnalyzer.Analyze() error = %v", err)
				return
			}

			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("TextAnalyzer.Analyze() = %+v, want %+v", got, tt.expected)
			}
		})
	}
}
