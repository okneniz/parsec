package parsec

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"
)

func TestSatisfy(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		c := byte('c')

		comb := Satisfy[byte](true, func(x byte) bool { return x != c })

		result, ok := ParseBytes([]byte("a"), comb)
		assert(t, ok, "expected true")
		assertEq(t, result, 'a')

		result, ok = ParseBytes([]byte("b"), comb)
		assert(t, ok, "expected true")
		assertEq(t, result, 'b')

		result, ok = ParseBytes([]byte("c"), comb)
		assert(t, !ok, "expected false")
		assertEq(t, result, 0)
	})

	t.Run("case 2", func(t *testing.T) {
		comb := Satisfy[byte](true, func(x byte) bool { return false })

		result, ok := ParseBytes([]byte{}, comb)
		assert(t, !ok, "expected false")
		assertEq(t, result, 0)
	})
}

func TestEq(t *testing.T) {
	c := byte('c')

	comb := Eq[byte](true, c)

	result, ok := ParseBytes([]byte("a"), comb)
	assert(t, !ok, "expected false")
	assertEq(t, result, 0)

	result, ok = ParseBytes([]byte("b"), comb)
	assert(t, !ok, "expected false")
	assertEq(t, result, 0)

	result, ok = ParseBytes([]byte("c"), comb)
	assert(t, ok, "expected true")
	assertEq(t, result, c)
}

func TestNotEq(t *testing.T) {
	c := byte('c')

	comb := NotEq[byte](true, c)

	result, ok := ParseBytes([]byte("a"), comb)
	assert(t, ok, "expected true")
	assertEq(t, result, byte('a'))

	result, ok = ParseBytes([]byte("b"), comb)
	assert(t, ok, "expected true")
	assertEq(t, result, byte('b'))

	result, ok = ParseBytes([]byte("abc"), comb)
	assert(t, ok, "expected true")
	assertEq(t, result, byte('a'))

	result, ok = ParseBytes([]byte("c"), comb)
	assert(t, !ok, "expected false")
	assertEq(t, result, 0)
}

func TestSequence(t *testing.T) {
	comb := Sequence(
		3,
		Eq[byte](true, byte('a')),
		Eq[byte](true, byte('b')),
		Satisfy[byte](true, func(x byte) bool { return x != byte('z') }),
	)

	result, ok := ParseBytes([]byte("abc"), comb)
	assert(t, ok, "expected true")
	assertSlice(t, result, []byte("abc"))

	result, ok = ParseBytes([]byte("abd"), comb)
	assert(t, ok, "expected true")
	assertSlice(t, result, []byte("abd"))

	result, ok = ParseBytes([]byte("abdasdasd"), comb)
	assert(t, ok, "expected true")
	assertSlice(t, result, []byte("abd"))

	result, ok = ParseBytes([]byte("xyz"), comb)
	assert(t, !ok, "expected false")
	assertSlice(t, result, nil)
}

func TestMany(t *testing.T) {
	comb := Many(0, Eq[byte](true, byte('a')))

	result, ok := ParseBytes([]byte("aaa"), comb)
	assert(t, ok, "expected true")
	assertSlice(t, result, []byte("aaa"))

	result, ok = ParseBytes([]byte("aaabc"), comb)
	assert(t, ok, "expected true")
	assertSlice(t, result, []byte("aaa"))

	result, ok = ParseBytes([]byte("xaaabc"), comb)
	assert(t, ok, "expected true")
	assertSlice(t, result, []byte{})
}

func TestSome(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		comb := Some(
			0,
			Eq[byte](true, byte('a')),
		)

		result, ok := ParseBytes([]byte("aaa"), comb)
		assert(t, ok, "expected true")
		assertSlice(t, result, []byte("aaa"))

		result, ok = ParseBytes([]byte("aaabc"), comb)
		assert(t, ok, "expected true")
		assertSlice(t, result, []byte("aaa"))

		result, ok = ParseBytes([]byte("xaaabc"), comb)
		assert(t, !ok, "expected false")
		assertSlice(t, result, nil)
	})

	t.Run("case 2", func(t *testing.T) {
		comb := Some(
			0,
			Satisfy[byte](true, func(x byte) bool { return false }),
		)

		result, ok := ParseBytes([]byte("abc"), comb)
		assert(t, !ok, "expected false")
		assertSlice(t, result, []byte{})
	})
}

