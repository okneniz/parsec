package strings

import (
	"fmt"
	"testing"

	p "github.com/okneniz/parsec/common"
	. "github.com/okneniz/parsec/testing"
)

func TestChainl(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		next := Satisfy(true, p.Anything[rune])

		comb := Chainl(
			func(buffer p.Buffer[rune, Position]) (string, error) {
				x, err := next(buffer)
				if err != nil {
					return "", err
				}

				return string(x), nil
			},
			func(buffer p.Buffer[rune, Position]) (func(string, string) string, error) {
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
		next := Satisfy(true, p.Anything[rune])
		c := 0

		comb := Chainl(
			func(buffer p.Buffer[rune, Position]) (string, error) {
				c++
				if c > 1 {
					return "-", fmt.Errorf("test error")
				}

				x, err := next(buffer)
				if err != nil {
					return "", err
				}

				return string(x), nil
			},
			func(buffer p.Buffer[rune, Position]) (func(string, string) string, error) {
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
		next := Satisfy(true, p.Anything[rune])

		comb := Chainl(
			func(buffer p.Buffer[rune, Position]) (string, error) {
				x, err := next(buffer)
				if err != nil {
					return "", err
				}

				return string(x), nil
			},
			func(buffer p.Buffer[rune, Position]) (func(string, string) string, error) {
				return func(x, y string) string { return "" }, fmt.Errorf("test error")
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
		next := Satisfy(true, p.Anything[rune])

		comb := Chainl1(
			func(buffer p.Buffer[rune, Position]) (string, error) {
				x, err := next(buffer)
				if err != nil {
					return "", err
				}

				return string(x), nil
			},
			func(buffer p.Buffer[rune, Position]) (func(string, string) string, error) {
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
		next := Satisfy(true, p.Anything[rune])
		c := 0

		comb := Chainl1(
			func(buffer p.Buffer[rune, Position]) (string, error) {
				c++
				if c > 1 {
					return "-", fmt.Errorf("test error")
				}

				x, err := next(buffer)
				if err != nil {
					return "", err
				}

				return string(x), nil
			},
			func(buffer p.Buffer[rune, Position]) (func(string, string) string, error) {
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
		next := Satisfy(true, p.Anything[rune])

		comb := Chainl1(
			func(buffer p.Buffer[rune, Position]) (string, error) {
				x, err := next(buffer)
				if err != nil {
					return "--", err
				}

				return string(x), nil
			},
			func(buffer p.Buffer[rune, Position]) (func(string, string) string, error) {
				return func(x, y string) string { return "++" }, fmt.Errorf("test error")
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
		next := Satisfy(true, p.Anything[rune])

		comb := Chainr(
			func(buffer p.Buffer[rune, Position]) (string, error) {
				x, err := next(buffer)
				if err != nil {
					return "", err
				}

				return string(x), nil
			},
			func(buffer p.Buffer[rune, Position]) (func(string, string) string, error) {
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
		next := Satisfy(true, p.Anything[rune])
		c := 0

		comb := Chainr(
			func(buffer p.Buffer[rune, Position]) (string, error) {
				c++
				if c > 1 {
					return "-", fmt.Errorf("test error")
				}

				x, err := next(buffer)
				if err != nil {
					return "", err
				}

				return string(x), nil
			},
			func(buffer p.Buffer[rune, Position]) (func(string, string) string, error) {
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
		next := Satisfy(true, p.Anything[rune])

		comb := Chainr(
			func(buffer p.Buffer[rune, Position]) (string, error) {
				x, err := next(buffer)
				if err != nil {
					return "", err
				}

				return string(x), nil
			},
			func(buffer p.Buffer[rune, Position]) (func(string, string) string, error) {
				return func(x, y string) string { return "" }, fmt.Errorf("test error")
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
		next := Satisfy(true, p.Anything[rune])

		comb := Chainr1(
			func(buffer p.Buffer[rune, Position]) (string, error) {
				x, err := next(buffer)
				if err != nil {
					return "", err
				}

				return string(x), nil
			},
			func(buffer p.Buffer[rune, Position]) (func(string, string) string, error) {
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
		next := Satisfy(true, p.Anything[rune])
		c := 0

		comb := Chainr1(
			func(buffer p.Buffer[rune, Position]) (string, error) {
				c++
				if c > 1 {
					return "-", fmt.Errorf("test error")
				}

				x, err := next(buffer)
				if err != nil {
					return "", err
				}

				return string(x), nil
			},
			func(buffer p.Buffer[rune, Position]) (func(string, string) string, error) {
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
		next := Satisfy(true, p.Anything[rune])

		comb := Chainr1(
			func(buffer p.Buffer[rune, Position]) (string, error) {
				x, err := next(buffer)
				if err != nil {
					return "", err
				}

				return string(x), nil
			},
			func(buffer p.Buffer[rune, Position]) (func(string, string) string, error) {
				return func(x, y string) string {
					return ""
				}, fmt.Errorf("test error")
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
			NotEq(','),
			Eq(','),
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
			Satisfy(true, p.Nothing[rune]),
			Eq(','),
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
			NotEq(','),
			Satisfy(true, p.Nothing[rune]),
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
			NotEq(','),
			Eq(','),
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
			Satisfy(true, p.Nothing[rune]),
			Eq(','),
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
			NotEq(','),
			Satisfy(true, p.Nothing[rune]),
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
			NotEq(','),
			Eq(','),
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
			Satisfy(true, p.Nothing[rune]),
			Eq(','),
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
			NotEq(','),
			Satisfy(true, p.Nothing[rune]),
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
			NotEq(','),
			Eq(','),
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
			Satisfy(true, p.Nothing[rune]),
			Eq(','),
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
			NotEq(','),
			Satisfy(true, p.Nothing[rune]),
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
			NotEq(','),
			Eq(','),
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
			Satisfy(true, p.Nothing[rune]),
			Eq(','),
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
			NotEq(','),
			Satisfy(true, p.Nothing[rune]),
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
			NotEq(','),
			Eq(','),
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
			Satisfy(true, p.Nothing[rune]),
			Eq(','),
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
			NotEq(','),
			Satisfy(true, p.Nothing[rune]),
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
			Satisfy(false, func(x rune) bool { return x == 'd' }),
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
			Satisfy(true, p.Nothing[rune]),
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
