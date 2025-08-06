package strings

import (
	"fmt"
	"testing"

	"github.com/okneniz/parsec/common"
	. "github.com/okneniz/parsec/testing"
)

func TestChainl(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		next := Satisfy("expected rune", true, common.Anything[rune])

		comb := Chainl(
			func(buffer common.Buffer[rune, Position]) (string, common.Error[Position]) {
				x, err := next(buffer)
				if err != nil {
					return "", err
				}

				return string(x), nil
			},
			func(
				buffer common.Buffer[rune, Position],
			) (func(string, string) string, common.Error[Position]) {
				return func(x, y string) string {
					return fmt.Sprintf("(%v %v)", x, y)
				}, nil
			},
			"default",
		)

		result, err := ParseString("abcd", comb)
		Check(t, err)
		AssertEq(t, result, "(((a b) c) d)")

		result, err = ParseString("ab", comb)
		Check(t, err)
		AssertEq(t, result, "(a b)")

		result, err = ParseString("a", comb)
		Check(t, err)
		AssertEq(t, result, "a")

		result, err = ParseString("", comb)
		Check(t, err)
		AssertEq(t, result, "default")
	})

	t.Run("case 2", func(t *testing.T) {
		next := Satisfy("expected byte", true, common.Anything[rune])
		c := 0

		comb := Chainl(
			func(buffer common.Buffer[rune, Position]) (string, common.Error[Position]) {
				c++
				if c > 1 {
					return "-", common.NewParseError(buffer.Position(), "test error")
				}

				x, err := next(buffer)
				if err != nil {
					return "", err
				}

				return string(x), nil
			},
			func(buffer common.Buffer[rune, Position]) (func(string, string) string, common.Error[Position]) {
				return func(x, y string) string {
					return fmt.Sprintf("(%v %v)", x, y)
				}, nil
			},
			"default",
		)

		result, err := ParseString("abcd", comb)
		Check(t, err)
		AssertEq(t, result, "a")

		c = 0

		result, err = ParseString("a", comb)
		Check(t, err)
		AssertEq(t, result, "a")

		c = 0

		result, err = ParseString("", comb)
		Check(t, err)
		AssertEq(t, result, "default")
	})

	t.Run("case 3", func(t *testing.T) {
		next := Satisfy("expected byte", true, common.Anything[rune])

		comb := Chainl(
			func(buffer common.Buffer[rune, Position]) (string, common.Error[Position]) {
				x, err := next(buffer)
				if err != nil {
					return "", err
				}

				return string(x), nil
			},
			func(buffer common.Buffer[rune, Position]) (func(string, string) string, common.Error[Position]) {
				return func(x, y string) string {
						return ""
					},
					common.NewParseError(buffer.Position(), "test error")
			},
			"default",
		)

		result, err := ParseString("abcd", comb)
		Check(t, err)
		AssertEq(t, result, "a")

		result, err = ParseString("a", comb)
		Check(t, err)
		AssertEq(t, result, "a")

		result, err = ParseString("", comb)
		Check(t, err)
		AssertEq(t, result, "default")
	})
}

