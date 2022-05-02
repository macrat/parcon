package parcon

import (
	"fmt"
)

// ConvertFunc is a function type to convert parsed value.
//
// Normally, it is used for Convert or Map parser.
type ConvertFunc[I any, O any] interface {
	func(input I) (output O, err error)
}

type Converter[I comparable, O1, O2 any, P Parser[I, O1], F ConvertFunc[O1, O2]] struct {
	parser P
	fun    F
}

// Convert converts the output of the given `parser` using given ConvertFunc.
func Convert[I comparable, O1, O2 any, P Parser[I, O1], F ConvertFunc[O1, O2]](parser P, fun F) Converter[I, O1, O2, P, F] {
	return Converter[I, O1, O2, P, F]{parser, fun}
}

func (c Converter[I, O1, O2, P, F]) Parse(input []I) (output O2, remain []I, err error) {
	var o O1
	o, remain, err = c.parser.Parse(input)
	if err != nil {
		return
	}

	output, err = c.fun(o)
	return
}

func (c Converter[I, O1, O2, P, F]) String() string {
	return fmt.Sprint(c.parser)
}

type Mapper[I comparable, O1, O2 any, P Parser[I, []O1], F ConvertFunc[O1, O2]] struct {
	parser P
	fun    F
}

// Map converts the all of outputs of the given `parser` using given ConvertFunc.
func Map[I comparable, O1, O2 any, P Parser[I, []O1], F ConvertFunc[O1, O2]](parser P, fun F) Mapper[I, O1, O2, P, F] {
	return Mapper[I, O1, O2, P, F]{parser, fun}
}

func (m Mapper[I, O1, O2, P, F]) Parse(input []I) (output []O2, remain []I, err error) {
	var o1s []O1
	o1s, remain, err = m.parser.Parse(input)
	if err != nil {
		return
	}

	for _, o1 := range o1s {
		var o2 O2
		o2, err = m.fun(o1)
		if err != nil {
			return
		}
		output = append(output, o2)
	}

	return
}

func (m Mapper[I, O1, O2, P, F]) String() string {
	return fmt.Sprint(m.parser)
}

type MatchOnlyParser[I comparable, O any, P Parser[I, O]] struct {
	parser P
}

// MatchOnly parses the input with the given `parser`, but returns a range of input string that parsed as is.
func MatchOnly[I comparable, O any, P Parser[I, O]](parser P) MatchOnlyParser[I, O, P] {
	return MatchOnlyParser[I, O, P]{parser}
}

func (m MatchOnlyParser[I, O, P]) Parse(input []I) (output []I, remain []I, err error) {
	_, remain, err = m.parser.Parse(input)
	if err != nil {
		return
	}
	l := len(input) - len(remain)
	return input[:l], remain, nil
}

func (m MatchOnlyParser[I, O, P]) String() string {
	return fmt.Sprint(m.parser)
}

type Replacer[I comparable, O1, O2 any, P Parser[I, O1]] struct {
	parser P
	value  O2
}

// Replace replaces result value with a fixed value.
func Replace[I comparable, O1, O2 any, P Parser[I, O1]](parser P, value O2) Replacer[I, O1, O2, P] {
	return Replacer[I, O1, O2, P]{parser, value}
}

func (r Replacer[I, O1, O2, P]) String() string {
	return fmt.Sprint(r.parser)
}

func (r Replacer[I, O1, O2, P]) Parse(input []I) (output O2, remain []I, err error) {
	_, remain, err = r.parser.Parse(input)
	if err != nil {
		return
	}
	return r.value, remain, nil
}
