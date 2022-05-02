package parcon_test

import (
	"fmt"

	"github.com/macrat/parcon"
)

func ExampleParserFunc() {
	parser := parcon.Repeat[rune, string](2, parcon.ParserFunc[rune, string](func(input []rune) (output string, remain []rune, err error) {
		return string(input[:2]), input[2:], nil
	}))

	output, remain, err := parser.Parse([]rune("hello"))
	fmt.Printf("output:%#v remain:%#v err:%v\n", output, string(remain), err)

	// OUTPUT:
	// output:[]string{"he", "ll"} remain:"o" err:<nil>
}

func ExampleFunc() {
	parser := parcon.Repeat(2, parcon.Func(func(input []rune) (output string, remain []rune, err error) {
		return string(input[:2]), input[2:], nil
	}))

	output, remain, err := parser.Parse([]rune("hello"))
	fmt.Printf("output:%#v remain:%#v err:%v\n", output, string(remain), err)

	// OUTPUT:
	// output:[]string{"he", "ll"} remain:"o" err:<nil>
}

func ExampleNamed() {
	raw := parcon.Or(
		parcon.WithPrefix(parcon.Tag("AT_SYMBOL", []rune("@")), parcon.MultiAlphas),
		parcon.WithPrefix(parcon.Tag("HASH_SYMBOL", []rune("#")), parcon.MultiAlphas),
	)
	_, _, err := raw.Parse([]rune("hello"))
	fmt.Println("raw:", err)

	named := parcon.Or(
		parcon.Named("MENTION", parcon.WithPrefix(parcon.Tag("AT_SYMBOL", []rune("@")), parcon.MultiAlphas)),
		parcon.Named("TAG", parcon.WithPrefix(parcon.Tag("HASH_SYMBOL", []rune("#")), parcon.MultiAlphas)),
	)
	_, _, err = named.Parse([]rune("hello"))
	fmt.Println("named:", err)

	// OUTPUT:
	// raw: expected one of [AT_SYMBOL, ALPHA, NOTHING] [HASH_SYMBOL, ALPHA, NOTHING] but got "hello"
	// named: expected one of [MENTION] [TAG] but got "hello"
}
