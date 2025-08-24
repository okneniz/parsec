package strings

import (
	"fmt"
	"math"
	"math/rand/v2"
	"testing"
	"time"

	ohsnap "github.com/okneniz/oh-snap"
	"github.com/okneniz/parsec/common"
	"golang.org/x/exp/constraints"
)

func TestParens(t *testing.T) {
	t.Parallel()

	runTests(t, []test[rune]{
		{
			comb: Parens(Try(NoneOf("expected none parens", '(', ')'))),
			cases: []testCase[rune]{
				{
					input: "",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected '('",
					),
				},
				{
					input:  "(x)",
					output: 'x',
				},
				{
					input: "(x",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 2,
							index:  2,
						},
						"expected ')'",
					),
				},
				{
					input: "x)",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected '('",
					),
				},
				{
					input: "()",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 1,
							index:  1,
						},
						"expected none parens",
					),
				},
				{
					input: "(",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 1,
							index:  1,
						},
						"expected none parens",
					),
				},
				{
					input: ")",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected '('",
					),
				},
			},
		},
	})
}

func TestBraces(t *testing.T) {
	t.Parallel()

	runTests(t, []test[rune]{
		{
			comb: Braces(Try(NoneOf("expected none braces", '{', '}'))),
			cases: []testCase[rune]{
				{
					input: "",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected '{'",
					),
				},
				{
					input:  "{x}",
					output: 'x',
				},
				{
					input: "{x",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 2,
							index:  2,
						},
						"expected '}'",
					),
				},
				{
					input: "x}",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected '{'",
					),
				},
				{
					input: "{}",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 1,
							index:  1,
						},
						"expected none braces",
					),
				},
				{
					input: "{",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 1,
							index:  1,
						},
						"expected none braces",
					),
				},
				{
					input: "}",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected '{'",
					),
				},
			},
		},
	})
}

func TestAngles(t *testing.T) {
	t.Parallel()

	runTests(t, []test[rune]{
		{
			comb: Angles(Try(NoneOf("expected none angles", '<', '>'))),
			cases: []testCase[rune]{
				{
					input: "",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected '<'",
					),
				},
				{
					input:  "<x>",
					output: 'x',
				},
				{
					input: "<x",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 2,
							index:  2,
						},
						"expected '>'",
					),
				},
				{
					input: "x>",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected '<'",
					),
				},
				{
					input: "<>",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 1,
							index:  1,
						},
						"expected none angles",
					),
				},
				{
					input: "<",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 1,
							index:  1,
						},
						"expected none angles",
					),
				},
				{
					input: ">",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected '<'",
					),
				},
			},
		},
	})
}

func TestSquares(t *testing.T) {
	t.Parallel()

	runTests(t, []test[rune]{
		{
			comb: Squares(Try(NoneOf("expected none squares", '[', ']'))),
			cases: []testCase[rune]{
				{
					input: "",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected '['",
					),
				},
				{
					input:  "[x]",
					output: 'x',
				},
				{
					input: "[x",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 2,
							index:  2,
						},
						"expected ']'",
					),
				},
				{
					input: "x]",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected '['",
					),
				},
				{
					input: "[]",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 1,
							index:  1,
						},
						"expected none squares",
					),
				},
				{
					input: "[",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 1,
							index:  1,
						},
						"expected none squares",
					),
				},
				{
					input: "]",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected '['",
					),
				},
			},
		},
	})
}

func TestSemi(t *testing.T) {
	t.Parallel()

	runTests(t, []test[rune]{
		{
			comb: Semi(),
			cases: []testCase[rune]{
				{
					input: "",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected ';'",
					),
				},
				{
					input:  ";",
					output: ';',
				},
				{
					input:  ";x",
					output: ';',
				},
				{
					input: "x;",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected ';'",
					),
				},
			},
		},
	})
}

