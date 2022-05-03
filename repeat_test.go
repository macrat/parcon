package parcon_test

import (
	"fmt"

	"github.com/macrat/parcon"
)

func ExampleSeparatedList() {
	parser := parcon.Map(parcon.SeparatedList(
		0,
		parcon.Tag("COMMA", []rune(",")),
		parcon.MultiDigits,
	), parcon.ToString)

	output, remain, err := parser.Parse([]rune("123,456,789"), true)
	fmt.Printf("output:%#v remain:%#v err:%v\n", output, string(remain), err)

	output, remain, err = parser.Parse([]rune("123,abc"), true)
	fmt.Printf("output:%#v remain:%#v err:%v\n", output, string(remain), err)

	output, remain, err = parser.Parse([]rune("abc"), true)
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

	output, remain, err := parser.Parse([]rune("123,456,789"), true)
	fmt.Printf("output:%#v remain:%#v err:%v\n", output, string(remain), err)

	// OUTPUT:
	// output:[]int{123, 456} remain:",789" err:<nil>
}

func ExampleMany() {
	parser := parcon.Many(0, parcon.TagStr("ITEM", "ab_"))

	output, remain, err := parser.Parse([]rune("ab_ab_ab_cd_"), true)
	fmt.Printf("output:%#v remain:%#v err:%v\n", output, string(remain), err)

	// OUTPUT:
	// output:[]string{"ab_", "ab_", "ab_"} remain:"cd_" err:<nil>
}

func ExampleManyLimited() {
	parser := parcon.ManyLimited(0, 2, parcon.TagStr("ITEM", "ab_"))

	output, remain, err := parser.Parse([]rune("ab_ab_ab_cd_"), true)
	fmt.Printf("output:%#v remain:%#v err:%v\n", output, string(remain), err)

	// OUTPUT:
	// output:[]string{"ab_", "ab_"} remain:"ab_cd_" err:<nil>
}

func ExampleRepeat() {
	parser := parcon.Repeat(3, parcon.SingleDigit)

	output, remain, err := parser.Parse([]rune("12345"), true)
	fmt.Printf("output:%#v remain:%#v err:%v\n", string(output), string(remain), err)

	// OUTPUT:
	// output:"123" remain:"45" err:<nil>
}
