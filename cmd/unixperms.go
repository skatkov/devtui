package cmd

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// permBit describes one permission bit.
type permBit struct {
	symbol string
	name   string
	value  int
}

var permBits = []permBit{
	{"r", "owner read", 0o400},
	{"w", "owner write", 0o200},
	{"x", "owner execute", 0o100},
	{"r", "group read", 0o040},
	{"w", "group write", 0o020},
	{"x", "group execute", 0o010},
	{"r", "other read", 0o004},
	{"w", "other write", 0o002},
	{"x", "other execute", 0o001},
}

// octalToSymbolic converts an octal permission mode to a 9-char symbolic string (e.g. "rwxr-xr-x").
func octalToSymbolic(mode int) string {
	var sb strings.Builder
	for _, bit := range permBits {
		if mode&bit.value != 0 {
			sb.WriteString(bit.symbol)
		} else {
			sb.WriteString("-")
		}
	}
	return sb.String()
}

// symbolicToOctal parses a 9-char symbolic string and returns the octal mode.
func symbolicToOctal(sym string) (int, error) {
	if len(sym) != 9 {
		return 0, errors.New("symbolic notation must be exactly 9 characters (e.g. rwxr-xr-x)")
	}
	mode := 0
	for i, bit := range permBits {
		ch := string(sym[i])
		if ch == bit.symbol {
			mode |= bit.value
		} else if ch != "-" {
			return 0, fmt.Errorf("invalid character %q at position %d (expected %q or \"-\")", ch, i+1, bit.symbol)
		}
	}
	return mode, nil
}

// explainMode prints a human-readable breakdown of the permission mode.
func explainMode(mode int) {
	symbolic := octalToSymbolic(mode)
	fmt.Printf("Octal:    %04o\n", mode)
	fmt.Printf("Symbolic: %s\n\n", symbolic)

	sections := []struct {
		label string
		bits  []permBit
		start int
	}{
		{"Owner", permBits[0:3], 0},
		{"Group", permBits[3:6], 3},
		{"Other", permBits[6:9], 6},
	}

	for _, sec := range sections {
		var granted []string
		for _, bit := range sec.bits {
			if mode&bit.value != 0 {
				granted = append(granted, bit.symbol)
			} else {
				granted = append(granted, "-")
			}
		}
		fmt.Printf("%s: %s  (%s)\n", sec.label, strings.Join(granted, ""), describeSection(mode, sec.start))
	}
}

func describeSection(mode, start int) string {
	bits := permBits[start : start+3]
	var parts []string
	for _, bit := range bits {
		if mode&bit.value != 0 {
			switch bit.symbol {
			case "r":
				parts = append(parts, "read")
			case "w":
				parts = append(parts, "write")
			case "x":
				parts = append(parts, "execute")
			}
		}
	}
	if len(parts) == 0 {
		return "no permissions"
	}
	return strings.Join(parts, ", ")
}

var unixpermsCmd = &cobra.Command{
	Use:   "unixperms <mode>",
	Short: "Explain Unix file permission modes",
	Long: `Explain Unix file permission modes.

Accepts either octal notation (e.g. 755, 0644) or 9-character symbolic
notation (e.g. rwxr-xr-x). Prints a human-readable breakdown of each
permission bit for owner, group, and other.`,
	Example: `  # Explain octal mode 755
  devtui unixperms 755

  # Explain octal mode with leading zero
  devtui unixperms 0644

  # Explain symbolic notation
  devtui unixperms rwxr-xr-x`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		input := strings.TrimSpace(args[0])

		var mode int

		// Detect whether input is octal (all digits, optional leading 0) or symbolic.
		if strings.ContainsAny(input, "rwx-") || len(input) == 9 {
			// Symbolic notation
			m, err := symbolicToOctal(input)
			if err != nil {
				return err
			}
			mode = m
		} else {
			// Octal notation
			stripped := strings.TrimPrefix(input, "0")
			if stripped == "" {
				stripped = "0"
			}
			m, err := strconv.ParseInt(stripped, 8, 64)
			if err != nil {
				return fmt.Errorf("invalid octal mode %q: %v", input, err)
			}
			if m < 0 || m > 0o777 {
				return fmt.Errorf("mode %04o out of range (must be 000-777)", m)
			}
			mode = int(m)
		}

		explainMode(mode)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(unixpermsCmd)
}
