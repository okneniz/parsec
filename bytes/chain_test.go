package bytes

import (
	"fmt"
	"testing"

	p "github.com/okneniz/parsec/common"
	. "github.com/okneniz/parsec/testing"
)

func TestChainl(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		next := Satisfy(true, p.Anything[byte])

		comb := Chainl(
			func(buffer p.Buffer[byte, int]) (string, error) {
				x, err := next(buffer)
				if err != nil {
					return "", err
				}

				return string(x), nil
			},
			func(buffer p.Buffer[byte, int]) (func(string, string) string, error) {
				return func(x, y string) string {
					return fmt.Sprintf("(%v %v)", x, y)
				}, nil
			},
			"default",
		)

		result, err := Parse([]byte("abcd"), comb)
		Check(t, err)
		AssertEq(t, result, "(((a b) c) d)")

		result, err = Parse([]byte("a"), comb)
		Check(t, err)
		AssertEq(t, result, "a")

		result, err = Parse([]byte(""), comb)
		Check(t, err)
		AssertEq(t, result, "default")
	})

	t.Run("case 2", func(t *testing.T) {
		next := Satisfy(true, p.Anything[byte])
		c := 0

		comb := Chainl(
			func(buffer p.Buffer[byte, int]) (string, error) {
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
			func(buffer p.Buffer[byte, int]) (func(string, string) string, error) {
				return func(x, y string) string {
					return fmt.Sprintf("(%v %v)", x, y)
				}, nil
			},
			"default",
		)

		result, err := Parse([]byte("abcd"), comb)
		Check(t, err)
		AssertEq(t, result, "a")

		c = 0

		result, err = Parse([]byte("a"), comb)
		Check(t, err)
		AssertEq(t, result, "a")

		c = 0

		result, err = Parse([]byte(""), comb)
		Check(t, err)
		AssertEq(t, result, "default")
	})

	t.Run("case 3", func(t *testing.T) {
		next := Satisfy(true, p.Anything[byte])

		comb := Chainl(
			func(buffer p.Buffer[byte, int]) (string, error) {
				x, err := next(buffer)
				if err != nil {
					return "", err
				}

				return string(x), nil
			},
			func(buffer p.Buffer[byte, int]) (func(string, string) string, error) {
				return func(x, y string) string { return "" }, fmt.Errorf("test error")
			},
			"default",
		)

		result, err := Parse([]byte("abcd"), comb)
		Check(t, err)
		AssertEq(t, result, "a")

		result, err = Parse([]byte("a"), comb)
		Check(t, err)
		AssertEq(t, result, "a")

		result, err = Parse([]byte(""), comb)
		Check(t, err)
		AssertEq(t, result, "default")
	})
}

func TestChainl1(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		next := Satisfy(true, p.Anything[byte])

		comb := Chainl1(
			func(buffer p.Buffer[byte, int]) (string, error) {
				x, err := next(buffer)
				if err != nil {
					return "", err
				}

				return string(x), nil
			},
			func(buffer p.Buffer[byte, int]) (func(string, string) string, error) {
				return func(x, y string) string {
					return fmt.Sprintf("(%v %v)", x, y)
				}, nil
			},
		)

		result, err := Parse([]byte("abcd"), comb)
		Check(t, err)
		AssertEq(t, result, "(((a b) c) d)")

		result, err = Parse([]byte("a"), comb)
		Check(t, err)
		AssertEq(t, result, "a")

		result, err = Parse([]byte(""), comb)
		AssertError(t, err)
		AssertEq(t, result, "")
	})

	t.Run("case 2", func(t *testing.T) {
		next := Satisfy(true, p.Anything[byte])
		c := 0

		comb := Chainl1(
			func(buffer p.Buffer[byte, int]) (string, error) {
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
			func(buffer p.Buffer[byte, int]) (func(string, string) string, error) {
				return func(x, y string) string {
					return fmt.Sprintf("(%v %v)", x, y)
				}, nil
			},
		)

		result, err := Parse([]byte("abcd"), comb)
		Check(t, err)
		AssertEq(t, result, "a")

		c = 0

		result, err = Parse([]byte("a"), comb)
		Check(t, err)
		AssertEq(t, result, "a")

		c = 0

		result, err = Parse([]byte(""), comb)
		AssertError(t, err)
		AssertEq(t, result, "")
	})

	t.Run("case 3", func(t *testing.T) {
		next := Satisfy(true, p.Anything[byte])

		comb := Chainl1(
			func(buffer p.Buffer[byte, int]) (string, error) {
				x, err := next(buffer)
				if err != nil {
					return "--", err
				}

				return string(x), nil
			},
			func(buffer p.Buffer[byte, int]) (func(string, string) string, error) {
				return func(x, y string) string { return "++" }, fmt.Errorf("test error")
			},
		)

		result, err := Parse([]byte("abcd"), comb)
		Check(t, err)
		AssertEq(t, result, "a")

		result, err = Parse([]byte("a"), comb)
		Check(t, err)
		AssertEq(t, result, "a")

		result, err = Parse([]byte(""), comb)
		AssertError(t, err)
		AssertEq(t, result, "")
	})
}

