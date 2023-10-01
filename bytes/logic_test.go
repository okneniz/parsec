package bytes

import (
	"testing"

	. "github.com/okneniz/parsec/testing"
)

func TestOr(t *testing.T) {
	comb := Or(
		Try(Eq('a')),
		Eq('b'),
	)

	result, err := Parse([]byte("a"), comb)
	Check(t, err)
	AssertEq(t, result, byte('a'))

	result, err = Parse([]byte("b"), comb)
	Check(t, err)
	AssertEq(t, result, byte('b'))

	result, err = Parse([]byte("c"), comb)
	AssertError(t, err)
	AssertEq(t, result, 0)
}

func TestAnd(t *testing.T) {
	comb := And(
		Eq('a'),
		Eq('b'),
		func(x, y byte) []byte { return []byte{x, y} },
	)

	result, err := Parse([]byte("abc"), comb)
	Check(t, err)
	AssertSlice(t, result, []byte{'a', 'b'})

	result, err = Parse([]byte("bca"), comb)
	AssertError(t, err)
	AssertSlice(t, result, nil)

	result, err = Parse([]byte("acb"), comb)
	AssertError(t, err)
	AssertSlice(t, result, nil)
}
