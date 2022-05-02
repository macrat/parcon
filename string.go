package parcon

import (
	"strconv"
)

// Pre-defined parsers for a single rune.
var (
	// A single space or tab character.
	SingleSpace = OneOf("SPACE", []rune(" \t"))

	// A single newline character.
	SingleNewline = OneOf("NEWLINE", []rune("\r\n"))

	// A single space, tab, or new line character.
	SingleSpaceOrNewline = OneOf("SPACE_OR_NEWLINE", []rune(" \t\r\n"))

	// A single latin alphabet.
	// This parser is case in-sensitive.
	SingleAlpha = TakeSingle("ALPHA", isAlpha)

	// A single character of decimal number.
	SingleDigit = TakeSingle("DIGIT", isDigit)

	// A single character of hex number.
	SingleHexDigit = TakeSingle("HEX_DIGIT", isHexDigit)

	// A single latin alphabet or a single decimal digit.
	SingleAlphaNum = TakeSingle("ALPHA_NUM", isAlphaNum)
)

// Pre-defined parsers for a slice of runes.
var (
	// A sequence of space or tab characters.
	MultiSpaces = OneOfList("SPACE", []rune(" \t"))

	// A sequence of new line characters.
	MultiNewline = OneOfList("NEWLINE", []rune("\r\n"))

	// A sequence of space, tab, or new line characters.
	MultiSpacesOrNewlines = OneOfList("SPACE_OR_NEWLINE", []rune(" \t\r\n"))

	// A sequence of latin alphabets.
	// This parser is case in-sensitive.
	MultiAlphas = TakeWhile("ALPHA", isAlpha)

	// A sequence of decimal digits.
	MultiDigits = TakeWhile("DIGIT", isDigit)

	// A sequence of hex number characters.
	MultiHexDigits = TakeWhile("HEX_DIGIT", isHexDigit)

	// A sequence of latin alphabets or a decimal digits.
	MultiAlphaNums = TakeWhile("ALPHA_NUM", isAlphaNum)
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

func isAlpha(c rune) bool {
	return ('A' <= c && c <= 'Z') || ('a' <= c && c <= 'z')
}

func isDigit(c rune) bool {
	return '0' <= c && c <= '9'
}

func isHexDigit(c rune) bool {
	return ('0' <= c && c <= '9') || ('a' <= c && c <= 'f') || ('A' <= c && c <= 'F')
}

func isAlphaNum(c rune) bool {
	return isAlpha(c) || isDigit(c)
}