func TestOptional(t *testing.T) {
	comb := Optional(Eq[byte](true, byte('a')), 0)

	result, ok := ParseBytes([]byte("aaa"), comb)
	assert(t, ok, "expected true")
	assertEq(t, result, byte('a'))

	result, ok = ParseBytes([]byte("bcd"), comb)
	assert(t, ok, "expected true")
	assertEq(t, result, 0)
}

func TestTry(t *testing.T) {
	comb := Try(
		Satisfy(true, func(x byte) bool { return x <= byte('b') }),
	)

	buf := BytesBuffer([]byte("abcd"))
	assertEq(t, buf.Position(), 0)

	result, ok := Parse[byte, byte](buf, comb)
	assert(t, ok, "expected true")
	assertEq(t, result, byte('a'))
	assertEq(t, buf.Position(), 1)

	result, ok = Parse[byte, byte](buf, comb)
	assert(t, ok, "expected true")
	assertEq(t, result, byte('b'))
	assertEq(t, buf.Position(), 2)

	result, ok = Parse[byte, byte](buf, comb)
	assert(t, !ok, "expected false")
	assertEq(t, result, 0)
	assertEq(t, buf.Position(), 2)

	result, ok = Parse[byte, byte](buf, comb)
	assert(t, !ok, "expected false")
	assertEq(t, result, 0)
	assertEq(t, buf.Position(), 2)
}

func TestOr(t *testing.T) {
	comb := Or(
		Try(Eq(true, byte('a'))),
		Eq(true, byte('b')),
	)

	result, ok := ParseBytes([]byte("a"), comb)
	assert(t, ok, "expected true")
	assertEq(t, result, byte('a'))

	result, ok = ParseBytes([]byte("b"), comb)
	assert(t, ok, "expected true")
	assertEq(t, result, byte('b'))

	result, ok = ParseBytes([]byte("c"), comb)
	assert(t, !ok, "expected false")
	assertEq(t, result, 0)
}

func TestAnd(t *testing.T) {
	comb := And(
		Eq(true, byte('a')),
		Eq(true, byte('b')),
		func(x, y byte) []byte { return []byte{x, y} },
	)

	result, ok := ParseBytes([]byte("abc"), comb)
	assert(t, ok, "expected true")
	assertSlice(t, result, []byte{'a', 'b'})

	result, ok = ParseBytes([]byte("bca"), comb)
	assert(t, !ok, "expected false")
	assertSlice(t, result, nil)

	result, ok = ParseBytes([]byte("acb"), comb)
	assert(t, !ok, "expected false")
	assertSlice(t, result, nil)
}

func TestBefore(t *testing.T) {
	comb := Before(
		Eq(true, byte('a')),
		Eq(true, byte('b')),
		func(x, y byte) []byte { return []byte{x, y} },
	)

	result, ok := ParseBytes([]byte("bac"), comb)
	assert(t, ok, "expected true")
	assertSlice(t, result, []byte{'b', 'a'})

	result, ok = ParseBytes([]byte("bca"), comb)
	assert(t, !ok, "expected false")
	assertSlice(t, result, nil)

	result, ok = ParseBytes([]byte("abc"), comb)
	assert(t, !ok, "expected false")
	assertSlice(t, result, nil)
}

func TestAfter(t *testing.T) {
	comb := After(
		Eq(true, byte('a')),
		Eq(true, byte('b')),
		func(x, y byte) []byte { return []byte{x, y} },
	)

	result, ok := ParseBytes([]byte("abc"), comb)
	assert(t, ok, "expected true")
	assertSlice(t, result, []byte{'a', 'b'})

	result, ok = ParseBytes([]byte("bac"), comb)
	assert(t, !ok, "expected false")
	assertSlice(t, result, nil)
}

