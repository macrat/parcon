package parcon

import (
	"errors"
	"fmt"
)

// ErrInvalidInput is a error when the parser found unexpected input.
var ErrInvalidInput = errors.New("invalid input")

// ErrInvalidInputVerbose is a error when the parser found unexpected input, with verbose information.
type ErrInvalidInputVerbose[I comparable] struct {
	Expected any
	Input    []I
}

// Unwrap always returns ErrInvalidInput.
func (e ErrInvalidInputVerbose[I]) Unwrap() error {
	return ErrInvalidInput
}

// Error returns human readable string.
func (e ErrInvalidInputVerbose[I]) Error() string {
	switch i := any(e.Input).(type) {
	case []rune:
		return fmt.Sprintf("invalid input: expected %v but got %#v", e.Expected, string(i))
	default:
		return fmt.Sprintf("invalid input: expected %v but got %v", e.Expected, e.Input)
	}
}
