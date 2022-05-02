// package parcon is a yet another parser combinator written in Go.
//
// Parcon uses Generics so you can parse non-string array like []byte.
package parcon

// Parser is the interface of parsers.
type Parser[I comparable, O any] interface {
	// Parse parses `input`, and returns parsed `output`, `remain` slice that not parsed with this parser, and error if it happened.
	Parse(input []I) (output O, remain []I, err error)
}

// ParseFunc is a function to parse input.
// This type implements Parser interface.
type ParseFunc[I comparable, O any] func(input []I) (output O, remain []I, err error)

// Parse parses input.
func (p ParseFunc[I, O]) Parse(input []I) (output O, remain []I, err error) {
	return p(input)
}
