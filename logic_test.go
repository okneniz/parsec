package parsec

import (
	"testing"
)

func TestOr(t *testing.T) {
	comb := Or(
		Try(Eq(true, byte('a'))),
		Eq(true, byte('b')),
	)

	result, ok := ParseBytes([]byte("a"), comb)
	assert(t, ok, "expected true")
	assertEq(t, result, byte('a'))

	result, ok = ParseBytes([]byte("b"), comb)
	assert(t, ok, "expected true")
	assertEq(t, result, byte('b'))

	result, ok = ParseBytes([]byte("c"), comb)
	assert(t, !ok, "expected false")
	assertEq(t, result, 0)
}

func TestAnd(t *testing.T) {
	comb := And(
		Eq(true, byte('a')),
		Eq(true, byte('b')),
		func(x, y byte) []byte { return []byte{x, y} },
	)

	result, ok := ParseBytes([]byte("abc"), comb)
	assert(t, ok, "expected true")
	assertSlice(t, result, []byte{'a', 'b'})

	result, ok = ParseBytes([]byte("bca"), comb)
	assert(t, !ok, "expected false")
	assertSlice(t, result, nil)

	result, ok = ParseBytes([]byte("acb"), comb)
	assert(t, !ok, "expected false")
	assertSlice(t, result, nil)
}
