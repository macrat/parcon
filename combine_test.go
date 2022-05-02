package parcon_test

import (
	"fmt"

	"github.com/macrat/parcon"
)

func ExampleOptional() {
	parser := parcon.Optional[rune, []rune](parcon.TagS("HELLO", "hello"))

	output, remain, err := parser.Parse([]rune("hello world"))
	fmt.Printf("output:%#v remain:%#v err:%v\n", string(output), string(remain), err)

	output, remain, err = parser.Parse([]rune("foo bar"))
	fmt.Printf("output:%#v remain:%#v err:%v\n", string(output), string(remain), err)

	// OUTPUT:
	// output:"hello" remain:" world" err:<nil>
	// output:"" remain:"foo bar" err:<nil>
}

func ExampleOptionalWithDefault() {
	parser := parcon.OptionalWithDefault[rune](parcon.TagS("HELLO", "hello"), []rune("not-found"))

	output, remain, err := parser.Parse([]rune("hello world"))
	fmt.Printf("output:%#v remain:%#v err:%v\n", string(output), string(remain), err)

	output, remain, err = parser.Parse([]rune("foo bar"))
	fmt.Printf("output:%#v remain:%#v err:%v\n", string(output), string(remain), err)

	// OUTPUT:
	// output:"hello" remain:" world" err:<nil>
	// output:"not-found" remain:"foo bar" err:<nil>
}

func ExampleOr() {
	parser := parcon.Or[rune, []rune](
		parcon.TagS("HELLO", "hello"),
		parcon.TagS("WORLD", "world"),
	)

	output, remain, err := parser.Parse([]rune("hello world"))
	fmt.Printf("output:%#v remain:%#v err:%v\n", string(output), string(remain), err)

	output, remain, err = parser.Parse([]rune("world hello"))
	fmt.Printf("output:%#v remain:%#v err:%v\n", string(output), string(remain), err)

	_, _, err = parser.Parse([]rune("foo bar"))
	fmt.Printf("err:%v\n", err)

	// OUTPUT:
	// output:"hello" remain:" world" err:<nil>
	// output:"world" remain:" hello" err:<nil>
	// err:expected one of [HELLO] [WORLD] but got "foo bar"
}

func ExampleSequence() {
	parser := parcon.Map[rune](parcon.Sequence[rune, []rune](
		parcon.TagS("HELLO", "hello"),
		parcon.MultiSpaces,
		parcon.TagS("WORLD", "world"),
	), parcon.ToString)

	output, remain, err := parser.Parse([]rune("hello   world"))
	fmt.Printf("output:%#v remain:%#v err:%v\n", output, string(remain), err)

	// OUTPUT:
	// output:[]string{"hello", "   ", "world"} remain:"" err:<nil>
}

func ExamplePair() {
	parser := parcon.Pair[rune, string, int](
		parcon.Convert[rune](parcon.MultiAlphas, parcon.ToString),
		parcon.Convert[rune](parcon.MultiDigits, parcon.ToInt),
	)

	output, _, _ := parser.Parse([]rune("hello123"))
	fmt.Printf("first: %#v\n", output.First)
	fmt.Printf("second: %#v\n", output.Second)

	// OUTPUT:
	// first: "hello"
	// second: 123
}

func ExampleSeparatedList() {
	parser := parcon.Map[rune](parcon.SeparatedList[rune, []rune, []rune](
		0,
		parcon.TagS("COMMA", ","),
		parcon.MultiDigits,
	), parcon.ToString)

	output, remain, err := parser.Parse([]rune("123,456,789"))
	fmt.Printf("output:%#v remain:%#v err:%v\n", output, string(remain), err)

	output, remain, err = parser.Parse([]rune("abc"))
	fmt.Printf("output:%#v remain:%#v err:%v\n", output, string(remain), err)

	// OUTPUT:
	// output:[]string{"123", "456", "789"} remain:"" err:<nil>
	// output:[]string(nil) remain:"abc" err:<nil>
}

func ExampleSeparatedListLimited() {
	parser := parcon.Map[rune](parcon.SeparatedListLimited[rune, []rune, []rune](
		0,
		2,
		parcon.TagS("COMMA", ","),
		parcon.MultiDigits,
	), parcon.ToInt)

	output, remain, err := parser.Parse([]rune("123,456,789"))
	fmt.Printf("output:%#v remain:%#v err:%v\n", output, string(remain), err)

	// OUTPUT:
	// output:[]int{123, 456} remain:",789" err:<nil>
}

