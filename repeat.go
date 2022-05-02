package parcon

import (
	"fmt"
)

type listParser[I comparable, O, D any] struct {
	Min       uint
	Max       uint
	Delimiter Parser[I, D]
	Parser    Parser[I, O]
}

// SeparatedList parses an array that separated by `delimiter` using the given `parser`.
// For example, you can use this parser for comma separated list.
//
// The output slice have at least `min` number of elements, or returns error.
// If you want to specify maximum number of elements, please use SeparatedListLimited.
func SeparatedList[I comparable, O, D any](min uint, delimiter Parser[I, D], parser Parser[I, O]) Parser[I, []O] {
	return listParser[I, O, D]{min, 0, delimiter, parser}
}

// SeparatedListLimited parses an array that separated by `delimiter` using the given `parser`.
//
// The output slice have `min` number of elements to `max` number of elements.
// If it did not find enough elements, it returns error. If found more than `max` number of elements, just remains them without error.
func SeparatedListLimited[I comparable, O, D any](min, max uint, delimiter Parser[I, D], parser Parser[I, O]) Parser[I, []O] {
	if max == 1 {
		return Convert(parser, func(o O) ([]O, error) {
			return []O{o}, nil
		})
	} else {
		return listParser[I, O, D]{min, max, delimiter, parser}
	}
}

// Many parses multiple values as a slice using the given `parser`.
//
// The output slice have at least `min` number of elements, or returns error.
// If you want to specify maximum number of elements, please use ManyLimited.
//
// This is a shorthand of SeparatedList that uses Nothing as a delimiter.
func Many[I comparable, O any](min uint, parser Parser[I, O]) Parser[I, []O] {
	return SeparatedList(min, Nothing[I](), parser)
}

// ManyLimited parses multiple values as a slice using the given `parser`.
//
// The output slice have `min` number of elements to `max` number of elements.
// If it did not find enough elements, it returns error. If found more than `max` number of elements, just remains them without error.
//
// This is a shorthand of SeparatedListLimited that uses Nothing as a delimiter.
func ManyLimited[I comparable, O any](min, max uint, parser Parser[I, O]) Parser[I, []O] {
	return SeparatedListLimited(min, max, Nothing[I](), parser)
}

// Repeat parses multiple values that exactly `num` number of elements using the given `parser`.
//
// This is a shorthand of `ManyLimited(num, num, parser)`.
func Repeat[I comparable, O any](num uint, parser Parser[I, O]) Parser[I, []O] {
	return ManyLimited(num, num, parser)
}

func (l listParser[I, O, D]) Parse(input []I) (output []O, remain []I, err error) {
	if l.Max == 0 {
		output = make([]O, 0)
	} else {
		output = make([]O, 0, l.Max)
	}

	var o O

	o, remain, err = l.Parser.Parse(input)
	if err != nil {
		if l.Min == 0 {
			err = nil
			remain = input
		}
		return
	}
	output = append(output, o)

	var count uint = 1

	for l.Max == 0 || count < l.Max {
		var r []I

		_, r, err = l.Delimiter.Parse(remain)
		if err != nil {
			break
		}

		o, r, err = l.Parser.Parse(r)
		if err != nil {
			break
		}

		remain = r
		output = append(output, o)
		count++
	}

	if l.Min <= count {
		err = nil
	}

	return
}

func (l listParser[I, O, D]) String() string {
	switch any(l.Delimiter).(type) {
	case nothing[I]:
		return fmt.Sprintf("multiple [%v]", l.Parser)
	default:
		return fmt.Sprintf("multiple [%v] separated by [%v]", l.Parser, l.Delimiter)
	}
}