func TestChainl1(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		next := Satisfy("anything", true, common.Anything[rune])

		comb := Chainl1(
			func(buffer common.Buffer[rune, Position]) (string, common.Error[Position]) {
				x, err := next(buffer)
				if err != nil {
					return "", err
				}

				return string(x), nil
			},
			func(buffer common.Buffer[rune, Position]) (func(string, string) string, common.Error[Position]) {
				return func(x, y string) string {
					return fmt.Sprintf("(%v %v)", x, y)
				}, nil
			},
		)

		result, err := ParseString("abcd", comb)
		Check(t, err)
		AssertEq(t, result, "(((a b) c) d)")

		result, err = ParseString("ab", comb)
		Check(t, err)
		AssertEq(t, result, "(a b)")

		result, err = ParseString("a", comb)
		Check(t, err)
		AssertEq(t, result, "a")

		result, err = ParseString("", comb)
		AssertError(t, err)
		AssertEq(t, result, "")
	})

	t.Run("case 2", func(t *testing.T) {
		next := Satisfy("anything", true, common.Anything[rune])
		c := 0

		comb := Chainl1(
			func(buffer common.Buffer[rune, Position]) (string, common.Error[Position]) {
				c++
				if c > 1 {
					return "-", common.NewParseError(buffer.Position(), "test error")
				}

				x, err := next(buffer)
				if err != nil {
					return "", err
				}

				return string(x), nil
			},
			func(buffer common.Buffer[rune, Position]) (func(string, string) string, common.Error[Position]) {
				return func(x, y string) string {
					return fmt.Sprintf("(%v %v)", x, y)
				}, nil
			},
		)

		result, err := ParseString("abcd", comb)
		Check(t, err)
		AssertEq(t, result, "a")

		c = 0

		result, err = ParseString("a", comb)
		Check(t, err)
		AssertEq(t, result, "a")

		c = 0

		result, err = ParseString("", comb)
		AssertError(t, err)
		AssertEq(t, result, "")
	})

	t.Run("case 3", func(t *testing.T) {
		next := Satisfy("anything", true, common.Anything[rune])

		comb := Chainl1(
			func(buffer common.Buffer[rune, Position]) (string, common.Error[Position]) {
				x, err := next(buffer)
				if err != nil {
					return "--", err
				}

				return string(x), nil
			},
			func(buffer common.Buffer[rune, Position]) (func(string, string) string, common.Error[Position]) {
				return func(x, y string) string { return "++" },
					common.NewParseError(buffer.Position(), "test error")
			},
		)

		result, err := ParseString("abcd", comb)
		Check(t, err)
		AssertEq(t, result, "a")

		result, err = ParseString("a", comb)
		Check(t, err)
		AssertEq(t, result, "a")

		result, err = ParseString("", comb)
		AssertError(t, err)
		AssertEq(t, result, "")
	})
}

func TestChainr(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		next := Satisfy("anything", true, common.Anything[rune])

		comb := Chainr(
			func(buffer common.Buffer[rune, Position]) (string, common.Error[Position]) {
				x, err := next(buffer)
				if err != nil {
					return "", err
				}

				return string(x), nil
			},
			func(buffer common.Buffer[rune, Position]) (func(string, string) string, common.Error[Position]) {
				return func(x, y string) string {
					return fmt.Sprintf("(%v %v)", x, y)
				}, nil
			},
			"default",
		)

		result, err := ParseString("abcd", comb)
		Check(t, err)
		AssertEq(t, result, "(a (b (c d)))")

		result, err = ParseString("ab", comb)
		Check(t, err)
		AssertEq(t, result, "(a b)")

		result, err = ParseString("a", comb)
		Check(t, err)
		AssertEq(t, result, "a")

		result, err = ParseString("", comb)
		Check(t, err)
		AssertEq(t, result, "default")
	})

	t.Run("case 2", func(t *testing.T) {
		next := Satisfy("anything", true, common.Anything[rune])
		c := 0

		comb := Chainr(
			func(buffer common.Buffer[rune, Position]) (string, common.Error[Position]) {
				c++
				if c > 1 {
					return "-", common.NewParseError(buffer.Position(), "test error")
				}

				x, err := next(buffer)
				if err != nil {
					return "", err
				}

				return string(x), nil
			},
			func(buffer common.Buffer[rune, Position]) (func(string, string) string, common.Error[Position]) {
				return func(x, y string) string {
					return fmt.Sprintf("(%v %v)", x, y)
				}, nil
			},
			"default",
		)

		result, err := ParseString("abcd", comb)
		Check(t, err)
		AssertEq(t, result, "a")

		c = 0

		result, err = ParseString("a", comb)
		Check(t, err)
		AssertEq(t, result, "a")

		c = 0

		result, err = ParseString("", comb)
		Check(t, err)
		AssertEq(t, result, "default")
	})

	t.Run("case 3", func(t *testing.T) {
		next := Satisfy("anything", true, common.Anything[rune])

		comb := Chainr(
			func(buffer common.Buffer[rune, Position]) (string, common.Error[Position]) {
				x, err := next(buffer)
				if err != nil {
					return "", err
				}

				return string(x), nil
			},
			func(buffer common.Buffer[rune, Position]) (func(string, string) string, common.Error[Position]) {
				return func(x, y string) string { return "" },
					common.NewParseError(buffer.Position(), "test error")
			},
			"default",
		)

		result, err := ParseString("abcd", comb)
		Check(t, err)
		AssertEq(t, result, "a")

		result, err = ParseString("a", comb)
		Check(t, err)
		AssertEq(t, result, "a")

		result, err = ParseString("", comb)
		Check(t, err)
		AssertEq(t, result, "default")
	})
}

