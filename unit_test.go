package parcon_test

import (
	"fmt"

	"github.com/macrat/parcon"
)

func ExampleTagAs() {
	parser := parcon.TagAs("HELLO", []rune("hello"), "greetings")

	output, remain, err := parser.Parse([]rune("hello world"), true)
	fmt.Printf("output:%#v remain:%#v err:%v\n", string(output), string(remain), err)

	// OUTPUT:
	// output:"greetings" remain:" world" err:<nil>
}

func ExampleTag() {
	parser := parcon.Tag("HELLO", []rune("hello"))

	output, remain, err := parser.Parse([]rune("hello world"), true)
	fmt.Printf("output:%#v remain:%#v err:%v\n", string(output), string(remain), err)

	// OUTPUT:
	// output:"hello" remain:" world" err:<nil>
}

func ExampleTagStr() {
	parser := parcon.TagStr("HELLO", "hello")

	output, remain, err := parser.Parse([]rune("hello world"), true)
	fmt.Printf("output:%#v remain:%#v err:%v\n", output, string(remain), err)

	// OUTPUT:
	// output:"hello" remain:" world" err:<nil>
}

func ExampleOneOf() {
	parser := parcon.OneOf("DIGIT", []rune("0123456789"))

	output, remain, err := parser.Parse([]rune("123 hello"), true)
	fmt.Printf("output:%#v remain:%#v err:%v\n", string(output), string(remain), err)

	// OUTPUT:
	// output:"1" remain:"23 hello" err:<nil>
}

func ExampleOneOfList() {
	parser := parcon.OneOfList("DIGITS", []rune("0123456789"))

	output, remain, err := parser.Parse([]rune("123 hello"), true)
	fmt.Printf("output:%#v remain:%#v err:%v\n", string(output), string(remain), err)

	// OUTPUT:
	// output:"123" remain:" hello" err:<nil>
}

func ExampleOneOfStr() {
	parser := parcon.OneOfStr("DIGITS", "0123456789")

	output, remain, err := parser.Parse([]rune("123 hello"), true)
	fmt.Printf("output:%#v remain:%#v err:%v\n", output, string(remain), err)

	// OUTPUT:
	// output:"123" remain:" hello" err:<nil>
}

func ExampleNoneOf() {
	parser := parcon.NoneOf("DIGIT", []rune("0123456789"))

	output, remain, err := parser.Parse([]rune("hello 123"), true)
	fmt.Printf("output:%#v remain:%#v err:%v\n", string(output), string(remain), err)

	// OUTPUT:
	// output:"h" remain:"ello 123" err:<nil>
}

func ExampleNoneOfList() {
	parser := parcon.NoneOfList("DIGITS", []rune("0123456789"))

	output, remain, err := parser.Parse([]rune("hello 123"), true)
	fmt.Printf("output:%#v remain:%#v err:%v\n", string(output), string(remain), err)

	// OUTPUT:
	// output:"hello " remain:"123" err:<nil>
}

func ExampleNoneOfStr() {
	parser := parcon.NoneOfStr("DIGITS", "0123456789")

	output, remain, err := parser.Parse([]rune("hello 123"), true)
	fmt.Printf("output:%#v remain:%#v err:%v\n", output, string(remain), err)

	// OUTPUT:
	// output:"hello " remain:"123" err:<nil>
}

func ExampleAnything() {
	parser := parcon.Anything[rune]()

	output, remain, err := parser.Parse([]rune("hello world"), true)
	fmt.Printf("output:%#v remain:%#v err:%v\n", string(output), string(remain), err)

	output, remain, err = parser.Parse(remain, true)
	fmt.Printf("output:%#v remain:%#v err:%v\n", string(output), string(remain), err)

	// OUTPUT:
	// output:"h" remain:"ello world" err:<nil>
	// output:"e" remain:"llo world" err:<nil>
}

func ExampleNothing() {
	parser := parcon.Nothing[rune]()

	output, remain, err := parser.Parse([]rune("hello world"), true)
	fmt.Printf("%#v\n", output)
	fmt.Printf("%#v\n", string(remain))
	fmt.Println(err)

	// OUTPUT:
	// struct {}{}
	// "hello world"
	// <nil>
}

func ExampleTakeSingle() {
	parser := parcon.TakeSingle("ABC", func(c rune) bool {
		return 'a' <= c && c <= 'c'
	})

	output, remain, err := parser.Parse([]rune("abc"), true)
	fmt.Printf("output:%#v remain:%#v err:%v\n", string(output), string(remain), err)

	output, remain, err = parser.Parse([]rune("bcd"), true)
	fmt.Printf("output:%#v remain:%#v err:%v\n", string(output), string(remain), err)

	_, _, err = parser.Parse([]rune("def"), true)
	fmt.Printf("err:%v\n", err)

	// OUTPUT:
	// output:"a" remain:"bc" err:<nil>
	// output:"b" remain:"cd" err:<nil>
	// err:invalid input: expected ABC but got "def"
}

func ExampleTakeWhile() {
	parser := parcon.TakeWhile("ABC", func(c rune) bool {
		return c == 'a' || c == 'b' || c == 'c'
	})

	output, remain, err := parser.Parse([]rune("abcdef"), true)
	fmt.Printf("output:%#v remain:%#v err:%v\n", string(output), string(remain), err)

	// OUTPUT:
	// output:"abc" remain:"def" err:<nil>
}
