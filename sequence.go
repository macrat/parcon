package parcon

import (
	"fmt"
	"strings"
)

type sequenceParser[I comparable, O any] []Parser[I, O]

// Sequence parses using all of `parsers` sequentially.
func Sequence[I comparable, O any](parsers ...Parser[I, O]) Parser[I, []O] {
	return sequenceParser[I, O](parsers)
}

func (s sequenceParser[I, O]) Parse(input []I, verbose bool) (output []O, remain []I, err error) {
	remain = input
	output = make([]O, len(s))
	for i, p := range s {
		output[i], remain, err = p.Parse(remain, verbose)
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

func (p pairParser[I, O1, O2]) Parse(input []I, verbose bool) (output PairValue[O1, O2], remain []I, err error) {
	output.First, remain, err = p.First.Parse(input, verbose)
	if err != nil {
		return
	}

	output.Second, remain, err = p.Second.Parse(remain, verbose)
	return
}

func (p pairParser[I, O1, O2]) String() string {
	return fmt.Sprintf("[%v, %v]", p.First, p.Second)
}

type enclosuredParser[I comparable, P, O, S any] struct {
	Prefix Parser[I, P]
	Body   Parser[I, O]
	Suffix Parser[I, S]
}

// WithEnclosure parses a value that have a prefix and a suffix.
// For example, you can use this parser for quoted string.
func WithEnclosure[I comparable, P, O, S any](prefix Parser[I, P], body Parser[I, O], suffix Parser[I, S]) Parser[I, O] {
	return enclosuredParser[I, P, O, S]{prefix, body, suffix}
}

// WithPrefix parses a value that have a prefix.
// For example, you can use this parser to parse a GitHub or Twitter style mention that have '@' prefix.
func WithPrefix[I comparable, P, O any](prefix Parser[I, P], body Parser[I, O]) Parser[I, O] {
	return enclosuredParser[I, P, O, struct{}]{prefix, body, Nothing[I]()}
}

// WithSuffix parses a value that have a suffix.
// For example, you can use this parser to parse a single line of C language that have semi-colon as a prefix at the end of lines.
func WithSuffix[I comparable, O, S any](body Parser[I, O], suffix Parser[I, S]) Parser[I, O] {
	return enclosuredParser[I, struct{}, O, S]{Nothing[I](), body, suffix}
}

func (d enclosuredParser[I, P, O, S]) Parse(input []I, verbose bool) (output O, remain []I, err error) {
	_, remain, err = d.Prefix.Parse(input, verbose)
	if err != nil {
		return
	}

	output, remain, err = d.Body.Parse(remain, verbose)
	if err != nil {
		return
	}

	_, remain, err = d.Suffix.Parse(remain, verbose)
	return
}

func (d enclosuredParser[I, P, O, S]) String() string {
	return fmt.Sprintf("%v, %v, %v", d.Prefix, d.Body, d.Suffix)
}