func TestBetween(t *testing.T) {
	notBrackets := Satisfy[byte](true, func(x byte) bool {
		return !(x == byte(')') || x == byte('('))
	})

	comb := Between(
		Eq(true, byte('(')),
		Some(0, Try(notBrackets)),
		Eq(true, byte(')')),
		func(x byte, y []byte, z byte) []byte { return y },
	)

	result, ok := ParseBytes([]byte("(abc)"), comb)
	assert(t, ok, "expected true")
	assertSlice(t, result, []byte{'a', 'b', 'c'})

	result, ok = ParseBytes([]byte("(abc)def"), comb)
	assert(t, ok, "expected true")
	assertSlice(t, result, []byte{'a', 'b', 'c'})

	result, ok = ParseBytes([]byte("(abc))"), comb)
	assert(t, ok, "expected true")
	assertSlice(t, result, []byte{'a', 'b', 'c'})

	result, ok = ParseBytes([]byte("(ab)"), comb)
	assert(t, ok, "expected true")
	assertSlice(t, result, []byte{'a', 'b'})

	result, ok = ParseBytes([]byte("x(abc)def"), comb)
	assert(t, !ok, "expected false")
	assertSlice(t, result, nil)

	result, ok = ParseBytes([]byte("()"), comb)
	assert(t, !ok, "expected false")
	assertSlice(t, result, nil)

	result, ok = ParseBytes([]byte("(()"), comb)
	assert(t, !ok, "expected false")
	assertSlice(t, result, nil)

	result, ok = ParseBytes([]byte("((1))"), comb)
	assert(t, !ok, "expected false")
	assertSlice(t, result, nil)

	result, ok = ParseBytes([]byte("(abc"), comb)
	assert(t, !ok, "expected false")
	assertSlice(t, result, nil)

	result, ok = ParseBytes([]byte("(abc("), comb)
	assert(t, !ok, "expected false")
	assertSlice(t, result, nil)

	result, ok = ParseBytes([]byte("((abc)"), comb)
	assert(t, !ok, "expected false")
	assertSlice(t, result, nil)
}

func TestOneOf(t *testing.T) {
	comb := OneOf(true, byte('a'), byte('b'), byte('c'))

	result, ok := ParseBytes([]byte("a"), comb)
	assert(t, ok, "expected true")
	assertEq(t, result, byte('a'))

	result, ok = ParseBytes([]byte("b"), comb)
	assert(t, ok, "expected true")
	assertEq(t, result, byte('b'))

	result, ok = ParseBytes([]byte("c"), comb)
	assert(t, ok, "expected true")
	assertEq(t, result, byte('c'))

	result, ok = ParseBytes([]byte("d"), comb)
	assert(t, !ok, "expected false")
	assertEq(t, result, 0)
}

func TestCount(t *testing.T) {
	comb := Count(2, Eq(true, byte('a')))

	result, ok := ParseBytes([]byte("aabbcc"), comb)
	assert(t, ok, "expected true")
	assertSlice(t, result, []byte{'a', 'a'})

	result, ok = ParseBytes([]byte("abbcc"), comb)
	assert(t, !ok, "expected false")
	assertSlice(t, result, nil)

	result, ok = ParseBytes([]byte("bbaacc"), comb)
	assert(t, !ok, "expected false")
	assertSlice(t, result, nil)
}

func TestSkip(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		comb := Skip(
			Optional(Eq(true, byte('a')), 0),
			Eq(true, byte('b')),
		)

		result, ok := ParseBytes([]byte("abc"), comb)
		assert(t, ok, "expected true")
		assertEq(t, result, byte('b'))

		result, ok = ParseBytes([]byte("cba"), comb)
		assert(t, ok, "expected true")
		assertEq(t, result, byte('b'))
	})

	t.Run("case 2", func(t *testing.T) {
		phrase := Sequence(
			3,
			Eq(true, byte('a')),
			Eq(true, byte('b')),
			Eq(true, byte('c')),
		)

		noice := Many(
			0,
			Try(
				NoneOf(
					true,
					byte('a'),
					byte('b'),
					byte('c'),
				),
			),
		)

		comb := Skip(noice, phrase)

		result, ok := ParseBytes([]byte("abc"), comb)
		assert(t, ok, "expected true")
		assertSlice(t, result, []byte{'a', 'b', 'c'})

		result, ok = ParseBytes([]byte("abc123"), comb)
		assert(t, ok, "expected true")
		assertSlice(t, result, []byte{'a', 'b', 'c'})

		result, ok = ParseBytes([]byte("123abc"), comb)
		assert(t, ok, "expected true")
		assertSlice(t, result, []byte{'a', 'b', 'c'})

		result, ok = ParseBytes([]byte("123abc123"), comb)
		assert(t, ok, "expected true")
		assertSlice(t, result, []byte{'a', 'b', 'c'})
	})
}

