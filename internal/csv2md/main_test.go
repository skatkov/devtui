package csv2md

import (
	"slices"
	"testing"
)

func TestConvert(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		header  string
		records [][]string
		aligned bool
		want    []string
	}{
		{
			name:    "table without header",
			records: [][]string{{"a", "b"}, {"c", "d"}, {"e", "f"}},
			want:    []string{"| a | b | ", "| --- | --- | ", "| c | d | ", "| e | f | "},
		},
		{
			name:    "header is trimmed",
			header:  " \tT\r\n",
			records: [][]string{{"a"}},
			want:    []string{"# T", "", "| a | ", "| --- | "},
		},
		{
			name:    "blank header is dropped",
			header:  " \t\r\n",
			records: [][]string{{"a"}},
			want:    []string{"| a | ", "| --- | "},
		},
		{
			name:   "header without records",
			header: "Title",
			want:   []string{"# Title", ""},
		},
		{
			name: "no header and no records",
			want: nil,
		},
		{
			name:    "aligned columns",
			records: [][]string{{"name", "x", "abc"}, {"alexander", "y", "z"}},
			aligned: true,
			want: []string{
				"| name      | x | abc | ",
				"| --------- | --- | --- | ",
				"| alexander | y | z   | ",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := Convert(tt.header, tt.records, tt.aligned)
			if !slices.Equal(got, tt.want) {
				t.Fatalf("Convert() = %q, want %q", got, tt.want)
			}
		})
	}
}