func TestComma(t *testing.T) {
	t.Parallel()

	runTests(t, []test[rune]{
		{
			comb: Comma(),
			cases: []testCase[rune]{
				{
					input: "",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected ','",
					),
				},
				{
					input:  ",",
					output: ',',
				},
				{
					input:  ",x",
					output: ',',
				},
				{
					input: "x,",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected ','",
					),
				},
			},
		},
	})
}

func TestColon(t *testing.T) {
	t.Parallel()

	runTests(t, []test[rune]{
		{
			comb: Colon(),
			cases: []testCase[rune]{
				{
					input: "",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected ':'",
					),
				},
				{
					input:  ":",
					output: ':',
				},
				{
					input:  ":x",
					output: ':',
				},
				{
					input: "x:",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected ':'",
					),
				},
			},
		},
	})
}

func TestDot(t *testing.T) {
	t.Parallel()

	runTests(t, []test[rune]{
		{
			comb: Dot(),
			cases: []testCase[rune]{
				{
					input: "",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected '.'",
					),
				},
				{
					input:  ".",
					output: '.',
				},
				{
					input:  ".x",
					output: '.',
				},
				{
					input: "x.",
					err: common.NewParseError(
						Position{
							line:   0,
							column: 0,
							index:  0,
						},
						"expected '.'",
					),
				},
			},
		},
	})
}

