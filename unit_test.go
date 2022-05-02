package parcon_test

import (
	"fmt"

	"github.com/macrat/parcon"
)

func ExampleTag() {
	parser := parcon.Tag("HELLO", []rune("hello"))

	output, remain, err := parser.Parse([]rune("hello world"))
	fmt.Printf("output:%#v remain:%#v err:%v\n", string(output), string(remain), err)

	// OUTPUT:
	// output:"hello" remain:" world" err:<nil>
}

func ExampleOneOf() {
	parser := parcon.OneOf("DIGIT", []rune("0123456789"))

	output, remain, err := parser.Parse([]rune("123 hello"))
	fmt.Printf("output:%#v remain:%#v err:%v\n", string(output), string(remain), err)

	// OUTPUT:
	// output:"1" remain:"23 hello" err:<nil>
}

func ExampleOneOfList() {
	parser := parcon.OneOfList("DIGITS", []rune("0123456789"))

	output, remain, err := parser.Parse([]rune("123 hello"))
	fmt.Printf("output:%#v remain:%#v err:%v\n", string(output), string(remain), err)

	// OUTPUT:
	// output:"123" remain:" hello" err:<nil>
}

func ExampleNoneOf() {
	parser := parcon.NoneOf("DIGIT", []rune("0123456789"))

	output, remain, err := parser.Parse([]rune("hello 123"))
	fmt.Printf("output:%#v remain:%#v err:%v\n", string(output), string(remain), err)

	// OUTPUT:
	// output:"h" remain:"ello 123" err:<nil>
}

func ExampleAnything() {
	parser := parcon.Anything[rune]()

	output, remain, err := parser.Parse([]rune("hello world"))
	fmt.Printf("output:%#v remain:%#v err:%v\n", string(output), string(remain), err)

	output, remain, err = parser.Parse(remain)
	fmt.Printf("output:%#v remain:%#v err:%v\n", string(output), string(remain), err)

	// OUTPUT:
	// output:"h" remain:"ello world" err:<nil>
	// output:"e" remain:"llo world" err:<nil>
}

func ExampleNothing() {
	parser := parcon.Nothing[rune]()

	output, remain, err := parser.Parse([]rune("hello world"))
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

	output, remain, err := parser.Parse([]rune("abc"))
	fmt.Printf("output:%#v remain:%#v err:%v\n", string(output), string(remain), err)

	output, remain, err = parser.Parse([]rune("bcd"))
	fmt.Printf("output:%#v remain:%#v err:%v\n", string(output), string(remain), err)

	_, _, err = parser.Parse([]rune("def"))
	fmt.Printf("err:%v\n", err)

	// OUTPUT:
	// output:"a" remain:"bc" err:<nil>
	// output:"b" remain:"cd" err:<nil>
	// err:expected ABC but got "def"
}

func ExampleTakeWhile() {
	parser := parcon.TakeWhile("ABC", func(c rune) bool {
		return c == 'a' || c == 'b' || c == 'c'
	})

	output, remain, err := parser.Parse([]rune("abcdef"))
	fmt.Printf("output:%#v remain:%#v err:%v\n", string(output), string(remain), err)

	// OUTPUT:
	// output:"abc" remain:"def" err:<nil>
}
