package parsec

import (
	"testing"
)

func TestConcat(t *testing.T) {
	t.Run("case 1", func(t *testing.T) {
		comb := Concat(
			3,
			Count(1, Eq[byte](true, byte('a'))),
			Count(1, Eq[byte](true, byte('b'))),
			Count(1, NotEq[byte](true, byte('z'))),
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
			Count(1, Eq[byte](true, byte('a'))),
			Count(1, Eq[byte](true, byte('b'))),
			Count(1, Satisfy(true, Nothing[byte])),
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
			Count(1, Satisfy(true, Nothing[byte])),
			Count(1, Satisfy(true, Nothing[byte])),
			Count(3, Satisfy(true, Nothing[byte])),
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
			Eq[byte](true, byte('a')),
			Eq[byte](true, byte('b')),
			NotEq[byte](true, byte('z')),
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
			Eq[byte](true, byte('a')),
			Eq[byte](true, byte('b')),
			Satisfy(true, Nothing[byte]),
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
			Try(Eq(true, byte('a'))),
			Try(Eq(true, byte('b'))),
			Try(Eq(true, byte('c'))),
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
				Try(Eq(true, byte('a'))),
				Try(Eq(true, byte('b'))),
				Try(Eq(true, byte('c'))),
				Try(NotEq(true, byte('z'))),
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
				Satisfy[byte](true, Nothing[byte]),
				Satisfy[byte](true, Nothing[byte]),
			),
		)

		result, err := ParseBytes([]byte("abcd"), comb)
		check(t, err)
		assertSlice(t, result, nil)
	})
}
