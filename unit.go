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

type tagParser[I comparable, O any] struct {
	Name  string
	Tag   []I
	Value O
}

// TagAs parses fixed string or array, and returns specified value.
// It is similar to Tag with Replace, but faster and simpler than that.
func TagAs[I comparable, O any](name string, tag []I, value O) Parser[I, O] {
	return tagParser[I, O]{name, tag, value}
}

// Tag parses fixed string or array, like keywords.
//
// The `name` in argument is used as human readable name in error messages.
func Tag[I comparable](name string, tag []I) Parser[I, []I] {
	return TagAs(name, tag, tag)
}

// TagStr is similar to Tag parser but it handles string.
func TagStr(name string, tag string) Parser[rune, string] {
	return TagAs(name, []rune(tag), tag)
}

func (t tagParser[I, O]) Parse(input []I, verbose bool) (output O, remain []I, err error) {
	if len(t.Tag) > len(input) {
		if verbose {
			err = ErrInvalidInputVerbose[I]{Expected: t.Name, Input: input}
		} else {
			err = ErrInvalidInput
		}
		return
	}
	for i := range t.Tag {
		if t.Tag[i] != input[i] {
			if verbose {
				err = ErrInvalidInputVerbose[I]{Expected: t.Name, Input: input}
			} else {
				err = ErrInvalidInput
			}
			return
		}
	}

	return t.Value, input[len(t.Tag):], nil
}

func (t tagParser[I, O]) String() string {
	return t.Name
}

type oneOfParser[T comparable] struct {
	Name string
	List []T
}

// OneOf parses a single value that listed in the `list`.
//
// If you want to parse two or more values, please use OneOfList.
//
// The `name` in argument is used as human readable name in error messages.
func OneOf[T comparable](name string, list []T) Parser[T, T] {
	return oneOfParser[T]{name, list}
}

func (o oneOfParser[T]) Parse(input []T, verbose bool) (output T, remain []T, err error) {
	if len(input) > 0 && contains(o.List, input[0]) {
		return input[0], input[1:], nil
	} else {
		if verbose {
			err = ErrInvalidInputVerbose[T]{Expected: o, Input: input}
		} else {
			err = ErrInvalidInput
		}
		return
	}
}

func (o oneOfParser[T]) String() string {
	return o.Name
}

type oneOfListParser[T comparable] struct {
	Name string
	List []T
}

// OneOfList parses one of more values that listed in the `list`.
//
// If you want to parse exact one value, please use OneOf.
//
// The `name` in argument is used as human readable name in error messages.
func OneOfList[T comparable](name string, list []T) Parser[T, []T] {
	return oneOfListParser[T]{name, list}
}

// OneOfStr is a similar parser to the OneOfList, but it parses string instead of generics type.
func OneOfStr(name string, list string) Parser[rune, string] {
	return Convert(OneOfList(name, []rune(list)), ToString)
}

func (o oneOfListParser[T]) Parse(input []T, verbose bool) (output []T, remain []T, err error) {
	if len(input) == 0 || !contains(o.List, input[0]) {
		if verbose {
			err = ErrInvalidInputVerbose[T]{Expected: o, Input: input}
		} else {
			err = ErrInvalidInput
		}
		return
	}

	i := 1
	for ; i < len(input); i++ {
		if !contains(o.List, input[i]) {
			break
		}
	}
	return input[:i], input[i:], nil
}

func (o oneOfListParser[T]) String() string {
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

func (n noneOfParser[T]) Parse(input []T, verbose bool) (output T, remain []T, err error) {
	if len(input) > 0 && !contains(n.List, input[0]) {
		return input[0], input[1:], nil
	} else {
		if verbose {
			err = ErrInvalidInputVerbose[T]{n.Name, input}
		} else {
			err = ErrInvalidInput
		}
		return
	}
}

func (n noneOfParser[T]) String() string {
	return n.Name
}

type noneOfListParser[T comparable] struct {
	Name string
	List []T
}

// NoneOfList parses one of more values that NOT listed in the `list`.
//
// If you want to parse exact one value, please use NoneOf.
//
// The `name` in argument is used as human readable name in error messages.
func NoneOfList[T comparable](name string, list []T) Parser[T, []T] {
	return noneOfListParser[T]{name, list}
}

// NoneOfStr is a similar parser to the NoneOfList, but it parses string instead of generics type.
func NoneOfStr(name string, list string) Parser[rune, string] {
	return Convert(NoneOfList(name, []rune(list)), ToString)
}

func (n noneOfListParser[T]) Parse(input []T, verbose bool) (output []T, remain []T, err error) {
	if len(input) == 0 || contains(n.List, input[0]) {
		if verbose {
			err = ErrInvalidInputVerbose[T]{n.Name, input}
		} else {
			err = ErrInvalidInput
		}
		return
	}

	i := 1
	for ; i < len(input); i++ {
		if contains(n.List, input[i]) {
			break
		}
	}
	return input[:i], input[i:], nil
}

func (n noneOfListParser[T]) String() string {
	return n.Name
}

type anything[T comparable] struct{}

// Anything parses any single value.
func Anything[T comparable]() Parser[T, T] {
	return anything[T]{}
}

func (a anything[T]) Parse(input []T, verbose bool) (output T, remain []T, err error) {
	if len(input) == 0 {
		if verbose {
			err = ErrInvalidInputVerbose[T]{"ANYTHING", input}
		} else {
			err = ErrInvalidInput
		}
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

func (n nothing[I]) Parse(input []I, verbose bool) (output struct{}, remain []I, err error) {
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

func (t takeSingleParser[I]) Parse(input []I, verbose bool) (output I, remain []I, err error) {
	if len(input) > 0 && t.Func(input[0]) {
		return input[0], input[1:], nil
	} else {
		if verbose {
			err = ErrInvalidInputVerbose[I]{t.Name, input}
		} else {
			err = ErrInvalidInput
		}
		return
	}
}

func (t takeSingleParser[I]) String() string {
	return t.Name
}

type takeWhileParser[I comparable] struct {
	Name string
	Func func(I) bool
}

// TakeWhile parses a sequence until the given function returns false.
// This parser expects at least one element that parseable.
func TakeWhile[I comparable](name string, fn func(I) bool) Parser[I, []I] {
	return takeWhileParser[I]{name, fn}
}

func (t takeWhileParser[I]) Parse(input []I, verbose bool) (output []I, remain []I, err error) {
	if len(input) == 0 || !t.Func(input[0]) {
		if verbose {
			err = ErrInvalidInputVerbose[I]{t.Name, input}
		} else {
			err = ErrInvalidInput
		}
		return
	}

	i := 1
	for ; i < len(input); i++ {
		if !t.Func(input[i]) {
			break
		}
	}
	return input[:i], input[i:], nil
}

func (t takeWhileParser[I]) String() string {
	return t.Name
}
