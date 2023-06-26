package parsec

import (
	"fmt"
	"testing"
)

func TestChainl(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		next := Satisfy(true, Anything[byte])

		comb := Chainl(
			func(buffer Buffer[byte]) (string, bool) {
				x, ok := next(buffer)
				if !ok {
					return "", false
				}

				return string(x), true
			},
			func(buffer Buffer[byte]) (func(string, string) string, bool) {
				return func(x, y string) string {
					return fmt.Sprintf("(%v %v)", x, y)
				}, true
			},
			"default",
		)

		result, ok := ParseBytes([]byte("abcd"), comb)
		assert(t, ok, "expected true")
		assertEq(t, result, "(((a b) c) d)")

		result, ok = ParseBytes([]byte("a"), comb)
		assert(t, ok, "expected true")
		assertEq(t, result, "a")

		result, ok = ParseBytes([]byte(""), comb)
		assert(t, ok, "expected true")
		assertEq(t, result, "default")
	})

	t.Run("case 2", func(t *testing.T) {
		next := Satisfy(true, Anything[byte])
		c := 0

		comb := Chainl(
			func(buffer Buffer[byte]) (string, bool) {
				c++
				if c > 1 {
					return "-", false
				}

				x, ok := next(buffer)
				if !ok {
					return "", false
				}

				return string(x), true
			},
			func(buffer Buffer[byte]) (func(string, string) string, bool) {
				return func(x, y string) string {
					return fmt.Sprintf("(%v %v)", x, y)
				}, true
			},
			"default",
		)

		result, ok := ParseBytes([]byte("abcd"), comb)
		assert(t, ok, "expected true")
		assertEq(t, result, "a")

		c = 0

		result, ok = ParseBytes([]byte("a"), comb)
		assert(t, ok, "expected true")
		assertEq(t, result, "a")

		c = 0

		result, ok = ParseBytes([]byte(""), comb)
		assert(t, ok, "expected true")
		assertEq(t, result, "default")
	})

	t.Run("case 3", func(t *testing.T) {
		next := Satisfy(true, Anything[byte])

		comb := Chainl(
			func(buffer Buffer[byte]) (string, bool) {
				x, ok := next(buffer)
				if !ok {
					return "", false
				}

				return string(x), true
			},
			func(buffer Buffer[byte]) (func(string, string) string, bool) {
				return func(x, y string) string { return "" }, false
			},
			"default",
		)

		result, ok := ParseBytes([]byte("abcd"), comb)
		assert(t, ok, "expected true")
		assertEq(t, result, "a")

		result, ok = ParseBytes([]byte("a"), comb)
		assert(t, ok, "expected true")
		assertEq(t, result, "a")

		result, ok = ParseBytes([]byte(""), comb)
		assert(t, ok, "expected true")
		assertEq(t, result, "default")
	})
}

func TestChainl1(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		next := Satisfy(true, Anything[byte])

		comb := Chainl1(
			func(buffer Buffer[byte]) (string, bool) {
				x, ok := next(buffer)
				if !ok {
					return "", false
				}

				return string(x), true
			},
			func(buffer Buffer[byte]) (func(string, string) string, bool) {
				return func(x, y string) string {
					return fmt.Sprintf("(%v %v)", x, y)
				}, true
			},
		)

		result, ok := ParseBytes([]byte("abcd"), comb)
		assert(t, ok, "expected true")
		assertEq(t, result, "(((a b) c) d)")

		result, ok = ParseBytes([]byte("a"), comb)
		assert(t, ok, "expected true")
		assertEq(t, result, "a")

		result, ok = ParseBytes([]byte(""), comb)
		assert(t, !ok, "expected false")
		assertEq(t, result, "")
	})

	t.Run("case 2", func(t *testing.T) {
		next := Satisfy(true, Anything[byte])
		c := 0

		comb := Chainl1(
			func(buffer Buffer[byte]) (string, bool) {
				c++
				if c > 1 {
					return "-", false
				}

				x, ok := next(buffer)
				if !ok {
					return "", false
				}

				return string(x), true
			},
			func(buffer Buffer[byte]) (func(string, string) string, bool) {
				return func(x, y string) string {
					return fmt.Sprintf("(%v %v)", x, y)
				}, true
			},
		)

		result, ok := ParseBytes([]byte("abcd"), comb)
		assert(t, ok, "expected true")
		assertEq(t, result, "a")

		c = 0

		result, ok = ParseBytes([]byte("a"), comb)
		assert(t, ok, "expected true")
		assertEq(t, result, "a")

		c = 0

		result, ok = ParseBytes([]byte(""), comb)
		assert(t, !ok, "expected false")
		assertEq(t, result, "")
	})

	t.Run("case 3", func(t *testing.T) {
		next := Satisfy(true, Anything[byte])

		comb := Chainl1(
			func(buffer Buffer[byte]) (string, bool) {
				x, ok := next(buffer)
				if !ok {
					return "", false
				}

				return string(x), true
			},
			func(buffer Buffer[byte]) (func(string, string) string, bool) {
				return func(x, y string) string { return "" }, false
			},
		)

		result, ok := ParseBytes([]byte("abcd"), comb)
		assert(t, ok, "expected true")
		assertEq(t, result, "a")

		result, ok = ParseBytes([]byte("a"), comb)
		assert(t, ok, "expected true")
		assertEq(t, result, "a")

		result, ok = ParseBytes([]byte(""), comb)
		assert(t, !ok, "expected false")
		assertEq(t, result, "")
	})
}

