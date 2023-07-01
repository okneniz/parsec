package parsec

import (
	"testing"
)

func TestRange(t *testing.T) {
	comb := Range(true, byte('a'), byte('c'))

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

	result, err = ParseBytes([]byte(""), comb)
	assertError(t, err)
	assertEq(t, result, 0)
}

func TestNotRange(t *testing.T) {
	comb := NotRange(true, byte('a'), byte('c'))

	result, err := ParseBytes([]byte("a"), comb)
	assertError(t, err)
	assertEq(t, result, 0)

	result, err = ParseBytes([]byte("b"), comb)
	assertError(t, err)
	assertEq(t, result, 0)

	result, err = ParseBytes([]byte("c"), comb)
	assertError(t, err)
	assertEq(t, result, 0)

	result, err = ParseBytes([]byte("d"), comb)
	check(t, err)
	assertEq(t, result, byte('d'))

	result, err = ParseBytes([]byte(""), comb)
	assertError(t, err)
	assertEq(t, result, 0)
}
