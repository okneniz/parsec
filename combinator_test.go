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

		comb := Satisfy[byte,int](true, func(x byte) bool { return x != c })

		result, err := ParseBytes([]byte("a"), comb)
		check(t, err)
		assertEq(t, result, 'a')

		result, err = ParseBytes([]byte("b"), comb)
		check(t, err)
		assertEq(t, result, 'b')

		result, err = ParseBytes([]byte("c"), comb)
		assertError(t, err)
		assertEq(t, result, 0)
	})

	t.Run("case 2", func(t *testing.T) {
		comb := Satisfy[byte,int](true, func(x byte) bool { return false })

		result, err := ParseBytes([]byte{}, comb)
		assertError(t, err)
		assertEq(t, result, 0)
	})
}

func TestAny(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		source := rand.New(rand.NewSource(time.Now().UnixNano()))
		comb := Any[byte,int]()

		for i := 0; i < 10000; i++ {
			b := byte(source.Intn(math.MaxUint8 + 1))

			result, err := ParseBytes([]byte{b}, comb)
			check(t, err)
			assertEq(t, result, b)
		}
	})

	t.Run("case 2", func(t *testing.T) {
		comb := Any[byte,int]()

		result, err := ParseBytes([]byte{}, comb)
		assertError(t, err)
		assertEq(t, result, 0)
	})
}

func TestTry(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		comb := Try(
			Satisfy[byte,int](true, func(x byte) bool { return x <= byte('b') }),
		)

		buf := BytesBuffer([]byte("abcd"))
		assertEq(t, buf.Position(), 0)

		result, err := Parse[byte, int, byte](buf, comb)
		check(t, err)
		assertEq(t, result, byte('a'))
		assertEq(t, buf.Position(), 1)

		result, err = Parse[byte, int, byte](buf, comb)
		check(t, err)
		assertEq(t, result, byte('b'))
		assertEq(t, buf.Position(), 2)

		result, err = Parse[byte, int, byte](buf, comb)
		assertError(t, err)
		assertEq(t, result, 0)
		assertEq(t, buf.Position(), 2)

		result, err = Parse[byte, int, byte](buf, comb)
		assertError(t, err)
		assertEq(t, result, 0)
		assertEq(t, buf.Position(), 2)
	})
}

func TestBetween(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		notBrackets := Satisfy[byte,int](true, func(x byte) bool {
			return !(x == byte(')') || x == byte('('))
		})

		comb := Between(
			Eq[byte,int]('('),
			Some[byte,int](0, Try(notBrackets)),
			Eq[byte,int](')'),
		)

		result, err := ParseBytes([]byte("(abc)"), comb)
		check(t, err)
		assertSlice(t, result, []byte{'a', 'b', 'c'})

		result, err = ParseBytes([]byte("(abc)def"), comb)
		check(t, err)
		assertSlice(t, result, []byte{'a', 'b', 'c'})

		result, err = ParseBytes([]byte("(abc))"), comb)
		check(t, err)
		assertSlice(t, result, []byte{'a', 'b', 'c'})

		result, err = ParseBytes([]byte("(ab)"), comb)
		check(t, err)
		assertSlice(t, result, []byte{'a', 'b'})

		result, err = ParseBytes([]byte("x(abc)def"), comb)
		assertError(t, err)
		assertSlice(t, result, nil)

		result, err = ParseBytes([]byte("()"), comb)
		assertError(t, err)
		assertSlice(t, result, nil)

		result, err = ParseBytes([]byte("(()"), comb)
		assertError(t, err)
		assertSlice(t, result, nil)

		result, err = ParseBytes([]byte("((1))"), comb)
		assertError(t, err)
		assertSlice(t, result, nil)

		result, err = ParseBytes([]byte("(abc"), comb)
		assertError(t, err)
		assertSlice(t, result, nil)

		result, err = ParseBytes([]byte("(abc("), comb)
		assertError(t, err)
		assertSlice(t, result, nil)

		result, err = ParseBytes([]byte("((abc)"), comb)
		assertError(t, err)
		assertSlice(t, result, nil)
	})
}

