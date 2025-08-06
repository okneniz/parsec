package bytes

import (
	"testing"

	"github.com/okneniz/parsec/common"
	. "github.com/okneniz/parsec/testing"
)

func TestConcat(t *testing.T) {
	t.Run("case 1", func(t *testing.T) {
		comb := Concat(
			3,
			Count(1, "expected 'a'", Eq("expected 'a'", 'a')),
			Count(1, "expected 'b'", Eq("expected 'b'", 'b')),
			Count(1, "expected 'z'", NotEq("expected 'z'", 'z')),
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
			Count(1, "expected 'a'", Eq("expected 'a'", 'a')),
			Count(1, "expected 'b'", Eq("expected 'a'", 'b')),
			Count(1, "expected any byte", Satisfy("anything", true, common.Nothing[byte])),
		)

		result, err := Parse([]byte("abc"), comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)

		result, err = Parse([]byte(""), comb)
		AssertError(t, err)
		AssertSlice(t, result, nil)
	})

	t.Run("case 3", func(t *testing.T) {
		oneByte := Satisfy("test", true, common.Nothing[byte])

		comb := Concat(
			0,
			Count(1, "expected 1 byte", oneByte),
			Count(1, "expected 1 byte", oneByte),
			Count(3, "expected 3 bytes", oneByte),
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
			Eq("expected 'a'", 'a'),
			Eq("expected 'b'", 'b'),
			NotEq("expected 'z'", 'z'),
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
			Eq("expected 'a'", byte('a')),
			Eq("expected 'b'", byte('b')),
			Satisfy("expected at least one byte", true, common.Nothing[byte]),
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
			Try(Eq("expected 'a'", 'a')),
			Try(Eq("expected 'b'", 'b')),
			Try(Eq("expected 'c'", 'c')),
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
				Try(Eq("expected 'a'", 'a')),
				Try(Eq("expected 'b'", 'b')),
				Try(Eq("expected 'c'", 'c')),
				Try(NotEq("expected not 'z'", 'z')),
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
				Satisfy("expected at least one byte", true, common.Nothing[byte]),
				Satisfy("expected at least one byte", true, common.Nothing[byte]),
			),
		)

		result, err := Parse([]byte("abcd"), comb)
		Check(t, err)
		AssertSlice(t, result, nil)
	})
}
