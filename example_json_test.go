package parcon_test

import (
	"fmt"
	"strconv"
	"unicode/utf16"

	pc "github.com/macrat/parcon"
)

func ToInterface[I any](x I) (interface{}, error) {
	return x, nil
}

var (
	quote = pc.Tag("DOUBLE_QUOTE", []rune(`"`))

	optionalSpaces = pc.Optional(pc.MultiSpacesOrNewlines)

	str = pc.Named("STRING_LITERAL", pc.Convert(pc.Delimited(
		quote,
		pc.Many(0, pc.Or(
			pc.WithPrefix(
				pc.Tag("ESCAPE", []rune{'\\'}),
				pc.Or(
					pc.TagAs("QUOTATION_MARK", []rune(`"`), '"'),
					pc.TagAs("REVERSE_SOLIDUS", []rune(`\`), '\\'),
					pc.TagAs("SOLIDUS", []rune(`/`), '/'),
					pc.TagAs("BACKSPACE", []rune(`b`), '\b'),
					pc.TagAs("FORM_FEED", []rune(`f`), '\f'),
					pc.TagAs("LINE_FEED", []rune(`n`), '\n'),
					pc.TagAs("CARRIAGE_RETURN", []rune(`r`), '\r'),
					pc.TagAs("TAB", []rune(`t`), '\t'),
				),
			),
			pc.Convert(
				pc.Repeat(2, pc.WithPrefix(
					pc.Tag("UNICODE", []rune(`\u`)),
					pc.Repeat(4, pc.SingleHexDigit),
				)),
				func(xs [][]rune) (rune, error) {
					a, err := strconv.ParseUint(string(xs[0]), 16, 32)
					if err != nil {
						return 0, err
					}
					b, err := strconv.ParseUint(string(xs[1]), 16, 32)
					if err != nil {
						return 0, err
					}
					if (0xD800 <= a && a <= 0xDBFF) && (0xDC00 <= b && b <= 0xDFFF) {
						return utf16.DecodeRune(rune(a), rune(b)), nil
					}
					return 0, fmt.Errorf("they are not a surrogate pair, go next!")
				},
			),
			pc.Convert(
				pc.WithPrefix(
					pc.Tag("UNICODE", []rune(`\u`)),
					pc.Repeat(4, pc.SingleHexDigit),
				),
				func(xs []rune) (rune, error) {
					i, err := strconv.ParseUint(string(xs), 16, 32)
					return rune(i), err
				},
			),
			pc.TakeSingle("UNESCAPED", func(c rune) bool {
				return (0x20 <= c && c <= 0x21) || (0x23 <= c && c <= 0x5B) || (0x5D <= c && c <= 0x10FFFF)
			}),
		)),
		quote,
	), pc.ToString))

	number = pc.Named("NUMBER_LITERAL", pc.Convert(pc.MatchOnly(pc.Sequence(
		pc.Optional(pc.Tag("MINUS", []rune("-"))),
		pc.Or(
			pc.Tag("ZERO", []rune("0")),
			pc.MatchOnly(pc.Pair(
				pc.OneOf("DIGIT_1-9", []rune("123456789")),
				pc.Optional(pc.MultiDigits),
			)),
		),
		pc.Optional(pc.MatchOnly(pc.Sequence(
			pc.Tag("PERIOD", []rune(".")),
			pc.MultiDigits,
		))),
		pc.Optional(pc.MatchOnly(pc.Sequence(
			pc.Sequence(
				pc.OneOf("E", []rune("eE")),
				pc.Optional(pc.OneOf("SIGN", []rune("+-"))),
			),
			pc.MultiDigits,
		))),
	)), pc.ToFloat))

	null = pc.TagAs("NULL", []rune("null"), (interface{})(nil))

	boolean = pc.Or(
		pc.TagAs("TRUE", []rune("true"), true),
		pc.TagAs("FALSE", []rune("false"), false),
	)

	beginArray = pc.WithSuffix(pc.Tag("BEGIN_ARRAY", []rune("[")), optionalSpaces)
	endArray   = pc.WithPrefix(optionalSpaces, pc.Tag("END_ARRAY", []rune("]")))

	beginObject = pc.WithSuffix(pc.Tag("BEGIN_OBJECT", []rune("{")), optionalSpaces)
	endObject   = pc.WithPrefix(optionalSpaces, pc.Tag("END_OBJECT", []rune("}")))

	nameSeparator  = pc.Delimited(optionalSpaces, pc.Tag("NAME_SEPARATOR", []rune(":")), optionalSpaces)
	valueSeparator = pc.Delimited(optionalSpaces, pc.Tag("VALUE_SEPARATOR", []rune(",")), optionalSpaces)

	keyValuePair = pc.Pair(
		pc.WithSuffix(
			str,
			nameSeparator,
		),
		jsonValue,
	)

	jsonValue = pc.Named("JSON_VALUE", pc.Delimited(
		optionalSpaces,
		pc.Or(
			null,
			pc.Convert(str, ToInterface[string]),
			pc.Convert(number, ToInterface[float64]),
			pc.Convert(boolean, ToInterface[bool]),
			pc.Convert[rune, []interface{}](Array{}, ToInterface[[]interface{}]),
			pc.Convert[rune, map[string]interface{}](Object{}, ToInterface[map[string]interface{}]),
		),
		optionalSpaces,
	))
)

type Array struct{}

func (a Array) String() string {
	return "ARRAY"
}

func (a Array) Parse(input []rune, verbose bool) ([]interface{}, []rune, error) {
	return pc.Delimited(
		beginArray,
		pc.SeparatedList(0, valueSeparator, jsonValue),
		endArray,
	).Parse(input, verbose)
}

type Object struct{}

func (o Object) String() string {
	return "OBJECT"
}

func (o Object) Parse(input []rune, verbose bool) (map[string]interface{}, []rune, error) {
	xs, remain, err := pc.Delimited(
		beginObject,
		pc.SeparatedList(0, valueSeparator, keyValuePair),
		endObject,
	).Parse(input, verbose)
	if err != nil {
		return nil, nil, err
	}

	result := make(map[string]interface{})
	for _, x := range xs {
		result[x.First] = x.Second
	}
	return result, remain, nil
}

// Parse JSON that defined in RFC8259
func ParseJson(s string) (interface{}, error) {
	output, remain, err := jsonValue.Parse([]rune(s), true)
	if err != nil {
		return nil, err
	}
	if len(remain) != 0 {
		return nil, fmt.Errorf("found extra string: %s", string(remain))
	}
	return output, nil
}

func Example_json() {
	input := `
		{
			"hello": "world",
			"foo": 123.456,
			"list": [
				1,
				2,
				"3"
			],
			"quoted": "hello\"world"
		}
	`

	output, err := ParseJson(input)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", output)

	// OUTPUT:
	// map[string]interface {}{"foo":123.456, "hello":"world", "list":[]interface {}{1, 2, "3"}, "quoted":"hello\"world"}
}
