package parcon

import (
	"fmt"
	"strings"
)

type OptionalParser[I comparable, O any, P Parser[I, O]] struct {
	parser   P
	default_ O
}

// Optional tries to parse with the given `parser`, and returns the `default_` value if failed to parse.
func OptionalWithDefault[I comparable, O any, P Parser[I, O]](parser P, default_ O) OptionalParser[I, O, P] {
	return OptionalParser[I, O, P]{parser, default_}
}

// Optional tries to parse with the given `parser`.
// It returns zero value if failure to parse the input.
func Optional[I comparable, O any, P Parser[I, O]](parser P) OptionalParser[I, O, P] {
	return OptionalParser[I, O, P]{parser: parser}
}

func (o OptionalParser[I, O, P]) Parse(input []I) (output O, remain []I, err error) {
	output, remain, err = o.parser.Parse(input)
	if err != nil {
		return o.default_, input, nil
	}
	return
}

func (o OptionalParser[I, O, P]) String() string {
	return fmt.Sprint(o.parser)
}

type OrParser[I comparable, O any] struct {
	parsers []Parser[I, O]
}

// Or parses using one of `parsers`, and returns the parsed value that first succeed.
func Or[I comparable, O any](parsers ...Parser[I, O]) OrParser[I, O] {
	return OrParser[I, O]{parsers}
}

func (o OrParser[I, O]) Parse(input []I) (output O, remain []I, err error) {
	for _, p := range o.parsers {
		output, remain, err = p.Parse(input)
		if err == nil {
			return
		}
	}
	err = ErrUnexpectedInput[I]{Name: o.String(), Input: input}
	return
}

func (o OrParser[I, O]) String() string {
	ss := make([]string, len(o.parsers))
	for i, p := range o.parsers {
		ss[i] = fmt.Sprintf("[%v]", p)
	}
	return fmt.Sprintf("one of %s", strings.Join(ss, " "))
}

type SequenceParser[I comparable, O any] struct {
	parsers []Parser[I, O]
}

// Sequence parses using all of `parsers` sequentially.
func Sequence[I comparable, O any](parsers ...Parser[I, O]) SequenceParser[I, O] {
	return SequenceParser[I, O]{parsers}
}

func (s SequenceParser[I, O]) Parse(input []I) (output []O, remain []I, err error) {
	remain = input
	output = make([]O, len(s.parsers))
	for i, p := range s.parsers {
		output[i], remain, err = p.Parse(remain)
		if err != nil {
			return
		}
	}
	return
}

func (s SequenceParser[I, O]) String() string {
	ss := make([]string, len(s.parsers))
	for i, p := range s.parsers {
		ss[i] = fmt.Sprint(p)
	}
	return fmt.Sprintf("[%s]", strings.Join(ss, ", "))
}

// PairValue is a pair of values.
type PairValue[F, S any] struct {
	First  F
	Second S
}

type PairParser[I comparable, O1, O2 any, P1 Parser[I, O1], P2 Parser[I, O2]] struct {
	first  P1
	second P2
}

// Pair parses a pair of elements that have different types.
func Pair[I comparable, O1, O2 any, P1 Parser[I, O1], P2 Parser[I, O2]](first P1, second P2) PairParser[I, O1, O2, P1, P2] {
	return PairParser[I, O1, O2, P1, P2]{first, second}
}

func (p PairParser[I, O1, O2, P1, P2]) Parse(input []I) (output PairValue[O1, O2], remain []I, err error) {
	output.First, remain, err = p.first.Parse(input)
	if err != nil {
		return
	}

	output.Second, remain, err = p.second.Parse(remain)
	return
}

type ListParser[I comparable, D, O any, DP Parser[I, D], OP Parser[I, O]] struct {
	min       uint
	max       uint
	delimiter DP
	parser    OP
}

// SeparatedList parses an array that separated by `delimiter` using the given `parser`.
// For example, you can use this parser for comma separated list.
//
// The output slice have at least `min` number of elements, or returns error.
// If you want to specify maximum number of elements, please use SeparatedListLimited.
func SeparatedList[I comparable, D, O any, DP Parser[I, D], OP Parser[I, O]](min uint, delimiter DP, parser OP) ListParser[I, D, O, DP, OP] {
	return ListParser[I, D, O, DP, OP]{min, 0, delimiter, parser}
}

// SeparatedListLimited parses an array that separated by `delimiter` using the given `parser`.
//
// The output slice have `min` number of elements to `max` number of elements.
// If it did not find enough elements, it returns error. If found more than `max` number of elements, just remains them without error.
func SeparatedListLimited[I comparable, D, O any, DP Parser[I, D], OP Parser[I, O]](min, max uint, delimiter DP, parser OP) ListParser[I, D, O, DP, OP] {
	return ListParser[I, D, O, DP, OP]{min, max, delimiter, parser}
}

