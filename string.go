package parcon

import (
	"strconv"
)

// Pre-defined parsers for a single rune.
var (
	// A single space or tab character.
	SingleSpace = OneOfS("SPACE", " \t")

	// A single newline character.
	SingleNewline = OneOfS("NEWLINE", "\r\n")

	// A single space, tab, or new line character.
	SingleSpaceOrNewline = OneOfS("SPACE_OR_NEWLINE", " \t\r\n")

	// A single latin alphabet.
	// This parser is case in-sensitive.
	SingleAlpha = TakeSingle("ALPHA", func(c rune) bool {
		return ('A' <= c && c <= 'Z') || ('a' <= c && c <= 'z')
	})

	// A single character of decimal number.
	SingleDigit = TakeSingle("DIGIT", func(c rune) bool {
		return '0' <= c && c <= '9'
	})

	// A single character of hex number.
	SingleHexDigit = TakeSingle("HEX_DIGIT", func(c rune) bool {
		return ('0' <= c && c <= '9') || ('a' <= c && c <= 'f') || ('A' <= c && c <= 'F')
	})

	// A single latin alphabet or a single decimal digit.
	SingleAlphaNum = TakeSingle("ALPHA_NUM", func(c rune) bool {
		return ('0' <= c && c <= '9') || ('A' <= c && c <= 'Z') || ('a' <= c && c <= 'z')
	})
)

// Pre-defined parsers for a slice of runes.
var (
	// A sequence of multiple SingleSpace.
	MultiSpaces = Many[rune, rune](1, SingleSpace)

	// A sequence of multiple SingleNewline.
	MultiNewline = Many[rune, rune](1, SingleNewline)

	// A sequence of multiple SingleSpaceOrNewline.
	MultiSpacesOrNewlines = Many[rune, rune](1, SingleSpaceOrNewline)

	// A sequence of multiple SingleAlpha.
	MultiAlphas = Many[rune, rune](1, SingleAlpha)

	// A sequence of multiple SingleDigit.
	MultiDigits = Many[rune, rune](1, SingleDigit)

	// A sequence of multiple SingleHexDigit.
	MultiHexDigits = Many[rune, rune](1, SingleHexDigit)

	// A sequence of SingleAlphaNum.
	MultiAlphaNums = Many[rune, rune](1, SingleAlphaNum)
)

// ToString is a ConvertFunc to convert []rune to string.
func ToString(input []rune) (string, error) {
	return string(input), nil
}

// ToInt is a ConvertFunc to convert []rune to int.
func ToInt(input []rune) (int, error) {
	return strconv.Atoi(string(input))
}

// ToFloat is a ConvertFunc to convert []rune to float64.
func ToFloat(input []rune) (float64, error) {
	return strconv.ParseFloat(string(input), 64)
}
