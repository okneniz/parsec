package bytes

import (
	p "github.com/okneniz/parsec/common"
)

// Range - succeeds for any bytes which include in input range.
// Returns the byte that is actually readed from input buffer.
// Greedy by default - keep position after reading.
func Range(from byte, to byte) p.Combinator[byte, int, byte] {
	return p.Range[byte, int](from, to)
}

// NotRange - succeeds for any byte which not included in input range.
// Returns the item that is actually readed from input buffer.
// Greedy by default - keep position after reading.
func NotRange(from byte, to byte) p.Combinator[byte, int, byte] {
	return p.NotRange[byte, int](from, to)
}

// Gt - succeeds for any byte which greater than input value.
// Returns the byte that is actually readed from input buffer.
// Greedy by default - keep position after reading.
func Gt(t byte) p.Combinator[byte, int, byte] {
	return Satisfy(true, func(x byte) bool {
		return x > t
	})
}

// Gte - succeeds for any byte which greater than or equal input value.
// Returns the byte that is actually readed from input buffer.
// Greedy by default - keep position after reading.
func Gte(t byte) p.Combinator[byte, int, byte] {
	return Satisfy(true, func(x byte) bool {
		return x >= t
	})
}

// Lt - succeeds for any byte which less than input value.
// Returns the byte that is actually readed from input buffer.
// Greedy by default - keep position after reading.
func Lt(t byte) p.Combinator[byte, int, byte] {
	return Satisfy(true, func(x byte) bool {
		return x < t
	})
}

// Lte - succeeds for any byte which less than or equal input byte.
// Returns the byte that is actually readed from input buffer.
// Greedy by default - keep position after reading.
func Lte(t byte) p.Combinator[byte, int, byte] {
	return Satisfy(true, func(x byte) bool {
		return x <= t
	})
}
