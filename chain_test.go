package parsec

import (
	"fmt"
	"testing"
)

func TestChainl(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		next := Satisfy[byte, int](true, Anything[byte])

		comb := Chainl(
			func(buffer Buffer[byte,int]) (string, error) {
				x, err := next(buffer)
				if err != nil {
					return "", err
				}

				return string(x), nil
			},
			func(buffer Buffer[byte,int]) (func(string, string) string, error) {
				return func(x, y string) string {
					return fmt.Sprintf("(%v %v)", x, y)
				}, nil
			},
			"default",
		)

		result, err := ParseBytes([]byte("abcd"), comb)
		check(t, err)
		assertEq(t, result, "(((a b) c) d)")

		result, err = ParseBytes([]byte("a"), comb)
		check(t, err)
		assertEq(t, result, "a")

		result, err = ParseBytes([]byte(""), comb)
		check(t, err)
		assertEq(t, result, "default")
	})

	t.Run("case 2", func(t *testing.T) {
		next := Satisfy[byte, int](true, Anything[byte])
		c := 0

		comb := Chainl(
			func(buffer Buffer[byte,int]) (string, error) {
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
			func(buffer Buffer[byte,int]) (func(string, string) string, error) {
				return func(x, y string) string {
					return fmt.Sprintf("(%v %v)", x, y)
				}, nil
			},
			"default",
		)

		result, err := ParseBytes([]byte("abcd"), comb)
		check(t, err)
		assertEq(t, result, "a")

		c = 0

		result, err = ParseBytes([]byte("a"), comb)
		check(t, err)
		assertEq(t, result, "a")

		c = 0

		result, err = ParseBytes([]byte(""), comb)
		check(t, err)
		assertEq(t, result, "default")
	})

	t.Run("case 3", func(t *testing.T) {
		next := Satisfy[byte, int](true, Anything[byte])

		comb := Chainl(
			func(buffer Buffer[byte,int]) (string, error) {
				x, err := next(buffer)
				if err != nil {
					return "", err
				}

				return string(x), nil
			},
			func(buffer Buffer[byte,int]) (func(string, string) string, error) {
				return func(x, y string) string { return "" }, fmt.Errorf("test error")
			},
			"default",
		)

		result, err := ParseBytes([]byte("abcd"), comb)
		check(t, err)
		assertEq(t, result, "a")

		result, err = ParseBytes([]byte("a"), comb)
		check(t, err)
		assertEq(t, result, "a")

		result, err = ParseBytes([]byte(""), comb)
		check(t, err)
		assertEq(t, result, "default")
	})
}

func TestChainl1(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		next := Satisfy[byte, int](true, Anything[byte])

		comb := Chainl1(
			func(buffer Buffer[byte,int]) (string, error) {
				x, err := next(buffer)
				if err != nil {
					return "", err
				}

				return string(x), nil
			},
			func(buffer Buffer[byte,int]) (func(string, string) string, error) {
				return func(x, y string) string {
					return fmt.Sprintf("(%v %v)", x, y)
				}, nil
			},
		)

		result, err := ParseBytes([]byte("abcd"), comb)
		check(t, err)
		assertEq(t, result, "(((a b) c) d)")

		result, err = ParseBytes([]byte("a"), comb)
		check(t, err)
		assertEq(t, result, "a")

		result, err = ParseBytes([]byte(""), comb)
		assertError(t, err)
		assertEq(t, result, "")
	})

	t.Run("case 2", func(t *testing.T) {
		next := Satisfy[byte, int](true, Anything[byte])
		c := 0

		comb := Chainl1(
			func(buffer Buffer[byte,int]) (string, error) {
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
			func(buffer Buffer[byte,int]) (func(string, string) string, error) {
				return func(x, y string) string {
					return fmt.Sprintf("(%v %v)", x, y)
				}, nil
			},
		)

		result, err := ParseBytes([]byte("abcd"), comb)
		check(t, err)
		assertEq(t, result, "a")

		c = 0

		result, err = ParseBytes([]byte("a"), comb)
		check(t, err)
		assertEq(t, result, "a")

		c = 0

		result, err = ParseBytes([]byte(""), comb)
		assertError(t, err)
		assertEq(t, result, "")
	})

	t.Run("case 3", func(t *testing.T) {
		next := Satisfy[byte, int](true, Anything[byte])

		comb := Chainl1(
			func(buffer Buffer[byte,int]) (string, error) {
				x, err := next(buffer)
				if err != nil {
					return "--", err
				}

				return string(x), nil
			},
			func(buffer Buffer[byte,int]) (func(string, string) string, error) {
				return func(x, y string) string { return "++" }, fmt.Errorf("test error")
			},
		)

		result, err := ParseBytes([]byte("abcd"), comb)
		check(t, err)
		assertEq(t, result, "a")

		result, err = ParseBytes([]byte("a"), comb)
		check(t, err)
		assertEq(t, result, "a")

		result, err = ParseBytes([]byte(""), comb)
		assertError(t, err)
		assertEq(t, result, "")
	})
}