// Many parses multiple values as a slice using the given `parser`.
//
// The output slice have at least `min` number of elements, or returns error.
// If you want to specify maximum number of elements, please use ManyLimited.
//
// This is a shorthand of SeparatedList that uses Nothing as a delimiter.
func Many[I comparable, O any, P Parser[I, O]](min uint, parser P) ListParser[I, struct{}, O, NothingParser[I], P] {
	return SeparatedList[I, struct{}, O](min, Nothing[I](), parser)
}

// ManyLimited parses multiple values as a slice using the given `parser`.
//
// The output slice have `min` number of elements to `max` number of elements.
// If it did not find enough elements, it returns error. If found more than `max` number of elements, just remains them without error.
//
// This is a shorthand of SeparatedListLimited that uses Nothing as a delimiter.
func ManyLimited[I comparable, O any, P Parser[I, O]](min, max uint, parser P) ListParser[I, struct{}, O, NothingParser[I], P] {
	return SeparatedListLimited[I, struct{}, O](min, max, Nothing[I](), parser)
}

// Repeat parses multiple values that exactly `num` number of elements using the given `parser`.
//
// This is a shorthand of `ManyLimited(num, num, parser)`.
func Repeat[I comparable, O any, P Parser[I, O]](num uint, parser P) ListParser[I, struct{}, O, NothingParser[I], P] {
	return ManyLimited[I, O](num, num, parser)
}

func (p ListParser[I, D, O, DP, OP]) Parse(input []I) (output []O, remain []I, err error) {
	remain = input

	if p.max != 0 {
		output = make([]O, 0, p.max)
	}

	var count uint
	for {
		var o O
		var r []I
		o, r, err = p.parser.Parse(remain)
		if err != nil {
			break
		}

		remain = r
		output = append(output, o)
		count++

		if p.max != 0 && count >= p.max {
			break
		}

		_, r, err = p.delimiter.Parse(remain)
		if err != nil {
			break
		}
		remain = r
	}

	if p.min <= count {
		err = nil
	}

	return
}

func (p ListParser[I, D, O, DP, OP]) String() string {
	switch any(p.delimiter).(type) {
	case NothingParser[I]:
		return fmt.Sprintf("multiple [%v]", p.parser)
	default:
		return fmt.Sprintf("multiple [%v] separated by [%v]", p.parser, p.delimiter)
	}
}

type DelimitedParser[I comparable, P, O, S any, PP Parser[I, P], OP Parser[I, O], SP Parser[I, S]] struct {
	prefix PP
	body   OP
	suffix SP
}

// Delimited parses a value that have a prefix and a suffix.
// For example, you can use this parser for quoted string.
func Delimited[I comparable, P, O, S any, PP Parser[I, P], OP Parser[I, O], SP Parser[I, S]](prefix PP, body OP, suffix SP) DelimitedParser[I, P, O, S, PP, OP, SP] {
	return DelimitedParser[I, P, O, S, PP, OP, SP]{prefix, body, suffix}
}

// WithPrefix parses a value that have a prefix.
// For example, you can use this parser to parse a GitHub or Twitter style mention that have '@' prefix.
//
// This is a shorthand for Delimited that uses Nothing as the suffix.
func WithPrefix[I comparable, P, O any, PP Parser[I, P], OP Parser[I, O]](prefix PP, body OP) DelimitedParser[I, P, O, struct{}, PP, OP, NothingParser[I]] {
	return DelimitedParser[I, P, O, struct{}, PP, OP, NothingParser[I]]{prefix, body, Nothing[I]()}
}

// WithSuffix parses a value that have a suffix.
// For example, you can use this parser to parse a single line of C language that have semi-colon as a prefix at the end of lines.
func WithSuffix[I comparable, O, S any, OP Parser[I, O], SP Parser[I, S]](body OP, suffix SP) DelimitedParser[I, struct{}, O, S, NothingParser[I], OP, SP] {
	return DelimitedParser[I, struct{}, O, S, NothingParser[I], OP, SP]{Nothing[I](), body, suffix}
}

func (d DelimitedParser[I, P, O, S, PP, OP, SP]) Parse(input []I) (output O, remain []I, err error) {
	_, remain, err = d.prefix.Parse(input)
	if err != nil {
		return
	}

	output, remain, err = d.body.Parse(remain)
	if err != nil {
		return
	}

	_, remain, err = d.suffix.Parse(remain)
	return
}

func (d DelimitedParser[I, P, O, S, PP, OP, SP]) String() string {
	return fmt.Sprintf("%v, %v, %v", d.prefix, d.body, d.suffix)
}

type NamedParser[I comparable, O any, P Parser[I, O]] struct {
	name   string
	parser P
}

// Named sets parser's name that shown in error message.
func Named[I comparable, O any, P Parser[I, O]](name string, parser P) NamedParser[I, O, P] {
	return NamedParser[I, O, P]{name, parser}
}

func (n NamedParser[I, O, P]) String() string {
	return n.name
}

func (n NamedParser[I, O, P]) Parse(input []I) (output O, remain []I, err error) {
	return n.parser.Parse(input)
}