func TestChainr1(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		next := Satisfy("anything", true, common.Anything[rune])

		comb := Chainr1(
			func(buffer common.Buffer[rune, Position]) (string, common.Error[Position]) {
				x, err := next(buffer)
				if err != nil {
					return "", err
				}

				return string(x), nil
			},
			func(buffer common.Buffer[rune, Position]) (func(string, string) string, common.Error[Position]) {
				return func(x, y string) string {
					return fmt.Sprintf("(%v %v)", x, y)
				}, nil
			},
		)

		result, err := ParseString("abcd", comb)
		Check(t, err)
		AssertEq(t, result, "(a (b (c d)))")

		result, err = ParseString("ab", comb)
		Check(t, err)
		AssertEq(t, result, "(a b)")

		result, err = ParseString("a", comb)
		Check(t, err)
		AssertEq(t, result, "a")

		result, err = ParseString("", comb)
		AssertError(t, err)
		AssertEq(t, result, "")
	})

	t.Run("case 2", func(t *testing.T) {
		next := Satisfy("expected any char", true, common.Anything[rune])
		c := 0

		comb := Chainr1(
			func(buffer common.Buffer[rune, Position]) (string, common.Error[Position]) {
				c++
				if c > 1 {
					return "-", common.NewParseError(buffer.Position(), "test error")
				}

				x, err := next(buffer)
				if err != nil {
					return "", err
				}

				return string(x), nil
			},
			func(buffer common.Buffer[rune, Position]) (func(string, string) string, common.Error[Position]) {
				return func(x, y string) string {
					return fmt.Sprintf("(%v %v)", x, y)
				}, nil
			},
		)

		result, err := ParseString("abcd", comb)
		Check(t, err)
		AssertEq(t, result, "a")

		c = 0

		result, err = ParseString("a", comb)
		Check(t, err)
		AssertEq(t, result, "a")

		c = 0

		result, err = ParseString("", comb)
		AssertError(t, err)
		AssertEq(t, result, "")
	})

	t.Run("case 3", func(t *testing.T) {
		next := Satisfy("anything", true, common.Anything[rune])

		comb := Chainr1(
			func(buffer common.Buffer[rune, Position]) (string, common.Error[Position]) {
				x, err := next(buffer)
				if err != nil {
					return "", err
				}

				return string(x), nil
			},
			func(buffer common.Buffer[rune, Position]) (func(string, string) string, common.Error[Position]) {
				return func(x, y string) string {
						return ""
					},
					common.NewParseError(buffer.Position(), "test error")
			},
		)

		result, err := ParseString("abcd", comb)
		Check(t, err)
		AssertEq(t, result, "a")

		result, err = ParseString("a", comb)
		Check(t, err)
		AssertEq(t, result, "a")

		result, err = ParseString("", comb)
		AssertError(t, err)
		AssertEq(t, result, "")
	})
}

