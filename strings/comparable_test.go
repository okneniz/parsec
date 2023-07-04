package strings

import (
	. "git.sr.ht/~okneniz/parsec/testing"
	"testing"
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
