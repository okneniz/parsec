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

		result, ok = ParseBytes([]byte(""), comb)
		assert(t, !ok, "expected false")
		assertSlice(t, result, nil)
	})

	t.Run("case 2", func(t *testing.T) {
		comb := Concat(
			3,
			Count(1, Eq[byte](true, byte('a'))),
			Count(1, Eq[byte](true, byte('b'))),
			Count(1, Satisfy(true, Nothing[byte])),
		)

		result, ok := ParseBytes([]byte("abc"), comb)
		assert(t, !ok, "expected false")
		assertSlice(t, result, nil)

		result, ok = ParseBytes([]byte(""), comb)
		assert(t, !ok, "expected false")
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

		result, ok = ParseBytes([]byte(""), comb)
		assert(t, !ok, "expected false")
		assertSlice(t, result, nil)
	})

	t.Run("case 2", func(t *testing.T) {
		comb := Sequence(
			3,
			Eq[byte](true, byte('a')),
			Eq[byte](true, byte('b')),
			Satisfy(true, Nothing[byte]),
		)

		result, ok := ParseBytes([]byte("abc"), comb)
		assert(t, !ok, "expected false")
		assertSlice(t, result, nil)

		result, ok = ParseBytes([]byte(""), comb)
		assert(t, !ok, "expected false")
		assertSlice(t, result, nil)
	})
}