func ExampleMany() {
	parser := parcon.Map[rune](
		parcon.Many[rune, []rune](0, parcon.TagS("ITEM", "ab_")),
		parcon.ToString,
	)

	output, remain, err := parser.Parse([]rune("ab_ab_ab_cd_"))
	fmt.Printf("output:%#v remain:%#v err:%v\n", output, string(remain), err)

	// OUTPUT:
	// output:[]string{"ab_", "ab_", "ab_"} remain:"cd_" err:<nil>
}

func ExampleManyLimited() {
	parser := parcon.Map[rune](
		parcon.ManyLimited[rune, []rune](0, 2, parcon.TagS("ITEM", "ab_")),
		parcon.ToString,
	)

	output, remain, err := parser.Parse([]rune("ab_ab_ab_cd_"))
	fmt.Printf("output:%#v remain:%#v err:%v\n", output, string(remain), err)

	// OUTPUT:
	// output:[]string{"ab_", "ab_"} remain:"ab_cd_" err:<nil>
}

func ExampleRepeat() {
	parser := parcon.Repeat[rune, rune](3, parcon.SingleDigit)

	output, remain, err := parser.Parse([]rune("12345"))
	fmt.Printf("output:%#v remain:%#v err:%v\n", string(output), string(remain), err)

	// OUTPUT:
	// output:"123" remain:"45" err:<nil>
}

func ExampleDelimited() {
	parser := parcon.Delimited[rune, []rune, []rune, []rune](
		parcon.TagS("OPEN_PAREN", "("),
		parcon.Many[rune, rune](0, parcon.NoneOfS("NOT_PAREN", "()")),
		parcon.TagS("CLOSE_PAREN", ")"),
	)

	output, remain, err := parser.Parse([]rune("(hello world)"))
	fmt.Printf("output:%#v remain:%#v err:%v\n", string(output), string(remain), err)

	// OUTPUT:
	// output:"hello world" remain:"" err:<nil>
}

func ExampleWithPrefix() {
	parser := parcon.WithPrefix[rune, []rune, []rune](
		parcon.TagS("AT_SYMBOL", "@"),
		parcon.MultiAlphas,
	)

	output, remain, err := parser.Parse([]rune("@user"))
	fmt.Printf("output:%#v remain:%#v err:%v\n", string(output), string(remain), err)

	// OUTPUT:
	// output:"user" remain:"" err:<nil>
}

func ExampleWithSuffix() {
	parser := parcon.WithSuffix[rune, []rune, []rune](
		parcon.Many[rune, rune](1, parcon.NoneOfS("NOT_SEMICOLON", ";")),
		parcon.TagS("SEMICOLON", ";"),
	)

	output, remain, err := parser.Parse([]rune("hello world; foo bar;"))
	fmt.Printf("output:%#v remain:%#v err:%v\n", string(output), string(remain), err)

	// OUTPUT:
	// output:"hello world" remain:" foo bar;" err:<nil>
}

func ExampleNamed() {
	raw := parcon.Or[rune, []rune](
		parcon.WithPrefix[rune, []rune, []rune](parcon.TagS("AT_SYMBOL", "@"), parcon.MultiAlphas),
		parcon.WithPrefix[rune, []rune, []rune](parcon.TagS("HASH_SYMBOL", "#"), parcon.MultiAlphas),
	)
	_, _, err := raw.Parse([]rune("hello"))
	fmt.Println("raw:", err)

	named := parcon.Or[rune, []rune](
		parcon.Named[rune, []rune]("MENTION", parcon.WithPrefix[rune, []rune, []rune](parcon.TagS("AT_SYMBOL", "@"), parcon.MultiAlphas)),
		parcon.Named[rune, []rune]("TAG", parcon.WithPrefix[rune, []rune, []rune](parcon.TagS("HASH_SYMBOL", "#"), parcon.MultiAlphas)),
	)
	_, _, err = named.Parse([]rune("hello"))
	fmt.Println("named:", err)

	// OUTPUT:
	// raw: expected one of [AT_SYMBOL, multiple [ALPHA], NOTHING] [HASH_SYMBOL, multiple [ALPHA], NOTHING] but got "hello"
	// named: expected one of [MENTION] [TAG] but got "hello"
}
