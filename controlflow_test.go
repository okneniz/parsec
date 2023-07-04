package parsec

import (
	"testing"
)

func TestConcat(t *testing.T) {
	t.Run("case 1", func(t *testing.T) {
		comb := Concat(
			3,
			Count(1, Eq[byte, int]('a')),
			Count(1, Eq[byte, int]('b')),
			Count(1, NotEq[byte, int]('z')),
		)

		result, err := ParseBytes([]byte("abc"), comb)
		check(t, err)
		assertSlice(t, result, []byte("abc"))

		result, err = ParseBytes([]byte("abd"), comb)
		check(t, err)
		assertSlice(t, result, []byte("abd"))

		result, err = ParseBytes([]byte("abdasdasd"), comb)
		check(t, err)
		assertSlice(t, result, []byte("abd"))

		result, err = ParseBytes([]byte("xyz"), comb)
		assertError(t, err)
		assertSlice(t, result, nil)

		result, err = ParseBytes([]byte(""), comb)
		assertError(t, err)
		assertSlice(t, result, nil)
	})

	t.Run("case 2", func(t *testing.T) {
		comb := Concat(
			3,
			Count(1, Eq[byte, int]('a')),
			Count(1, Eq[byte, int]('b')),
			Count(1, Satisfy[byte,int](true, Nothing[byte])),
		)

		result, err := ParseBytes([]byte("abc"), comb)
		assertError(t, err)
		assertSlice(t, result, nil)

		result, err = ParseBytes([]byte(""), comb)
		assertError(t, err)
		assertSlice(t, result, nil)
	})

	t.Run("case 3", func(t *testing.T) {
		comb := Concat(
			0,
			Count(1, Satisfy[byte,int](true, Nothing[byte])),
			Count(1, Satisfy[byte,int](true, Nothing[byte])),
			Count(3, Satisfy[byte,int](true, Nothing[byte])),
		)

		result, err := ParseBytes([]byte("abc"), comb)
		assertError(t, err)
		assertSlice(t, result, nil)

		result, err = ParseBytes([]byte(""), comb)
		assertError(t, err)
		assertSlice(t, result, nil)
	})
}

func TestSequence(t *testing.T) {
	t.Run("case 1", func(t *testing.T) {
		comb := Sequence(
			3,
			Eq[byte,int]('a'),
			Eq[byte, int]('b'),
			NotEq[byte, int]('z'),
		)

		result, err := ParseBytes([]byte("abc"), comb)
		check(t, err)
		assertSlice(t, result, []byte("abc"))

		result, err = ParseBytes([]byte("abd"), comb)
		check(t, err)
		assertSlice(t, result, []byte("abd"))

		result, err = ParseBytes([]byte("abdasdasd"), comb)
		check(t, err)
		assertSlice(t, result, []byte("abd"))

		result, err = ParseBytes([]byte("xyz"), comb)
		assertError(t, err)
		assertSlice(t, result, nil)

		result, err = ParseBytes([]byte(""), comb)
		assertError(t, err)
		assertSlice(t, result, nil)
	})

	t.Run("case 2", func(t *testing.T) {
		comb := Sequence(
			3,
			Eq[byte, int](byte('a')),
			Eq[byte, int](byte('b')),
			Satisfy[byte,int](true, Nothing[byte]),
		)

		result, err := ParseBytes([]byte("abc"), comb)
		assertError(t, err)
		assertSlice(t, result, nil)

		result, err = ParseBytes([]byte(""), comb)
		assertError(t, err)
		assertSlice(t, result, nil)
	})
}

func TestChoice(t *testing.T) {
	t.Run("case 1", func(t *testing.T) {
		comb := Choice(
			Try(Eq[byte,int]('a')),
			Try(Eq[byte,int]('b')),
			Try(Eq[byte,int]('c')),
		)

		result, err := ParseBytes([]byte("a"), comb)
		check(t, err)
		assertEq(t, result, byte('a'))

		result, err = ParseBytes([]byte("b"), comb)
		check(t, err)
		assertEq(t, result, byte('b'))

		result, err = ParseBytes([]byte("c"), comb)
		check(t, err)
		assertEq(t, result, byte('c'))
	})

	t.Run("case 2", func(t *testing.T) {
		comb := Many(
			4,
			Choice(
				Try(Eq[byte,int]('a')),
				Try(Eq[byte,int]('b')),
				Try(Eq[byte,int]('c')),
				Try(NotEq[byte,int]('z')),
			),
		)

		result, err := ParseBytes([]byte("abcd"), comb)
		check(t, err)
		assertSlice(t, result, []byte{
			byte('a'),
			byte('b'),
			byte('c'),
			byte('d'),
		})
	})

	t.Run("case 3", func(t *testing.T) {
		comb := Many(
			0,
			Choice(
				Satisfy[byte, int](true, Nothing[byte]),
				Satisfy[byte, int](true, Nothing[byte]),
			),
		)

		result, err := ParseBytes([]byte("abcd"), comb)
		check(t, err)
		assertSlice(t, result, nil)
	})
}
