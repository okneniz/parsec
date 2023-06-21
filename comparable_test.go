package parsec

import (
	"testing"
)

func TestEq(t *testing.T) {
	c := byte('c')

	comb := Eq[byte](true, c)

	result, ok := ParseBytes([]byte("a"), comb)
	assert(t, !ok, "expected false")
	assertEq(t, result, 0)

	result, ok = ParseBytes([]byte("b"), comb)
	assert(t, !ok, "expected false")
	assertEq(t, result, 0)

	result, ok = ParseBytes([]byte("c"), comb)
	assert(t, ok, "expected true")
	assertEq(t, result, c)
}

func TestNotEq(t *testing.T) {
	c := byte('c')

	comb := NotEq[byte](true, c)

	result, ok := ParseBytes([]byte("a"), comb)
	assert(t, ok, "expected true")
	assertEq(t, result, byte('a'))

	result, ok = ParseBytes([]byte("b"), comb)
	assert(t, ok, "expected true")
	assertEq(t, result, byte('b'))

	result, ok = ParseBytes([]byte("abc"), comb)
	assert(t, ok, "expected true")
	assertEq(t, result, byte('a'))

	result, ok = ParseBytes([]byte("c"), comb)
	assert(t, !ok, "expected false")
	assertEq(t, result, 0)
}

func TestOneOf(t *testing.T) {
	comb := OneOf(true, byte('a'), byte('b'), byte('c'))

	result, ok := ParseBytes([]byte("a"), comb)
	assert(t, ok, "expected true")
	assertEq(t, result, byte('a'))

	result, ok = ParseBytes([]byte("b"), comb)
	assert(t, ok, "expected true")
	assertEq(t, result, byte('b'))

	result, ok = ParseBytes([]byte("c"), comb)
	assert(t, ok, "expected true")
	assertEq(t, result, byte('c'))

	result, ok = ParseBytes([]byte("d"), comb)
	assert(t, !ok, "expected false")
	assertEq(t, result, 0)
}
