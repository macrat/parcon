package parcon_test

import (
	"fmt"

	"github.com/macrat/parcon"
)

func ExampleSequence() {
	parser := parcon.Map(parcon.Sequence(
		parcon.Tag("HELLO", []rune("hello")),
		parcon.MultiSpaces,
		parcon.Tag("WORLD", []rune("world")),
	), parcon.ToString)

	output, remain, err := parser.Parse([]rune("hello   world"))
	fmt.Printf("output:%#v remain:%#v err:%v\n", output, string(remain), err)

	// OUTPUT:
	// output:[]string{"hello", "   ", "world"} remain:"" err:<nil>
}

func ExamplePair() {
	parser := parcon.Pair(
		parcon.Convert(parcon.MultiAlphas, parcon.ToString),
		parcon.Convert(parcon.MultiDigits, parcon.ToInt),
	)

	output, _, _ := parser.Parse([]rune("hello123"))
	fmt.Printf("first: %#v\n", output.First)
	fmt.Printf("second: %#v\n", output.Second)

	// OUTPUT:
	// first: "hello"
	// second: 123
}

func ExampleDelimited() {
	parser := parcon.Delimited(
		parcon.Tag("OPEN_PAREN", []rune("(")),
		parcon.NoneOfStr("NOT_PAREN", "()"),
		parcon.Tag("CLOSE_PAREN", []rune(")")),
	)

	output, remain, err := parser.Parse([]rune("(hello world)"))
	fmt.Printf("output:%#v remain:%#v err:%v\n", output, string(remain), err)

	// OUTPUT:
	// output:"hello world" remain:"" err:<nil>
}

func ExampleWithPrefix() {
	parser := parcon.WithPrefix(
		parcon.Tag("AT_SYMBOL", []rune("@")),
		parcon.MultiAlphas,
	)

	output, remain, err := parser.Parse([]rune("@user"))
	fmt.Printf("output:%#v remain:%#v err:%v\n", string(output), string(remain), err)

	// OUTPUT:
	// output:"user" remain:"" err:<nil>
}

func ExampleWithSuffix() {
	parser := parcon.WithSuffix(
		parcon.NoneOfStr("NOT_SEMICOLON", ";"),
		parcon.Tag("SEMICOLON", []rune(";")),
	)

	output, remain, err := parser.Parse([]rune("hello world; foo bar;"))
	fmt.Printf("output:%#v remain:%#v err:%v\n", output, string(remain), err)

	// OUTPUT:
	// output:"hello world" remain:" foo bar;" err:<nil>
}
