package parsec

import (
	"testing"
)

func TestBuffer(t *testing.T) {
	b := BytesBuffer([]byte("foo"))

	assertEq(t, b.Position(), 0)
	assertEq(t, b.IsEOF(), false)

	x, err := b.Read(false)
	check(t, err)

	assertEq(t, x, byte('f'))
	assertEq(t, b.Position(), 0)
	assertEq(t, b.IsEOF(), false)

	x, err = b.Read(false)
	check(t, err)

	assertEq(t, x, byte('f'))
	assertEq(t, b.Position(), 0)
	assertEq(t, b.IsEOF(), false)

	x, err = b.Read(true)
	check(t, err)

	assertEq(t, x, byte('f'))
	assertEq(t, b.Position(), 1)
	assertEq(t, b.IsEOF(), false)

	x, err = b.Read(true)
	check(t, err)

	assertEq(t, x, byte('o'))
	assertEq(t, b.Position(), 2)
	assertEq(t, b.IsEOF(), false)

	x, err = b.Read(true)
	check(t, err)

	assertEq(t, x, byte('o'))
	assertEq(t, b.Position(), 3)
	assertEq(t, b.IsEOF(), true)

	b.Seek(1)

	x, err = b.Read(true)
	check(t, err)

	assertEq(t, x, byte('o'))
	assertEq(t, b.Position(), 2)
	assertEq(t, b.IsEOF(), false)

	x, err = b.Read(true)
	check(t, err)

	assertEq(t, x, byte('o'))
	assertEq(t, b.Position(), 3)
	assertEq(t, b.IsEOF(), true)
}
