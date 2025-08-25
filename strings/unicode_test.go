package strings

import (
	"math/rand/v2"
	"testing"
	"time"
	"unicode"

	ohsnap "github.com/okneniz/oh-snap"
	"github.com/okneniz/parsec/common"
)

func TestUnicodeHelpers(t *testing.T) {
	type test struct {
		name   string
		parser common.Combinator[rune, Position, rune]
		arb    ohsnap.Arbitrary[rune]
	}

	seed := time.Now().UnixNano()
	t.Logf("seed: %v", seed)
	rnd := rand.New(rand.NewPCG(0, uint64(seed)))

	tests := []test{
		// skip control chars because don't know how to make valid strings with them
		// {
		// 	name:   "control",
		// 	parser: Control("expected control char"),
		// 	arb:    ohsnap.RuneFromTable(rnd, unicode.C),
		// },
		{
			name:   "digit",
			parser: Digit("expected digit"),
			arb:    ohsnap.RuneFromTable(rnd, unicode.Digit),
		},
		{
			name:   "Graphic",
			parser: Graphic("expected graphic"),
			arb:    ohsnap.RuneFromTable(rnd, unicode.L),
		},
		{
			name:   "Letter",
			parser: Letter("expected letter"),
			arb:    ohsnap.RuneFromTable(rnd, unicode.Letter),
		},
		{
			name:   "Lower",
			parser: Lower("expected lower"),
			arb:    ohsnap.RuneFromTable(rnd, unicode.Lower),
		},
		{
			name:   "Mark",
			parser: Mark("expected mark"),
			arb:    ohsnap.RuneFromTable(rnd, unicode.Mark),
		},
		{
			name:   "Number",
			parser: Number("expected number"),
			arb:    ohsnap.RuneFromTable(rnd, unicode.Number),
		},
		{
			name:   "Punct",
			parser: Punct("expected punct"),
			arb:    ohsnap.RuneFromTable(rnd, unicode.Punct),
		},
		{
			name:   "Space",
			parser: Space("expected space"),
			arb:    ohsnap.RuneFromTable(rnd, unicode.Space),
		},
		{
			name:   "Print",
			parser: Print("expected print"),
			arb:    ohsnap.RuneFromTable(rnd, unicode.Title),
		},
		{
			name:   "Symbol",
			parser: Symbol("expected symbol"),
			arb:    ohsnap.RuneFromTable(rnd, unicode.Symbol),
		},
		{
			name:   "Title",
			parser: Title("expected title"),
			arb:    ohsnap.RuneFromTable(rnd, unicode.Title),
		},
		{
			name:   "Upper",
			parser: Upper("expected upper"),
			arb:    ohsnap.RuneFromTable(rnd, unicode.Upper),
		},
	}

	const iterations = 100_000

	for _, x := range tests {
		test := x

		t.Run(test.name, func(t *testing.T) {
			ohsnap.Check(t, iterations, test.arb, func(r rune) bool {
				input := string(r)

				result, err := ParseString(input, test.parser)
				if err != nil {
					t.Log("input", r, string(r))
					t.Log("result", result)
					t.Error(err)
					return false
				}

				return r == result
			})
		})
	}
}

func TestRangeTable(t *testing.T) {
	t.Parallel()

	seed := time.Now().UnixNano()
	t.Logf("seed: %v", seed)

	rnd := rand.New(rand.NewPCG(0, uint64(seed)))

	for name, tbl := range unicode.Categories {
		// skip control chars because don't know how to make valid strings with them
		if name[0] == 'C' {
			continue
		}

		t.Run(name, checkRangeTable(t, rnd, tbl))
	}

	for name, tbl := range unicode.Properties {
		t.Run(name, checkRangeTable(t, rnd, tbl))
	}
}

func checkRangeTable(t *testing.T, rnd *rand.Rand, tbl *unicode.RangeTable) func(t *testing.T) {
	t.Helper()

	const iterations = 100_000
	arb := ohsnap.RuneFromTable(rnd, tbl)

	parseChar := RangeTable("expected char from table", tbl)

	return func(t *testing.T) {
		ohsnap.Check(t, iterations, arb, func(r rune) bool {
			result, err := ParseString(string(r), parseChar)
			if err != nil {
				t.Log("input", r, string(r))
				t.Log("result", result)
				t.Error(err)
				return false
			}

			return r == result
		})
	}
}