func TestChainr(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		next := Satisfy(true, Anything[byte])

		comb := Chainr(
			func(buffer Buffer[byte]) (string, bool) {
				x, ok := next(buffer)
				if !ok {
					return "", false
				}

				return string(x), true
			},
			func(buffer Buffer[byte]) (func(string, string) string, bool) {
				return func(x, y string) string {
					return fmt.Sprintf("(%v %v)", x, y)
				}, true
			},
			"default",
		)

		result, ok := ParseBytes([]byte("abcd"), comb)
		assert(t, ok, "expected true")
		assertEq(t, result, "(a (b (c d)))")

		result, ok = ParseBytes([]byte("a"), comb)
		assert(t, ok, "expected true")
		assertEq(t, result, "a")

		result, ok = ParseBytes([]byte(""), comb)
		assert(t, ok, "expected true")
		assertEq(t, result, "default")
	})

	t.Run("case 2", func(t *testing.T) {
		next := Satisfy(true, Anything[byte])
		c := 0

		comb := Chainr(
			func(buffer Buffer[byte]) (string, bool) {
				c++

				if c > 1 {
					return "-", c <= 1
				}

				x, ok := next(buffer)
				if !ok {
					return "", false
				}

				return string(x), true
			},
			func(buffer Buffer[byte]) (func(string, string) string, bool) {
				return func(x, y string) string {
					return fmt.Sprintf("(%v %v)", x, y)
				}, true
			},
			"default",
		)

		result, ok := ParseBytes([]byte("abcd"), comb)
		assert(t, ok, "expected true")
		assertEq(t, result, "a")

		c = 0

		result, ok = ParseBytes([]byte("a"), comb)
		assert(t, ok, "expected true")
		assertEq(t, result, "a")

		c = 0

		result, ok = ParseBytes([]byte(""), comb)
		assert(t, ok, "expected true")
		assertEq(t, result, "default")
	})

	t.Run("case 3", func(t *testing.T) {
		next := Satisfy(true, Anything[byte])

		comb := Chainr(
			func(buffer Buffer[byte]) (string, bool) {
				x, ok := next(buffer)
				if !ok {
					return "", false
				}

				return string(x), true
			},
			func(buffer Buffer[byte]) (func(string, string) string, bool) {
				return func(x, y string) string { return "" }, false
			},
			"default",
		)

		result, ok := ParseBytes([]byte("abcd"), comb)
		assert(t, ok, "expected true")
		assertEq(t, result, "a")

		result, ok = ParseBytes([]byte("a"), comb)
		assert(t, ok, "expected true")
		assertEq(t, result, "a")

		result, ok = ParseBytes([]byte(""), comb)
		assert(t, ok, "expected true")
		assertEq(t, result, "default")
	})
}

