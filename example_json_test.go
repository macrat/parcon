package parcon_test

import (
	"fmt"

	pc "github.com/macrat/parcon"
)

func ToInterface[I any](x I) (interface{}, error) {
	return x, nil
}

var (
	quote = pc.TagS("DOUBLE_QUOTE", `"`)

	optionalSpaces = pc.Optional[rune, []rune](pc.MultiSpacesOrNewlines)

	str = pc.Named[rune, string]("STRING_LITERAL", pc.Convert[rune](pc.Delimited[rune, []rune, []rune, []rune](
		quote,
		pc.Many[rune, rune](0, pc.Or[rune, rune](
			pc.Replace[rune, []rune, rune](pc.TagS("ESCAPE", `\"`), '"'),
			pc.NoneOfS("CHARACTER", `"`),
		)),
		quote,
	), pc.ToString))

	number = pc.Named[rune, float64]("NUMBER_LITERAL", pc.Convert[rune](pc.MatchOnly[rune, [][]rune](pc.Sequence[rune, []rune](
		pc.MultiDigits,
		pc.Optional[rune, []rune](pc.MatchOnly[rune, [][]rune](pc.Sequence[rune, []rune](
			pc.TagS("PERIOD", "."),
			pc.MultiDigits,
		))),
	)), pc.ToFloat))

	listSeparator = pc.TagS("LIST_SEPARATOR", ",")
	listStart     = pc.TagS("LIST_START", "[")
	listEnd       = pc.TagS("LIST_END", "]")

	objectSeparator = pc.TagS("OBJECT_SEPARATOR", ":")
	objectStart     = pc.TagS("OBJECT_START", "{")
	objectEnd       = pc.TagS("OBJECT_END", "}")
)

type List struct{}

func (l List) String() string {
	return "LIST"
}

func (l List) Parse(input []rune) ([]interface{}, []rune, error) {
	return pc.Delimited[rune, [][]rune, []interface{}, [][]rune](
		pc.Sequence[rune, []rune](listStart, optionalSpaces),
		pc.SeparatedList[rune, [][]rune, interface{}](
			0,
			pc.Sequence[rune, []rune](optionalSpaces, listSeparator, optionalSpaces),
			JsonValue{},
		),
		pc.Sequence[rune, []rune](optionalSpaces, listEnd),
	).Parse(input)
}

type Object struct{}

func (o Object) String() string {
	return "OBJECT"
}

func (o Object) Parse(input []rune) (map[string]interface{}, []rune, error) {
	keyValue := pc.Pair[rune, string, interface{}](
		pc.WithSuffix[rune, string, [][]rune](
			str,
			pc.Sequence[rune, []rune](optionalSpaces, objectSeparator, optionalSpaces),
		),
		JsonValue{},
	)

	xs, remain, err := pc.Delimited[rune, [][]rune, []pc.PairValue[string, interface{}], [][]rune](
		pc.Sequence[rune, []rune](objectStart, optionalSpaces),
		pc.SeparatedList[rune, [][]rune, pc.PairValue[string, interface{}]](
			0,
			pc.Sequence[rune, []rune](optionalSpaces, listSeparator, optionalSpaces),
			keyValue,
		),
		pc.Sequence[rune, []rune](optionalSpaces, objectEnd),
	).Parse(input)
	if err != nil {
		return nil, nil, err
	}

	result := make(map[string]interface{})
	for _, x := range xs {
		result[x.First] = x.Second
	}
	return result, remain, nil
}

type JsonValue struct{}

func (j JsonValue) String() string {
	return "JSON_VALUE"
}

func (j JsonValue) Parse(input []rune) (interface{}, []rune, error) {
	return pc.Delimited[rune, []rune, interface{}, []rune](
		optionalSpaces,
		pc.Or[rune, interface{}](
			pc.Convert[rune](str, ToInterface[string]),
			pc.Convert[rune](number, ToInterface[float64]),
			pc.Convert[rune](List{}, ToInterface[[]interface{}]),
			pc.Convert[rune](Object{}, ToInterface[map[string]interface{}]),
		),
		optionalSpaces,
	).Parse(input)
}

func Example_json() {
	output, remain, err := JsonValue{}.Parse([]rune(`
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
	`))
	if err != nil {
		panic(err)
	}
	if len(remain) != 0 {
		panic(string(remain))
	}

	fmt.Printf("%#v\n", output)

	// OUTPUT:
	// map[string]interface {}{"foo":123.456, "hello":"world", "list":[]interface {}{1, 2, "3"}, "quoted":"hello\"world"}
}
