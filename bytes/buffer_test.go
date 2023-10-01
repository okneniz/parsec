package bytes

import (
	"testing"

	. "github.com/okneniz/parsec/testing"
)

func TestBuffer(t *testing.T) {
	b := Buffer([]byte("foo"))

	AssertEq(t, b.Position(), 0)
	AssertEq(t, b.IsEOF(), false)

	x, err := b.Read(false)
	Check(t, err)

	AssertEq(t, x, byte('f'))
	AssertEq(t, b.Position(), 0)
	AssertEq(t, b.IsEOF(), false)

	x, err = b.Read(false)
	Check(t, err)

	AssertEq(t, x, byte('f'))
	AssertEq(t, b.Position(), 0)
	AssertEq(t, b.IsEOF(), false)

	x, err = b.Read(true)
	Check(t, err)

	AssertEq(t, x, byte('f'))
	AssertEq(t, b.Position(), 1)
	AssertEq(t, b.IsEOF(), false)

	x, err = b.Read(true)
	Check(t, err)

	AssertEq(t, x, byte('o'))
	AssertEq(t, b.Position(), 2)
	AssertEq(t, b.IsEOF(), false)

	x, err = b.Read(true)
	Check(t, err)

	AssertEq(t, x, byte('o'))
	AssertEq(t, b.Position(), 3)
	AssertEq(t, b.IsEOF(), true)

	b.Seek(1)

	x, err = b.Read(true)
	Check(t, err)

	AssertEq(t, x, byte('o'))
	AssertEq(t, b.Position(), 2)
	AssertEq(t, b.IsEOF(), false)

	x, err = b.Read(true)
	Check(t, err)

	AssertEq(t, x, byte('o'))
	AssertEq(t, b.Position(), 3)
	AssertEq(t, b.IsEOF(), true)
}
