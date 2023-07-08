package strings

import (
	"testing"

	. "git.sr.ht/~okneniz/parsec/testing"
)

func TestEq(t *testing.T) {
	comb := Eq('c')

	result, err := ParseString("a", comb)
	AssertError(t, err)
	AssertEq(t, result, 0)

	result, err = ParseString("b", comb)
	AssertError(t, err)
	AssertEq(t, result, 0)

	result, err = ParseString("c", comb)
	Check(t, err)
	AssertEq(t, result, 'c')
}

func TestNotEq(t *testing.T) {
	comb := NotEq('c')

	result, err := ParseString("a", comb)
	Check(t, err)
	AssertEq(t, result, 'a')

	result, err = ParseString("b", comb)
	Check(t, err)
	AssertEq(t, result, 'b')

	result, err = ParseString("abc", comb)
	Check(t, err)
	AssertEq(t, result, 'a')

	result, err = ParseString("c", comb)
	AssertError(t, err)
	AssertEq(t, result, 0)
}

func TestOneOf(t *testing.T) {
	comb := OneOf('a', 'b', 'c')

	result, err := ParseString("a", comb)
	Check(t, err)
	AssertEq(t, result, 'a')

	result, err = ParseString("b", comb)
	Check(t, err)
	AssertEq(t, result, 'b')

	result, err = ParseString("c", comb)
	Check(t, err)
	AssertEq(t, result, 'c')

	result, err = ParseString("d", comb)
	AssertError(t, err)
	AssertEq(t, result, 0)
}

func TestSequenceOf(t *testing.T) {
	comb := SequenceOf('f', 'o', 'o')

	result, err := ParseString("foo", comb)
	Check(t, err)
	AssertSlice(t, result, []rune{'f', 'o', 'o'})

	result, err = ParseString("foobar", comb)
	Check(t, err)
	AssertSlice(t, result, []rune{'f', 'o', 'o'})

	result, err = ParseString("fo", comb)
	AssertError(t, err)
	AssertSlice(t, result, nil)

	result, err = ParseString(" foobar", comb)
	AssertError(t, err)
	AssertSlice(t, result, nil)

	result, err = ParseString(" ", comb)
	AssertError(t, err)
	AssertSlice(t, result, nil)

	result, err = ParseString("", comb)
	AssertError(t, err)
	AssertSlice(t, result, nil)
}

func TestMap(t *testing.T) {
	cases := map[rune]int{'a': 1, 'b': 2, 'c': 3}
	noice := Try(NoneOf('a', 'b', 'c'))

	comb := Some(
		1,
		Skip(
			Many(0, noice),
			Map(cases, Any()),
		),
	)

	result, err := ParseString("a", comb)
	Check(t, err)
	AssertSlice(t, result, []int{1})

	result, err = ParseString("..a//b++c**d,,e--a", comb)
	Check(t, err)
	AssertSlice(t, result, []int{1, 2, 3, 1})

	result, err = ParseString("bb", comb)
	Check(t, err)
	AssertSlice(t, result, []int{2, 2})

	result, err = ParseString("", comb)
	AssertError(t, err)
	AssertSlice(t, result, nil)
}

func TestString(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		comb := String("foo")

		result, err := ParseString("foo", comb)
		Check(t, err)
		AssertEq(t, result, "foo")

		result, err = ParseString("foobar", comb)
		Check(t, err)
		AssertEq(t, result, "foo")

		result, err = ParseString("bar", comb)
		AssertError(t, err)
		AssertEq(t, result, "")

		result, err = ParseString("baz", comb)
		AssertError(t, err)
		AssertEq(t, result, "")

		result, err = ParseString(" foo", comb)
		AssertError(t, err)
		AssertEq(t, result, "")

		result, err = ParseString(" foobar", comb)
		AssertError(t, err)
		AssertEq(t, result, "")

		result, err = ParseString("", comb)
		AssertError(t, err)
		AssertEq(t, result, "")
	})
}

func TestOneOfStrings(t *testing.T) {
	t.Parallel()

	t.Run("case 1", func(t *testing.T) {
		comb := OneOfStrings("foo", "bar")

		result, err := ParseString("foo", comb)
		Check(t, err)
		AssertEq(t, result, "foo")

		result, err = ParseString("foobar", comb)
		Check(t, err)
		AssertEq(t, result, "foo")

		result, err = ParseString("barbaz", comb)
		Check(t, err)
		AssertEq(t, result, "bar")

		result, err = ParseString("baz", comb)
		AssertError(t, err)
		AssertEq(t, result, "")

		result, err = ParseString("", comb)
		AssertError(t, err)
		AssertEq(t, result, "")
	})

	t.Run("case 2", func(t *testing.T) {
		noice := Many(0, Try(NoneOf('f', 'o', 'b', 'a', 'r')))

		comb := Many(
			0,
			Skip(
				noice,
				OneOfStrings("foo", "bar"),
			),
		)

		result, err := ParseString("foo", comb)
		Check(t, err)
		AssertSlice(t, result, []string{"foo"})

		result, err = ParseString("barfoo", comb)
		Check(t, err)
		AssertSlice(t, result, []string{"bar", "foo"})

		result, err = ParseString("bar12334foo123", comb)
		Check(t, err)
		AssertSlice(t, result, []string{"bar", "foo"})

		result, err = ParseString("bar12334foo123baz", comb)
		Check(t, err)
		AssertSlice(t, result, []string{"bar", "foo"})

		result, err = ParseString("12311231820398", comb)
		Check(t, err)
		AssertSlice(t, result, nil)
	})
}
