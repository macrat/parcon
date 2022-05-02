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

	output, remain, err := parser.Parse([]rune("123"))
	fmt.Printf("output:%#v remain:%#v err:%v\n", output, string(remain), err)

	// OUTPUT:
	// output:123 remain:"" err:<nil>
}

func ExampleMap() {
	parser := parcon.Map(
		parcon.SeparatedList(
			0,
			parcon.MultiSpaces,
			parcon.MultiDigits,
		),
		func(input []rune) (int, error) {
			return strconv.Atoi(string(input))
		},
	)

	output, remain, err := parser.Parse([]rune("123 456 789"))
	fmt.Printf("output:%#v remain:%#v err:%v\n", output, string(remain), err)

	// OUTPUT:
	// output:[]int{123, 456, 789} remain:"" err:<nil>
}

func ExampleMatchOnly() {
	parser := parcon.MatchOnly(parcon.Sequence(
		parcon.Many(1, parcon.MultiAlphas),
		parcon.Many(0, parcon.MultiDigits),
	))

	output, remain, err := parser.Parse([]rune("hello123"))
	fmt.Printf("output:%#v remain:%#v err:%v\n", string(output), string(remain), err)

	output, remain, err = parser.Parse([]rune("hello"))
	fmt.Printf("output:%#v remain:%#v err:%v\n", string(output), string(remain), err)

	output, remain, err = parser.Parse([]rune("123hello"))
	fmt.Printf("err:%v\n", err)

	// OUTPUT:
	// output:"hello123" remain:"" err:<nil>
	// output:"hello" remain:"" err:<nil>
	// err:expected ALPHA but got "123hello"
}

func ExampleReplace() {
	parser := parcon.Many(0, parcon.Or(
		parcon.Replace(parcon.Tag("NEWLINE", []rune(`\n`)), '\n'),
		parcon.Anything[rune](),
	))

	output, _, _ := parser.Parse([]rune(`hello\nworld`))
	fmt.Println(string(output))

	// OUTPUT:
	// hello
	// world
}
