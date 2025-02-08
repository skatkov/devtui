package main

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
)

func TestExtractUUIDData(t *testing.T) {
	// Test cases
	tests := []struct {
		name     string
		uuid     uuid.UUID
		expected [][]string
	}{
		{
			name: "Valid UUID v4",
			uuid: uuid.MustParse("4326ff5f-774d-4506-a18c-4bc50c761863"),
			expected: [][]string{
				{"Standard String Format", "123e4567-e89b-12d3-a456-426614174000"},
				{"Single Integer Value", "89260762576260186387968344354060638307"},
				{"Version", "4"},
				{"Variant", "DCE 1.1, ISO/IEC 11578:1996"},
				{"Contents", "43:26:FF:5F:77:4D:05:06:21:8C:4B:C5:0C:76:18:63"},
			},
		},
		{
			name: "Valid UUID v1",
			uuid: uuid.MustParse("550e8400-e29b-11d4-a716-446655440000"),
			expected: [][]string{
				{"Standard String Format", "550e8400-e29b-11d4-a716-446655440000"},
				{"Single Integer Value", "113059749145936098728763079434011148288"},
				{"Version", "1"},
				{"Variant", "DCE 1.1, ISO/IEC 11578:1996"},
				{"Contents - Time", "2001-01-04 23:43:07.540992 UTC"},
				{"Contents - Clock", "10006"},
				{"Contents - Node", "44:66:55:44:00:00"},
			},
		},
		// {
		// 	name: "UUID v2",
		// 	uuid: uuid.MustParse("000003e8-e64c-21ef-ac00-325096b39f47"),
		// 	expected: [][]string{
		// 		{"Standard String Format", "000003e8-e64c-21ef-ac00-325096b39f47"},
		// 		{"Single Integer Value", "79299436105144797400979033268039"},
		// 		{"Version", "2"},
		// 		{"Variant", "DCE 1.1, ISO/IEC 11578:1996"},
		// 		{"Contents", "00:00:03:E8:E6:4C:01:EF:2C:00:32:50:96:B3:9F:47"},
		// 	},
		// },
		{
			name:     "Nil UUID",
			uuid:     uuid.Nil,
			expected: [][]string{},
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractUUIDData(tt.uuid)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("extractUUIDData() = %v, want %v", got, tt.expected)
			}
		})
	}
}