func TestSepBy(t *testing.T) {
	comb := NotEq(true, byte(','))
	sep := Eq(true, byte(','))

	result, ok := ParseBytes([]byte("a,b,c"), SepBy(0, comb, sep))
	assert(t, ok, "expected true")
	assertSlice(t, result, []byte{'a', 'b', 'c'})

	result, ok = ParseBytes([]byte(""), SepBy(0, comb, sep))
	assert(t, ok, "expected true")
	assertSlice(t, result, []byte{})

	result, ok = ParseBytes([]byte(","), SepBy(0, comb, sep))
	assert(t, ok, "expected true")
	assertSlice(t, result, []byte{})

	result, ok = ParseBytes([]byte(",a,b,c"), SepBy(0, comb, sep))
	assert(t, ok, "expected true")
	assertSlice(t, result, []byte{})

	result, ok = ParseBytes([]byte("a,b,c,"), SepBy(0, comb, sep))
	assert(t, ok, "expected true")
	assertSlice(t, result, []byte{'a', 'b', 'c'})

	result, ok = ParseBytes([]byte("a,b,"), SepBy(0, comb, sep))
	assert(t, ok, "expected true")
	assertSlice(t, result, []byte{'a', 'b'})

	result, ok = ParseBytes([]byte("abc"), SepBy(0, comb, sep))
	assert(t, ok, "expected true")
	assertSlice(t, result, []byte{'a'})
}

func TestSepBy1(t *testing.T) {
	comb := NotEq(true, byte(','))
	sep := Eq(true, byte(','))

	result, ok := ParseBytes([]byte("a,b,c"), SepBy1(0, comb, sep))
	assert(t, ok, "expected true")
	assertSlice(t, result, []byte{'a', 'b', 'c'})

	result, ok = ParseBytes([]byte(""), SepBy1(0, comb, sep))
	assert(t, !ok, "expected false")
	assertSlice(t, result, nil)

	result, ok = ParseBytes([]byte(","), SepBy1(0, comb, sep))
	assert(t, !ok, "expected false")
	assertSlice(t, result, nil)

	result, ok = ParseBytes([]byte(",a,b,c"), SepBy1(0, comb, sep))
	assert(t, !ok, "expected false")
	assertSlice(t, result, nil)

	result, ok = ParseBytes([]byte("a,b,c,"), SepBy1(0, comb, sep))
	assert(t, ok, "expected true")
	assertSlice(t, result, []byte{'a', 'b', 'c'})

	result, ok = ParseBytes([]byte("a,b,"), SepBy1(0, comb, sep))
	assert(t, ok, "expected true")
	assertSlice(t, result, []byte{'a', 'b'})

	result, ok = ParseBytes([]byte("abc"), SepBy1(0, comb, sep))
	assert(t, ok, "expected true")
	assertSlice(t, result, []byte{'a'})
}

func TestEndBy(t *testing.T) {
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
}

func TestEndBy1(t *testing.T) {
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
}

func TestSepEndBy(t *testing.T) {
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
}

func TestSepEndBy1(t *testing.T) {
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
}

func TestChainl(t *testing.T) {
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

	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		result, ok := ParseBytes([]byte("abcd"), comb)
		assert(t, ok, "expected true")
		assertEq(t, result, "(((a b) c) d)")
	})

	t.Run("case 2", func(t *testing.T) {
		result, ok := ParseBytes([]byte("a"), comb)
		assert(t, ok, "expected true")
		assertEq(t, result, "a")
	})

	t.Run("case 3", func(t *testing.T) {
		result, ok := ParseBytes([]byte(""), comb)
		assert(t, ok, "expected true")
		assertEq(t, result, "default")
	})
}