func TestChainr(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		next := Satisfy[byte, int](true, Anything[byte])

		comb := Chainr(
			func(buffer Buffer[byte,int]) (string, error) {
				x, err := next(buffer)
				if err != nil {
					return "", err
				}

				return string(x), nil
			},
			func(buffer Buffer[byte,int]) (func(string, string) string, error) {
				return func(x, y string) string {
					return fmt.Sprintf("(%v %v)", x, y)
				}, nil
			},
			"default",
		)

		result, err := ParseBytes([]byte("abcd"), comb)
		check(t, err)
		assertEq(t, result, "(a (b (c d)))")

		result, err = ParseBytes([]byte("a"), comb)
		check(t, err)
		assertEq(t, result, "a")

		result, err = ParseBytes([]byte(""), comb)
		check(t, err)
		assertEq(t, result, "default")
	})

	t.Run("case 2", func(t *testing.T) {
		next := Satisfy[byte, int](true, Anything[byte])
		c := 0

		comb := Chainr(
			func(buffer Buffer[byte,int]) (string, error) {
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
			func(buffer Buffer[byte,int]) (func(string, string) string, error) {
				return func(x, y string) string {
					return fmt.Sprintf("(%v %v)", x, y)
				}, nil
			},
			"default",
		)

		result, err := ParseBytes([]byte("abcd"), comb)
		check(t, err)
		assertEq(t, result, "a")

		c = 0

		result, err = ParseBytes([]byte("a"), comb)
		check(t, err)
		assertEq(t, result, "a")

		c = 0

		result, err = ParseBytes([]byte(""), comb)
		check(t, err)
		assertEq(t, result, "default")
	})

	t.Run("case 3", func(t *testing.T) {
		next := Satisfy[byte, int](true, Anything[byte])

		comb := Chainr(
			func(buffer Buffer[byte,int]) (string, error) {
				x, err := next(buffer)
				if err != nil {
					return "", err
				}

				return string(x), nil
			},
			func(buffer Buffer[byte,int]) (func(string, string) string, error) {
				return func(x, y string) string { return "" }, fmt.Errorf("test error")
			},
			"default",
		)

		result, err := ParseBytes([]byte("abcd"), comb)
		check(t, err)
		assertEq(t, result, "a")

		result, err = ParseBytes([]byte("a"), comb)
		check(t, err)
		assertEq(t, result, "a")

		result, err = ParseBytes([]byte(""), comb)
		check(t, err)
		assertEq(t, result, "default")
	})
}