func TestChainr(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		next := Satisfy(true, p.Anything[byte])

		comb := Chainr(
			func(buffer p.Buffer[byte, int]) (string, error) {
				x, err := next(buffer)
				if err != nil {
					return "", err
				}

				return string(x), nil
			},
			func(buffer p.Buffer[byte, int]) (func(string, string) string, error) {
				return func(x, y string) string {
					return fmt.Sprintf("(%v %v)", x, y)
				}, nil
			},
			"default",
		)

		result, err := Parse([]byte("abcd"), comb)
		Check(t, err)
		AssertEq(t, result, "(a (b (c d)))")

		result, err = Parse([]byte("a"), comb)
		Check(t, err)
		AssertEq(t, result, "a")

		result, err = Parse([]byte(""), comb)
		Check(t, err)
		AssertEq(t, result, "default")
	})

	t.Run("case 2", func(t *testing.T) {
		next := Satisfy(true, p.Anything[byte])
		c := 0

		comb := Chainr(
			func(buffer p.Buffer[byte, int]) (string, error) {
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
			func(buffer p.Buffer[byte, int]) (func(string, string) string, error) {
				return func(x, y string) string {
					return fmt.Sprintf("(%v %v)", x, y)
				}, nil
			},
			"default",
		)

		result, err := Parse([]byte("abcd"), comb)
		Check(t, err)
		AssertEq(t, result, "a")

		c = 0

		result, err = Parse([]byte("a"), comb)
		Check(t, err)
		AssertEq(t, result, "a")

		c = 0

		result, err = Parse([]byte(""), comb)
		Check(t, err)
		AssertEq(t, result, "default")
	})

	t.Run("case 3", func(t *testing.T) {
		next := Satisfy(true, p.Anything[byte])

		comb := Chainr(
			func(buffer p.Buffer[byte, int]) (string, error) {
				x, err := next(buffer)
				if err != nil {
					return "", err
				}

				return string(x), nil
			},
			func(buffer p.Buffer[byte, int]) (func(string, string) string, error) {
				return func(x, y string) string { return "" }, fmt.Errorf("test error")
			},
			"default",
		)

		result, err := Parse([]byte("abcd"), comb)
		Check(t, err)
		AssertEq(t, result, "a")

		result, err = Parse([]byte("a"), comb)
		Check(t, err)
		AssertEq(t, result, "a")

		result, err = Parse([]byte(""), comb)
		Check(t, err)
		AssertEq(t, result, "default")
	})
}

