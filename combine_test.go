package parcon_test

import (
	"fmt"

	"github.com/macrat/parcon"
)

func ExampleOptional() {
	parser := parcon.Optional(parcon.TagStr("HELLO", "hello"))

	output, remain, err := parser.Parse([]rune("hello world"))
	fmt.Printf("output:%#v remain:%#v err:%v\n", output, string(remain), err)

	output, remain, err = parser.Parse([]rune("foo bar"))
	fmt.Printf("output:%#v remain:%#v err:%v\n", output, string(remain), err)

	// OUTPUT:
	// output:"hello" remain:" world" err:<nil>
	// output:"" remain:"foo bar" err:<nil>
}

func ExampleOptionalWithDefault() {
	parser := parcon.OptionalWithDefault(parcon.TagStr("HELLO", "hello"), "not-found")

	output, remain, err := parser.Parse([]rune("hello world"))
	fmt.Printf("output:%#v remain:%#v err:%v\n", output, string(remain), err)

	output, remain, err = parser.Parse([]rune("foo bar"))
	fmt.Printf("output:%#v remain:%#v err:%v\n", output, string(remain), err)

	// OUTPUT:
	// output:"hello" remain:" world" err:<nil>
	// output:"not-found" remain:"foo bar" err:<nil>
}

func ExampleOr() {
	parser := parcon.Or(
		parcon.TagStr("HELLO", "hello"),
		parcon.TagStr("WORLD", "world"),
	)

	output, remain, err := parser.Parse([]rune("hello world"))
	fmt.Printf("output:%#v remain:%#v err:%v\n", output, string(remain), err)

	output, remain, err = parser.Parse([]rune("world hello"))
	fmt.Printf("output:%#v remain:%#v err:%v\n", output, string(remain), err)

	_, _, err = parser.Parse([]rune("foo bar"))
	fmt.Printf("err:%v\n", err)

	// OUTPUT:
	// output:"hello" remain:" world" err:<nil>
	// output:"world" remain:" hello" err:<nil>
	// err:expected one of [HELLO] [WORLD] but got "foo bar"
}

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

func ExampleSeparatedList() {
	parser := parcon.Map(parcon.SeparatedList(
		0,
		parcon.Tag("COMMA", []rune(",")),
		parcon.MultiDigits,
	), parcon.ToString)

	output, remain, err := parser.Parse([]rune("123,456,789"))
	fmt.Printf("output:%#v remain:%#v err:%v\n", output, string(remain), err)

	output, remain, err = parser.Parse([]rune("123,abc"))
	fmt.Printf("output:%#v remain:%#v err:%v\n", output, string(remain), err)

	output, remain, err = parser.Parse([]rune("abc"))
	fmt.Printf("output:%#v remain:%#v err:%v\n", output, string(remain), err)

	// OUTPUT:
	// output:[]string{"123", "456", "789"} remain:"" err:<nil>
	// output:[]string{"123"} remain:",abc" err:<nil>
	// output:[]string(nil) remain:"abc" err:<nil>
}

func ExampleSeparatedListLimited() {
	parser := parcon.Map(parcon.SeparatedListLimited(
		0,
		2,
		parcon.Tag("COMMA", []rune(",")),
		parcon.MultiDigits,
	), parcon.ToInt)

	output, remain, err := parser.Parse([]rune("123,456,789"))
	fmt.Printf("output:%#v remain:%#v err:%v\n", output, string(remain), err)

	// OUTPUT:
	// output:[]int{123, 456} remain:",789" err:<nil>
}

func ExampleMany() {
	parser := parcon.Many(0, parcon.TagStr("ITEM", "ab_"))

	output, remain, err := parser.Parse([]rune("ab_ab_ab_cd_"))
	fmt.Printf("output:%#v remain:%#v err:%v\n", output, string(remain), err)

	// OUTPUT:
	// output:[]string{"ab_", "ab_", "ab_"} remain:"cd_" err:<nil>
}

func ExampleManyLimited() {
	parser := parcon.ManyLimited(0, 2, parcon.TagStr("ITEM", "ab_"))

	output, remain, err := parser.Parse([]rune("ab_ab_ab_cd_"))
	fmt.Printf("output:%#v remain:%#v err:%v\n", output, string(remain), err)

	// OUTPUT:
	// output:[]string{"ab_", "ab_"} remain:"ab_cd_" err:<nil>
}

func ExampleRepeat() {
	parser := parcon.Repeat(3, parcon.SingleDigit)

	output, remain, err := parser.Parse([]rune("12345"))
	fmt.Printf("output:%#v remain:%#v err:%v\n", string(output), string(remain), err)

	// OUTPUT:
	// output:"123" remain:"45" err:<nil>
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