func TestSkip(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		comb := Skip(
			Optional(Eq[byte,int]('a'), 0),
			Eq[byte,int]('b'),
		)

		result, err := ParseBytes([]byte("abc"), comb)
		check(t, err)
		assertEq(t, result, byte('b'))

		result, err = ParseBytes([]byte("cba"), comb)
		check(t, err)
		assertEq(t, result, byte('b'))
	})

	t.Run("case 2", func(t *testing.T) {
		phrase := SequenceOf[byte,int]('a', 'b', 'c')
		noice := Many(0, Try(NoneOf[byte,int]('a', 'b', 'c')))
		comb := Skip(noice, phrase)

		result, err := ParseBytes([]byte("abc"), comb)
		check(t, err)
		assertSlice(t, result, []byte{'a', 'b', 'c'})

		result, err = ParseBytes([]byte("abc123"), comb)
		check(t, err)
		assertSlice(t, result, []byte{'a', 'b', 'c'})

		result, err = ParseBytes([]byte("123abc"), comb)
		check(t, err)
		assertSlice(t, result, []byte{'a', 'b', 'c'})

		result, err = ParseBytes([]byte("123abc123"), comb)
		check(t, err)
		assertSlice(t, result, []byte{'a', 'b', 'c'})
	})

	t.Run("case 3", func(t *testing.T) {
		comb := Skip(
			NotEq[byte,int]('a'),
			Eq[byte,int]('a'),
		)

		result, err := ParseBytes([]byte("abc"), comb)
		assertError(t, err)
		assertEq(t, result, 0)
	})
}

func TestSkipAfter(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		comb := SkipAfter(
			Eq[byte,int]('b'),
			Eq[byte,int]('a'),
		)

		result, err := ParseBytes([]byte("abc"), comb)
		check(t, err)
		assertEq(t, result, byte('a'))

		result, err = ParseBytes([]byte("ab"), comb)
		check(t, err)
		assertEq(t, result, byte('a'))

		result, err = ParseBytes([]byte("a"), comb)
		assertError(t, err)
		assertEq(t, result, 0)
	})

	t.Run("case 2", func(t *testing.T) {
		comb := SkipAfter(
			Eq[byte,int]('b'),
			Satisfy[byte,int](true, Nothing[byte]),
		)

		result, err := ParseBytes([]byte("abc"), comb)
		assertError(t, err)
		assertEq(t, result, 0)
	})

	t.Run("case 3", func(t *testing.T) {
		comb := SkipAfter(
			Satisfy[byte,int](true, Nothing[byte]),
			Eq[byte,int]('a'),
		)

		result, err := ParseBytes([]byte("abc"), comb)
		assertError(t, err)
		assertEq(t, result, 0)
	})
}

func TestPadded(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		comb := Padded(
			Eq[byte,int]('.'),
			Range[byte,int]('0', '9'),
		)

		result, err := ParseBytes([]byte("1"), comb)
		check(t, err)
		assertEq(t, result, '1')

		result, err = ParseBytes([]byte(".1"), comb)
		check(t, err)
		assertEq(t, result, '1')

		result, err = ParseBytes([]byte("...1"), comb)
		check(t, err)
		assertEq(t, result, '1')

		result, err = ParseBytes([]byte("..1..."), comb)
		check(t, err)
		assertEq(t, result, '1')
	})
}

func TestEOF(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		result, err := ParseBytes([]byte("abcd"), EOF[byte,int]())
		check(t, err)
		assertEq(t, result, false)
	})

	t.Run("case 2", func(t *testing.T) {
		result, err := ParseBytes([]byte(""), EOF[byte,int]())
		check(t, err)
		assertEq(t, result, true)
	})
}

func TestCast(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		comb := Cast[byte, int, byte](
			Satisfy[byte,int](true, Anything[byte]),
			func(x byte) (int, error) { return int(x), nil },
		)

		result, err := ParseBytes([]byte("a"), comb)
		check(t, err)
		assertEq(t, result, 97)

		result, err = ParseBytes([]byte("b"), comb)
		check(t, err)
		assertEq(t, result, 98)

		result, err = ParseBytes([]byte(""), comb)
		assertError(t, err)
		assertEq(t, result, 0)
	})

	t.Run("case 2", func(t *testing.T) {
		comb := Cast[byte, int, byte](
			Any[byte,int](),
			func(x byte) (int, error) { return -1, fmt.Errorf("test error") },
		)

		result, err := ParseBytes([]byte("abc"), comb)
		assertError(t, err)
		assertEq(t, result, 0)
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

func check(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatal(err)
	}
}

func assertError(t *testing.T, err error) {
	t.Helper()

	if err == nil {
		t.Fatal("expected error")
	} else {
		t.Log("catch error: ", err)
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
