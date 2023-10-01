package bytes

import (
	"testing"

	p "github.com/okneniz/parsec/common"
	. "github.com/okneniz/parsec/testing"
)

func TestConcat(t *testing.T) {
	t.Run("case 1", func(t *testing.T) {
		comb := Concat(
			3,
			Count(1, Eq('a')),
			Count(1, Eq('b')),
			Count(1, NotEq('z')),
		)

		result, err := Parse([]byte("abc"), comb)
		Check(t, err)
		AssertSlice(t, result, []byte("abc"))

		result, err = Parse([]byte("abd"), comb)
		Check(t, err)
		AssertSlice(t, result, []byte("abd"))

		result, err = Parse([]byte("abdasdasd"), comb)
		Check(t, err)
		AssertSlice(t, result, []byte("abd"))

		result, err = Parse([]byte("xyz"), comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)

		result, err = Parse([]byte(""), comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})

	t.Run("case 2", func(t *testing.T) {
		comb := Concat(
			3,
			Count(1, Eq('a')),
			Count(1, Eq('b')),
			Count(1, Satisfy(true, p.Nothing[byte])),
		)

		result, err := Parse([]byte("abc"), comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)

		result, err = Parse([]byte(""), comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})

	t.Run("case 3", func(t *testing.T) {
		comb := Concat(
			0,
			Count(1, Satisfy(true, p.Nothing[byte])),
			Count(1, Satisfy(true, p.Nothing[byte])),
			Count(3, Satisfy(true, p.Nothing[byte])),
		)

		result, err := Parse([]byte("abc"), comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)

		result, err = Parse([]byte(""), comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})
}

func TestSequence(t *testing.T) {
	t.Run("case 1", func(t *testing.T) {
		comb := Sequence(
			3,
			Eq('a'),
			Eq('b'),
			NotEq('z'),
		)

		result, err := Parse([]byte("abc"), comb)
		Check(t, err)
		AssertSlice(t, result, []byte("abc"))

		result, err = Parse([]byte("abd"), comb)
		Check(t, err)
		AssertSlice(t, result, []byte("abd"))

		result, err = Parse([]byte("abdasdasd"), comb)
		Check(t, err)
		AssertSlice(t, result, []byte("abd"))

		result, err = Parse([]byte("xyz"), comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)

		result, err = Parse([]byte(""), comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})

	t.Run("case 2", func(t *testing.T) {
		comb := Sequence(
			3,
			Eq(byte('a')),
			Eq(byte('b')),
			Satisfy(true, p.Nothing[byte]),
		)

		result, err := Parse([]byte("abc"), comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)

		result, err = Parse([]byte(""), comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})
}

func TestChoice(t *testing.T) {
	t.Run("case 1", func(t *testing.T) {
		comb := Choice(
			Try(Eq('a')),
			Try(Eq('b')),
			Try(Eq('c')),
		)

		result, err := Parse([]byte("a"), comb)
		Check(t, err)
		AssertEq(t, result, byte('a'))

		result, err = Parse([]byte("b"), comb)
		Check(t, err)
		AssertEq(t, result, byte('b'))

		result, err = Parse([]byte("c"), comb)
		Check(t, err)
		AssertEq(t, result, byte('c'))
	})

	t.Run("case 2", func(t *testing.T) {
		comb := Many(
			4,
			Choice(
				Try(Eq('a')),
				Try(Eq('b')),
				Try(Eq('c')),
				Try(NotEq('z')),
			),
		)

		result, err := Parse([]byte("abcd"), comb)
		Check(t, err)
		AssertSlice(t, result, []byte{
			byte('a'),
			byte('b'),
			byte('c'),
			byte('d'),
		})
	})

	t.Run("case 3", func(t *testing.T) {
		comb := Many(
			0,
			Choice(
				Satisfy(true, p.Nothing[byte]),
				Satisfy(true, p.Nothing[byte]),
			),
		)

		result, err := Parse([]byte("abcd"), comb)
		Check(t, err)
		AssertSlice(t, result, nil)
	})
}
