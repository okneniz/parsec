package strings

import (
	"testing"

	p "github.com/okneniz/parsec/common"
	. "github.com/okneniz/parsec/testing"
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

func TestSkip(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		comb := Skip(
			Optional(Eq('a'), 0),
			Eq('b'),
		)

		result, err := ParseString("abc", comb)
		Check(t, err)
		AssertEq(t, result, 'b')

		result, err = ParseString("cba", comb)
		Check(t, err)
		AssertEq(t, result, 'b')
	})

	t.Run("case 2", func(t *testing.T) {
		phrase := SequenceOf('a', 'b', 'c')
		noice := Many(0, Try(NoneOf('a', 'b', 'c')))
		comb := Skip(noice, phrase)

		result, err := ParseString("abc", comb)
		Check(t, err)
		AssertSlice(t, result, []rune{'a', 'b', 'c'})

		result, err = ParseString("abc123", comb)
		Check(t, err)
		AssertSlice(t, result, []rune{'a', 'b', 'c'})

		result, err = ParseString("123abc", comb)
		Check(t, err)
		AssertSlice(t, result, []rune{'a', 'b', 'c'})

		result, err = ParseString("123abc123", comb)
		Check(t, err)
		AssertSlice(t, result, []rune{'a', 'b', 'c'})
	})

	t.Run("case 3", func(t *testing.T) {
		comb := Skip(
			NotEq('a'),
			Eq('a'),
		)

		result, err := ParseString("abc", comb)
		AssertError(t, err)
		AssertEq(t, result, 0)
	})
}

func TestSkipAfter(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		comb := SkipAfter(
			Eq('b'),
			Eq('a'),
		)

		result, err := ParseString("abc", comb)
		Check(t, err)
		AssertEq(t, result, 'a')

		result, err = ParseString("ab", comb)
		Check(t, err)
		AssertEq(t, result, 'a')

		result, err = ParseString("a", comb)
		AssertError(t, err)
		AssertEq(t, result, 0)
	})

	t.Run("case 2", func(t *testing.T) {
		comb := SkipAfter(
			Eq('b'),
			Satisfy(true, p.Nothing[rune]),
		)

		result, err := ParseString("abc", comb)
		AssertError(t, err)
		AssertEq(t, result, 0)
	})

	t.Run("case 3", func(t *testing.T) {
		comb := SkipAfter(
			Satisfy(true, p.Nothing[rune]),
			Eq('a'),
		)

		result, err := ParseString("abc", comb)
		AssertError(t, err)
		AssertEq(t, result, 0)
	})
}
