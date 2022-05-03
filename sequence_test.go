package parcon_test

import (
	"fmt"

	"github.com/macrat/parcon"
)

func ExampleSequence() {
	parser := parcon.Sequence(
		parcon.TagStr("HELLO", "hello"),
		parcon.Convert(parcon.MultiSpaces, parcon.ToString),
		parcon.TagStr("WORLD", "world"),
	)

	output, remain, err := parser.Parse([]rune("hello   world"), true)
	fmt.Printf("output:%#v remain:%#v err:%v\n", output, string(remain), err)

	// OUTPUT:
	// output:[]string{"hello", "   ", "world"} remain:"" err:<nil>
}

func ExamplePair() {
	parser := parcon.Pair(
		parcon.Convert(parcon.MultiAlphas, parcon.ToString),
		parcon.Convert(parcon.MultiDigits, parcon.ToInt),
	)

	output, _, _ := parser.Parse([]rune("hello123"), true)
	fmt.Printf("first: %#v\n", output.First)
	fmt.Printf("second: %#v\n", output.Second)

	// OUTPUT:
	// first: "hello"
	// second: 123
}

func ExampleWithEnclosure() {
	parser := parcon.WithEnclosure(
		parcon.TagStr("OPEN_PAREN", "("),
		parcon.NoneOfStr("NOT_PAREN", "()"),
		parcon.TagStr("CLOSE_PAREN", ")"),
	)

	output, remain, err := parser.Parse([]rune("(hello world)"), true)
	fmt.Printf("output:%#v remain:%#v err:%v\n", output, string(remain), err)

	// OUTPUT:
	// output:"hello world" remain:"" err:<nil>
}

func ExampleWithPrefix() {
	parser := parcon.WithPrefix(
		parcon.TagStr("AT_SYMBOL", "@"),
		parcon.MultiAlphas,
	)

	output, remain, err := parser.Parse([]rune("@user"), true)
	fmt.Printf("output:%#v remain:%#v err:%v\n", string(output), string(remain), err)

	// OUTPUT:
	// output:"user" remain:"" err:<nil>
}

func ExampleWithSuffix() {
	parser := parcon.WithSuffix(
		parcon.NoneOfStr("NOT_SEMICOLON", ";"),
		parcon.TagStr("SEMICOLON", ";"),
	)

	output, remain, err := parser.Parse([]rune("hello world; foo bar;"), true)
	fmt.Printf("output:%#v remain:%#v err:%v\n", output, string(remain), err)

	// OUTPUT:
	// output:"hello world" remain:" foo bar;" err:<nil>
}
