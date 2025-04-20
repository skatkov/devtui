package csv2md

import (
	"fmt"
	"strings"
)

// Convert formats data from file or stdin as markdown
func Convert(header string, records [][]string, aligned bool) []string {
	// Preallocate result slice to reduce allocations
	// Size estimate: header (if exists) + empty line + records + separator line
	resultSize := len(records) + 1 // +1 for the separator after header row
	if len(header) > 0 {
		resultSize += 2 // Add 2 for header title and blank line
	}
	result := make([]string, 0, resultSize)

	// add h1 if passed
	header = strings.Trim(header, "\t\r\n ")
	if len(header) != 0 {
		result = append(result, "# "+header)
		result = append(result, "")
	}

	// if user wants aligned columns width then we
	// count max length of every value in every column
	widths := make(map[int]int)
	if aligned {
		for _, row := range records {
			for col_idx, col := range row {
				length := len(col)
				if len(widths) == 0 || widths[col_idx] < length {
					widths[col_idx] = length
				}
			}
		}
	}

	// build markdown table
	for row_idx, row := range records {
		// table content
		str := "| "
		for col_idx, col := range row {
			if aligned {
				str += fmt.Sprintf("%-*s | ", widths[col_idx], col)
			} else {
				str += col + " | "
			}
		}
		result = append(result, str)

		// content separator only after first row (header)
		if row_idx == 0 {
			str := "| "
			for col_idx := range row {
				if !aligned || widths[col_idx] < 3 {
					str += "--- | "
				} else {
					str += strings.Repeat("-", widths[col_idx]) + " | "
				}
			}
			result = append(result, str)
		}
	}
	return result
}

func Print(data []string) {
	for _, row := range data {
		fmt.Println(row)
	}
}