func TestChainr1(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		next := Satisfy(true, p.Anything[byte])

		comb := Chainr1(
			func(buffer p.Buffer[byte, int]) (string, error) {
				x, err := next(buffer)
				if err != nil {
					return "", err
				}

				return string(x), nil
			},
			func(buffer p.Buffer[byte, int]) (func(string, string) string, error) {
				return func(x, y string) string {
					return fmt.Sprintf("(%v %v)", x, y)
				}, nil
			},
		)

		result, err := Parse([]byte("abcd"), comb)
		Check(t, err)
		AssertEq(t, result, "(a (b (c d)))")

		result, err = Parse([]byte("a"), comb)
		Check(t, err)
		AssertEq(t, result, "a")

		result, err = Parse([]byte(""), comb)
		AssertError(t, err)
		AssertEq(t, result, "")
	})

	t.Run("case 2", func(t *testing.T) {
		next := Satisfy(true, p.Anything[byte])
		c := 0

		comb := Chainr1(
			func(buffer p.Buffer[byte, int]) (string, error) {
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
			func(buffer p.Buffer[byte, int]) (func(string, string) string, error) {
				return func(x, y string) string {
					return fmt.Sprintf("(%v %v)", x, y)
				}, nil
			},
		)

		result, err := Parse([]byte("abcd"), comb)
		Check(t, err)
		AssertEq(t, result, "a")

		c = 0

		result, err = Parse([]byte("a"), comb)
		Check(t, err)
		AssertEq(t, result, "a")

		c = 0

		result, err = Parse([]byte(""), comb)
		AssertError(t, err)
		AssertEq(t, result, "")
	})

	t.Run("case 3", func(t *testing.T) {
		next := Satisfy(true, p.Anything[byte])

		comb := Chainr1(
			func(buffer p.Buffer[byte, int]) (string, error) {
				x, err := next(buffer)
				if err != nil {
					return "", err
				}

				return string(x), nil
			},
			func(buffer p.Buffer[byte, int]) (func(string, string) string, error) {
				return func(x, y string) string {
					return ""
				}, fmt.Errorf("test error")
			},
		)

		result, err := Parse([]byte("abcd"), comb)
		Check(t, err)
		AssertEq(t, result, "a")

		result, err = Parse([]byte("a"), comb)
		Check(t, err)
		AssertEq(t, result, "a")

		result, err = Parse([]byte(""), comb)
		AssertError(t, err)
		AssertEq(t, result, "")
	})
}

func TestSepBy(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		comb := SepBy(0, NotEq(','), Eq(','))

		result, err := Parse([]byte("a,b,c"), comb)
		Check(t, err)
		AssertSlice(t, result, []byte{'a', 'b', 'c'})

		result, err = Parse([]byte(""), comb)
		Check(t, err)
		AssertSlice(t, result, []byte{})

		result, err = Parse([]byte(","), comb)
		Check(t, err)
		AssertSlice(t, result, []byte{})

		result, err = Parse([]byte(",a,b,c"), comb)
		Check(t, err)
		AssertSlice(t, result, []byte{})

		result, err = Parse([]byte("a,b,c,"), comb)
		Check(t, err)
		AssertSlice(t, result, []byte{'a', 'b', 'c'})

		result, err = Parse([]byte("a,b,"), comb)
		Check(t, err)
		AssertSlice(t, result, []byte{'a', 'b'})

		result, err = Parse([]byte("abc"), comb)
		Check(t, err)
		AssertSlice(t, result, []byte{'a'})
	})

	t.Run("case 2", func(t *testing.T) {
		comb := SepBy(
			0,
			Satisfy(true, p.Nothing[byte]),
			Eq(','),
		)

		result, err := Parse([]byte("a,b,c"), comb)
		Check(t, err)
		AssertSlice(t, result, nil)

		result, err = Parse([]byte(""), comb)
		Check(t, err)
		AssertSlice(t, result, nil)
	})

	t.Run("case 3", func(t *testing.T) {
		comb := SepBy(
			0,
			NotEq(','),
			Satisfy(true, p.Nothing[byte]),
		)

		result, err := Parse([]byte("a,b,c"), comb)
		Check(t, err)
		AssertSlice(t, result, []byte{'a'})

		result, err = Parse([]byte(""), comb)
		Check(t, err)
		AssertSlice(t, result, nil)
	})
}