func TestChainr1(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		next := Satisfy[byte, int](true, Anything[byte])

		comb := Chainr1(
			func(buffer Buffer[byte,int]) (string, error) {
				x, err := next(buffer)
				if err != nil {
					return "", err
				}

				return string(x), nil
			},
			func(buffer Buffer[byte,int]) (func(string, string) string, error) {
				return func(x, y string) string {
					return fmt.Sprintf("(%v %v)", x, y)
				}, nil
			},
		)

		result, err := ParseBytes([]byte("abcd"), comb)
		check(t, err)
		assertEq(t, result, "(a (b (c d)))")

		result, err = ParseBytes([]byte("a"), comb)
		check(t, err)
		assertEq(t, result, "a")

		result, err = ParseBytes([]byte(""), comb)
		assertError(t, err)
		assertEq(t, result, "")
	})

	t.Run("case 2", func(t *testing.T) {
		next := Satisfy[byte, int](true, Anything[byte])
		c := 0

		comb := Chainr1(
			func(buffer Buffer[byte,int]) (string, error) {
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
			func(buffer Buffer[byte,int]) (func(string, string) string, error) {
				return func(x, y string) string {
					return fmt.Sprintf("(%v %v)", x, y)
				}, nil
			},
		)

		result, err := ParseBytes([]byte("abcd"), comb)
		check(t, err)
		assertEq(t, result, "a")

		c = 0

		result, err = ParseBytes([]byte("a"), comb)
		check(t, err)
		assertEq(t, result, "a")

		c = 0

		result, err = ParseBytes([]byte(""), comb)
		assertError(t, err)
		assertEq(t, result, "")
	})

	t.Run("case 3", func(t *testing.T) {
		next := Satisfy[byte, int](true, Anything[byte])

		comb := Chainr1(
			func(buffer Buffer[byte,int]) (string, error) {
				x, err := next(buffer)
				if err != nil {
					return "", err
				}

				return string(x), nil
			},
			func(buffer Buffer[byte,int]) (func(string, string) string, error) {
				return func(x, y string) string {
					return ""
				}, fmt.Errorf("test error")
			},
		)

		result, err := ParseBytes([]byte("abcd"), comb)
		check(t, err)
		assertEq(t, result, "a")

		result, err = ParseBytes([]byte("a"), comb)
		check(t, err)
		assertEq(t, result, "a")

		result, err = ParseBytes([]byte(""), comb)
		assertError(t, err)
		assertEq(t, result, "")
	})
}

func TestSepBy(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		comb := SepBy(
			0,
			NotEq[byte, int](','),
			Eq[byte, int](','),
		)

		result, err := ParseBytes([]byte("a,b,c"), comb)
		check(t, err)
		assertSlice(t, result, []byte{'a', 'b', 'c'})

		result, err = ParseBytes([]byte(""), comb)
		check(t, err)
		assertSlice(t, result, []byte{})

		result, err = ParseBytes([]byte(","), comb)
		check(t, err)
		assertSlice(t, result, []byte{})

		result, err = ParseBytes([]byte(",a,b,c"), comb)
		check(t, err)
		assertSlice(t, result, []byte{})

		result, err = ParseBytes([]byte("a,b,c,"), comb)
		check(t, err)
		assertSlice(t, result, []byte{'a', 'b', 'c'})

		result, err = ParseBytes([]byte("a,b,"), comb)
		check(t, err)
		assertSlice(t, result, []byte{'a', 'b'})

		result, err = ParseBytes([]byte("abc"), comb)
		check(t, err)
		assertSlice(t, result, []byte{'a'})
	})

	t.Run("case 2", func(t *testing.T) {
		comb := SepBy(
			0,
			Satisfy[byte, int](true, Nothing[byte]),
			Eq[byte, int](','),
		)

		result, err := ParseBytes([]byte("a,b,c"), comb)
		check(t, err)
		assertSlice(t, result, nil)

		result, err = ParseBytes([]byte(""), comb)
		check(t, err)
		assertSlice(t, result, nil)
	})

	t.Run("case 3", func(t *testing.T) {
		comb := SepBy(
			0,
			NotEq[byte,int](','),
			Satisfy[byte, int](true, Nothing[byte]),
		)

		result, err := ParseBytes([]byte("a,b,c"), comb)
		check(t, err)
		assertSlice(t, result, []byte{'a'})

		result, err = ParseBytes([]byte(""), comb)
		check(t, err)
		assertSlice(t, result, nil)
	})
}

