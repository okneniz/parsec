package parsec

import (
	"testing"
)

func TestRange(t *testing.T) {
	comb := Range(true, byte('a'), byte('c'))

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

	result, ok = ParseBytes([]byte(""), comb)
	assert(t, !ok, "expected false")
	assertEq(t, result, 0)
}

func TestNotRange(t *testing.T) {
	comb := NotRange(true, byte('a'), byte('c'))

	result, ok := ParseBytes([]byte("a"), comb)
	assert(t, !ok, "expected false")
	assertEq(t, result, 0)

	result, ok = ParseBytes([]byte("b"), comb)
	assert(t, !ok, "expected false")
	assertEq(t, result, 0)

	result, ok = ParseBytes([]byte("c"), comb)
	assert(t, !ok, "expected false")
	assertEq(t, result, 0)

	result, ok = ParseBytes([]byte("d"), comb)
	assert(t, ok, "expected true")
	assertEq(t, result, byte('d'))

	result, ok = ParseBytes([]byte(""), comb)
	assert(t, !ok, "expected false")
	assertEq(t, result, 0)
}
