package parcon_test

import (
	"fmt"
	"strings"
	"testing"

	pc "github.com/macrat/parcon"
)

func generateSimpleList() ([]rune, int) {
	xs := make([]string, 1000)
	for i := range xs {
		xs[i] = fmt.Sprint(i)
	}
	return []rune(strings.Join(xs, ",")), len(xs)
}

func Benchmark_simpleListWithoutParcon(b *testing.B) {
	parser := func(input []rune) [][]rune {
		result := make([][]rune, 0, 1000)
		var buf []rune
		for _, x := range input {
			if x == ',' {
				result = append(result, buf)
			} else {
				buf = append(buf, x)
			}
		}
		return append(result, buf)
	}

	input, l := generateSimpleList()
	b.SetBytes(int64(len(input)))

	output := parser(input)
	if len(output) != l {
		b.Fatalf("found unexpected length of array: expected %d but got %d", l, len(output))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser(input)
	}
}

func Benchmark_simpleList(b *testing.B) {
	parser := pc.SeparatedListLimited(
		0,
		1000,
		pc.TagS("COMMA", ","),
		pc.ManyLimited(1, 5, pc.NoneOfS("NOT_COMMA", ",")),
	)

	input, l := generateSimpleList()
	b.SetBytes(int64(len(input)))

	output, _, err := parser.Parse(input)
	if err != nil {
		b.Fatalf("failed to parse: %s", err)
	}
	if len(output) != l {
		b.Fatalf("found unexpected length of array: expected %d but got %d", l, len(output))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser.Parse(input)
	}
}

func Benchmark_jsonParser(b *testing.B) {
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
	b.SetBytes(int64(len(input)))

	_, _, err := parser.Parse(input)
	if err != nil {
		b.Fatalf("failed to parse: %s", err)
	}

	b.ResetTimer()
	for i := 0; i <= b.N; i++ {
		parser.Parse(input)
	}
}
