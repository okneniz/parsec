package parsec

import (
	"testing"
)

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
