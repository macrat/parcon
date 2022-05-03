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

func (o optionalParser[I, O]) Parse(input []I, verbose bool) (output O, remain []I, err error) {
	output, remain, err = o.Parser.Parse(input, false)
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

func (o orParser[I, O]) Parse(input []I, verbose bool) (output O, remain []I, err error) {
	for _, p := range o {
		output, remain, err = p.Parse(input, false)
		if err == nil {
			return
		}
	}
	if verbose {
		err = ErrInvalidInputVerbose[I]{Expected: o, Input: input}
	} else {
		err = ErrInvalidInput
	}
	return
}

func (o orParser[I, O]) String() string {
	var ss []string
	for _, p := range o {
		ss = append(ss, fmt.Sprintf("[%v]", p))
	}
	return fmt.Sprintf("one of %s", strings.Join(ss, " "))
}
