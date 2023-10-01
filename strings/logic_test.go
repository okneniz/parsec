package strings

import (
	"testing"

	. "github.com/okneniz/parsec/testing"
)

func TestOr(t *testing.T) {
	comb := Or(Try(Eq('a')), Eq('b'))

	result, err := ParseString("a", comb)
	Check(t, err)
	AssertEq(t, result, 'a')

	result, err = ParseString("b", comb)
	Check(t, err)
	AssertEq(t, result, 'b')

	result, err = ParseString("c", comb)
	AssertError(t, err)
	AssertEq(t, result, 0)
}

func TestAnd(t *testing.T) {
	comb := And(Eq('a'), Eq('b'), func(x, y rune) []rune {
		return []rune{x, y}
	})

	result, err := ParseString("abc", comb)
	Check(t, err)
	AssertSlice(t, result, []rune{'a', 'b'})

	result, err = ParseString("bca", comb)
	AssertError(t, err)
	AssertSlice(t, result, nil)

	result, err = ParseString("acb", comb)
	AssertError(t, err)
	AssertSlice(t, result, nil)
}