func TestSepBy1(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		comb := SepBy1(0, NotEq(','), Eq(','))

		result, err := Parse([]byte("a,b,c"), comb)
		Check(t, err)
		AssertSlice(t, result, []byte{'a', 'b', 'c'})

		result, err = Parse([]byte(""), comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)

		result, err = Parse([]byte(","), comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)

		result, err = Parse([]byte(",a,b,c"), comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)

		result, err = Parse([]byte("a,b,c,"), comb)
		Check(t, err)
		AssertSlice(t, result, []byte{'a', 'b', 'c'})

		result, err = Parse([]byte("a,b,"), comb)
		Check(t, err)
		AssertSlice(t, result, []byte{'a', 'b'})

		result, err = Parse([]byte("abc"), comb)
		Check(t, err)
		AssertSlice(t, result, []byte{'a'})
	})

	t.Run("case 2", func(t *testing.T) {
		comb := SepBy1(
			0,
			Satisfy(true, p.Nothing[byte]),
			Eq(','),
		)

		result, err := Parse([]byte("a,b,c"), comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)

		result, err = Parse([]byte(""), comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})

	t.Run("case 3", func(t *testing.T) {
		comb := SepBy1(
			0,
			NotEq(','),
			Satisfy(true, p.Nothing[byte]),
		)

		result, err := Parse([]byte("a,b,c"), comb)
		Check(t, err)
		AssertSlice(t, result, []byte{'a'})

		result, err = Parse([]byte(""), comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})
}

func TestEndBy(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		comb := EndBy(0, NotEq(','), Eq(','))

		result, err := Parse([]byte("a,b,c"), comb)
		Check(t, err)
		AssertSlice(t, result, []byte{'a', 'b'})

		result, err = Parse([]byte("a,b,c,"), comb)
		Check(t, err)
		AssertSlice(t, result, []byte{'a', 'b', 'c'})

		result, err = Parse([]byte(""), comb)
		Check(t, err)
		AssertSlice(t, result, nil)

		result, err = Parse([]byte(","), comb)
		Check(t, err)
		AssertSlice(t, result, nil)

		result, err = Parse([]byte("a"), comb)
		Check(t, err)
		AssertSlice(t, result, nil)

		result, err = Parse([]byte("a,"), comb)
		Check(t, err)
		AssertSlice(t, result, []byte{'a'})

		result, err = Parse([]byte(",a"), comb)
		Check(t, err)
		AssertSlice(t, result, nil)

		result, err = Parse([]byte("a,,"), comb)
		Check(t, err)
		AssertSlice(t, result, []byte{'a'})

		result, err = Parse([]byte(",a,"), comb)
		Check(t, err)
		AssertSlice(t, result, nil)
	})

	t.Run("case 2", func(t *testing.T) {
		comb := EndBy(
			0,
			Satisfy(true, p.Nothing[byte]),
			Eq(','),
		)

		result, err := Parse([]byte("a,b,c"), comb)
		Check(t, err)
		AssertSlice(t, result, nil)

		result, err = Parse([]byte(""), comb)
		Check(t, err)
		AssertSlice(t, result, nil)
	})

	t.Run("case 3", func(t *testing.T) {
		comb := EndBy(
			0,
			NotEq(','),
			Satisfy(true, p.Nothing[byte]),
		)

		result, err := Parse([]byte("a,b,c"), comb)
		Check(t, err)
		AssertSlice(t, result, nil)

		result, err = Parse([]byte(""), comb)
		Check(t, err)
		AssertSlice(t, result, nil)
	})
}

func TestEndBy1(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		comb := EndBy1(0, NotEq(','), Eq(','))

		result, err := Parse([]byte("a,b,c"), comb)
		Check(t, err)
		AssertSlice(t, result, []byte{'a', 'b'})

		result, err = Parse([]byte("a,b,c,"), comb)
		Check(t, err)
		AssertSlice(t, result, []byte{'a', 'b', 'c'})

		result, err = Parse([]byte(""), comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)

		result, err = Parse([]byte(","), comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)

		result, err = Parse([]byte("a"), comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)

		result, err = Parse([]byte("a,"), comb)
		Check(t, err)
		AssertSlice(t, result, []byte{'a'})

		result, err = Parse([]byte(",a"), comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)

		result, err = Parse([]byte("a,,"), comb)
		Check(t, err)
		AssertSlice(t, result, []byte{'a'})

		result, err = Parse([]byte(",a,"), comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})

	t.Run("case 2", func(t *testing.T) {
		comb := EndBy1(
			0,
			Satisfy(true, p.Nothing[byte]),
			Eq(','),
		)

		result, err := Parse([]byte("a,b,c"), comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)

		result, err = Parse([]byte(""), comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})

	t.Run("case 3", func(t *testing.T) {
		comb := EndBy1(
			0,
			NotEq(','),
			Satisfy(true, p.Nothing[byte]),
		)

		result, err := Parse([]byte("a,b,c"), comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)

		result, err = Parse([]byte(""), comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})
}

func TestSepEndBy(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		comb := SepEndBy(0, NotEq(','), Eq(','))

		result, err := Parse([]byte("a,b,c"), comb)
		Check(t, err)
		AssertSlice(t, result, []byte{'a', 'b', 'c'})

		result, err = Parse([]byte("a,b,c,"), comb)
		Check(t, err)
		AssertSlice(t, result, []byte{'a', 'b', 'c'})

		result, err = Parse([]byte("a,b,c,,"), comb)
		Check(t, err)
		AssertSlice(t, result, []byte{'a', 'b', 'c'})

		result, err = Parse([]byte(""), comb)
		Check(t, err)
		AssertSlice(t, result, nil)

		result, err = Parse([]byte(","), comb)
		Check(t, err)
		AssertSlice(t, result, nil)

		result, err = Parse([]byte(",a,b,c"), comb)
		Check(t, err)
		AssertSlice(t, result, nil)
	})

	t.Run("case 2", func(t *testing.T) {
		comb := SepEndBy(
			0,
			Satisfy(true, p.Nothing[byte]),
			Eq(','),
		)

		result, err := Parse([]byte("a,b,c"), comb)
		Check(t, err)
		AssertSlice(t, result, nil)

		result, err = Parse([]byte(","), comb)
		Check(t, err)
		AssertSlice(t, result, nil)

		result, err = Parse([]byte(""), comb)
		Check(t, err)
		AssertSlice(t, result, nil)
	})

	t.Run("case 3", func(t *testing.T) {
		comb := SepEndBy(
			0,
			NotEq(','),
			Satisfy(true, p.Nothing[byte]),
		)

		result, err := Parse([]byte("a,b,c"), comb)
		Check(t, err)
		AssertSlice(t, result, []byte{'a'})

		result, err = Parse([]byte(","), comb)
		Check(t, err)
		AssertSlice(t, result, nil)

		result, err = Parse([]byte(""), comb)
		Check(t, err)
		AssertSlice(t, result, nil)
	})
}

