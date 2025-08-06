package bytes

import (
	"testing"

	. "github.com/okneniz/parsec/testing"
)

func TestRange(t *testing.T) {
	comb := Range("expected byte between 'a' and 'c'", 'a', 'c')

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
	comb := NotRange("expecte byte between 'a' and 'c'", 'a', 'c')

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

func TestGt(t *testing.T) {
	comb := Gt("expected byte greater than 'c'", 'c')

	result, err := Parse([]byte("d"), comb)
	Check(t, err)
	AssertEq(t, result, 'd')

	result, err = Parse([]byte("e"), comb)
	Check(t, err)
	AssertEq(t, result, 'e')

	result, err = Parse([]byte("a"), comb)
	AssertError(t, err)
	AssertEq(t, result, 0)

	result, err = Parse([]byte("b"), comb)
	AssertError(t, err)
	AssertEq(t, result, 0)

	result, err = Parse([]byte("c"), comb)
	AssertError(t, err)
	AssertEq(t, result, 0)
}

func TestGte(t *testing.T) {
	comb := Gte("expected byte greater or equal than 'c'", 'c')

	result, err := Parse([]byte("d"), comb)
	Check(t, err)
	AssertEq(t, result, 'd')

	result, err = Parse([]byte("e"), comb)
	Check(t, err)
	AssertEq(t, result, 'e')

	result, err = Parse([]byte("a"), comb)
	AssertError(t, err)
	AssertEq(t, result, 0)

	result, err = Parse([]byte("b"), comb)
	AssertError(t, err)
	AssertEq(t, result, 0)

	result, err = Parse([]byte("c"), comb)
	Check(t, err)
	AssertEq(t, result, 'c')
}

func TestLt(t *testing.T) {
	comb := Lt("expected byte less than 'c'", 'c')

	result, err := Parse([]byte("a"), comb)
	Check(t, err)
	AssertEq(t, result, 'a')

	result, err = Parse([]byte("b"), comb)
	Check(t, err)
	AssertEq(t, result, 'b')

	result, err = Parse([]byte("c"), comb)
	AssertError(t, err)
	AssertEq(t, result, 0)

	result, err = Parse([]byte("d"), comb)
	AssertError(t, err)
	AssertEq(t, result, 0)

	result, err = Parse([]byte("e"), comb)
	AssertError(t, err)
	AssertEq(t, result, 0)
}

func TestLte(t *testing.T) {
	comb := Lte("expecte byte less than or equal 'c'", 'c')

	result, err := Parse([]byte("a"), comb)
	Check(t, err)
	AssertEq(t, result, 'a')

	result, err = Parse([]byte("b"), comb)
	Check(t, err)
	AssertEq(t, result, 'b')

	result, err = Parse([]byte("c"), comb)
	Check(t, err)
	AssertEq(t, result, 'c')

	result, err = Parse([]byte("d"), comb)
	AssertError(t, err)
	AssertEq(t, result, 0)

	result, err = Parse([]byte("e"), comb)
	AssertError(t, err)
	AssertEq(t, result, 0)
}