func TestChainr1(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		next := Satisfy(true, Anything[byte])

		comb := Chainr1(
			func(buffer Buffer[byte]) (string, bool) {
				x, ok := next(buffer)
				if !ok {
					return "", false
				}

				return string(x), true
			},
			func(buffer Buffer[byte]) (func(string, string) string, bool) {
				return func(x, y string) string {
					return fmt.Sprintf("(%v %v)", x, y)
				}, true
			},
		)

		result, ok := ParseBytes([]byte("abcd"), comb)
		assert(t, ok, "expected true")
		assertEq(t, result, "(a (b (c d)))")

		result, ok = ParseBytes([]byte("a"), comb)
		assert(t, ok, "expected true")
		assertEq(t, result, "a")

		result, ok = ParseBytes([]byte(""), comb)
		assert(t, !ok, "expected false")
		assertEq(t, result, "")
	})

	t.Run("case 2", func(t *testing.T) {
		next := Satisfy(true, Anything[byte])
		c := 0

		comb := Chainr1(
			func(buffer Buffer[byte]) (string, bool) {
				c++
				if c > 1 {
					return "-", false
				}

				x, ok := next(buffer)
				if !ok {
					return "", false
				}

				return string(x), true
			},
			func(buffer Buffer[byte]) (func(string, string) string, bool) {
				return func(x, y string) string {
					return fmt.Sprintf("(%v %v)", x, y)
				}, true
			},
		)

		result, ok := ParseBytes([]byte("abcd"), comb)
		assert(t, ok, "expected true")
		assertEq(t, result, "a")

		c = 0

		result, ok = ParseBytes([]byte("a"), comb)
		assertEq(t, result, "a")
		assert(t, ok, "expected true")

		c = 0

		result, ok = ParseBytes([]byte(""), comb)
		assert(t, !ok, "expected false")
		assertEq(t, result, "")
	})

	t.Run("case 3", func(t *testing.T) {
		next := Satisfy(true, Anything[byte])

		comb := Chainr1(
			func(buffer Buffer[byte]) (string, bool) {
				x, ok := next(buffer)
				if !ok {
					return "", false
				}

				return string(x), true
			},
			func(buffer Buffer[byte]) (func(string, string) string, bool) {
				return func(x, y string) string { return "" }, false
			},
		)

		result, ok := ParseBytes([]byte("abcd"), comb)
		assert(t, ok, "expected true")
		assertEq(t, result, "a")

		result, ok = ParseBytes([]byte("a"), comb)
		assert(t, ok, "expected true")
		assertEq(t, result, "a")

		result, ok = ParseBytes([]byte(""), comb)
		assert(t, !ok, "expected false")
		assertEq(t, result, "")
	})
}

func TestSepBy(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		comb := SepBy(
			0,
			NotEq(true, byte(',')),
			Eq(true, byte(',')),
		)

		result, ok := ParseBytes([]byte("a,b,c"), comb)
		assert(t, ok, "expected true")
		assertSlice(t, result, []byte{'a', 'b', 'c'})

		result, ok = ParseBytes([]byte(""), comb)
		assert(t, ok, "expected true")
		assertSlice(t, result, []byte{})

		result, ok = ParseBytes([]byte(","), comb)
		assert(t, ok, "expected true")
		assertSlice(t, result, []byte{})

		result, ok = ParseBytes([]byte(",a,b,c"), comb)
		assert(t, ok, "expected true")
		assertSlice(t, result, []byte{})

		result, ok = ParseBytes([]byte("a,b,c,"), comb)
		assert(t, ok, "expected true")
		assertSlice(t, result, []byte{'a', 'b', 'c'})

		result, ok = ParseBytes([]byte("a,b,"), comb)
		assert(t, ok, "expected true")
		assertSlice(t, result, []byte{'a', 'b'})

		result, ok = ParseBytes([]byte("abc"), comb)
		assert(t, ok, "expected true")
		assertSlice(t, result, []byte{'a'})
	})

	t.Run("case 2", func(t *testing.T) {
		comb := SepBy(
			0,
			Satisfy(true, Nothing[byte]),
			Eq(true, byte(',')),
		)

		result, ok := ParseBytes([]byte("a,b,c"), comb)
		assert(t, ok, "expected true")
		assertSlice(t, result, nil)

		result, ok = ParseBytes([]byte(""), comb)
		assert(t, ok, "expected true")
		assertSlice(t, result, nil)
	})

	t.Run("case 3", func(t *testing.T) {
		comb := SepBy(
			0,
			NotEq(true, byte(',')),
			Satisfy(true, Nothing[byte]),
		)

		result, ok := ParseBytes([]byte("a,b,c"), comb)
		assert(t, ok, "expected true")
		assertSlice(t, result, []byte{'a'})

		result, ok = ParseBytes([]byte(""), comb)
		assert(t, ok, "expected true")
		assertSlice(t, result, nil)
	})
}

