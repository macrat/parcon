package parcon

import (
	"fmt"
	"strings"
)

type optionalParser[I comparable, O any] struct {
	Parser  Parser[I, O]
	Default O
}

// Optional tries to parse with the given `parser`, and returns the `default_` value if failed to parse.
func OptionalWithDefault[I comparable, O any](parser Parser[I, O], default_ O) Parser[I, O] {
	return optionalParser[I, O]{
		Parser:  parser,
		Default: default_,
	}
}

// Optional tries to parse with the given `parser`.
// It returns zero value if failure to parse the input.
func Optional[I comparable, O any](parser Parser[I, O]) Parser[I, O] {
	return optionalParser[I, O]{Parser: parser}
}

func (o optionalParser[I, O]) Parse(input []I) (output O, remain []I, err error) {
	output, remain, err = o.Parser.Parse(input)
	if err != nil {
		return o.Default, input, nil
	}
	return
}

func (o optionalParser[I, O]) String() string {
	return fmt.Sprintf("%v", o.Parser)
}

type orParser[I comparable, O any] []Parser[I, O]

// Or parses using one of `parsers`, and returns the parsed value that first succeed.
func Or[I comparable, O any](parsers ...Parser[I, O]) Parser[I, O] {
	return orParser[I, O](parsers)
}

func (o orParser[I, O]) Parse(input []I) (output O, remain []I, err error) {
	for _, p := range o {
		output, remain, err = p.Parse(input)
		if err == nil {
			return
		}
	}
	err = ErrUnexpectedInput[I]{Name: o.String(), Input: input}
	return
}

func (o orParser[I, O]) String() string {
	var ss []string
	for _, p := range o {
		ss = append(ss, fmt.Sprintf("[%v]", p))
	}
	return fmt.Sprintf("one of %s", strings.Join(ss, " "))
}

type sequenceParser[I comparable, O any] []Parser[I, O]

// Sequence parses using all of `parsers` sequentially.
func Sequence[I comparable, O any](parsers ...Parser[I, O]) Parser[I, []O] {
	return sequenceParser[I, O](parsers)
}

func (s sequenceParser[I, O]) Parse(input []I) (output []O, remain []I, err error) {
	remain = input
	output = make([]O, len(s))
	for i, p := range s {
		output[i], remain, err = p.Parse(remain)
		if err != nil {
			return
		}
	}
	return
}

func (s sequenceParser[I, O]) String() string {
	var ss []string
	for _, p := range s {
		ss = append(ss, fmt.Sprint(p))
	}
	return fmt.Sprintf("[%s]", strings.Join(ss, ", "))
}

// PairValue is a pair of values.
type PairValue[F, S any] struct {
	First  F
	Second S
}

type pairParser[I comparable, O1, O2 any] struct {
	First  Parser[I, O1]
	Second Parser[I, O2]
}

// Pair parses a pair of elements that have different types.
func Pair[I comparable, O1, O2 any](first Parser[I, O1], second Parser[I, O2]) Parser[I, PairValue[O1, O2]] {
	return pairParser[I, O1, O2]{first, second}
}

func (p pairParser[I, O1, O2]) Parse(input []I) (output PairValue[O1, O2], remain []I, err error) {
	output.First, remain, err = p.First.Parse(input)
	if err != nil {
		return
	}

	output.Second, remain, err = p.Second.Parse(remain)
	return
}

type separatedListParser[I comparable, O, D any] struct {
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
	return separatedListParser[I, O, D]{min, 0, delimiter, parser}
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
		return separatedListParser[I, O, D]{min, max, delimiter, parser}
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

func (s separatedListParser[I, O, D]) Parse(input []I) (output []O, remain []I, err error) {
	remain = input

	if s.Max != 0 {
		output = make([]O, 0, s.Max)
	}

	var count uint
	for {
		var o O
		var r []I
		o, r, err = s.Parser.Parse(remain)
		if err != nil {
			break
		}

		remain = r
		output = append(output, o)
		count++

		if s.Max != 0 && count >= s.Max {
			break
		}

		_, r, err = s.Delimiter.Parse(remain)
		if err != nil {
			break
		}
		remain = r
	}

	if s.Min <= count {
		err = nil
	}

	return
}

func (s separatedListParser[I, O, D]) String() string {
	switch any(s.Delimiter).(type) {
	case nothing[I]:
		return fmt.Sprintf("multiple [%v]", s.Parser)
	default:
		return fmt.Sprintf("multiple [%v] separated by [%v]", s.Parser, s.Delimiter)
	}
}

type delimitedParser[I comparable, P, O, S any] struct {
	Prefix Parser[I, P]
	Body   Parser[I, O]
	Suffix Parser[I, S]
}

// Delimited parses a value that have a prefix and a suffix.
// For example, you can use this parser for quoted string.
func Delimited[I comparable, P, O, S any](prefix Parser[I, P], body Parser[I, O], suffix Parser[I, S]) Parser[I, O] {
	return delimitedParser[I, P, O, S]{prefix, body, suffix}
}

// WithPrefix parses a value that have a prefix.
// For example, you can use this parser to parse a GitHub or Twitter style mention that have '@' prefix.
//
// This is a shorthand for Delimited that uses Nothing as the suffix.
func WithPrefix[I comparable, P, O any](prefix Parser[I, P], body Parser[I, O]) Parser[I, O] {
	return delimitedParser[I, P, O, struct{}]{prefix, body, Nothing[I]()}
}

// WithSuffix parses a value that have a suffix.
// For example, you can use this parser to parse a single line of C language that have semi-colon as a prefix at the end of lines.
func WithSuffix[I comparable, O, S any](body Parser[I, O], suffix Parser[I, S]) Parser[I, O] {
	return delimitedParser[I, struct{}, O, S]{Nothing[I](), body, suffix}
}

func (d delimitedParser[I, P, O, S]) Parse(input []I) (output O, remain []I, err error) {
	_, remain, err = d.Prefix.Parse(input)
	if err != nil {
		return
	}

	output, remain, err = d.Body.Parse(remain)
	if err != nil {
		return
	}

	_, remain, err = d.Suffix.Parse(remain)
	return
}

func (d delimitedParser[I, P, O, S]) String() string {
	return fmt.Sprintf("%v, %v, %v", d.Prefix, d.Body, d.Suffix)
}

type named[I comparable, O any] struct {
	Name   string
	Parser Parser[I, O]
}

// Named sets parser's name that shown in error message.
func Named[I comparable, O any](name string, parser Parser[I, O]) Parser[I, O] {
	return named[I, O]{name, parser}
}

func (n named[I, O]) String() string {
	return n.Name
}

func (n named[I, O]) Parse(input []I) (output O, remain []I, err error) {
	return n.Parser.Parse(input)
}
