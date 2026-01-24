package numbers

import (
	"fmt"
	"strconv"
)

// Base describes a number base label.
type Base struct {
	Label string
	Base  int
}

// Conversion represents a converted value for a base.
type Conversion struct {
	Label string `json:"label"`
	Base  int    `json:"base"`
	Value string `json:"value"`
}

// Result contains conversion results for a number.
type Result struct {
	Input       string       `json:"input"`
	Base        int          `json:"base"`
	Value       int64        `json:"value"`
	Conversions []Conversion `json:"conversions"`
}

// Bases lists supported base conversions.
var Bases = []Base{
	{Label: "Base 2 (binary)", Base: 2},
	{Label: "Base 8 (octal)", Base: 8},
	{Label: "Base 10 (decimal)", Base: 10},
	{Label: "Base 16 (hexadecimal)", Base: 16},
}

// DefaultBase returns the default base (decimal).
func DefaultBase() Base {
	return Bases[2]
}

// Parse parses the input number using the provided base.
func Parse(input string, base int) (int64, error) {
	if !isSupportedBase(base) {
		return 0, fmt.Errorf("unsupported base: %d (supported: 2, 8, 10, 16)", base)
	}

	return strconv.ParseInt(input, base, 64)
}

// Convert converts the input into supported bases.
func Convert(input string, base int) (Result, error) {
	value, err := Parse(input, base)
	if err != nil {
		return Result{}, err
	}

	conversions := make([]Conversion, 0, len(Bases))
	for _, info := range Bases {
		conversions = append(conversions, Conversion{
			Label: info.Label,
			Base:  info.Base,
			Value: strconv.FormatInt(value, info.Base),
		})
	}

	return Result{
		Input:       input,
		Base:        base,
		Value:       value,
		Conversions: conversions,
	}, nil
}

func isSupportedBase(base int) bool {
	for _, info := range Bases {
		if info.Base == base {
			return true
		}
	}
	return false
}