func TestSepBy1(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		comb := SepBy1(
			0,
			NotEq[byte,int](','),
			Eq[byte,int](','),
		)

		result, err := ParseBytes([]byte("a,b,c"), comb)
		check(t, err)
		assertSlice(t, result, []byte{'a', 'b', 'c'})

		result, err = ParseBytes([]byte(""), comb)
		assertError(t, err)
		assertSlice(t, result, nil)

		result, err = ParseBytes([]byte(","), comb)
		assertError(t, err)
		assertSlice(t, result, nil)

		result, err = ParseBytes([]byte(",a,b,c"), comb)
		assertError(t, err)
		assertSlice(t, result, nil)

		result, err = ParseBytes([]byte("a,b,c,"), comb)
		check(t, err)
		assertSlice(t, result, []byte{'a', 'b', 'c'})

		result, err = ParseBytes([]byte("a,b,"), comb)
		check(t, err)
		assertSlice(t, result, []byte{'a', 'b'})

		result, err = ParseBytes([]byte("abc"), comb)
		check(t, err)
		assertSlice(t, result, []byte{'a'})
	})

	t.Run("case 2", func(t *testing.T) {
		comb := SepBy1(
			0,
			Satisfy[byte,int](true, Nothing[byte]),
			Eq[byte,int](','),
		)

		result, err := ParseBytes([]byte("a,b,c"), comb)
		assertError(t, err)
		assertSlice(t, result, nil)

		result, err = ParseBytes([]byte(""), comb)
		assertError(t, err)
		assertSlice(t, result, nil)
	})

	t.Run("case 3", func(t *testing.T) {
		comb := SepBy1(
			0,
			NotEq[byte,int](','),
			Satisfy[byte,int](true, Nothing[byte]),
		)

		result, err := ParseBytes([]byte("a,b,c"), comb)
		check(t, err)
		assertSlice(t, result, []byte{'a'})

		result, err = ParseBytes([]byte(""), comb)
		assertError(t, err)
		assertSlice(t, result, nil)
	})
}

func TestEndBy(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		comb := EndBy(
			0,
			NotEq[byte,int](','),
			Eq[byte,int](','),
		)

		result, err := ParseBytes([]byte("a,b,c"), comb)
		check(t, err)
		assertSlice(t, result, []byte{'a', 'b'})

		result, err = ParseBytes([]byte("a,b,c,"), comb)
		check(t, err)
		assertSlice(t, result, []byte{'a', 'b', 'c'})

		result, err = ParseBytes([]byte(""), comb)
		check(t, err)
		assertSlice(t, result, nil)

		result, err = ParseBytes([]byte(","), comb)
		check(t, err)
		assertSlice(t, result, nil)

		result, err = ParseBytes([]byte("a"), comb)
		check(t, err)
		assertSlice(t, result, nil)

		result, err = ParseBytes([]byte("a,"), comb)
		check(t, err)
		assertSlice(t, result, []byte{'a'})

		result, err = ParseBytes([]byte(",a"), comb)
		check(t, err)
		assertSlice(t, result, nil)

		result, err = ParseBytes([]byte("a,,"), comb)
		check(t, err)
		assertSlice(t, result, []byte{'a'})

		result, err = ParseBytes([]byte(",a,"), comb)
		check(t, err)
		assertSlice(t, result, nil)
	})

	t.Run("case 2", func(t *testing.T) {
		comb := EndBy(
			0,
			Satisfy[byte,int](true, Nothing[byte]),
			Eq[byte,int](','),
		)

		result, err := ParseBytes([]byte("a,b,c"), comb)
		check(t, err)
		assertSlice(t, result, nil)

		result, err = ParseBytes([]byte(""), comb)
		check(t, err)
		assertSlice(t, result, nil)
	})

	t.Run("case 3", func(t *testing.T) {
		comb := EndBy(
			0,
			NotEq[byte,int](','),
			Satisfy[byte,int](true, Nothing[byte]),
		)

		result, err := ParseBytes([]byte("a,b,c"), comb)
		check(t, err)
		assertSlice(t, result, nil)

		result, err = ParseBytes([]byte(""), comb)
		check(t, err)
		assertSlice(t, result, nil)
	})
}

