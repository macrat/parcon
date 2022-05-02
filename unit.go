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

type TagParser[T comparable] struct {
	name string
	tag  []T
}

// Tag parses fixed string or array, like keywords.
//
// The `name` in argument is used as human readable name in error messages.
func Tag[T comparable](name string, tag []T) TagParser[T] {
	return TagParser[T]{name, tag}
}

// TagS is a shorthand for Tag using string as argument.
func TagS(name string, tag string) TagParser[rune] {
	return Tag(name, []rune(tag))
}

func (t TagParser[T]) Parse(input []T) (output []T, remain []T, err error) {
	if len(t.tag) > len(input) {
		return nil, nil, ErrUnexpectedInput[T]{t.name, input}
	}
	for i := range t.tag {
		if t.tag[i] != input[i] {
			return nil, nil, ErrUnexpectedInput[T]{t.name, input}
		}
	}

	return t.tag, input[len(t.tag):], nil
}

func (t TagParser[T]) String() string {
	return t.name
}

type OneOfParser[T comparable] struct {
	name string
	list []T
}

// OneOf parses a single value that listed in the `list`.
//
// If you want to parse two or more values, please use Many parser.
//
// The `name` in argument is used as human readable name in error messages.
func OneOf[T comparable](name string, list []T) OneOfParser[T] {
	return OneOfParser[T]{name, list}
}

// OneOfS is a shorthand for OneOf using a string as `list`.
func OneOfS(name string, list string) OneOfParser[rune] {
	return OneOf(name, []rune(list))
}

func (o OneOfParser[T]) Parse(input []T) (output T, remain []T, err error) {
	if len(input) > 0 && contains(o.list, input[0]) {
		return input[0], input[1:], nil
	} else {
		err = ErrUnexpectedInput[T]{o.name, input}
		return
	}
}

func (o OneOfParser[T]) String() string {
	return o.name
}

type NoneOfParser[T comparable] struct {
	name string
	list []T
}

// NoneOf is almost the same as OneOf, but parses a single value that NOT listed in the `list`.
//
// The `name` in argument is used as human readable name in error messages.
func NoneOf[T comparable](name string, list []T) NoneOfParser[T] {
	return NoneOfParser[T]{name, list}
}

// NoneOfS is a shorthand for NoneOf using a string as `list`.
func NoneOfS(name string, list string) NoneOfParser[rune] {
	return NoneOf(name, []rune(list))
}

func (n NoneOfParser[T]) Parse(input []T) (output T, remain []T, err error) {
	if len(input) > 0 && !contains(n.list, input[0]) {
		return input[0], input[1:], nil
	} else {
		err = ErrUnexpectedInput[T]{n.name, input}
		return
	}
}

func (n NoneOfParser[T]) String() string {
	return n.name
}

type AnythingParser[T comparable] struct{}

// Anything parses any single value.
func Anything[T comparable]() AnythingParser[T] {
	return AnythingParser[T]{}
}

func (a AnythingParser[T]) Parse(input []T) (output T, remain []T, err error) {
	if len(input) == 0 {
		err = ErrUnexpectedInput[T]{"ANYTHING", input}
		return
	}
	return input[0], input[1:], nil
}

func (a AnythingParser[T]) String() string {
	return "ANYTHING"
}

type NothingParser[I comparable] struct{}

// Nothing parses nothing, just leave all of inputs as `remain` and returns `interface{}` as an output.
func Nothing[I comparable]() NothingParser[I] {
	return NothingParser[I]{}
}

func (n NothingParser[I]) Parse(input []I) (output struct{}, remain []I, err error) {
	return struct{}{}, input, nil
}

func (n NothingParser[I]) String() string {
	return "NOTHING"
}

type TakeSingleParser[I comparable] struct {
	name string
	fun  func(I) bool
}

// TakeSingle parses input using the given function.
// The parse will be succeed if the function returns true, otherwise failure with error.
func TakeSingle[I comparable](name string, fun func(I) bool) TakeSingleParser[I] {
	return TakeSingleParser[I]{name, fun}
}

// TakeWhile parses a sequence until the given function returns false.
// This parser expects at least one element that parseable.
func TakeWhile[I comparable](name string, fun func(I) bool) ListParser[I, struct{}, I, NothingParser[I], TakeSingleParser[I]] {
	return Many[I, I](1, TakeSingle(name, fun))
}

func (t TakeSingleParser[I]) Parse(input []I) (output I, remain []I, err error) {
	if len(input) > 0 && t.fun(input[0]) {
		return input[0], input[1:], nil
	} else {
		err = ErrUnexpectedInput[I]{t.name, input}
		return
	}
}

func (t TakeSingleParser[I]) String() string {
	return t.name
}