func TestSepBy1(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		comb := SepBy1(
			0,
			NotEq(true, byte(',')),
			Eq(true, byte(',')),
		)

		result, ok := ParseBytes([]byte("a,b,c"), comb)
		assert(t, ok, "expected true")
		assertSlice(t, result, []byte{'a', 'b', 'c'})

		result, ok = ParseBytes([]byte(""), comb)
		assert(t, !ok, "expected false")
		assertSlice(t, result, nil)

		result, ok = ParseBytes([]byte(","), comb)
		assert(t, !ok, "expected false")
		assertSlice(t, result, nil)

		result, ok = ParseBytes([]byte(",a,b,c"), comb)
		assert(t, !ok, "expected false")
		assertSlice(t, result, nil)

		result, ok = ParseBytes([]byte("a,b,c,"), comb)
		assert(t, ok, "expected true")
		assertSlice(t, result, []byte{'a', 'b', 'c'})

		result, ok = ParseBytes([]byte("a,b,"), comb)
		assert(t, ok, "expected true")
		assertSlice(t, result, []byte{'a', 'b'})

		result, ok = ParseBytes([]byte("abc"), comb)
		assert(t, ok, "expected true")
		assertSlice(t, result, []byte{'a'})
	})

	t.Run("case 2", func(t *testing.T) {
		comb := SepBy1(
			0,
			Satisfy(true, Nothing[byte]),
			Eq(true, byte(',')),
		)

		result, ok := ParseBytes([]byte("a,b,c"), comb)
		assert(t, !ok, "expected false")
		assertSlice(t, result, nil)

		result, ok = ParseBytes([]byte(""), comb)
		assert(t, !ok, "expected false")
		assertSlice(t, result, nil)
	})

	t.Run("case 3", func(t *testing.T) {
		comb := SepBy1(
			0,
			NotEq(true, byte(',')),
			Satisfy(true, Nothing[byte]),
		)

		result, ok := ParseBytes([]byte("a,b,c"), comb)
		assert(t, ok, "expected true")
		assertSlice(t, result, []byte{'a'})

		result, ok = ParseBytes([]byte(""), comb)
		assert(t, !ok, "expected false")
		assertSlice(t, result, nil)
	})
}

func TestEndBy(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		comb := EndBy(
			0,
			NotEq(true, byte(',')),
			Eq(true, byte(',')),
		)

		result, ok := ParseBytes([]byte("a,b,c"), comb)
		assert(t, ok, "expected true")
		assertSlice(t, result, []byte{'a', 'b'})

		result, ok = ParseBytes([]byte("a,b,c,"), comb)
		assert(t, ok, "expected true")
		assertSlice(t, result, []byte{'a', 'b', 'c'})

		result, ok = ParseBytes([]byte(""), comb)
		assert(t, ok, "expected true")
		assertSlice(t, result, nil)

		result, ok = ParseBytes([]byte(","), comb)
		assert(t, ok, "expected true")
		assertSlice(t, result, nil)

		result, ok = ParseBytes([]byte("a"), comb)
		assert(t, ok, "expected true")
		assertSlice(t, result, nil)

		result, ok = ParseBytes([]byte("a,"), comb)
		assert(t, ok, "expected true")
		assertSlice(t, result, []byte{'a'})

		result, ok = ParseBytes([]byte(",a"), comb)
		assert(t, ok, "expected true")
		assertSlice(t, result, nil)

		result, ok = ParseBytes([]byte("a,,"), comb)
		assert(t, ok, "expected true")
		assertSlice(t, result, []byte{'a'})

		result, ok = ParseBytes([]byte(",a,"), comb)
		assert(t, ok, "expected true")
		assertSlice(t, result, nil)
	})

	t.Run("case 2", func(t *testing.T) {
		comb := EndBy(
			0,
			Satisfy(true, Nothing[byte]),
			Eq(true, byte(',')),
		)

		result, ok := ParseBytes([]byte("a,b,c"), comb)
		assert(t, ok, "expected true")
		assertSlice(t, result, nil)

		result, ok = ParseBytes([]byte(""), comb)
		assert(t, ok, "expected true")
		assertSlice(t, result, nil)
	})

	t.Run("case 3", func(t *testing.T) {
		comb := EndBy(
			0,
			NotEq(true, byte(',')),
			Satisfy(true, Nothing[byte]),
		)

		result, ok := ParseBytes([]byte("a,b,c"), comb)
		assert(t, ok, "expected true")
		assertSlice(t, result, nil)

		result, ok = ParseBytes([]byte(""), comb)
		assert(t, ok, "expected true")
		assertSlice(t, result, nil)
	})
}

