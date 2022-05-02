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
