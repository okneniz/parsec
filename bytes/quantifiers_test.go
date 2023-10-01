package bytes

import (
	"testing"

	. "github.com/okneniz/parsec/testing"
)

func TestMany(t *testing.T) {
	comb := Many(0, Eq('a'))

	result, err := Parse([]byte("aaa"), comb)
	Check(t, err)
	AssertSlice(t, result, []byte("aaa"))

	result, err = Parse([]byte("aaabc"), comb)
	Check(t, err)
	AssertSlice(t, result, []byte("aaa"))

	result, err = Parse([]byte("xaaabc"), comb)
	Check(t, err)
	AssertSlice(t, result, []byte{})
}

func TestSome(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		comb := Some(0, Eq('a'))

		result, err := Parse([]byte("aaa"), comb)
		Check(t, err)
		AssertSlice(t, result, []byte("aaa"))

		result, err = Parse([]byte("aaabc"), comb)
		Check(t, err)
		AssertSlice(t, result, []byte("aaa"))

		result, err = Parse([]byte("xaaabc"), comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})

	t.Run("case 2", func(t *testing.T) {
		comb := Some(
			0,
			Satisfy(true, func(x byte) bool { return false }),
		)

		result, err := Parse([]byte("abc"), comb)
		AssertError(t, err)
		AssertSlice(t, result, []byte{})
	})
}

func TestOptional(t *testing.T) {
	comb := Optional(Eq('a'), 0)

	result, err := Parse([]byte("aaa"), comb)
	Check(t, err)
	AssertEq(t, result, byte('a'))

	result, err = Parse([]byte("bcd"), comb)
	Check(t, err)
	AssertEq(t, result, 0)
}

func TestCount(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		comb := Count(2, Eq('a'))

		result, err := Parse([]byte("aabbcc"), comb)
		Check(t, err)
		AssertSlice(t, result, []byte{'a', 'a'})

		result, err = Parse([]byte("abbcc"), comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)

		result, err = Parse([]byte("bbaacc"), comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})

	t.Run("case 2", func(t *testing.T) {
		comb := Count(2, EOF())

		result, err := Parse([]byte("aab"), comb)
		Check(t, err)
		AssertSlice(t, result, []bool{false, false})

		result, err = Parse([]byte(""), comb)
		Check(t, err)
		AssertSlice(t, result, []bool{true, true})
	})
}