func TestEndBy1(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		comb := EndBy1(
			0,
			NotEq[byte,int](','),
			Eq[byte,int](','),
		)

		result, err := ParseBytes([]byte("a,b,c"), comb)
		check(t, err)
		assertSlice(t, result, []byte{'a', 'b'})

		result, err = ParseBytes([]byte("a,b,c,"), comb)
		check(t, err)
		assertSlice(t, result, []byte{'a', 'b', 'c'})

		result, err = ParseBytes([]byte(""), comb)
		assertError(t, err)
		assertSlice(t, result, nil)

		result, err = ParseBytes([]byte(","), comb)
		assertError(t, err)
		assertSlice(t, result, nil)

		result, err = ParseBytes([]byte("a"), comb)
		assertError(t, err)
		assertSlice(t, result, nil)

		result, err = ParseBytes([]byte("a,"), comb)
		check(t, err)
		assertSlice(t, result, []byte{'a'})

		result, err = ParseBytes([]byte(",a"), comb)
		assertError(t, err)
		assertSlice(t, result, nil)

		result, err = ParseBytes([]byte("a,,"), comb)
		check(t, err)
		assertSlice(t, result, []byte{'a'})

		result, err = ParseBytes([]byte(",a,"), comb)
		assertError(t, err)
		assertSlice(t, result, nil)
	})

	t.Run("case 2", func(t *testing.T) {
		comb := EndBy1(
			0,
			Satisfy[byte,int](true, Nothing[byte]),
			Eq[byte,int](','),
		)

		result, err := ParseBytes([]byte("a,b,c"), comb)
		assertError(t, err)
		assertSlice(t, result, nil)

		result, err = ParseBytes([]byte(""), comb)
		assertError(t, err)
		assertSlice(t, result, nil)
	})

	t.Run("case 3", func(t *testing.T) {
		comb := EndBy1(
			0,
			NotEq[byte,int](','),
			Satisfy[byte,int](true, Nothing[byte]),
		)

		result, err := ParseBytes([]byte("a,b,c"), comb)
		assertError(t, err)
		assertSlice(t, result, nil)

		result, err = ParseBytes([]byte(""), comb)
		assertError(t, err)
		assertSlice(t, result, nil)
	})
}

func TestSepEndBy(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		comb := SepEndBy(
			0,
			NotEq[byte,int](','),
			Eq[byte,int](','),
		)

		result, err := ParseBytes([]byte("a,b,c"), comb)
		check(t, err)
		assertSlice(t, result, []byte{'a', 'b', 'c'})

		result, err = ParseBytes([]byte("a,b,c,"), comb)
		check(t, err)
		assertSlice(t, result, []byte{'a', 'b', 'c'})

		result, err = ParseBytes([]byte("a,b,c,,"), comb)
		check(t, err)
		assertSlice(t, result, []byte{'a', 'b', 'c'})

		result, err = ParseBytes([]byte(""), comb)
		check(t, err)
		assertSlice(t, result, nil)

		result, err = ParseBytes([]byte(","), comb)
		check(t, err)
		assertSlice(t, result, nil)

		result, err = ParseBytes([]byte(",a,b,c"), comb)
		check(t, err)
		assertSlice(t, result, nil)
	})

	t.Run("case 2", func(t *testing.T) {
		comb := SepEndBy(
			0,
			Satisfy[byte,int](true, Nothing[byte]),
			Eq[byte,int](','),
		)

		result, err := ParseBytes([]byte("a,b,c"), comb)
		check(t, err)
		assertSlice(t, result, nil)

		result, err = ParseBytes([]byte(","), comb)
		check(t, err)
		assertSlice(t, result, nil)

		result, err = ParseBytes([]byte(""), comb)
		check(t, err)
		assertSlice(t, result, nil)
	})

	t.Run("case 3", func(t *testing.T) {
		comb := SepEndBy(
			0,
			NotEq[byte,int](','),
			Satisfy[byte,int](true, Nothing[byte]),
		)

		result, err := ParseBytes([]byte("a,b,c"), comb)
		check(t, err)
		assertSlice(t, result, []byte{'a'})

		result, err = ParseBytes([]byte(","), comb)
		check(t, err)
		assertSlice(t, result, nil)

		result, err = ParseBytes([]byte(""), comb)
		check(t, err)
		assertSlice(t, result, nil)
	})
}

