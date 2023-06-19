package parsec

import (
	"testing"
	"encoding/json"
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
		assert(t, !ok, "expected true")
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

func TestSlice(t *testing.T) {
	comb := Slice(
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
		assert(t, !ok, "expected true")
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
	comb := Optional(Eq[byte](true, byte('a')))

	result, ok := ParseBytes([]byte("aaa"), comb)
	assert(t, ok, "expected true")
	assertEqPointer(t, result, pointer(byte('a')))

	result, ok = ParseBytes([]byte("bcd"), comb)
	assert(t, ok, "expected true")
	assertEqPointer(t, result, nil)
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

func TestNot(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		comb := Not(Satisfy(true, Nothing[byte]))

		result, ok := ParseBytes([]byte("abc"), comb)
		assert(t, ok, "expected true")
		assertEq(t, result, 0) // ?
	})

	t.Run("case 2", func(t *testing.T) {
		comb := Not(Satisfy(true, Anything[byte]))

		result, ok := ParseBytes([]byte("abc"), comb)
		assert(t, !ok, "expected false")
		assertEq(t, result, 0)
	})
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
		func(x, y byte) []byte { return []byte{x,y} },
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
		func(x, y byte) []byte { return []byte{x,y} },
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
		func(x, y byte) []byte { return []byte{x,y} },
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
	assert(t, !ok, "expected true")
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
			Optional(Eq(true, byte('a'))),
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
		phrase := Slice(
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
	type node struct {
		Value byte `json:"v"`
		Next *node `json:"n,omitempty"`
	}

	next := Satisfy(true, Anything[byte])

	comb := Chainl(
		func(buffer Buffer[byte]) (*node, bool) {
			x, ok := next(buffer)
			if !ok {
				return nil, false
			}

			return &node{Value: x}, true
		},
		func(buffer Buffer[byte]) (func(*node, *node) *node, bool) {
			return func(x, y *node) *node {
				x.Next = y
				return y
			}, true
		},
	)

	result, ok := ParseBytes([]byte("abc"), comb)
	assert(t, ok, "expected true")
	assertEqDump(
		t,
		result,
		&node{
			Value: byte('a'),
			Next: &node{
				Value: byte('b'),
				Next: &node{
					Value: byte('c'),
				},
			},
		},
	)

	result, ok = ParseBytes([]byte("a"), comb)
	assert(t, ok, "expected true")
	assertEqDump(
		t,
		result,
		&node{Value: byte('a')},
	)

	result, ok = ParseBytes([]byte(""), comb)
	assert(t, !ok, "expected false")
	assertEqDump(t, result, nil)
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