func TestChainl1(t *testing.T) {
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

	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		result, ok := ParseBytes([]byte("abcd"), comb)
		assert(t, ok, "expected true")
		assertEq(t, result, "(((a b) c) d)")
	})

	t.Run("case 2", func(t *testing.T) {
		result, ok := ParseBytes([]byte("a"), comb)
		assert(t, ok, "expected true")
		assertEq(t, result, "a")
	})

	t.Run("case 3", func(t *testing.T) {
		result, ok := ParseBytes([]byte(""), comb)
		assert(t, !ok, "expected false")
		assertEq(t, result, "")
	})
}

func TestChainr(t *testing.T) {
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

	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		result, ok := ParseBytes([]byte("abcd"), comb)
		assert(t, ok, "expected true")
		assertEq(t, result, "(a (b (c d)))")
	})

	t.Run("case 2", func(t *testing.T) {
		result, ok := ParseBytes([]byte("a"), comb)
		assert(t, ok, "expected true")
		assertEq(t, result, "a")
	})

	t.Run("case 3", func(t *testing.T) {
		result, ok := ParseBytes([]byte(""), comb)
		assert(t, ok, "expected true")
		assertEq(t, result, "default")
	})
}

func TestChainr1(t *testing.T) {
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

	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		result, ok := ParseBytes([]byte("abcd"), comb)
		assert(t, ok, "expected true")
		assertEq(t, result, "(a (b (c d)))")
	})

	t.Run("case 2", func(t *testing.T) {
		result, ok := ParseBytes([]byte("a"), comb)
		assert(t, ok, "expected true")
		assertEq(t, result, "a")
	})

	t.Run("case 3", func(t *testing.T) {
		result, ok := ParseBytes([]byte(""), comb)
		assert(t, !ok, "expected false")
		assertEq(t, result, "")
	})
}

func TestAny(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		source := rand.New(rand.NewSource(time.Now().UnixNano()))

		comb := Any[byte](true)

		for i := 0; i < 10000; i++ {
			b := byte(source.Intn(math.MaxUint8 + 1))

			result, ok := ParseBytes([]byte{b}, comb)
			assert(t, ok, "expected true")
			assertEq(t, result, b)
		}
	})

	t.Run("case 2", func(t *testing.T) {
		comb := Any[byte](true)

		result, ok := ParseBytes([]byte{}, comb)
		assert(t, !ok, "expected false")
		assertEq(t, result, 0)
	})
}

func TestEOF(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		_, ok := ParseBytes([]byte("abcd"), EOF[byte]())
		assert(t, !ok, "expected false")
	})

	t.Run("case 2", func(t *testing.T) {
		_, ok := ParseBytes([]byte(""), EOF[byte]())
		assert(t, ok, "expected true")
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
}

func assertEqDump[T any](t *testing.T, actual, expected T) {
	t.Helper()

	ex, err := json.Marshal(expected)
	if err != nil {
		t.Fatal(err)
	}

	ac, err := json.Marshal(actual)
	if err != nil {
		t.Fatal(err)
	}

	if string(ex) != string(ac) {
		t.Errorf("expected %v", string(ex))
		t.Errorf("actual %v", string(ac))
		t.Fatal("invalid result")
	}
}

func assert(t *testing.T, x bool, m string) {
	t.Helper()

	if !x {
		t.Fatal(m)
	}
}

func assertEq[T comparable](t *testing.T, x, y T) {
	t.Helper()

	if x != y {
		t.Fatalf("%v != %v", x, y)
	}
}

func assertEqPointer[T comparable](t *testing.T, x, y *T) {
	t.Helper()

	if x == nil && y == nil {
		return
	}

	if (x == nil && y != nil) || (x != nil && y == nil) {
		t.Fatalf("%v != %v", x, y)
	}

	if *x != *y {
		t.Fatalf("%v != %v", *x, *y)
	}
}

func assertSlice[T comparable](t *testing.T, xs, ys []T) {
	t.Helper()

	if len(xs) != len(ys) {
		t.Fatalf("%v != %v", xs, ys)
	}

	for i, x := range xs {
		if x != ys[i] {
			t.Fatalf("%v != %v", xs, ys)
		}
	}
}

func pointer[T any](t T) *T {
	return &t
}
