package bytes

import (
	"github.com/okneniz/parsec"
)

// Eq succeeds when the next byte is equal to t and returns it.
// Greedy: consumes the byte even on failure (see Try).
func Eq[S any](
	errMessage string,
	t byte,
) parsec.Combinator[byte, int, byte, S] {
	return parsec.Eq[byte, int, S](errMessage, t)
}

// NotEq succeeds when the next byte is not equal to t and returns it.
// Greedy: consumes the byte even on failure (see Try).
func NotEq[S any](
	errMessage string,
	t byte,
) parsec.Combinator[byte, int, byte, S] {
	return parsec.NotEq[byte, int, S](errMessage, t)
}

// OneOf succeeds when the next byte is one of data and returns it.
// Greedy: consumes the byte even on failure (see Try).
func OneOf[S any](
	errMessage string,
	data ...byte,
) parsec.Combinator[byte, int, byte, S] {
	return parsec.OneOf[byte, int, S](errMessage, data...)
}

// NoneOf succeeds when the next byte is none of data and returns it.
// Greedy: consumes the byte even on failure (see Try).
func NoneOf[S any](
	errMessage string,
	data ...byte,
) parsec.Combinator[byte, int, byte, S] {
	return parsec.NoneOf[byte, int, S](errMessage, data...)
}

// SequenceOf expects the next bytes to be equal to data in the same order
// and returns them as a slice.
func SequenceOf[S any](
	errMessage string,
	data ...byte,
) parsec.Combinator[byte, int, []byte, S] {
	return parsec.SequenceOf[byte, int, S](errMessage, data...)
}

// Map parses a key with the c combinator, looks the key up in cases
// and returns the mapped value. It fails with errMessage
// when the key is not found.
func Map[K comparable, O any, S any](
	errMessage string,
	cases map[K]O,
	c parsec.Combinator[byte, int, K, S],
) parsec.Combinator[byte, int, O, S] {
	return parsec.Map[byte, int, K, O, S](errMessage, cases, c)
}