func TestSepBy(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		comb := SepBy(
			0,
			NotEq("expected not ','", ','),
			Eq("expected ',' as separator", ','),
		)

		result, err := ParseString("a,b,c", comb)
		Check(t, err)
		AssertSlice(t, result, []rune{'a', 'b', 'c'})

		result, err = ParseString("", comb)
		Check(t, err)
		AssertSlice(t, result, nil)

		result, err = ParseString(",", comb)
		Check(t, err)
		AssertSlice(t, result, nil)

		result, err = ParseString(",a,b,c", comb)
		Check(t, err)
		AssertSlice(t, result, nil)

		result, err = ParseString("a,b,c,", comb)
		Check(t, err)
		AssertSlice(t, result, []rune{'a', 'b', 'c'})

		result, err = ParseString("a,b,", comb)
		Check(t, err)
		AssertSlice(t, result, []rune{'a', 'b'})

		result, err = ParseString("abc", comb)
		Check(t, err)
		AssertSlice(t, result, []rune{'a'})
	})

	t.Run("case 2", func(t *testing.T) {
		comb := SepBy(
			0,
			Satisfy("expected anything", true, common.Nothing[rune]),
			Eq("expected separator", ','),
		)

		result, err := ParseString("a,b,c", comb)
		Check(t, err)
		AssertSlice(t, result, nil)

		result, err = ParseString("", comb)
		Check(t, err)
		AssertSlice(t, result, nil)
	})

	t.Run("case 3", func(t *testing.T) {
		comb := SepBy(
			0,
			NotEq("expected not ','", ','),
			Satisfy("expected anything", true, common.Nothing[rune]),
		)

		result, err := ParseString("a,b,c", comb)
		Check(t, err)
		AssertSlice(t, result, []rune{'a'})

		result, err = ParseString("", comb)
		Check(t, err)
		AssertSlice(t, result, nil)
	})
}

func TestSepBy1(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		comb := SepBy1(
			0,
			"expected at least one rune separated by ','",
			NotEq("expected not ','", ','),
			Eq("expected ','", ','),
		)

		result, err := ParseString("a,b,c", comb)
		Check(t, err)
		AssertSlice(t, result, []rune{'a', 'b', 'c'})

		result, err = ParseString("", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)

		result, err = ParseString(",", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)

		result, err = ParseString(",a,b,c", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)

		result, err = ParseString("a,b,c,", comb)
		Check(t, err)
		AssertSlice(t, result, []rune{'a', 'b', 'c'})

		result, err = ParseString("a,b,", comb)
		Check(t, err)
		AssertSlice(t, result, []rune{'a', 'b'})

		result, err = ParseString("abc", comb)
		Check(t, err)
		AssertSlice(t, result, []rune{'a'})
	})

	t.Run("case 2", func(t *testing.T) {
		comb := SepBy1(
			0,
			"expected at least one rune",
			Satisfy("expected anything", true, common.Nothing[rune]),
			Eq("expected ','", ','),
		)

		result, err := ParseString("a,b,c", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)

		result, err = ParseString("", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})

	t.Run("case 3", func(t *testing.T) {
		comb := SepBy1(
			0,
			"expected at least one rune",
			NotEq("expecte not ','", ','),
			Satisfy("expected anything", true, common.Nothing[rune]),
		)

		result, err := ParseString("a,b,c", comb)
		Check(t, err)
		AssertSlice(t, result, []rune{'a'})

		result, err = ParseString("", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})
}

func TestEndBy(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		comb := EndBy(
			0,
			NotEq("expected not eq ','", ','),
			Eq("expected ','", ','),
		)

		result, err := ParseString("a,b,c", comb)
		Check(t, err)
		AssertSlice(t, result, []rune{'a', 'b'})

		result, err = ParseString("a,b,c,", comb)
		Check(t, err)
		AssertSlice(t, result, []rune{'a', 'b', 'c'})

		result, err = ParseString("", comb)
		Check(t, err)
		AssertSlice(t, result, nil)

		result, err = ParseString(",", comb)
		Check(t, err)
		AssertSlice(t, result, nil)

		result, err = ParseString("a", comb)
		Check(t, err)
		AssertSlice(t, result, nil)

		result, err = ParseString("a,", comb)
		Check(t, err)
		AssertSlice(t, result, []rune{'a'})

		result, err = ParseString(",a", comb)
		Check(t, err)
		AssertSlice(t, result, nil)

		result, err = ParseString("a,,", comb)
		Check(t, err)
		AssertSlice(t, result, []rune{'a'})

		result, err = ParseString(",a,", comb)
		Check(t, err)
		AssertSlice(t, result, nil)
	})

	t.Run("case 2", func(t *testing.T) {
		comb := EndBy(
			0,
			Satisfy("expected something", true, common.Nothing[rune]),
			Eq("expected separator", ','),
		)

		result, err := ParseString("a,b,c", comb)
		Check(t, err)
		AssertSlice(t, result, nil)

		result, err = ParseString("", comb)
		Check(t, err)
		AssertSlice(t, result, nil)
	})

	t.Run("case 3", func(t *testing.T) {
		comb := EndBy(
			0,
			NotEq("expected not ','", ','),
			Satisfy("expected something", true, common.Nothing[rune]),
		)

		result, err := ParseString("a,b,c", comb)
		Check(t, err)
		AssertSlice(t, result, nil)

		result, err = ParseString("", comb)
		Check(t, err)
		AssertSlice(t, result, nil)
	})
}

