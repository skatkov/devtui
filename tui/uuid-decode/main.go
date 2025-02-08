package main

import (
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

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
	var i big.Int
	i.SetString(strings.Replace(id.String(), "-", "", 4), 16)

	var result [][]string

	result = [][]string{
		{"Standard String Format", id.String()},
		{"Single Integer Value", i.String()},
		{"Version", fmt.Sprintf("%d", id.Version())},
		{"Variant", fmt.Sprintf("%s", mapVariant(id.Variant()))},
	}
	switch id.Version() {
	case uuid.Version(1):
		t := id.Time()
		sec, nsec := t.UnixTime()
		timeStamp := time.Unix(sec, nsec)
		node := id.NodeID()
		clockSeq := id.ClockSequence()

		result = append(result, []string{"Contents - Time", timeStamp.UTC().Format("2006-01-02 15:04:05.999999999 UTC")})
		result = append(result, []string{"Contents - Clock", fmt.Sprintf("%d", clockSeq)})
		result = append(result, []string{"Contents - Node", fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", node[0], node[1], node[2], node[3], node[4], node[5])})

	default:
		formatted := strings.ToUpper(strings.ReplaceAll(id.String(), "-", ""))
		var pairs []string
		for i := 0; i < len(formatted); i += 2 {
			if i < len(formatted)-1 {
				pairs = append(pairs, formatted[i:i+2])
			}
		}
		result = append(result, []string{"Contents", strings.Join(pairs, ":")})
	}

	return result
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
