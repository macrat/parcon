package parcon

import (
	"fmt"
)

// ErrUnexpectedInput is a error when the parser found unexpected input.
type ErrUnexpectedInput[I comparable] struct {
	Name  string
	Input []I
}

// Error returns human readable string.
func (e ErrUnexpectedInput[I]) Error() string {
	switch i := any(e.Input).(type) {
	case []rune:
		return fmt.Sprintf("expected %s but got %#v", e.Name, string(i))
	default:
		return fmt.Sprintf("expected %s but got %v", e.Name, e.Input)
	}
}
