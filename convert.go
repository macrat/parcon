package parcon

import (
	"fmt"
)

// ConvertFunc is a function type to convert parsed value.
//
// Normally, it is used for Convert or Map parser.
type ConvertFunc[I any, O any] func(input I) (output O, err error)

type converter[I comparable, O1, O2 any] struct {
	Parser Parser[I, O1]
	Func   ConvertFunc[O1, O2]
}

// Convert converts the output of the given `parser` using given ConvertFunc.
func Convert[I comparable, O1, O2 any](parser Parser[I, O1], fn ConvertFunc[O1, O2]) Parser[I, O2] {
	return converter[I, O1, O2]{parser, fn}
}

func (c converter[I, O1, O2]) Parse(input []I, verbose bool) (output O2, remain []I, err error) {
	var o O1
	o, remain, err = c.Parser.Parse(input, verbose)
	if err != nil {
		return
	}

	output, err = c.Func(o)
	return
}

func (c converter[I, O1, O2]) String() string {
	return fmt.Sprint(c.Parser)
}

type mapper[I comparable, O1, O2 any] struct {
	Parser Parser[I, []O1]
	Func   ConvertFunc[O1, O2]
}

// Map converts the all of outputs of the given `parser` using given ConvertFunc.
func Map[I comparable, O1, O2 any](parser Parser[I, []O1], fn ConvertFunc[O1, O2]) Parser[I, []O2] {
	return mapper[I, O1, O2]{parser, fn}
}

func (m mapper[I, O1, O2]) Parse(input []I, verbose bool) (output []O2, remain []I, err error) {
	var o1s []O1
	o1s, remain, err = m.Parser.Parse(input, verbose)
	if err != nil {
		return
	}

	for _, o1 := range o1s {
		var o2 O2
		o2, err = m.Func(o1)
		if err != nil {
			return
		}
		output = append(output, o2)
	}

	return
}

func (m mapper[I, O1, O2]) String() string {
	return fmt.Sprint(m.Parser)
}

type matchOnly[I comparable, O any] struct {
	Parser Parser[I, O]
}

// MatchOnly parses the input with the given `parser`, but returns a range of input string that parsed as is.
func MatchOnly[I comparable, O any](parser Parser[I, O]) Parser[I, []I] {
	return matchOnly[I, O]{parser}
}

func (m matchOnly[I, O]) Parse(input []I, verbose bool) (output []I, remain []I, err error) {
	_, remain, err = m.Parser.Parse(input, verbose)
	if err != nil {
		return
	}
	l := len(input) - len(remain)
	return input[:l], remain, nil
}

func (m matchOnly[I, O]) String() string {
	return fmt.Sprint(m.Parser)
}

type replace[I comparable, O1, O2 any] struct {
	Parser Parser[I, O1]
	Value  O2
}

// Replace replaces result value with a fixed value.
func Replace[I comparable, O1, O2 any](parser Parser[I, O1], value O2) Parser[I, O2] {
	return replace[I, O1, O2]{parser, value}
}

func (r replace[I, O1, O2]) String() string {
	return fmt.Sprint(r.Parser)
}

func (r replace[I, O1, O2]) Parse(input []I, verbose bool) (output O2, remain []I, err error) {
	_, remain, err = r.Parser.Parse(input, verbose)
	if err != nil {
		return
	}
	return r.Value, remain, nil
}
