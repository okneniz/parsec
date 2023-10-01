package bytes

import (
	"testing"

	"github.com/okneniz/parsec/common"
	. "github.com/okneniz/parsec/testing"
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

func TestMap(t *testing.T) {
	cases := map[byte]string{0: "foo", 1: "bar", 2: "baz"}

	comb := Some(
		1,
		common.SkipMany(
			NoneOf(0, 1, 2),
			Map(cases, Any()),
		),
	)

	result, err := Parse([]byte{1}, comb)
	Check(t, err)
	AssertSlice(t, result, []string{"bar"})

	result, err = Parse([]byte{10, 0, 4, 1, 6, 8, 2, 5, 5, 3}, comb)
	Check(t, err)
	AssertSlice(t, result, []string{"foo", "bar", "baz"})

	result, err = Parse([]byte{10}, comb)
	AssertError(t, err)
	AssertSlice(t, result, nil)

	result, err = Parse([]byte{}, comb)
	AssertError(t, err)
	AssertSlice(t, result, nil)
}
