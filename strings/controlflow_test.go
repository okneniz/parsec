package strings

import (
	p "git.sr.ht/~okneniz/parsec/common"
	. "git.sr.ht/~okneniz/parsec/testing"
	"testing"
)

func TestConcat(t *testing.T) {
	t.Run("case 1", func(t *testing.T) {
		comb := Concat(
			3,
			Count(1, Eq('a')),
			Count(1, Eq('b')),
			Count(1, NotEq('z')),
		)

		result, err := ParseString("abc", comb)
		Check(t, err)
		AssertSlice(t, result, []rune("abc"))

		result, err = ParseString("abd", comb)
		Check(t, err)
		AssertSlice(t, result, []rune("abd"))

		result, err = ParseString("abdasdasd", comb)
		Check(t, err)
		AssertSlice(t, result, []rune("abd"))

		result, err = ParseString("xyz", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)

		result, err = ParseString("", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})

	t.Run("case 2", func(t *testing.T) {
		comb := Concat(
			3,
			Count(1, Eq('a')),
			Count(1, Eq('b')),
			Count(1, Satisfy(true, p.Nothing[rune])),
		)

		result, err := ParseString("abc", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)

		result, err = ParseString("", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})

	t.Run("case 3", func(t *testing.T) {
		comb := Concat(
			0,
			Count(1, Satisfy(true, p.Nothing[rune])),
			Count(1, Satisfy(true, p.Nothing[rune])),
			Count(3, Satisfy(true, p.Nothing[rune])),
		)

		result, err := ParseString("abc", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)

		result, err = ParseString("", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})
}

func TestSequence(t *testing.T) {
	t.Run("case 1", func(t *testing.T) {
		comb := Sequence(
			3,
			Eq('a'),
			Eq('b'),
			NotEq('z'),
		)

		result, err := ParseString("abc", comb)
		Check(t, err)
		AssertSlice(t, result, []rune("abc"))

		result, err = ParseString("abd", comb)
		Check(t, err)
		AssertSlice(t, result, []rune("abd"))

		result, err = ParseString("abdasdasd", comb)
		Check(t, err)
		AssertSlice(t, result, []rune("abd"))

		result, err = ParseString("xyz", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)

		result, err = ParseString("", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})

	t.Run("case 2", func(t *testing.T) {
		comb := Sequence(
			3,
			Eq(rune('a')),
			Eq(rune('b')),
			Satisfy(true, p.Nothing[rune]),
		)

		result, err := ParseString("abc", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)

		result, err = ParseString("", comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})
}

func TestChoice(t *testing.T) {
	t.Run("case 1", func(t *testing.T) {
		comb := Choice(
			Try(Eq('a')),
			Try(Eq('b')),
			Try(Eq('c')),
		)

		result, err := ParseString("a", comb)
		Check(t, err)
		AssertEq(t, result, 'a')

		result, err = ParseString("b", comb)
		Check(t, err)
		AssertEq(t, result, 'b')

		result, err = ParseString("c", comb)
		Check(t, err)
		AssertEq(t, result, 'c')
	})

	t.Run("case 2", func(t *testing.T) {
		comb := Many(
			4,
			Choice(
				Try(Eq('a')),
				Try(Eq('b')),
				Try(Eq('c')),
				Try(NotEq('z')),
			),
		)

		result, err := ParseString("abcd", comb)
		Check(t, err)
		AssertSlice(t, result, []rune("abcd"))
	})

	t.Run("case 3", func(t *testing.T) {
		comb := Many(
			0,
			Choice(
				Satisfy(true, p.Nothing[rune]),
				Satisfy(true, p.Nothing[rune]),
			),
		)

		result, err := ParseString("abcd", comb)
		Check(t, err)
		AssertSlice(t, result, nil)
	})
}