func TestEndBy1(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		comb := EndBy1(
			0,
			"expecte chars separated by ','",
			NotEq("expected not ','", ','),
			Eq("expected ','", ','),
		)

		result, err := ParseString("a,b,c", comb)
		Check(t, err)
		AssertSlice(t, result, []rune{'a', 'b'})

		result, err = ParseString("a,b,c,", comb)
		Check(t, err)
		AssertSlice(t, result, []rune{'a', 'b', 'c'})

		result, err = ParseString("", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)

		result, err = ParseString(",", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)

		result, err = ParseString("a", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)

		result, err = ParseString("a,", comb)
		Check(t, err)
		AssertSlice(t, result, []rune{'a'})

		result, err = ParseString(",a", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)

		result, err = ParseString("a,,", comb)
		Check(t, err)
		AssertSlice(t, result, []rune{'a'})

		result, err = ParseString(",a,", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})

	t.Run("case 2", func(t *testing.T) {
		comb := EndBy1(
			0,
			"expected char separated by ','",
			Satisfy("expected char", true, common.Nothing[rune]),
			Eq("expected ','", ','),
		)

		result, err := ParseString("a,b,c", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)

		result, err = ParseString("", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})

	t.Run("case 3", func(t *testing.T) {
		comb := EndBy1(
			0,
			"expecte char separated by ','",
			NotEq("expected not ','", ','),
			Satisfy("expected char", true, common.Nothing[rune]),
		)

		result, err := ParseString("a,b,c", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)

		result, err = ParseString("", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})
}

func TestSepEndBy(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		comb := SepEndBy(
			0,
			NotEq("expected not ','", ','),
			Eq("expected ','", ','),
		)

		result, err := ParseString("a,b,c", comb)
		Check(t, err)
		AssertSlice(t, result, []rune{'a', 'b', 'c'})

		result, err = ParseString("a,b,c,", comb)
		Check(t, err)
		AssertSlice(t, result, []rune{'a', 'b', 'c'})

		result, err = ParseString("a,b,c,,", comb)
		Check(t, err)
		AssertSlice(t, result, []rune{'a', 'b', 'c'})

		result, err = ParseString("", comb)
		Check(t, err)
		AssertSlice(t, result, nil)

		result, err = ParseString(",", comb)
		Check(t, err)
		AssertSlice(t, result, nil)

		result, err = ParseString(",a,b,c", comb)
		Check(t, err)
		AssertSlice(t, result, nil)
	})

	t.Run("case 2", func(t *testing.T) {
		comb := SepEndBy(
			0,
			Satisfy("expected char", true, common.Nothing[rune]),
			Eq("expected ','", ','),
		)

		result, err := ParseString("a,b,c", comb)
		Check(t, err)
		AssertSlice(t, result, nil)

		result, err = ParseString(",", comb)
		Check(t, err)
		AssertSlice(t, result, nil)

		result, err = ParseString("", comb)
		Check(t, err)
		AssertSlice(t, result, nil)
	})

	t.Run("case 3", func(t *testing.T) {
		comb := SepEndBy(
			0,
			NotEq("expected not ','", ','),
			Satisfy("expected any char", true, common.Nothing[rune]),
		)

		result, err := ParseString("a,b,c", comb)
		Check(t, err)
		AssertSlice(t, result, []rune{'a'})

		result, err = ParseString(",", comb)
		Check(t, err)
		AssertSlice(t, result, nil)

		result, err = ParseString("", comb)
		Check(t, err)
		AssertSlice(t, result, nil)
	})
}