func TestEndBy1(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		comb := EndBy1(
			0,
			NotEq(true, byte(',')),
			Eq(true, byte(',')),
		)

		result, ok := ParseBytes([]byte("a,b,c"), comb)
		assert(t, ok, "expected true")
		assertSlice(t, result, []byte{'a', 'b'})

		result, ok = ParseBytes([]byte("a,b,c,"), comb)
		assert(t, ok, "expected true")
		assertSlice(t, result, []byte{'a', 'b', 'c'})

		result, ok = ParseBytes([]byte(""), comb)
		assert(t, !ok, "expected false")
		assertSlice(t, result, nil)

		result, ok = ParseBytes([]byte(","), comb)
		assert(t, !ok, "expected false")
		assertSlice(t, result, nil)

		result, ok = ParseBytes([]byte("a"), comb)
		assert(t, !ok, "expected false")
		assertSlice(t, result, nil)

		result, ok = ParseBytes([]byte("a,"), comb)
		assert(t, ok, "expected true")
		assertSlice(t, result, []byte{'a'})

		result, ok = ParseBytes([]byte(",a"), comb)
		assert(t, !ok, "expected false")
		assertSlice(t, result, nil)

		result, ok = ParseBytes([]byte("a,,"), comb)
		assert(t, ok, "expected true")
		assertSlice(t, result, []byte{'a'})

		result, ok = ParseBytes([]byte(",a,"), comb)
		assert(t, !ok, "expected false")
		assertSlice(t, result, nil)
	})

	t.Run("case 2", func(t *testing.T) {
		comb := EndBy1(
			0,
			Satisfy(true, Nothing[byte]),
			Eq(true, byte(',')),
		)

		result, ok := ParseBytes([]byte("a,b,c"), comb)
		assert(t, !ok, "expected false")
		assertSlice(t, result, nil)

		result, ok = ParseBytes([]byte(""), comb)
		assert(t, !ok, "expected false")
		assertSlice(t, result, nil)
	})

	t.Run("case 3", func(t *testing.T) {
		comb := EndBy1(
			0,
			NotEq(true, byte(',')),
			Satisfy(true, Nothing[byte]),
		)

		result, ok := ParseBytes([]byte("a,b,c"), comb)
		assert(t, !ok, "expected false")
		assertSlice(t, result, nil)

		result, ok = ParseBytes([]byte(""), comb)
		assert(t, !ok, "expected false")
		assertSlice(t, result, nil)
	})
}

func TestSepEndBy(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		comb := SepEndBy(
			0,
			NotEq(true, byte(',')),
			Eq(true, byte(',')),
		)

		result, ok := ParseBytes([]byte("a,b,c"), comb)
		assert(t, ok, "expected true")
		assertSlice(t, result, []byte{'a', 'b', 'c'})

		result, ok = ParseBytes([]byte("a,b,c,"), comb)
		assert(t, ok, "expected true")
		assertSlice(t, result, []byte{'a', 'b', 'c'})

		result, ok = ParseBytes([]byte("a,b,c,,"), comb)
		assert(t, ok, "expected true")
		assertSlice(t, result, []byte{'a', 'b', 'c'})

		result, ok = ParseBytes([]byte(""), comb)
		assert(t, ok, "expected true")
		assertSlice(t, result, nil)

		result, ok = ParseBytes([]byte(","), comb)
		assert(t, ok, "expected true")
		assertSlice(t, result, nil)

		result, ok = ParseBytes([]byte(",a,b,c"), comb)
		assert(t, ok, "expected true")
		assertSlice(t, result, nil)
	})

	t.Run("case 2", func(t *testing.T) {
		comb := SepEndBy(
			0,
			Satisfy(true, Nothing[byte]),
			Eq(true, byte(',')),
		)

		result, ok := ParseBytes([]byte("a,b,c"), comb)
		assert(t, ok, "expected true")
		assertSlice(t, result, nil)

		result, ok = ParseBytes([]byte(","), comb)
		assert(t, ok, "expected true")
		assertSlice(t, result, nil)

		result, ok = ParseBytes([]byte(""), comb)
		assert(t, ok, "expected true")
		assertSlice(t, result, nil)
	})

	t.Run("case 3", func(t *testing.T) {
		comb := SepEndBy(
			0,
			NotEq(true, byte(',')),
			Satisfy(true, Nothing[byte]),
		)

		result, ok := ParseBytes([]byte("a,b,c"), comb)
		assert(t, ok, "expected true")
		assertSlice(t, result, []byte{'a'})

		result, ok = ParseBytes([]byte(","), comb)
		assert(t, ok, "expected true")
		assertSlice(t, result, nil)

		result, ok = ParseBytes([]byte(""), comb)
		assert(t, ok, "expected true")
		assertSlice(t, result, nil)
	})
}

