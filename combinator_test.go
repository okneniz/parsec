package parsec

import (
	"encoding/json"
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

func TestTry(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
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
	})
}

func TestBefore(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
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
	})
}

func TestAfter(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
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
	})
}

func TestBetween(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
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
	})
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

	t.Run("case 3", func(t *testing.T) {
		comb := Skip(
			NotEq(true, byte('a')),
			Eq(true, byte('a')),
		)

		result, ok := ParseBytes([]byte("abc"), comb)
		assert(t, !ok, "expected false")
		assertEq(t, result, 0)
	})
}

func TestSkipAfter(t *testing.T) {
	t.Run("case 1", func(t *testing.T) {
		comb := SkipAfter(
			Eq(true, byte('b')),
			Eq(true, byte('a')),
		)

		result, ok := ParseBytes([]byte("abc"), comb)
		assert(t, ok, "expected true")
		assertEq(t, result, byte('a'))

		result, ok = ParseBytes([]byte("ab"), comb)
		assert(t, ok, "expected true")
		assertEq(t, result, byte('a'))

		result, ok = ParseBytes([]byte("a"), comb)
		assert(t, !ok, "expected false")
		assertEq(t, result, 0)
	})

	t.Run("case 2", func(t *testing.T) {
		comb := SkipAfter(
			Eq(true, byte('b')),
			Satisfy[byte](true, Nothing[byte]),
		)

		result, ok := ParseBytes([]byte("abc"), comb)
		assert(t, !ok, "expected false")
		assertEq(t, result, 0)
	})

	t.Run("case 3", func(t *testing.T) {
		comb := SkipAfter(
			Satisfy[byte](true, Nothing[byte]),
			Eq(true, byte('a')),
		)

		result, ok := ParseBytes([]byte("abc"), comb)
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

func TestCast(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		comb := Cast[byte, byte, int](
			Satisfy[byte](true, Anything[byte]),
			func(x byte) int { return int(x) },
		)

		result, ok := ParseBytes([]byte("a"), comb)
		assert(t, ok, "expected true")
		assertEq(t, result, 97)

		result, ok = ParseBytes([]byte("b"), comb)
		assert(t, ok, "expected true")
		assertEq(t, result, 98)

		result, ok = ParseBytes([]byte(""), comb)
		assert(t, !ok, "expected false")
		assertEq(t, result, 0)
	})

	// TODO : more cases
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
