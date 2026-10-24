package csv2md

import (
	"fmt"
	"strings"
)

// Convert formats data from file or stdin as markdown
func Convert(header string, records [][]string, aligned bool) []string {
	var result []string

	// add h1 if passed
	header = strings.Trim(header, "\t\r\n ")
	if len(header) != 0 {
		result = append(result, "# "+header, "")
	}

	if len(records) == 0 {
		return result
	}

	// if user wants aligned columns width then we
	// count max length of every value in every column
	widths := make(map[int]int)
	if aligned {
		for _, row := range records {
			for colIdx, col := range row {
				widths[colIdx] = max(widths[colIdx], len(col))
			}
		}
	}

	// content separator only after first row (header)
	result = append(result, tableRow(records[0], widths, aligned), separatorRow(records[0], widths, aligned))
	for _, row := range records[1:] {
		result = append(result, tableRow(row, widths, aligned))
	}
	return result
}

func tableRow(row []string, widths map[int]int, aligned bool) string {
	var builder strings.Builder
	builder.WriteString("| ")
	for colIdx, col := range row {
		if aligned {
			fmt.Fprintf(&builder, "%-*s | ", widths[colIdx], col)
		} else {
			builder.WriteString(col + " | ")
		}
	}
	return builder.String()
}

func separatorRow(row []string, widths map[int]int, aligned bool) string {
	var builder strings.Builder
	builder.WriteString("| ")
	for colIdx := range row {
		if aligned {
			builder.WriteString(strings.Repeat("-", max(widths[colIdx], 3)) + " | ")
		} else {
			builder.WriteString("--- | ")
		}
	}
	return builder.String()
}

func Print(data []string) {
	for _, row := range data {
		fmt.Println(row)
	}
}
