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
	Type  NumberBase
}

type NumberBase struct {
	title string
	base  int
}

var (
	Base2    = NumberBase{title: "Base 2 (binary)", base: 2}
	Base8    = NumberBase{title: "Base 8 (octal)", base: 8}
	Base10   = NumberBase{title: "Base 10 (decimal)", base: 10}
	Base16   = NumberBase{title: "Base 16 (hexadecimal)", base: 16}
	BaseList = []NumberBase{Base2, Base8, Base10, Base16}
)

func main() {
	number := Number{}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[NumberBase]().
				Options(
					huh.NewOption(Base2.title, Base2),
					huh.NewOption(Base8.title, Base8),
					huh.NewOption(Base10.title, Base10),
					huh.NewOption(Base16.title, Base16),
				).
				Title("Select Base").Value(&number.Type),
			huh.NewInput().
				Placeholder(fmt.Sprintf("Enter a %s number", number.Type.title)).
				Title("Enter a number").
				Validate(func(s string) error {
					if len(s) == 0 {
						return errors.New("number cannot be empty")
					}
					val, err := strconv.ParseInt(s, number.Type.base, 64)
					if err != nil {
						return errors.New(fmt.Sprintf("please enter a valid %s number", number.Type.title))
					}
					number.Value = val
					return nil
				}).Value(new(string)),
		),
	).WithTheme(huh.ThemeCharm())

	err := form.Run()

	if err != nil {
		fmt.Println("Uh oh:", err)
		os.Exit(1)
	}
	rows := make([][]string, len(BaseList))
	for i, numberBase := range BaseList {
		rows[i] = []string{
			numberBase.title,
			strconv.FormatInt(number.Value, numberBase.base),
		}
	}
	t := table.New().
		Border(lipgloss.RoundedBorder()).
		Width(120).
		Headers("Base", "Value").
		Rows(rows...)

	fmt.Println(t)
}
