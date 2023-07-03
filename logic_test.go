package parsec

import (
	"testing"
)

func TestOr(t *testing.T) {
	comb := Or(
		Try(Eq(byte('a'))),
		Eq(byte('b')),
	)

	result, err := ParseBytes([]byte("a"), comb)
	check(t, err)
	assertEq(t, result, byte('a'))

	result, err = ParseBytes([]byte("b"), comb)
	check(t, err)
	assertEq(t, result, byte('b'))

	result, err = ParseBytes([]byte("c"), comb)
	assertError(t, err)
	assertEq(t, result, 0)
}

func TestAnd(t *testing.T) {
	comb := And(
		Eq(byte('a')),
		Eq(byte('b')),
		func(x, y byte) []byte { return []byte{x, y} },
	)

	result, err := ParseBytes([]byte("abc"), comb)
	check(t, err)
	assertSlice(t, result, []byte{'a', 'b'})

	result, err = ParseBytes([]byte("bca"), comb)
	assertError(t, err)
	assertSlice(t, result, nil)

	result, err = ParseBytes([]byte("acb"), comb)
	assertError(t, err)
	assertSlice(t, result, nil)
}
