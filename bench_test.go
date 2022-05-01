package parcon_test

import (
	"strings"
	"fmt"
	"testing"

	pc "github.com/macrat/parcon"
)

func BenchmarkSimpleList(b *testing.B) {
	parser := pc.SeparatedListLimited(
		0,
		1000,
		pc.TagS("COMMA", ","),
		pc.ManyLimited(1, 5, pc.NoneOfS("NOT_COMMA", ",")),
	)

	xs := make([]string, 1000)
	for i := range xs {
		xs[i] = fmt.Sprint(i)
	}
	input := []rune(strings.Join(xs, ","))

	output, _, err := parser.Parse(input)
	if err != nil {
		b.Fatalf("failed to parse: %s", err)
	}
	if len(output) != len(xs) {
		b.Fatalf("found unexpected length of array: expected %d but got %d", len(xs), len(output))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser.Parse(input)
	}
}

func BenchmarkJsonParser(b *testing.B) {
	parser := JsonValue{}

	xs := make([]string, 100)
	for i := range xs {
		ys := make([]string, 100)
		for j := range ys {
			ys[j] = fmt.Sprintf(`"%d"`, j)
		}
		xs[i] = fmt.Sprintf(`"%d": [%s]`, i, strings.Join(ys, ", "))
	}
	input := []rune(fmt.Sprintf(`{%s}`, strings.Join(xs, ", ")))

	_, _, err := parser.Parse(input)
	if err != nil {
		b.Fatalf("failed to parse: %s", err)
	}

	b.ResetTimer()
	for i := 0; i <= b.N; i++ {
		parser.Parse(input)
	}
}
