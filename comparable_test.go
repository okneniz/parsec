package parsec

import (
	"testing"
)

func TestEq(t *testing.T) {
	c := byte('c')

	comb := Eq[byte, int](c)

	result, err := ParseBytes([]byte("a"), comb)
	assertError(t, err)
	assertEq(t, result, 0)

	result, err = ParseBytes([]byte("b"), comb)
	assertError(t, err)
	assertEq(t, result, 0)

	result, err = ParseBytes([]byte("c"), comb)
	check(t, err)
	assertEq(t, result, c)
}

func TestNotEq(t *testing.T) {
	c := byte('c')

	comb := NotEq[byte, int](c)

	result, err := ParseBytes([]byte("a"), comb)
	check(t, err)
	assertEq(t, result, byte('a'))

	result, err = ParseBytes([]byte("b"), comb)
	check(t, err)
	assertEq(t, result, byte('b'))

	result, err = ParseBytes([]byte("abc"), comb)
	check(t, err)
	assertEq(t, result, byte('a'))

	result, err = ParseBytes([]byte("c"), comb)
	assertError(t, err)
	assertEq(t, result, 0)
}

func TestOneOf(t *testing.T) {
	comb := OneOf[byte, int]('a', 'b', 'c')

	result, err := ParseBytes([]byte("a"), comb)
	check(t, err)
	assertEq(t, result, byte('a'))

	result, err = ParseBytes([]byte("b"), comb)
	check(t, err)
	assertEq(t, result, byte('b'))

	result, err = ParseBytes([]byte("c"), comb)
	check(t, err)
	assertEq(t, result, byte('c'))

	result, err = ParseBytes([]byte("d"), comb)
	assertError(t, err)
	assertEq(t, result, 0)
}

func TestSequenceOf(t *testing.T) {
	comb := SequenceOf[byte, int]('f', 'o', 'o')

	result, err := ParseBytes([]byte("foo"), comb)
	check(t, err)
	assertSlice(t, result, []byte{'f', 'o', 'o'})

	result, err = ParseBytes([]byte("foobar"), comb)
	check(t, err)
	assertSlice(t, result, []byte{'f', 'o', 'o'})

	result, err = ParseBytes([]byte("fo"), comb)
	assertError(t, err)
	assertSlice(t, result, nil)

	result, err = ParseBytes([]byte(" foobar"), comb)
	assertError(t, err)
	assertSlice(t, result, nil)

	result, err = ParseBytes([]byte(" "), comb)
	assertError(t, err)
	assertSlice(t, result, nil)

	result, err = ParseBytes([]byte(""), comb)
	assertError(t, err)
	assertSlice(t, result, nil)
}
