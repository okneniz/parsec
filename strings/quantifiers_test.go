package strings

import (
	"testing"

	. "github.com/okneniz/parsec/testing"
)

func TestMany(t *testing.T) {
	comb := Many(0, Eq('a'))

	result, err := ParseString("aaa", comb)
	Check(t, err)
	AssertSlice(t, result, []rune("aaa"))

	result, err = ParseString("aaabc", comb)
	Check(t, err)
	AssertSlice(t, result, []rune("aaa"))

	result, err = ParseString("xaaabc", comb)
	Check(t, err)
	AssertSlice(t, result, []rune{})
}

func TestSome(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		comb := Some(0, Eq('a'))

		result, err := ParseString("aaa", comb)
		Check(t, err)
		AssertSlice(t, result, []rune("aaa"))

		result, err = ParseString("aaabc", comb)
		Check(t, err)
		AssertSlice(t, result, []rune("aaa"))

		result, err = ParseString("xaaabc", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})

	t.Run("case 2", func(t *testing.T) {
		comb := Some(
			0,
			Satisfy(true, func(x rune) bool { return false }),
		)

		result, err := ParseString("abc", comb)
		AssertError(t, err)
		AssertSlice(t, result, []rune{})
	})
}

func TestOptional(t *testing.T) {
	comb := Optional(Eq('a'), 0)

	result, err := ParseString("aaa", comb)
	Check(t, err)
	AssertEq(t, result, rune('a'))

	result, err = ParseString("bcd", comb)
	Check(t, err)
	AssertEq(t, result, 0)
}

func TestCount(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		comb := Count(2, Eq('a'))

		result, err := ParseString("aabbcc", comb)
		Check(t, err)
		AssertSlice(t, result, []rune{'a', 'a'})

		result, err = ParseString("abbcc", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)

		result, err = ParseString("bbaacc", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})

	t.Run("case 2", func(t *testing.T) {
		comb := Count(2, EOF())

		result, err := ParseString("aab", comb)
		Check(t, err)
		AssertSlice(t, result, []bool{false, false})

		result, err = ParseString("", comb)
		Check(t, err)
		AssertSlice(t, result, []bool{true, true})
	})
}
