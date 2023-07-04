package bytes

import (
	. "git.sr.ht/~okneniz/parsec/testing"
	"testing"
)

func TestRange(t *testing.T) {
	comb := Range('a', 'c')

	result, err := Parse([]byte("a"), comb)
	Check(t, err)
	AssertEq(t, result, byte('a'))

	result, err = Parse([]byte("b"), comb)
	Check(t, err)
	AssertEq(t, result, byte('b'))

	result, err = Parse([]byte("c"), comb)
	Check(t, err)
	AssertEq(t, result, byte('c'))

	result, err = Parse([]byte("d"), comb)
	AssertError(t, err)
	AssertEq(t, result, 0)

	result, err = Parse([]byte(""), comb)
	AssertError(t, err)
	AssertEq(t, result, 0)
}

func TestNotRange(t *testing.T) {
	comb := NotRange('a', 'c')

	result, err := Parse([]byte("a"), comb)
	AssertError(t, err)
	AssertEq(t, result, 0)

	result, err = Parse([]byte("b"), comb)
	AssertError(t, err)
	AssertEq(t, result, 0)

	result, err = Parse([]byte("c"), comb)
	AssertError(t, err)
	AssertEq(t, result, 0)

	result, err = Parse([]byte("d"), comb)
	Check(t, err)
	AssertEq(t, result, byte('d'))

	result, err = Parse([]byte(""), comb)
	AssertError(t, err)
	AssertEq(t, result, 0)
}