func TestUnsigned(t *testing.T) {
	seed := time.Now().UnixNano()
	t.Logf("seed: %v", seed)

	rnd := rand.New(rand.NewPCG(0, uint64(seed)))
	arb := ohsnap.NewBuilder(rnd)

	t.Run("uint8", func(t *testing.T) {
		parseUnsigned := Unsigned[uint8]()

		ohsnap.Check(t, math.MaxUint8, arb.Uint8(), func(data uint8) bool {
			str := fmt.Sprintf("%d", data)

			result, err := ParseString(str, parseUnsigned)
			if err != nil {
				t.Log("number", data)
				t.Log("number string", str)
				t.Log("parsed data", result)
				t.Error(err)
				return false
			}

			if data != result {
				t.Log("data", data)
				t.Log("json", str)
				t.Log("parsed data", result)
				return false
			}

			return true
		})
	})

	t.Run("uint16", func(t *testing.T) {
		parseUnsigned := Unsigned[uint16]()

		ohsnap.Check(t, math.MaxUint16, arb.Uint16(), func(data uint16) bool {
			str := fmt.Sprintf("%d", data)

			result, err := ParseString(str, parseUnsigned)
			if err != nil {
				t.Log("number", data)
				t.Log("number string", str)
				t.Log("parsed data", result)
				t.Error(err)
				return false
			}

			if data != result {
				t.Log("data", data)
				t.Log("json", str)
				t.Log("parsed data", result)
				return false
			}

			return true
		})
	})

	t.Run("uint32", func(t *testing.T) {
		parseUnsigned := Unsigned[uint32]()

		ohsnap.Check(t, 50000, arb.Uint32(), func(data uint32) bool {
			str := fmt.Sprintf("%d", data)

			result, err := ParseString(str, parseUnsigned)
			if err != nil {
				t.Log("number", data)
				t.Log("number string", str)
				t.Log("parsed data", result)
				t.Error(err)
				return false
			}

			if data != result {
				t.Log("data", data)
				t.Log("json", str)
				t.Log("parsed data", result)
				return false
			}

			return true
		})
	})

	t.Run("uint64", func(t *testing.T) {
		parseUnsigned := Unsigned[uint64]()

		ohsnap.Check(t, 50000, arb.Uint64(), func(data uint64) bool {
			str := fmt.Sprintf("%d", data)

			result, err := ParseString(str, parseUnsigned)
			if err != nil {
				t.Log("number", data)
				t.Log("number string", str)
				t.Log("parsed data", result)
				t.Error(err)
				return false
			}

			if data != result {
				t.Log("data", data)
				t.Log("json", str)
				t.Log("parsed data", result)
				return false
			}

			return true
		})
	})

	t.Run("uint", func(t *testing.T) {
		parseUnsigned := Unsigned[uint]()

		ohsnap.Check(t, 50000, arb.Uint(), func(data uint) bool {
			str := fmt.Sprintf("%d", data)

			result, err := ParseString(str, parseUnsigned)
			if err != nil {
				t.Log("number", data)
				t.Log("number string", str)
				t.Log("parsed data", result)
				t.Error(err)
				return false
			}

			if data != result {
				t.Log("data", data)
				t.Log("json", str)
				t.Log("parsed data", result)
				return false
			}

			return true
		})
	})

	t.Run("positive int8", func(t *testing.T) {
		parseUnsigned := Unsigned[int8]()

		ohsnap.Check(t, math.MaxInt8, arb.MinInt8(0).Int8(), func(data int8) bool {
			str := fmt.Sprintf("%d", data)

			result, err := ParseString(str, parseUnsigned)
			if err != nil {
				t.Log("number", data)
				t.Log("number string", str)
				t.Log("parsed data", result)
				t.Error(err)
				return false
			}

			if data != result {
				t.Log("data", data)
				t.Log("json", str)
				t.Log("parsed data", result)
				return false
			}

			return true
		})
	})

	t.Run("positive int16", func(t *testing.T) {
		parseUnsigned := Unsigned[int16]()

		ohsnap.Check(t, math.MaxInt16, arb.MinInt16(0).Int16(), func(data int16) bool {
			str := fmt.Sprintf("%d", data)

			result, err := ParseString(str, parseUnsigned)
			if err != nil {
				t.Log("number", data)
				t.Log("number string", str)
				t.Log("parsed data", result)
				t.Error(err)
				return false
			}

			if data != result {
				t.Log("data", data)
				t.Log("json", str)
				t.Log("parsed data", result)
				return false
			}

			return true
		})
	})

	t.Run("positive int32", func(t *testing.T) {
		parseUnsigned := Unsigned[int32]()

		ohsnap.Check(t, 50000, arb.MinInt32(0).Int32(), func(data int32) bool {
			str := fmt.Sprintf("%d", data)

			result, err := ParseString(str, parseUnsigned)
			if err != nil {
				t.Log("number", data)
				t.Log("number string", str)
				t.Log("parsed data", result)
				t.Error(err)
				return false
			}

			if data != result {
				t.Log("data", data)
				t.Log("json", str)
				t.Log("parsed data", result)
				return false
			}

			return true
		})
	})

	t.Run("positive int64", func(t *testing.T) {
		parseUnsigned := Unsigned[int64]()

		ohsnap.Check(t, 50000, arb.MinInt64(0).Int64(), func(data int64) bool {
			str := fmt.Sprintf("%d", data)

			result, err := ParseString(str, parseUnsigned)
			if err != nil {
				t.Log("number", data)
				t.Log("number string", str)
				t.Log("parsed data", result)
				t.Error(err)
				return false
			}

			if data != result {
				t.Log("data", data)
				t.Log("json", str)
				t.Log("parsed data", result)
				return false
			}

			return true
		})
	})

	t.Run("positive int", func(t *testing.T) {
		parseUnsigned := Unsigned[int]()

		ohsnap.Check(t, 50000, arb.MinInt(0).Int(), func(data int) bool {
			str := fmt.Sprintf("%d", data)

			result, err := ParseString(str, parseUnsigned)
			if err != nil {
				t.Log("number", data)
				t.Log("number string", str)
				t.Log("parsed data", result)
				t.Error(err)
				return false
			}

			if data != result {
				t.Log("data", data)
				t.Log("json", str)
				t.Log("parsed data", result)
				return false
			}

			return true
		})
	})
}

