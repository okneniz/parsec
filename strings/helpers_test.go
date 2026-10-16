package strings

import (
	"fmt"
	"math"
	"math/rand/v2"
	"testing"
	"time"

	ohsnap "github.com/okneniz/oh-snap"
	"github.com/okneniz/oh-snap/arbitrary"
	"github.com/okneniz/parsec"
)

func TestParens(t *testing.T) {
	t.Parallel()

	runTests(t, []test[rune]{
		{
			comb: Parens(Try(NoneOf[parsec.Stateless]("expected none parens", '(', ')'))),
			cases: []testCase[rune]{
				{
					input: "",
					err: parsec.NewParseError(
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
					err: parsec.NewParseError(
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
					err: parsec.NewParseError(
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
					err: parsec.NewParseError(
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
					err: parsec.NewParseError(
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
					err: parsec.NewParseError(
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
			comb: Braces(Try(NoneOf[parsec.Stateless]("expected none braces", '{', '}'))),
			cases: []testCase[rune]{
				{
					input: "",
					err: parsec.NewParseError(
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
					err: parsec.NewParseError(
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
					err: parsec.NewParseError(
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
					err: parsec.NewParseError(
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
					err: parsec.NewParseError(
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
					err: parsec.NewParseError(
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
			comb: Angles(Try(NoneOf[parsec.Stateless]("expected none angles", '<', '>'))),
			cases: []testCase[rune]{
				{
					input: "",
					err: parsec.NewParseError(
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
					err: parsec.NewParseError(
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
					err: parsec.NewParseError(
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
					err: parsec.NewParseError(
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
					err: parsec.NewParseError(
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
					err: parsec.NewParseError(
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
			comb: Squares(Try(NoneOf[parsec.Stateless]("expected none squares", '[', ']'))),
			cases: []testCase[rune]{
				{
					input: "",
					err: parsec.NewParseError(
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
					err: parsec.NewParseError(
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
					err: parsec.NewParseError(
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
					err: parsec.NewParseError(
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
					err: parsec.NewParseError(
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
					err: parsec.NewParseError(
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
			comb: Semi[parsec.Stateless](),
			cases: []testCase[rune]{
				{
					input: "",
					err: parsec.NewParseError(
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
					err: parsec.NewParseError(
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
			comb: Comma[parsec.Stateless](),
			cases: []testCase[rune]{
				{
					input: "",
					err: parsec.NewParseError(
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
					err: parsec.NewParseError(
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
			comb: Colon[parsec.Stateless](),
			cases: []testCase[rune]{
				{
					input: "",
					err: parsec.NewParseError(
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
					err: parsec.NewParseError(
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
			comb: Dot[parsec.Stateless](),
			cases: []testCase[rune]{
				{
					input: "",
					err: parsec.NewParseError(
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
					err: parsec.NewParseError(
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

	t.Run("uint8", func(t *testing.T) {
		parseUnsigned := Unsigned[uint8, parsec.Stateless]()

		ohsnap.Check(t, math.MaxUint8, arbitrary.Uint8(rnd, 0, math.MaxUint8), func(data uint8) bool {
			str := fmt.Sprintf("%d", data)

			result, err := ParseString(parsec.Stateless{}, str, parseUnsigned)
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
		parseUnsigned := Unsigned[uint16, parsec.Stateless]()

		ohsnap.Check(t, math.MaxUint16, arbitrary.Uint16(rnd, 0, math.MaxUint16), func(data uint16) bool {
			str := fmt.Sprintf("%d", data)

			result, err := ParseString(parsec.Stateless{}, str, parseUnsigned)
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
		parseUnsigned := Unsigned[uint32, parsec.Stateless]()

		ohsnap.Check(t, 50000, arbitrary.Uint32(rnd, 0, math.MaxUint32), func(data uint32) bool {
			str := fmt.Sprintf("%d", data)

			result, err := ParseString(parsec.Stateless{}, str, parseUnsigned)
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
		parseUnsigned := Unsigned[uint64, parsec.Stateless]()

		ohsnap.Check(t, 50000, arbitrary.Uint64(rnd, 0, math.MaxUint64), func(data uint64) bool {
			str := fmt.Sprintf("%d", data)

			result, err := ParseString(parsec.Stateless{}, str, parseUnsigned)
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
		parseUnsigned := Unsigned[uint, parsec.Stateless]()

		ohsnap.Check(t, 50000, arbitrary.Uint(rnd, 0, math.MaxUint), func(data uint) bool {
			str := fmt.Sprintf("%d", data)

			result, err := ParseString(parsec.Stateless{}, str, parseUnsigned)
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
		parseUnsigned := Unsigned[int8, parsec.Stateless]()

		ohsnap.Check(t, math.MaxInt8, arbitrary.Int8(rnd, 0, math.MaxInt8), func(data int8) bool {
			str := fmt.Sprintf("%d", data)

			result, err := ParseString(parsec.Stateless{}, str, parseUnsigned)
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
		parseUnsigned := Unsigned[int16, parsec.Stateless]()

		ohsnap.Check(t, math.MaxInt16, arbitrary.Int16(rnd, 0, math.MaxInt16), func(data int16) bool {
			str := fmt.Sprintf("%d", data)

			result, err := ParseString(parsec.Stateless{}, str, parseUnsigned)
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
		parseUnsigned := Unsigned[int32, parsec.Stateless]()

		ohsnap.Check(t, 50000, arbitrary.Int32(rnd, 0, math.MaxInt32), func(data int32) bool {
			str := fmt.Sprintf("%d", data)

			result, err := ParseString(parsec.Stateless{}, str, parseUnsigned)
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
		parseUnsigned := Unsigned[int64, parsec.Stateless]()

		ohsnap.Check(t, 50000, arbitrary.Int64(rnd, 0, math.MaxInt64), func(data int64) bool {
			str := fmt.Sprintf("%d", data)

			result, err := ParseString(parsec.Stateless{}, str, parseUnsigned)
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
		parseUnsigned := Unsigned[int, parsec.Stateless]()

		ohsnap.Check(t, 50000, arbitrary.Int(rnd, 0, math.MaxInt), func(data int) bool {
			str := fmt.Sprintf("%d", data)

			result, err := ParseString(parsec.Stateless{}, str, parseUnsigned)
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

func checkUnsignedN[T parsec.Integer](t *testing.T, size int, arb ohsnap.Arbitrary[T]) func(t *testing.T) {
	t.Helper()

	return func(t *testing.T) {
		t.Helper()
		t.Parallel()

		parseUnsigned := UnsignedN[T, parsec.Stateless](
			size,
			fmt.Sprintf("expected digit (size=%d)", size),
		)

		ohsnap.Check(t, 1000, arb, func(data T) bool {
			str := fmt.Sprintf("%d", data)

			result, err := ParseString(parsec.Stateless{}, str, parseUnsigned)
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

	t.Run("uint8", func(t *testing.T) {
		t.Parallel()

		sizes := map[int]ohsnap.Arbitrary[uint8]{
			1: arbitrary.Uint8(rnd, 0, 9),
			2: arbitrary.Uint8(rnd, 10, 99),
			3: arbitrary.Uint8(rnd, 100, math.MaxUint8),
		}

		for size, arb := range sizes {
			t.Run(fmt.Sprintf("size %d", size), checkUnsignedN(t, size, arb))
		}
	})

	t.Run("uint16", func(t *testing.T) {
		t.Parallel()

		sizes := map[int]ohsnap.Arbitrary[uint16]{
			1: arbitrary.Uint16(rnd, 0, 9),
			2: arbitrary.Uint16(rnd, 10, 99),
			3: arbitrary.Uint16(rnd, 100, 999),
			4: arbitrary.Uint16(rnd, 1000, 9999),
			5: arbitrary.Uint16(rnd, 10000, math.MaxUint16),
		}

		for size, arb := range sizes {
			t.Run(fmt.Sprintf("size %d", size), checkUnsignedN(t, size, arb))
		}
	})

	t.Run("uint32", func(t *testing.T) {
		t.Parallel()

		sizes := map[int]ohsnap.Arbitrary[uint32]{
			1: arbitrary.Uint32(rnd, 0, 9),
			2: arbitrary.Uint32(rnd, 10, 99),
			3: arbitrary.Uint32(rnd, 100, 999),
			4: arbitrary.Uint32(rnd, 1000, 9999),
			5: arbitrary.Uint32(rnd, 10000, 99999),
			6: arbitrary.Uint32(rnd, 100000, 999999),
			7: arbitrary.Uint32(rnd, 1000000, 9999999),
		}

		for size, arb := range sizes {
			t.Run(fmt.Sprintf("size %d", size), checkUnsignedN(t, size, arb))
		}
	})

	t.Run("uint64", func(t *testing.T) {
		t.Parallel()

		sizes := map[int]ohsnap.Arbitrary[uint64]{
			1: arbitrary.Uint64(rnd, 0, 9),
			2: arbitrary.Uint64(rnd, 10, 99),
			3: arbitrary.Uint64(rnd, 100, 999),
			4: arbitrary.Uint64(rnd, 1000, 9999),
			5: arbitrary.Uint64(rnd, 10000, 99999),
			6: arbitrary.Uint64(rnd, 100000, 999999),
			7: arbitrary.Uint64(rnd, 1000000, 9999999),
		}

		for size, arb := range sizes {
			t.Run(fmt.Sprintf("size %d", size), checkUnsignedN(t, size, arb))
		}
	})

	t.Run("uint", func(t *testing.T) {
		t.Parallel()

		sizes := map[int]ohsnap.Arbitrary[uint]{
			1: arbitrary.Uint(rnd, 0, 9),
			2: arbitrary.Uint(rnd, 10, 99),
			3: arbitrary.Uint(rnd, 100, 999),
			4: arbitrary.Uint(rnd, 1000, 9999),
			5: arbitrary.Uint(rnd, 10000, 99999),
			6: arbitrary.Uint(rnd, 100000, 999999),
			7: arbitrary.Uint(rnd, 1000000, 9999999),
		}

		for size, arb := range sizes {
			t.Run(fmt.Sprintf("size %d", size), checkUnsignedN(t, size, arb))
		}
	})
}
