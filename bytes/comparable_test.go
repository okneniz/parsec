package bytes

import (
	. "git.sr.ht/~okneniz/parsec/testing"
	"testing"
)

func TestEq(t *testing.T) {
	comb := Eq('c')

	result, err := Parse([]byte("a"), comb)
	AssertError(t, err)
	AssertEq(t, result, 0)

	result, err = Parse([]byte("b"), comb)
	AssertError(t, err)
	AssertEq(t, result, 0)

	result, err = Parse([]byte("c"), comb)
	Check(t, err)
	AssertEq(t, result, 'c')
}

func TestNotEq(t *testing.T) {
	comb := NotEq('c')

	result, err := Parse([]byte("a"), comb)
	Check(t, err)
	AssertEq(t, result, byte('a'))

	result, err = Parse([]byte("b"), comb)
	Check(t, err)
	AssertEq(t, result, byte('b'))

	result, err = Parse([]byte("abc"), comb)
	Check(t, err)
	AssertEq(t, result, byte('a'))

	result, err = Parse([]byte("c"), comb)
	AssertError(t, err)
	AssertEq(t, result, 0)
}

func TestOneOf(t *testing.T) {
	comb := OneOf('a', 'b', 'c')

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
}

func TestSequenceOf(t *testing.T) {
	comb := SequenceOf('f', 'o', 'o')

	result, err := Parse([]byte("foo"), comb)
	Check(t, err)
	AssertSlice(t, result, []byte{'f', 'o', 'o'})

	result, err = Parse([]byte("foobar"), comb)
	Check(t, err)
	AssertSlice(t, result, []byte{'f', 'o', 'o'})

	result, err = Parse([]byte("fo"), comb)
	AssertError(t, err)
	AssertSlice(t, result, nil)

	result, err = Parse([]byte(" foobar"), comb)
	AssertError(t, err)
	AssertSlice(t, result, nil)

	result, err = Parse([]byte(" "), comb)
	AssertError(t, err)
	AssertSlice(t, result, nil)

	result, err = Parse([]byte(""), comb)
	AssertError(t, err)
	AssertSlice(t, result, nil)
}
