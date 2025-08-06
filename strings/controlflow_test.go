package strings

import (
	"testing"

	"github.com/okneniz/parsec/common"
	. "github.com/okneniz/parsec/testing"
)

func TestConcat(t *testing.T) {
	t.Run("case 1", func(t *testing.T) {
		comb := Concat(
			3,
			Count(1, "expected one 'a'", Eq("expected 'a'", 'a')),
			Count(1, "expected one 'b'", Eq("expected 'b'", 'b')),
			Count(1, "expected one 'z'", NotEq("expected 'z'", 'z')),
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
			Count(1, "expected one 'a'", Eq("expected 'a'", 'a')),
			Count(1, "expected one 'b'", Eq("expected 'b'", 'b')),
			Count(1, "expected one char", Satisfy("expected any char", true, common.Nothing[rune])),
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
			Count(1, "expected one char", Satisfy("expected any char", true, common.Nothing[rune])),
			Count(1, "expected one char", Satisfy("expected any char", true, common.Nothing[rune])),
			Count(3, "expected three char", Satisfy("expected any char", true, common.Nothing[rune])),
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
			Eq("expected 'a'", 'a'),
			Eq("expected 'b'", 'b'),
			NotEq("expected char not equal 'z'", 'z'),
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
			Eq("expected 'a'", 'a'),
			Eq("expected 'b'", 'b'),
			Satisfy("expected any char", true, common.Nothing[rune]),
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
			Try(Eq("expected 'a'", 'a')),
			Try(Eq("expected 'b'", 'b')),
			Try(Eq("expected 'c'", 'c')),
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
				Try(Eq("expected 'a'", 'a')),
				Try(Eq("expected 'b'", 'b')),
				Try(Eq("expected 'v'", 'c')),
				Try(NotEq("expected 'z'", 'z')),
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
				Satisfy("expected any char", true, common.Nothing[rune]),
				Satisfy("expected any char", true, common.Nothing[rune]),
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
			Optional(Eq("expected 'a'", 'a'), 0),
			Eq("expected 'b'", 'b'),
		)

		result, err := ParseString("abc", comb)
		Check(t, err)
		AssertEq(t, result, 'b')

		result, err = ParseString("cba", comb)
		Check(t, err)
		AssertEq(t, result, 'b')
	})

	t.Run("case 2", func(t *testing.T) {
		phrase := SequenceOf("expected 'abc'", 'a', 'b', 'c')
		noice := Many(0, Try(NoneOf("expected not 'a', 'b' or 'c'", 'a', 'b', 'c')))
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
			NotEq("expected 'a'", 'a'),
			Eq("expected 'a'", 'a'),
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
			Eq("expected 'b'", 'b'),
			Eq("expected 'a'", 'a'),
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
			Eq("expected 'b'", 'b'),
			Satisfy("expected char", true, common.Nothing[rune]),
		)

		result, err := ParseString("abc", comb)
		AssertError(t, err)
		AssertEq(t, result, 0)
	})

	t.Run("case 3", func(t *testing.T) {
		comb := SkipAfter(
			Satisfy("expected any char", true, common.Nothing[rune]),
			Eq("expected 'a'", 'a'),
		)

		result, err := ParseString("abc", comb)
		AssertError(t, err)
		AssertEq(t, result, 0)
	})
}