func TestSepEndBy1(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		comb := SepEndBy1(
			0,
			NotEq(true, byte(',')),
			Eq(true, byte(',')),
		)

		result, ok := ParseBytes([]byte("a,b,c"), comb)
		assert(t, ok, "expected true")
		assertSlice(t, result, []byte{'a', 'b', 'c'})

		result, ok = ParseBytes([]byte("a,b,c,"), comb)
		assert(t, ok, "expected true")
		assertSlice(t, result, []byte{'a', 'b', 'c'})

		result, ok = ParseBytes([]byte("a,b,c,,"), comb)
		assert(t, ok, "expected true")
		assertSlice(t, result, []byte{'a', 'b', 'c'})

		result, ok = ParseBytes([]byte(""), comb)
		assert(t, !ok, "expected false")
		assertSlice(t, result, nil)

		result, ok = ParseBytes([]byte(","), comb)
		assert(t, !ok, "expected false")
		assertSlice(t, result, nil)

		result, ok = ParseBytes([]byte(",a,b,c"), comb)
		assert(t, !ok, "expected false")
		assertSlice(t, result, nil)
	})

	t.Run("case 2", func(t *testing.T) {
		comb := SepEndBy1(
			0,
			Satisfy(true, Nothing[byte]),
			Eq(true, byte(',')),
		)

		result, ok := ParseBytes([]byte("a,b,c"), comb)
		assert(t, !ok, "expected false")
		assertSlice(t, result, nil)

		result, ok = ParseBytes([]byte(""), comb)
		assert(t, !ok, "expected false")
		assertSlice(t, result, nil)
	})

	t.Run("case 3", func(t *testing.T) {
		comb := SepEndBy1(
			0,
			NotEq(true, byte(',')),
			Satisfy(true, Nothing[byte]),
		)

		result, ok := ParseBytes([]byte("a,b,c"), comb)
		assert(t, ok, "expected true")
		assertSlice(t, result, []byte{'a'})

		result, ok = ParseBytes([]byte(""), comb)
		assert(t, !ok, "expected false")
		assertSlice(t, result, nil)
	})
}

func TestManyTill(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		comb := ManyTill(
			0,
			Any[byte](true),
			EOF[byte](),
		)

		result, ok := ParseBytes([]byte("abcd"), comb)
		assert(t, ok, "expected true")
		assertSlice(t, result, []byte{
			byte('a'),
			byte('b'),
			byte('c'),
			byte('d'),
		})
	})

	t.Run("case 2", func(t *testing.T) {
		comb := ManyTill(
			0,
			Any[byte](true),
			Satisfy[byte](false, func(x byte) bool { return x == byte('d') }),
		)

		result, ok := ParseBytes([]byte("abcd"), comb)
		assert(t, ok, "expected true")
		assertSlice(t, result, []byte{
			byte('a'),
			byte('b'),
			byte('c'),
		})

		result, ok = ParseBytes([]byte(""), comb)
		assert(t, !ok, "expected false")
		assertSlice(t, result, nil)
	})

	t.Run("case 3", func(t *testing.T) {
		comb := ManyTill(
			0,
			Satisfy(true, Nothing[byte]),
			EOF[byte](),
		)

		result, ok := ParseBytes([]byte("abcd"), comb)
		assert(t, !ok, "expected false")
		assertSlice(t, result, nil)
	})

	t.Run("case 4", func(t *testing.T) {
		comb := ManyTill(
			0,
			Any[byte](true),
			Satisfy(true, Nothing[byte]),
		)

		result, ok := ParseBytes([]byte("abcd"), comb)
		assert(t, !ok, "expected false")
		assertSlice(t, result, nil)
	})
}
