package uuidutil

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Field represents a decoded UUID field.
type Field struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Decode extracts human-readable fields from a UUID.
func Decode(id uuid.UUID) []Field {
	if id == uuid.Nil {
		return nil
	}

	var i big.Int
	i.SetString(strings.ReplaceAll(id.String(), "-", ""), 16)

	fields := []Field{
		{Name: "Standard String Format", Value: id.String()},
		{Name: "Single Integer Value", Value: i.String()},
		{Name: "Version", Value: fmt.Sprintf("%d", id.Version())},
		{Name: "Variant", Value: mapVariant(id.Variant())},
	}

	switch id.Version() {
	case uuid.Version(1):
		timestamp := id.Time()
		sec, nsec := timestamp.UnixTime()
		timeStamp := time.Unix(sec, nsec)
		node := id.NodeID()
		clockSeq := id.ClockSequence()

		fields = append(fields,
			Field{Name: "Contents - Time", Value: timeStamp.UTC().Format("2006-01-02 15:04:05.999999999 UTC")},
			Field{Name: "Contents - Clock", Value: strconv.Itoa(clockSeq)},
			Field{Name: "Contents - Node", Value: fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", node[0], node[1], node[2], node[3], node[4], node[5])},
		)
	default:
		formatted := strings.ToUpper(strings.ReplaceAll(id.String(), "-", ""))
		pairs := make([]string, 0, len(formatted)/2)
		for index := 0; index < len(formatted); index += 2 {
			if index < len(formatted)-1 {
				pairs = append(pairs, formatted[index:index+2])
			}
		}
		fields = append(fields, Field{Name: "Contents", Value: strings.Join(pairs, ":")})
	}

	return fields
}

// FieldsToRows converts decoded fields into table rows.
func FieldsToRows(fields []Field) [][]string {
	rows := make([][]string, 0, len(fields))
	for _, field := range fields {
		rows = append(rows, []string{field.Name, field.Value})
	}
	return rows
}

// Generate creates a UUID for the requested version.
func Generate(version int, namespace string) (uuid.UUID, error) {
	switch version {
	case 1:
		return uuid.NewUUID()
	case 2:
		return uuid.NewDCEGroup()
	case 3:
		return uuid.NewMD5(uuid.NameSpaceURL, []byte(namespace)), nil
	case 4:
		return uuid.NewRandom()
	case 5:
		return uuid.NewSHA1(uuid.NameSpaceURL, []byte(namespace)), nil
	case 6:
		return uuid.NewV6()
	case 7:
		return uuid.NewV7()
	default:
		return uuid.Nil, fmt.Errorf("unsupported uuid version: %d", version)
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
