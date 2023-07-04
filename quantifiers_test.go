package parsec

import (
	"testing"
)

func TestMany(t *testing.T) {
	comb := Many(0, Eq[byte, int]('a'))

	result, err := ParseBytes([]byte("aaa"), comb)
	check(t, err)
	assertSlice(t, result, []byte("aaa"))

	result, err = ParseBytes([]byte("aaabc"), comb)
	check(t, err)
	assertSlice(t, result, []byte("aaa"))

	result, err = ParseBytes([]byte("xaaabc"), comb)
	check(t, err)
	assertSlice(t, result, []byte{})
}

func TestSome(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		comb := Some(0, Eq[byte, int]('a'))

		result, err := ParseBytes([]byte("aaa"), comb)
		check(t, err)
		assertSlice(t, result, []byte("aaa"))

		result, err = ParseBytes([]byte("aaabc"), comb)
		check(t, err)
		assertSlice(t, result, []byte("aaa"))

		result, err = ParseBytes([]byte("xaaabc"), comb)
		assertError(t, err)
		assertSlice(t, result, nil)
	})

	t.Run("case 2", func(t *testing.T) {
		comb := Some(
			0,
			Satisfy[byte,int](true, func(x byte) bool { return false }),
		)

		result, err := ParseBytes([]byte("abc"), comb)
		assertError(t, err)
		assertSlice(t, result, []byte{})
	})
}

func TestOptional(t *testing.T) {
	comb := Optional(Eq[byte, int]('a'), 0)

	result, err := ParseBytes([]byte("aaa"), comb)
	check(t, err)
	assertEq(t, result, byte('a'))

	result, err = ParseBytes([]byte("bcd"), comb)
	check(t, err)
	assertEq(t, result, 0)
}

func TestCount(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		comb := Count(2, Eq[byte, int]('a'))

		result, err := ParseBytes([]byte("aabbcc"), comb)
		check(t, err)
		assertSlice(t, result, []byte{'a', 'a'})

		result, err = ParseBytes([]byte("abbcc"), comb)
		assertError(t, err)
		assertSlice(t, result, nil)

		result, err = ParseBytes([]byte("bbaacc"), comb)
		assertError(t, err)
		assertSlice(t, result, nil)
	})

	t.Run("case 2", func(t *testing.T) {
		comb := Count(2, EOF[byte,int]())

		result, err := ParseBytes([]byte("aab"), comb)
		check(t, err)
		assertSlice(t, result, []bool{false, false})

		result, err = ParseBytes([]byte(""), comb)
		check(t, err)
		assertSlice(t, result, []bool{true, true})
	})
}
