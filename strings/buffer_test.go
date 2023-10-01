package strings

import (
	"testing"

	. "github.com/okneniz/parsec/testing"
)

func TestBuffer(t *testing.T) {
	t.Run("case 1", func(t *testing.T) {
		b := Buffer([]rune("foo\nbar\nbaz"))

		AssertEq(t, b.Position().Line(), 0)
		AssertEq(t, b.Position().Column(), 0)
		AssertEq(t, b.IsEOF(), false)

		x, err := b.Read(false)
		Check(t, err)

		AssertEq(t, x, 'f')
		AssertEq(t, b.Position().Line(), 0)
		AssertEq(t, b.Position().Column(), 0)
		AssertEq(t, b.IsEOF(), false)

		x, err = b.Read(false)
		Check(t, err)

		AssertEq(t, x, 'f')
		AssertEq(t, b.Position().Line(), 0)
		AssertEq(t, b.Position().Column(), 0)
		AssertEq(t, b.IsEOF(), false)

		x, err = b.Read(true)
		Check(t, err)

		AssertEq(t, x, 'f')
		AssertEq(t, b.Position().Line(), 0)
		AssertEq(t, b.Position().Column(), 1)
		AssertEq(t, b.IsEOF(), false)

		x, err = b.Read(true)
		Check(t, err)

		AssertEq(t, x, 'o')
		AssertEq(t, b.Position().Line(), 0)
		AssertEq(t, b.Position().Column(), 2)
		AssertEq(t, b.IsEOF(), false)

		x, err = b.Read(true)
		Check(t, err)

		AssertEq(t, x, 'o')
		AssertEq(t, b.Position().Line(), 0)
		AssertEq(t, b.Position().Column(), 3)
		AssertEq(t, b.IsEOF(), false)

		x, err = b.Read(true)
		Check(t, err)

		AssertEq(t, x, '\n')
		AssertEq(t, b.Position().Line(), 1)
		AssertEq(t, b.Position().Column(), 0)
		AssertEq(t, b.IsEOF(), false)

		x, err = b.Read(true)
		Check(t, err)

		AssertEq(t, x, 'b')
		AssertEq(t, b.Position().Line(), 1)
		AssertEq(t, b.Position().Column(), 1)
		AssertEq(t, b.IsEOF(), false)

		x, err = b.Read(true)
		Check(t, err)

		AssertEq(t, x, 'a')
		AssertEq(t, b.Position().Line(), 1)
		AssertEq(t, b.Position().Column(), 2)
		AssertEq(t, b.IsEOF(), false)

		x, err = b.Read(true)
		Check(t, err)

		AssertEq(t, x, 'r')
		AssertEq(t, b.Position().Line(), 1)
		AssertEq(t, b.Position().Column(), 3)
		AssertEq(t, b.IsEOF(), false)

		x, err = b.Read(true)
		Check(t, err)

		AssertEq(t, x, '\n')
		AssertEq(t, b.Position().Line(), 2)
		AssertEq(t, b.Position().Column(), 0)
		AssertEq(t, b.IsEOF(), false)

		pos := b.Position()

		x, err = b.Read(true)
		Check(t, err)

		AssertEq(t, x, 'b')
		AssertEq(t, b.Position().Line(), 2)
		AssertEq(t, b.Position().Column(), 1)
		AssertEq(t, b.IsEOF(), false)

		x, err = b.Read(true)
		Check(t, err)

		AssertEq(t, x, 'a')
		AssertEq(t, b.Position().Line(), 2)
		AssertEq(t, b.Position().Column(), 2)
		AssertEq(t, b.IsEOF(), false)

		b.Seek(pos)

		x, err = b.Read(true)
		Check(t, err)

		AssertEq(t, x, 'b')
		AssertEq(t, b.Position().Line(), 2)
		AssertEq(t, b.Position().Column(), 1)
		AssertEq(t, b.IsEOF(), false)

		x, err = b.Read(true)
		Check(t, err)

		AssertEq(t, x, 'a')
		AssertEq(t, b.Position().Line(), 2)
		AssertEq(t, b.Position().Column(), 2)
		AssertEq(t, b.IsEOF(), false)

		x, err = b.Read(true)
		Check(t, err)

		AssertEq(t, x, 'z')
		AssertEq(t, b.Position().Line(), 2)
		AssertEq(t, b.Position().Column(), 3)
		AssertEq(t, b.IsEOF(), true)

		x, err = b.Read(true)
		AssertError(t, err)

		AssertEq(t, x, 0)
		AssertEq(t, b.Position().Line(), 2)
		AssertEq(t, b.Position().Column(), 3)
		AssertEq(t, b.IsEOF(), true)
	})

	t.Run("case 2", func(t *testing.T) {
		b := Buffer([]rune("12a3a"), 'a')

		x, err := b.Read(true)
		Check(t, err)

		AssertEq(t, x, '1')
		AssertEq(t, b.Position().Line(), 0)
		AssertEq(t, b.Position().Column(), 1)
		AssertEq(t, b.IsEOF(), false)

		x, err = b.Read(true)
		Check(t, err)

		AssertEq(t, x, '2')
		AssertEq(t, b.Position().Line(), 0)
		AssertEq(t, b.Position().Column(), 2)
		AssertEq(t, b.IsEOF(), false)

		x, err = b.Read(true)
		Check(t, err)

		AssertEq(t, x, 'a')
		AssertEq(t, b.Position().Line(), 1)
		AssertEq(t, b.Position().Column(), 0)
		AssertEq(t, b.IsEOF(), false)

		x, err = b.Read(true)
		Check(t, err)

		AssertEq(t, x, '3')
		AssertEq(t, b.Position().Line(), 1)
		AssertEq(t, b.Position().Column(), 1)
		AssertEq(t, b.IsEOF(), false)

		x, err = b.Read(true)
		Check(t, err)

		AssertEq(t, x, 'a')
		AssertEq(t, b.Position().Line(), 2)
		AssertEq(t, b.Position().Column(), 0)
		AssertEq(t, b.IsEOF(), true)
	})
}
