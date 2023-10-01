package bytes

import (
	p "github.com/okneniz/parsec/common"
)

// Eq - succeeds for any byte which equal input t.
// Returns the byte that is actually readed from input buffer.
// Greedy by default - keep position after reading.
func Eq(t byte) p.Combinator[byte, int, byte] {
	return p.Eq[byte, int](t)
}

// NotEq - succeeds for any byte which not equal input t.
// Returns the item that is actually readed from input buffer.
// Greedy by default - keep position after reading.
func NotEq(t byte) p.Combinator[byte, int, byte] {
	return p.NotEq[byte, int](t)
}

// OneOf - succeeds for any byte which included in input data.
// Returns the item that is actually readed from input buffer.
// Greedy by default - keep position after reading.
func OneOf(data ...byte) p.Combinator[byte, int, byte] {
	return p.OneOf[byte, int](data...)
}

// NoneOf - succeeds for any byte which not included in input data.
// Returns the item that is actually readed from input buffer.
// Greedy by default - keep position after reading.
func NoneOf(data ...byte) p.Combinator[byte, int, byte] {
	return p.NoneOf[byte, int](data...)
}

// SequenceOf - expects a sequence of bytes in the buffer
// equal to the input data sequence. If expectations are not met,
// returns NothingMatched error.
func SequenceOf(data ...byte) p.Combinator[byte, int, []byte] {
	return p.SequenceOf[byte, int](data...)
}

// Map - Reads one element from the bytes buffer using the combinator,
// then uses the resulting element to obtain a value from the map cases and returns it.
// If the value is not found then it returns NothingMatched error.
func Map[K comparable, V any](
	cases map[K]V,
	c p.Combinator[byte, int, K],
) p.Combinator[byte, int, V] {
	return p.Map[byte, int, K, V](cases, c)
}