func TestSepEndBy1(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		comb := SepEndBy1(
			0,
			"expected at least one char separated and ended by ','",
			NotEq("expected not ','", ','),
			Eq("expected ','", ','),
		)

		result, err := ParseString("a,b,c", comb)
		Check(t, err)
		AssertSlice(t, result, []rune{'a', 'b', 'c'})

		result, err = ParseString("a,b,c,", comb)
		Check(t, err)
		AssertSlice(t, result, []rune{'a', 'b', 'c'})

		result, err = ParseString("a,b,c,,", comb)
		Check(t, err)
		AssertSlice(t, result, []rune{'a', 'b', 'c'})

		result, err = ParseString("", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)

		result, err = ParseString(",", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)

		result, err = ParseString(",a,b,c", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})

	t.Run("case 2", func(t *testing.T) {
		comb := SepEndBy1(
			0,
			"expected at least one char separated and ended by ','",
			Satisfy("expected any char", true, common.Nothing[rune]),
			Eq("expected ','", ','),
		)

		result, err := ParseString("a,b,c", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)

		result, err = ParseString("", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})

	t.Run("case 3", func(t *testing.T) {
		comb := SepEndBy1(
			0,
			"expected at least one char separated and ended by ','",
			NotEq("expected not ','", ','),
			Satisfy("expected any char", true, common.Nothing[rune]),
		)

		result, err := ParseString("a,b,c", comb)
		Check(t, err)
		AssertSlice(t, result, []rune{'a'})

		result, err = ParseString("", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})
}

func TestManyTill(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		comb := ManyTill(
			0,
			Any(),
			Satisfy("expected 'd'", false, func(x rune) bool { return x == 'd' }),
		)

		result, err := ParseString("abcd", comb)
		Check(t, err)
		AssertSlice(t, result, []rune{
			rune('a'),
			rune('b'),
			rune('c'),
		})

		result, err = ParseString("", comb)
		Check(t, err)
		AssertSlice(t, result, nil)
	})

	t.Run("case 2", func(t *testing.T) {
		comb := ManyTill(
			0,
			Any(),
			Satisfy("expected any char", true, common.Nothing[rune]),
		)

		result, err := ParseString("abcd", comb)
		Check(t, err)
		AssertSlice(t, result, []rune{
			rune('a'),
			rune('b'),
			rune('c'),
			rune('d'),
		})

		result, err = ParseString("", comb)
		Check(t, err)
		AssertSlice(t, result, nil)
	})

	t.Run("case 3", func(t *testing.T) {
		comb := ManyTill(
			0,
			Any(),
			Any(),
		)

		result, err := ParseString("abcd", comb)
		Check(t, err)
		AssertSlice(t, result, nil)

		result, err = ParseString("", comb)
		Check(t, err)
		AssertSlice(t, result, nil)
	})
}