func checkUnsignedN[T constraints.Integer](t *testing.T, size int, arb ohsnap.Arbitrary[T]) func(t *testing.T) {
	t.Helper()

	return func(t *testing.T) {
		t.Helper()
		t.Parallel()

		parseUnsigned := UnsignedN[T](
			size,
			fmt.Sprintf("expected digit (size=%d)", size),
		)

		ohsnap.Check(t, 1000, arb, func(data T) bool {
			str := fmt.Sprintf("%d", data)

			result, err := ParseString(str, parseUnsigned)
			if err != nil {
				t.Log("number", data)
				t.Log("number string", str)
				t.Log("parsed data", result)
				t.Error(err)
				return false
			}

			if data != result {
				t.Log("data", data)
				t.Log("json", str)
				t.Log("parsed data", result)
				return false
			}

			return true
		})
	}
}

func TestUnsignedN(t *testing.T) {
	seed := time.Now().UnixNano()
	t.Logf("seed: %v", seed)

	rnd := rand.New(rand.NewPCG(0, uint64(seed)))
	builder := ohsnap.NewBuilder(rnd)

	t.Run("uint8", func(t *testing.T) {
		t.Parallel()

		sizes := map[int]*ohsnap.Builder{
			1: builder.MaxUint8(9),
			2: builder.MinUint8(10).MaxUint8(99),
			3: builder.MinUint8(100).MaxUint8(255),
		}

		for size, arb := range sizes {
			t.Run(fmt.Sprintf("size %d", size), checkUnsignedN(t, size, arb.Uint8()))
		}
	})

	t.Run("uint16", func(t *testing.T) {
		t.Parallel()

		sizes := map[int]*ohsnap.Builder{
			1: builder.MaxUint16(9),
			2: builder.MinUint16(10).MaxUint16(99),
			3: builder.MinUint16(100).MaxUint16(999),
			4: builder.MinUint16(1000).MaxUint16(9999),
			5: builder.MinUint16(10000).MaxUint16(math.MaxUint16),
		}

		for size, arb := range sizes {
			t.Run(fmt.Sprintf("size %d", size), checkUnsignedN(t, size, arb.Uint16()))
		}
	})

	t.Run("uint32", func(t *testing.T) {
		t.Parallel()

		sizes := map[int]*ohsnap.Builder{
			1: builder.MaxUint32(9),
			2: builder.MinUint32(10).MaxUint32(99),
			3: builder.MinUint32(100).MaxUint32(999),
			4: builder.MinUint32(1000).MaxUint32(9999),
			5: builder.MinUint32(10000).MaxUint32(99999),
			6: builder.MinUint32(100000).MaxUint32(999999),
			7: builder.MinUint32(1000000).MaxUint32(9999999),
		}

		for size, arb := range sizes {
			t.Run(fmt.Sprintf("size %d", size), checkUnsignedN(t, size, arb.Uint32()))
		}
	})

	t.Run("uint64", func(t *testing.T) {
		t.Parallel()

		sizes := map[int]*ohsnap.Builder{
			1: builder.MaxUint64(9),
			2: builder.MinUint64(10).MaxUint64(99),
			3: builder.MinUint64(100).MaxUint64(999),
			4: builder.MinUint64(1000).MaxUint64(9999),
			5: builder.MinUint64(10000).MaxUint64(99999),
			6: builder.MinUint64(100000).MaxUint64(999999),
			7: builder.MinUint64(1000000).MaxUint64(9999999),
		}

		for size, arb := range sizes {
			t.Run(fmt.Sprintf("size %d", size), checkUnsignedN(t, size, arb.Uint64()))
		}
	})

	t.Run("uint", func(t *testing.T) {
		t.Parallel()

		sizes := map[int]*ohsnap.Builder{
			1: builder.MaxUint(9),
			2: builder.MinUint(10).MaxUint(99),
			3: builder.MinUint(100).MaxUint(999),
			4: builder.MinUint(1000).MaxUint(9999),
			5: builder.MinUint(10000).MaxUint(99999),
			6: builder.MinUint(100000).MaxUint(999999),
			7: builder.MinUint(1000000).MaxUint(9999999),
		}

		for size, arb := range sizes {
			t.Run(fmt.Sprintf("size %d", size), checkUnsignedN(t, size, arb.Uint()))
		}
	})
}