func TestSepEndBy1(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		comb := SepEndBy1(0, NotEq(','), Eq(','))

		result, err := Parse([]byte("a,b,c"), comb)
		Check(t, err)
		AssertSlice(t, result, []byte{'a', 'b', 'c'})

		result, err = Parse([]byte("a,b,c,"), comb)
		Check(t, err)
		AssertSlice(t, result, []byte{'a', 'b', 'c'})

		result, err = Parse([]byte("a,b,c,,"), comb)
		Check(t, err)
		AssertSlice(t, result, []byte{'a', 'b', 'c'})

		result, err = Parse([]byte(""), comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)

		result, err = Parse([]byte(","), comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)

		result, err = Parse([]byte(",a,b,c"), comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})

	t.Run("case 2", func(t *testing.T) {
		comb := SepEndBy1(
			0,
			Satisfy(true, p.Nothing[byte]),
			Eq(','),
		)

		result, err := Parse([]byte("a,b,c"), comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)

		result, err = Parse([]byte(""), comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})

	t.Run("case 3", func(t *testing.T) {
		comb := SepEndBy1(
			0,
			NotEq(','),
			Satisfy(true, p.Nothing[byte]),
		)

		result, err := Parse([]byte("a,b,c"), comb)
		Check(t, err)
		AssertSlice(t, result, []byte{'a'})

		result, err = Parse([]byte(""), comb)
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
			Satisfy(false, func(x byte) bool { return x == byte('d') }),
		)

		result, err := Parse([]byte("abcd"), comb)
		Check(t, err)
		AssertSlice(t, result, []byte{
			byte('a'),
			byte('b'),
			byte('c'),
		})

		result, err = Parse([]byte(""), comb)
		Check(t, err)
		AssertSlice(t, result, nil)
	})

	t.Run("case 2", func(t *testing.T) {
		comb := ManyTill(
			0,
			Any(),
			Satisfy(true, p.Nothing[byte]),
		)

		result, err := Parse([]byte("abcd"), comb)
		Check(t, err)
		AssertSlice(t, result, []byte{
			byte('a'),
			byte('b'),
			byte('c'),
			byte('d'),
		})

		result, err = Parse([]byte(""), comb)
		Check(t, err)
		AssertSlice(t, result, nil)
	})

	t.Run("case 3", func(t *testing.T) {
		comb := ManyTill(
			0,
			Any(),
			Any(),
		)

		result, err := Parse([]byte("abcd"), comb)
		Check(t, err)
		AssertSlice(t, result, nil)

		result, err = Parse([]byte(""), comb)
		Check(t, err)
		AssertSlice(t, result, nil)
	})
}
