package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

type Number struct {
	Value int64
	Type  NumberType
}

type NumberType string

const (
	purple    = lipgloss.Color("99")
	gray      = lipgloss.Color("245")
	lightGray = lipgloss.Color("241")
)
const (
	Base2  NumberType = "Base 2 (binary)"
	Base8  NumberType = "Base 8 (octal)"
	Base10 NumberType = "Base 10 (decimal)"
	Base16 NumberType = "Base 16 (hexadecimal)"
)

var (
	re = lipgloss.NewRenderer(os.Stdout)
	// HeaderStyle is the lipgloss style used for the table headers.
	HeaderStyle = re.NewStyle().Foreground(purple).Bold(true).Align(lipgloss.Center)
	// CellStyle is the base lipgloss style used for the table rows.
	CellStyle = re.NewStyle().Padding(0, 1).Width(14)
	// OddRowStyle is the lipgloss style used for odd-numbered table rows.
	OddRowStyle = CellStyle.Foreground(gray)
	// EvenRowStyle is the lipgloss style used for even-numbered table rows.
	EvenRowStyle = CellStyle.Foreground(lightGray)
	// BorderStyle is the lipgloss style used for the table border.
	BorderStyle = lipgloss.NewStyle().Foreground(purple)
)

func main() {
	number := Number{Type: Base10}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Placeholder("Enter a number").
				Validate(func(s string) error {
					if len(s) == 0 {
						return errors.New("number cannot be empty")
					}
					val, err := strconv.ParseInt(s, 10, 64)
					if err != nil {
						return errors.New("please enter a valid integer")
					}
					number.Value = val
					return nil
				}).Value(new(string)),
		),
	)

	err := form.Run()

	if err != nil {
		fmt.Println("Uh oh:", err)
		os.Exit(1)
	}

	rows := [][]string{
		{string(number.Type), strconv.FormatInt(number.Value, 10)},
		{string(Base2), strconv.FormatInt(number.Value, 2)},
		{string(Base8), strconv.FormatInt(number.Value, 8)},
		{string(Base16), strconv.FormatInt(number.Value, 16)},
	}
	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(BorderStyle).
		Width(120).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == table.HeaderRow:
				return HeaderStyle
			case row%2 == 0:
				return EvenRowStyle
			default:
				return OddRowStyle
			}
		}).
		Headers("Base", "Value").
		Rows(rows...)

	fmt.Println(t)
}
