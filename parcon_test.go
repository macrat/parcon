package parcon_test

import (
	"fmt"

	"github.com/macrat/parcon"
)

func ExampleParserFunc() {
	fun := func(input []rune, verbose bool) (output string, remain []rune, err error) {
		return string(input[:2]), input[2:], nil
	}
	parser := parcon.Repeat[rune, string](2, parcon.ParserFunc[rune, string](fun))

	output, remain, err := parser.Parse([]rune("hello"), true)
	fmt.Printf("output:%#v remain:%#v err:%v\n", output, string(remain), err)

	// OUTPUT:
	// output:[]string{"he", "ll"} remain:"o" err:<nil>
}

func ExampleFunc() {
	parser := parcon.Repeat(2, parcon.Func(func(input []rune, verbose bool) (output string, remain []rune, err error) {
		return string(input[:2]), input[2:], nil
	}))

	output, remain, err := parser.Parse([]rune("hello"), true)
	fmt.Printf("output:%#v remain:%#v err:%v\n", output, string(remain), err)

	// OUTPUT:
	// output:[]string{"he", "ll"} remain:"o" err:<nil>
}

func ExampleNamed() {
	raw := parcon.Or(
		parcon.WithPrefix(parcon.TagStr("AT_SYMBOL", "@"), parcon.MultiAlphas),
		parcon.WithPrefix(parcon.TagStr("HASH_SYMBOL", "#"), parcon.MultiAlphas),
	)
	_, _, err := raw.Parse([]rune("hello"), true)
	fmt.Println("raw:", err)

	named := parcon.Or(
		parcon.Named("MENTION", parcon.WithPrefix(parcon.TagStr("AT_SYMBOL", "@"), parcon.MultiAlphas)),
		parcon.Named("TAG", parcon.WithPrefix(parcon.TagStr("HASH_SYMBOL", "#"), parcon.MultiAlphas)),
	)
	_, _, err = named.Parse([]rune("hello"), true)
	fmt.Println("named:", err)

	// OUTPUT:
	// raw: invalid input: expected one of [AT_SYMBOL, ALPHA, NOTHING] [HASH_SYMBOL, ALPHA, NOTHING] but got "hello"
	// named: invalid input: expected one of [MENTION] [TAG] but got "hello"
}