func TestSepEndBy1(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		comb := SepEndBy1(
			0,
			NotEq[byte,int](','),
			Eq[byte,int](','),
		)

		result, err := ParseBytes([]byte("a,b,c"), comb)
		check(t, err)
		assertSlice(t, result, []byte{'a', 'b', 'c'})

		result, err = ParseBytes([]byte("a,b,c,"), comb)
		check(t, err)
		assertSlice(t, result, []byte{'a', 'b', 'c'})

		result, err = ParseBytes([]byte("a,b,c,,"), comb)
		check(t, err)
		assertSlice(t, result, []byte{'a', 'b', 'c'})

		result, err = ParseBytes([]byte(""), comb)
		assertError(t, err)
		assertSlice(t, result, nil)

		result, err = ParseBytes([]byte(","), comb)
		assertError(t, err)
		assertSlice(t, result, nil)

		result, err = ParseBytes([]byte(",a,b,c"), comb)
		assertError(t, err)
		assertSlice(t, result, nil)
	})

	t.Run("case 2", func(t *testing.T) {
		comb := SepEndBy1(
			0,
			Satisfy[byte,int](true, Nothing[byte]),
			Eq[byte,int](','),
		)

		result, err := ParseBytes([]byte("a,b,c"), comb)
		assertError(t, err)
		assertSlice(t, result, nil)

		result, err = ParseBytes([]byte(""), comb)
		assertError(t, err)
		assertSlice(t, result, nil)
	})

	t.Run("case 3", func(t *testing.T) {
		comb := SepEndBy1(
			0,
			NotEq[byte,int](','),
			Satisfy[byte,int](true, Nothing[byte]),
		)

		result, err := ParseBytes([]byte("a,b,c"), comb)
		check(t, err)
		assertSlice(t, result, []byte{'a'})

		result, err = ParseBytes([]byte(""), comb)
		assertError(t, err)
		assertSlice(t, result, nil)
	})
}

func TestManyTill(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		comb := ManyTill(
			0,
			Any[byte,int](),
			Satisfy[byte,int](false, func(x byte) bool { return x == byte('d') }),
		)

		result, err := ParseBytes([]byte("abcd"), comb)
		check(t, err)
		assertSlice(t, result, []byte{
			byte('a'),
			byte('b'),
			byte('c'),
		})

		result, err = ParseBytes([]byte(""), comb)
		check(t, err)
		assertSlice(t, result, nil)
	})

	t.Run("case 2", func(t *testing.T) {
		comb := ManyTill(
			0,
			Any[byte,int](),
			Satisfy[byte,int](true, Nothing[byte]),
		)

		result, err := ParseBytes([]byte("abcd"), comb)
		check(t, err)
		assertSlice(t, result, []byte{
			byte('a'),
			byte('b'),
			byte('c'),
			byte('d'),
		})

		result, err = ParseBytes([]byte(""), comb)
		check(t, err)
		assertSlice(t, result, nil)
	})

	t.Run("case 3", func(t *testing.T) {
		comb := ManyTill(
			0,
			Any[byte,int](),
			Any[byte,int](),
		)

		result, err := ParseBytes([]byte("abcd"), comb)
		check(t, err)
		assertSlice(t, result, nil)

		result, err = ParseBytes([]byte(""), comb)
		check(t, err)
		assertSlice(t, result, nil)
	})
}
