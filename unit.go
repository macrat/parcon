package parcon

// contains checks if `slice` contains `item` or not.
func contains[T comparable](slice []T, item T) bool {
	for _, x := range slice {
		if x == item {
			return true
		}
	}
	return false
}

type tagParser[T comparable] struct {
	Name string
	Tag  []T
}

// Tag parses fixed string or array, like keywords.
//
// The `name` in argument is used as human readable name in error messages.
func Tag[T comparable](name string, tag []T) Parser[T, []T] {
	return tagParser[T]{name, tag}
}

// TagS is a shorthand for Tag using string as argument.
func TagS(name string, tag string) Parser[rune, []rune] {
	return Tag(name, []rune(tag))
}

func (t tagParser[T]) Parse(input []T) (output []T, remain []T, err error) {
	if len(t.Tag) > len(input) {
		return nil, nil, ErrUnexpectedInput[T]{t.Name, input}
	}
	for i := range t.Tag {
		if t.Tag[i] != input[i] {
			return nil, nil, ErrUnexpectedInput[T]{t.Name, input}
		}
	}

	return t.Tag, input[len(t.Tag):], nil
}

func (t tagParser[T]) String() string {
	return t.Name
}

type oneOfParser[T comparable] struct {
	Name string
	List []T
}

// OneOf parses a single value that listed in the `list`.
//
// If you want to parse two or more values, please use Many parser.
//
// The `name` in argument is used as human readable name in error messages.
func OneOf[T comparable](name string, list []T) Parser[T, T] {
	return oneOfParser[T]{name, list}
}

// OneOfS is a shorthand for OneOf using a string as `list`.
func OneOfS(name string, list string) Parser[rune, rune] {
	return OneOf(name, []rune(list))
}

func (o oneOfParser[T]) Parse(input []T) (output T, remain []T, err error) {
	if len(input) > 0 && contains(o.List, input[0]) {
		return input[0], input[1:], nil
	} else {
		err = ErrUnexpectedInput[T]{o.Name, input}
		return
	}
}

func (o oneOfParser[T]) String() string {
	return o.Name
}

type noneOfParser[T comparable] struct {
	Name string
	List []T
}

// NoneOf is almost the same as OneOf, but parses a single value that NOT listed in the `list`.
//
// The `name` in argument is used as human readable name in error messages.
func NoneOf[T comparable](name string, list []T) Parser[T, T] {
	return noneOfParser[T]{name, list}
}

// NoneOfS is a shorthand for NoneOf using a string as `list`.
func NoneOfS(name string, list string) Parser[rune, rune] {
	return NoneOf(name, []rune(list))
}

func (n noneOfParser[T]) Parse(input []T) (output T, remain []T, err error) {
	if len(input) > 0 && !contains(n.List, input[0]) {
		return input[0], input[1:], nil
	} else {
		err = ErrUnexpectedInput[T]{n.Name, input}
		return
	}
}

func (n noneOfParser[T]) String() string {
	return n.Name
}

type anything[T comparable] struct{}

// Anything parses any single value.
func Anything[T comparable]() Parser[T, T] {
	return anything[T]{}
}

func (a anything[T]) Parse(input []T) (output T, remain []T, err error) {
	if len(input) == 0 {
		err = ErrUnexpectedInput[T]{"ANYTHING", input}
		return
	}
	return input[0], input[1:], nil
}

func (a anything[T]) String() string {
	return "ANYTHING"
}

type nothing[I comparable] struct{}

// Nothing parses nothing, just leave all of inputs as `remain` and returns `interface{}` as an output.
func Nothing[I comparable]() Parser[I, struct{}] {
	return nothing[I]{}
}

func (n nothing[I]) Parse(input []I) (output struct{}, remain []I, err error) {
	return struct{}{}, input, nil
}

func (n nothing[I]) String() string {
	return "NOTHING"
}

type takeSingleParser[I comparable] struct {
	Name string
	Func func(I) bool
}

// TakeSingle parses input using the given function.
// The parse will be succeed if the function returns true, otherwise failure with error.
func TakeSingle[I comparable](name string, fn func(I) bool) Parser[I, I] {
	return takeSingleParser[I]{name, fn}
}

// TakeWhile parses a sequence until the given function returns false.
// This parser expects at least one element that parseable.
func TakeWhile[I comparable](name string, fn func(I) bool) Parser[I, []I] {
	return Many(1, TakeSingle(name, fn))
}

func (t takeSingleParser[I]) Parse(input []I) (output I, remain []I, err error) {
	if len(input) > 0 && t.Func(input[0]) {
		return input[0], input[1:], nil
	} else {
		err = ErrUnexpectedInput[I]{t.Name, input}
		return
	}
}

func (t takeSingleParser[I]) String() string {
	return t.Name
}
