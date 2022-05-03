package parcon_test

import (
	"fmt"
	"strconv"

	"github.com/macrat/parcon"
)

func ExampleConvert() {
	parser := parcon.Convert(
		parcon.MultiDigits,
		func(input []rune) (int, error) {
			return strconv.Atoi(string(input))
		},
	)

	output, remain, err := parser.Parse([]rune("123"), true)
	fmt.Printf("output:%#v remain:%#v err:%v\n", output, string(remain), err)

	// OUTPUT:
	// output:123 remain:"" err:<nil>
}

func ExampleMatchOnly() {
	parser := parcon.MatchOnly(parcon.Sequence(
		parcon.Many(1, parcon.MultiAlphas),
		parcon.Many(0, parcon.MultiDigits),
	))

	output, remain, err := parser.Parse([]rune("hello123"), true)
	fmt.Printf("output:%#v remain:%#v err:%v\n", string(output), string(remain), err)

	output, remain, err = parser.Parse([]rune("hello"), true)
	fmt.Printf("output:%#v remain:%#v err:%v\n", string(output), string(remain), err)

	output, remain, err = parser.Parse([]rune("123hello"), true)
	fmt.Printf("err:%v\n", err)

	// OUTPUT:
	// output:"hello123" remain:"" err:<nil>
	// output:"hello" remain:"" err:<nil>
	// err:invalid input: expected ALPHA but got "123hello"
}

func ExampleReplace() {
	parser := parcon.Or(
		// "yes" or "true" are true values.
		parcon.Replace(parcon.Or(
			parcon.TagStr("YES", "yes"),
			parcon.TagStr("TRUE", "true"),
		), true),
		parcon.Replace(parcon.MultiAlphas, false), // otherwise, parse as false value.
	)

	output, _, _ := parser.Parse([]rune("yes"), true)
	fmt.Println(output)

	output, _, _ = parser.Parse([]rune("no"), true)
	fmt.Println(output)

	// OUTPUT:
	// true
	// false
}
