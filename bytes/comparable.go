package bytes

import (
	"github.com/okneniz/parsec/common"
)

// Eq - succeeds for any byte which equal input t.
// Returns the byte that is actually readed from input buffer.
// Greedy by default - keep position after reading.
func Eq(
	errMessage string,
	t byte,
) common.Combinator[byte, int, byte] {
	return common.Eq[byte, int](errMessage, t)
}

// NotEq - succeeds for any byte which not equal input t.
// Returns the item that is actually readed from input buffer.
// Greedy by default - keep position after reading.
func NotEq(
	errMessage string,
	t byte,
) common.Combinator[byte, int, byte] {
	return common.NotEq[byte, int](errMessage, t)
}

// OneOf - succeeds for any byte which included in input data.
// Returns the item that is actually readed from input buffer.
// Greedy by default - keep position after reading.
func OneOf(
	errMessage string,
	data ...byte,
) common.Combinator[byte, int, byte] {
	return common.OneOf[byte, int](errMessage, data...)
}

// NoneOf - succeeds for any byte which not included in input data.
// Returns the item that is actually readed from input buffer.
// Greedy by default - keep position after reading.
func NoneOf(
	errMessage string,
	data ...byte,
) common.Combinator[byte, int, byte] {
	return common.NoneOf[byte, int](errMessage, data...)
}

// SequenceOf - expects a sequence of bytes in the buffer
// equal to the input data sequence. If expectations are not met,
// returns ParseError.
func SequenceOf(
	errMessage string,
	data ...byte,
) common.Combinator[byte, int, []byte] {
	return common.SequenceOf[byte, int](errMessage, data...)
}

// Map - Reads one element from the bytes buffer using the combinator,
// then uses the resulting element to obtain a value from the map cases and returns it.
// If the value is not found then it returns ParseError error.
func Map[K comparable, V any](
	errMessage string,
	cases map[K]V,
	c common.Combinator[byte, int, K],
) common.Combinator[byte, int, V] {
	return common.Map[byte, int, K, V](errMessage, cases, c)
}
