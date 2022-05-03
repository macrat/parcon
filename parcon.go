// package parcon is a yet another parser combinator written in Go.
//
// Parcon uses Generics so you can parse non-string array like []byte.
package parcon

// Parser is the interface of parsers.
type Parser[I comparable, O any] interface {
	// Parse parses `input`, and returns parsed `output`, `remain` slice that not parsed with this parser, and error if it happened.
	//
	// This method tries to record detail of error if `verbose` is true, otherwise, it can only returns ErrInvalidInput.
	// ErrInvalidInput has no any detail, and faster than detailed errors like ErrInvalidInputVerbose.
	Parse(input []I, verbose bool) (output O, remain []I, err error)
}

// ParserFunc is a function to parse input.
// This type implements Parser interface.
//
// See also: Func
type ParserFunc[I comparable, O any] func(input []I, verbose bool) (output O, remain []I, err error)

// Parse parses input.
func (p ParserFunc[I, O]) Parse(input []I, verbose bool) (output O, remain []I, err error) {
	return p(input, verbose)
}

type parserFuncType[I comparable, O any] interface {
	~func(input []I, verbose bool) (output O, remain []I, err error)
}

// Func makes a Parser by a function.
// It is a shorthand for ParserFunc.
func Func[I comparable, O any, F parserFuncType[I, O]](fun F) Parser[I, O] {
	return ParserFunc[I, O](fun)
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

func (n named[I, O]) Parse(input []I, verbose bool) (output O, remain []I, err error) {
	return n.Parser.Parse(input, verbose)
}
