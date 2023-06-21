package parsec

import (
	"testing"
)

func TestBuffer(t *testing.T) {
	b := BytesBuffer([]byte("foo"))

	assertEq(t, b.Position(), 0)
	assertEq(t, b.IsEOF(), false)

	x, ok := b.Read(false)

	assertEq(t, x, byte('f'))
	assertEq(t, ok, true)
	assertEq(t, b.Position(), 0)
	assertEq(t, b.IsEOF(), false)

	x, ok = b.Read(false)

	assertEq(t, x, byte('f'))
	assertEq(t, ok, true)
	assertEq(t, b.Position(), 0)
	assertEq(t, b.IsEOF(), false)

	x, ok = b.Read(true)

	assertEq(t, x, byte('f'))
	assertEq(t, ok, true)
	assertEq(t, b.Position(), 1)
	assertEq(t, b.IsEOF(), false)

	x, ok = b.Read(true)

	assertEq(t, x, byte('o'))
	assertEq(t, ok, true)
	assertEq(t, b.Position(), 2)
	assertEq(t, b.IsEOF(), false)

	x, ok = b.Read(true)

	assertEq(t, x, byte('o'))
	assertEq(t, ok, true)
	assertEq(t, b.Position(), 3)
	assertEq(t, b.IsEOF(), true)

	b.Seek(1)

	x, ok = b.Read(true)

	assertEq(t, x, byte('o'))
	assertEq(t, ok, true)
	assertEq(t, b.Position(), 2)
	assertEq(t, b.IsEOF(), false)

	x, ok = b.Read(true)

	assertEq(t, x, byte('o'))
	assertEq(t, ok, true)
	assertEq(t, b.Position(), 3)
	assertEq(t, b.IsEOF(), true)
}
