package bytes

import (
	"github.com/okneniz/parsec/common"
)

// Eq succeeds when the next byte is equal to t and returns it.
// Greedy: consumes the byte even on failure (see Try).
func Eq(
	errMessage string,
	t byte,
) common.Combinator[byte, int, byte] {
	return common.Eq[byte, int](errMessage, t)
}

// NotEq succeeds when the next byte is not equal to t and returns it.
// Greedy: consumes the byte even on failure (see Try).
func NotEq(
	errMessage string,
	t byte,
) common.Combinator[byte, int, byte] {
	return common.NotEq[byte, int](errMessage, t)
}

// OneOf succeeds when the next byte is one of data and returns it.
// Greedy: consumes the byte even on failure (see Try).
func OneOf(
	errMessage string,
	data ...byte,
) common.Combinator[byte, int, byte] {
	return common.OneOf[byte, int](errMessage, data...)
}

// NoneOf succeeds when the next byte is none of data and returns it.
// Greedy: consumes the byte even on failure (see Try).
func NoneOf(
	errMessage string,
	data ...byte,
) common.Combinator[byte, int, byte] {
	return common.NoneOf[byte, int](errMessage, data...)
}

// SequenceOf expects the next bytes to be equal to data in the same order
// and returns them as a slice.
func SequenceOf(
	errMessage string,
	data ...byte,
) common.Combinator[byte, int, []byte] {
	return common.SequenceOf[byte, int](errMessage, data...)
}

// Map parses a key with the c combinator, looks the key up in cases
// and returns the mapped value. It fails with errMessage
// when the key is not found.
func Map[K comparable, V any](
	errMessage string,
	cases map[K]V,
	c common.Combinator[byte, int, K],
) common.Combinator[byte, int, V] {
	return common.Map[byte, int, K, V](errMessage, cases, c)
}
