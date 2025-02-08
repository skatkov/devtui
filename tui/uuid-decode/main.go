package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/google/uuid"
)

func main() {
	var parsedUUID string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("UUID").
				Placeholder("Enter a UUID").
				Validate(func(value string) error {
					_, err := uuid.Parse(value)
					return err
				}).Value(&parsedUUID),
		),
	)

	err := form.Run()

	if err != nil {
		fmt.Println("Uh oh:", err)
		os.Exit(1)
	}

	result, _ := uuid.Parse(parsedUUID)
	//t := result.Time()
	// sec, nsec := t.UnixTime()
	//timeStamp := time.Unix(sec, nsec)

	tableOutput := table.New().
		Border(lipgloss.RoundedBorder()).
		Width(100).
		Rows(extractUUIDData(result)...)

	fmt.Println(tableOutput)
}

func extractUUIDData(id uuid.UUID) [][]string {
	if id == uuid.Nil {
		return [][]string{}
	}

	return [][]string{
		{"Standard String Format", id.String()},
		{"Raw Content", fmt.Sprintf("%x", id[:])},
		{"Version", fmt.Sprintf("%d", id.Version())},
		{"Variant", fmt.Sprintf("%s", mapVariant(id.Variant()))},
	}
}

func mapVariant(v uuid.Variant) string {
	switch v {
	case uuid.Invalid:
		return "Invalid UUID"
	case uuid.RFC4122:
		return "DCE 1.1, ISO/IEC 11578:1996"
	case uuid.Reserved:
		return "Reserved (NCS backward compatibility)"
	case uuid.Microsoft:
		return "Reserved (Microsoft GUID)"
	case uuid.Future:
		return "Reserved (future use)"
	default:
		return "Unknown"
	}
}
