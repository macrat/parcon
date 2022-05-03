package parcon_test

import (
	"fmt"

	pc "github.com/macrat/parcon"
)

var QuotedString = pc.WithEnclosure(
	pc.Tag("QUOTE", []rune(`"`)),
	pc.Convert(
		pc.Many(0, pc.Or(
			pc.TagAs("ESCAPED_QUOTE", []rune(`\"`), '"'),
			pc.TagAs("ESCAPED_SLASH", []rune(`\\`), '\\'),
			pc.NoneOf("CHARACTER", []rune(`\"`)),
		)),
		pc.ToString,
	),
	pc.Tag("QUOTE", []rune(`"`)),
)

func ParseQuotedString(input string) (string, error) {
	output, remain, err := QuotedString.Parse([]rune(input), true)
	if err != nil {
		return "", err
	}
	if len(remain) != 0 {
		return "", fmt.Errorf("found extra string: %#v", string(remain))
	}
	return output, nil
}

func Example_quotedStringWithEscape() {
	output, err := ParseQuotedString(`"hello \" world!"`)
	if err != nil {
		panic(err)
	}

	fmt.Println(output)

	// OUTPUT:
	// hello " world!
}
