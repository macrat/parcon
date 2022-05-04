Parcon
======

[![Go Reference](https://pkg.go.dev/badge/github.com/macrat/parcon.svg)](https://pkg.go.dev/github.com/macrat/parcon)

Parcon is a yet another parser combinator written in Go. (This is a experimental implementation and it is not effective at all)

``` golang
package main

import (
	"fmt"
	"strconv"

	pc "github.com/macrat/parcon"
)

type Color struct {
	Red   uint8
	Green uint8
	Blue  uint8
}

var HexNumber = pc.Convert(
	pc.Repeat(2, pc.SingleHexDigit),
	func(input []rune) (uint8, error) {
		i, err := strconv.ParseUint(string(input), 16, 8)
		return uint8(i), err
	},
)

var ColorParser = pc.Convert(
	pc.WithPrefix(
		pc.TagStr("HASH", "#"),
		pc.Repeat(3, HexNumber),
	),
	func(input []uint8) (Color, error) {
		return Color{
			Red:   input[0],
			Green: input[1],
			Blue:  input[2],
		}, nil
	},
)

func ParseColor(input string) (Color, error) {
	output, remain, err := ColorParser.Parse([]rune(input), true)
	if err != nil {
		return Color{}, err
	}
	if len(remain) != 0 {
		return Color{}, fmt.Errorf("found extra string: %#v", string(remain))
	}
	return output, nil
}

func main() {
	color, err := ParseColor("#2F14DF")
	if err != nil {
		panic(err)
	}

	fmt.Printf("red:%d green:%d blue:%d\n", color.Red, color.Green, color.Blue)

	// OUTPUT:
	// red:47 green:20 blue:223
}
```
