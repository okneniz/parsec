package strings

import (
	"testing"

	. "github.com/okneniz/parsec/testing"
)

func TestRange(t *testing.T) {
	comb := Range('a', 'c')

	result, err := ParseString("a", comb)
	Check(t, err)
	AssertEq(t, result, 'a')

	result, err = ParseString("b", comb)
	Check(t, err)
	AssertEq(t, result, 'b')

	result, err = ParseString("c", comb)
	Check(t, err)
	AssertEq(t, result, 'c')

	result, err = ParseString("d", comb)
	AssertError(t, err)
	AssertEq(t, result, 0)

	result, err = ParseString("", comb)
	AssertError(t, err)
	AssertEq(t, result, 0)
}

func TestNotRange(t *testing.T) {
	comb := NotRange('a', 'c')

	result, err := ParseString("a", comb)
	AssertError(t, err)
	AssertEq(t, result, 0)

	result, err = ParseString("b", comb)
	AssertError(t, err)
	AssertEq(t, result, 0)

	result, err = ParseString("c", comb)
	AssertError(t, err)
	AssertEq(t, result, 0)

	result, err = ParseString("d", comb)
	Check(t, err)
	AssertEq(t, result, 'd')

	result, err = ParseString("", comb)
	AssertError(t, err)
	AssertEq(t, result, 0)
}

func TestGt(t *testing.T) {
	comb := Gt('c')

	result, err := ParseString("d", comb)
	Check(t, err)
	AssertEq(t, result, 'd')

	result, err = ParseString("e", comb)
	Check(t, err)
	AssertEq(t, result, 'e')

	result, err = ParseString("a", comb)
	AssertError(t, err)
	AssertEq(t, result, 0)

	result, err = ParseString("b", comb)
	AssertError(t, err)
	AssertEq(t, result, 0)

	result, err = ParseString("c", comb)
	AssertError(t, err)
	AssertEq(t, result, 0)
}

func TestGte(t *testing.T) {
	comb := Gte('c')

	result, err := ParseString("d", comb)
	Check(t, err)
	AssertEq(t, result, 'd')

	result, err = ParseString("e", comb)
	Check(t, err)
	AssertEq(t, result, 'e')

	result, err = ParseString("a", comb)
	AssertError(t, err)
	AssertEq(t, result, 0)

	result, err = ParseString("b", comb)
	AssertError(t, err)
	AssertEq(t, result, 0)

	result, err = ParseString("c", comb)
	Check(t, err)
	AssertEq(t, result, 'c')
}

func TestLt(t *testing.T) {
	comb := Lt('c')

	result, err := ParseString("a", comb)
	Check(t, err)
	AssertEq(t, result, 'a')

	result, err = ParseString("b", comb)
	Check(t, err)
	AssertEq(t, result, 'b')

	result, err = ParseString("c", comb)
	AssertError(t, err)
	AssertEq(t, result, 0)

	result, err = ParseString("d", comb)
	AssertError(t, err)
	AssertEq(t, result, 0)

	result, err = ParseString("e", comb)
	AssertError(t, err)
	AssertEq(t, result, 0)
}

func TestLte(t *testing.T) {
	comb := Lte('c')

	result, err := ParseString("a", comb)
	Check(t, err)
	AssertEq(t, result, 'a')

	result, err = ParseString("b", comb)
	Check(t, err)
	AssertEq(t, result, 'b')

	result, err = ParseString("c", comb)
	Check(t, err)
	AssertEq(t, result, 'c')

	result, err = ParseString("d", comb)
	AssertError(t, err)
	AssertEq(t, result, 0)

	result, err = ParseString("e", comb)
	AssertError(t, err)
	AssertEq(t, result, 0)
}
